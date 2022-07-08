package cmd

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"whatsappproxy/aliyunoss"
	"whatsappproxy/discovery"
	"whatsappproxy/rabbitmq"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"os"
	"whatsappproxy/config"
	"whatsappproxy/middleware"
	"whatsappproxy/routers"
	"whatsappproxy/services"
	"whatsappproxy/utils"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short: "Whatsapp API",
	Long:  `you can send whatsapp over http api but your whatsapp account have to be multi device version`,
	Run:   runRest,
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVarP(&config.AppPort, "port", "p", config.AppPort, "change port number with --port <number> | example: --port=8080")
	rootCmd.PersistentFlags().BoolVarP(&config.AppDebug, "debug", "d", config.AppDebug, "hide or displaying log with --debug <true/false> | example: --debug=true")
	rootCmd.PersistentFlags().StringVarP(&config.WhatsappAutoReplyMessage, "autoreply", "", config.WhatsappAutoReplyMessage, `auto reply when received message --autoreply <string>`)
}

func runRest(cmd *cobra.Command, args []string) {
	if config.AppDebug {
		config.WhatsappLogLevel = "DEBUG"
	}

	//preparing folder if not exist
	err := utils.CreateFolder(config.PathQrCode, config.PathSendItems)
	if err != nil {
		logger.Error(err)
	}

	app := fiber.New()

	app.Use(middleware.Recovery())
	if config.AppDebug {
		app.Use(logger.GetLogger())
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// register service
	discovery.InitNacos()

	// init apollo
	//utils.ApolloClient = config.InitApollo()

	// init rabbitMq
	imRabbitMq := rabbitmq.InitBossRabbitMq()
	imRabbitMq.MqConnect()

	//value := utils.ApolloClient.GetConfig("DevCenter.atta-rabbitmq")
	//logger.Info(value.GetValue("spring.rabbitmq.host"))

	//init aliyunOss
	utils.Bucket = aliyunoss.InitOssClient()

	//init db
	db := utils.InitWaDB()

	// Service
	appService := services.NewAppService(db)
	sendService := services.NewSendService(db)
	userService := services.NewUserService(db)

	// Controller
	appController := routers.NewAppController(appService)
	sendController := routers.NewSendController(sendService)
	userController := routers.NewUserController(userService)

	appController.Route(app)
	sendController.Route(app)
	userController.Route(app)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{
			"AppHost":      fmt.Sprintf("%s://%s", ctx.Protocol(), ctx.Hostname()),
			"AppVersion":   config.AppVersion,
			"MaxFileSize":  humanize.Bytes(uint64(config.WhatsappSettingMaxFileSize)),
			"MaxVideoSize": humanize.Bytes(uint64(config.WhatsappSettingMaxVideoSize)),
		})
	})

	err = app.Listen(":" + config.AppPort)
	logger.Infof("xxx")
	if err != nil {
		logger.Errorf("Failed to start: ", err.Error())
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
