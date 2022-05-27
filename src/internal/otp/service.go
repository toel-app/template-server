package otp

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/toel-app/common-utils/logger"
	"github.com/toel-app/common-utils/slack"
	"github.com/toel-app/template-server/src/pkg/utils"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type IService interface {
	Send(payload CreateOtp) *utils.ErrResponse
	VerifyOtp(payload VerifyOtp) *utils.ErrResponse
}

func NewService(repo IRepository, slackWebHook slack.WebHook, twilio *twilio.RestClient) IService {
	return service{repo, slackWebHook, twilio}
}

type service struct {
	repo         IRepository
	slackWebHook slack.WebHook
	twilioClient *twilio.RestClient
}

func (s service) Send(payload CreateOtp) *utils.ErrResponse {
	code, err := s.repo.Create(payload)
	if err != nil {
		return utils.NewError(err.Code)
	}

	go s.slackWebHook.SendOtpToSlack(payload.Addressee, code)
	params := &openapi.CreateMessageParams{}
	smsTemplate := fmt.Sprintf(viper.GetString("twilio.smsTemplate"), code)
	addressee := FormatPhoneNumberIfNotEmail(payload.Addressee)
	params.SetTo(addressee)
	params.SetFrom(viper.GetString("twilio.phoneNumber"))
	params.SetBody(smsTemplate)

	_, errSend := s.twilioClient.Api.CreateMessage(params)
	if errSend != nil {
		logger.Error("error send otp", errSend)
	}

	return nil
}

func (s service) VerifyOtp(payload VerifyOtp) *utils.ErrResponse {
	addressee := FormatPhoneNumberIfNotEmail(payload.Addressee)
	findOtpData := &FindOtp{
		Addressee: addressee,
		Method:    payload.Method,
		Code:      payload.Code,
	}

	existingOtp, err := s.repo.FindValidOtp(*findOtpData)
	if err != nil {
		return utils.NewError(err.Code)
	}

	if existingOtp == nil {
		return utils.NewError(utils.ERR_INVALID_OTP_CODE)
	}

	ok, err := s.repo.InvalidateOtpById(existingOtp.Id)
	if err != nil || !ok {
		return utils.NewError(utils.ERR_INTERNAL_SERVER_ERROR)
	}

	return nil
}
