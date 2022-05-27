package otp

import (
	"context"
	"database/sql"
	"time"

	"github.com/toel-app/common-utils/logger"
	"github.com/toel-app/template-server/src/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRepository interface {
	Create(data CreateOtp) (code string, errCreate *utils.ErrResponse)
	FindValidOtp(data FindOtp) (*Otp, *utils.ErrResponse)
	InvalidateOtpById(id primitive.ObjectID) (ok bool, err *utils.ErrResponse)
}

type repository struct {
	collection *mongo.Collection
	sql        *sql.DB
}

func NewRepository(collection *mongo.Collection, sql *sql.DB) IRepository {
	return &repository{collection, sql}
}

func (r *repository) Create(data CreateOtp) (code string, errCreate *utils.ErrResponse) {
	timeNow := time.Now()
	expiredAt := timeNow.Add(time.Minute * OTP_EXPIRATION_IN_MINUTES)
	otpCode := GenerateOtpCode()

	if _, err := r.collection.InsertOne(context.TODO(), Otp{
		Addressee: data.Addressee,
		Method:    data.Method,
		Code:      otpCode,
		IsExpired: false,
		ExpiredAt: expiredAt,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}); err != nil {
		logger.Error("Error create OTP", err)
		return "", utils.NewError(utils.ERR_INTERNAL_SERVER_ERROR)
	}

	q := "INSERT INTO otps(address, code) VALUES (?, ?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.sql.PrepareContext(ctx, q)
	if err != nil {
		logger.Error("Error %s when preparing SQL statement", err)
		return "", utils.NewError(utils.ERR_INTERNAL_SERVER_ERROR)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, data.Addressee, otpCode)
	if err != nil {
		logger.Error("Error %s when inserting row into products table", err)
		return "", utils.NewError(utils.ERR_INTERNAL_SERVER_ERROR)
	}
	_, err = res.RowsAffected()
	if err != nil {
		logger.Error("Error %s when finding rows affected", err)
		return "", utils.NewError(utils.ERR_INTERNAL_SERVER_ERROR)
	}

	return otpCode, nil
}

func (r *repository) FindValidOtp(data FindOtp) (*Otp, *utils.ErrResponse) {
	var otp Otp

	condition := bson.D{
		primitive.E{
			Key:   "addressee",
			Value: data.Addressee,
		},
		primitive.E{
			Key:   "code",
			Value: data.Code,
		},
		primitive.E{
			Key:   "method",
			Value: data.Method,
		},
		primitive.E{
			Key:   "isExpired",
			Value: false,
		},
		primitive.E{
			Key: "expiredAt",
			Value: bson.M{
				"$gte": time.Now(),
			},
		},
	}

	err := r.collection.FindOne(context.TODO(), condition).Decode(&otp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		logger.Error("error find valid otp", err)
		return nil, utils.NewError(utils.ERR_INTERNAL_SERVER_ERROR)
	}

	return &otp, nil
}

func (r *repository) InvalidateOtpById(id primitive.ObjectID) (ok bool, err *utils.ErrResponse) {
	updateOperator := bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{
					Key:   "isExpired",
					Value: true,
				},
			},
		},
	}
	_, errorUpdate := r.collection.UpdateByID(context.TODO(), id, updateOperator)

	if errorUpdate != nil {
		return false, utils.NewError(utils.ERR_INTERNAL_SERVER_ERROR)
	}

	return true, nil
}
