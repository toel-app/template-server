package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/toel-app/common-utils/config"
	"github.com/toel-app/common-utils/logger"
	"github.com/toel-app/common-utils/redis"
	"github.com/toel-app/common-utils/slack"
	"github.com/toel-app/template-server/src/internal/otp"
	"github.com/toel-app/template-server/src/internal/ping"
	"github.com/toel-app/template-server/src/pkg/db"
	"github.com/toel-app/template-server/src/pkg/db/mysql"
	"github.com/toel-app/template-server/src/pkg/utils"
	"github.com/twilio/twilio-go"
)

var (
	router = gin.Default()
)

func StartApplication() {
	// Initialize
	logger.Initialize()
	config.Initialize()
	redis.Initialize()

	logger.Info("Application starting...")

	db.Connect()
	mysql.Connect()

	registerApiHandlers()
	router.Use(utils.Recovery())
	gin.SetMode(gin.ReleaseMode)
	if err := router.Run(); err != nil {
		logger.Error("Server not start", err)
	}
}

func registerApiHandlers() {
	// Setup dependency injection
	accountSid := viper.GetString("twilio.sid")
	authToken := viper.GetString("twilio.token")

	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	webHookUrl := viper.GetString("slackWebhookUrl")
	slackWebHook := slack.NewWebhook(webHookUrl)

	otpCollection := db.Database.Collection(db.COLLECTION_OTP)
	otpRepository := otp.NewRepository(otpCollection, mysql.Database)
	otpService := otp.NewService(otpRepository, slackWebHook, twilioClient)

	// Setup API handlers
	otp.RegisterRoute(router, otpService)
	ping.RegisterRoute(router)
}

func WrapWithGracefulShutdown(fn func()) {
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	go fn()
	<-cancelChan
	fmt.Println("SIGTERM detected")
	db.Close()
}
