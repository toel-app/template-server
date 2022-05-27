package otp

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/toel-app/template-server/src/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	OTP_METHOD_SMS   = "SMS"
	OTP_METHOD_EMAIL = "EMAIL"
)

type Otp struct {
	Id        primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Addressee string             `json:"addressee" bson:"addressee"`
	Code      string             `json:"code" bson:"code"`
	Method    string             `json:"method" bson:"method"`
	ExpiredAt time.Time          `json:"expiredAt" bson:"expiredAt"`
	IsExpired bool               `json:"isExpired" bson:"isExpired"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type CreateOtp struct {
	Addressee string `json:"addressee" bson:"addressee" validate:"required"`
	Method    string `json:"method" bson:"method" validate:"required"`
}

type VerifyOtp struct {
	Addressee string `json:"addressee" bson:"addressee" validate:"required"`
	Method    string `json:"method" bson:"method" validate:"required"`
	Code      string `json:"code" bson:"code" validate:"required"`
}

type FindOtp struct {
	Addressee string `json:"addressee" bson:"addressee" validate:"required"`
	Method    string `json:"method" bson:"method" validate:"required"`
	Code      string `json:"code" bson:"code" validate:"required"`
}

func (model VerifyOtp) Validate() *utils.ErrResponse {
	err := validator.New().Struct(model)

	if err != nil {
		return utils.NewError(utils.ERR_INVALID_PAYLOAD)
	}

	isOtpValid := ValidateOtpMethod(model.Method)

	if !isOtpValid {
		return utils.NewError(utils.ERR_INVALID_OTP_METHOD)
	}

	return nil
}

func (model CreateOtp) Validate() *utils.ErrResponse {
	err := validator.New().Struct(model)

	if err != nil {
		return utils.NewError(utils.ERR_INVALID_PAYLOAD)
	}

	isOtpValid := ValidateOtpMethod(model.Method)

	if !isOtpValid {
		return utils.NewError(utils.ERR_INVALID_OTP_METHOD)
	}

	return nil
}
