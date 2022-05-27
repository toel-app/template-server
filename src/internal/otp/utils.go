package otp

import (
	"crypto/rand"
	"io"

	"github.com/rbrahul/gofp"
	"github.com/spf13/viper"
	"github.com/toel-app/common-utils/string_utils"
)

func GenerateOtpCode() string {
	otpCodeLength := viper.GetInt("otp_code_length")
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, otpCodeLength)
	n, err := io.ReadAtLeast(rand.Reader, b, otpCodeLength)
	if n != otpCodeLength {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func ValidateOtpMethod(method string) (isValid bool) {
	validOtpMethods := []interface{}{
		OTP_METHOD_EMAIL,
		OTP_METHOD_SMS,
	}

	return gofp.Contains(validOtpMethods, method)
}

func FormatPhoneNumberIfNotEmail(phoneOrEmail string) string {
	_, isEmail := string_utils.ValidMailAddress(phoneOrEmail)
	if isEmail {
		return phoneOrEmail
	} else {
		return string_utils.FormatPhoneNumber(phoneOrEmail)
	}
}
