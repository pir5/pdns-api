package cmd

import (
	"context"
	"fmt"
	stdLog "log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pir5/pdns-api/controller"
	"github.com/pir5/pdns-api/docs"
	"github.com/pir5/pdns-api/pdns_api"
)

var CmdServer = &Command{
	Run:       runServer,
	UsageLine: "server",
	Short:     "Start API Server",
	Long: `
Start API Server
	`,
}

func init() {
	// Set your flag here like below.
	// cmdServer.Flag.BoolVar(&flagA, "a", false, "")

}

// runServer executes sub command and return exit code.
func runServer(cmdFlags *GlobalFlags, args []string) error {
	config, err := pdns_api.NewConfig(*cmdFlags.ConfPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	logger := log.New("pdns-api")
	e := echo.New()
	e.Logger = logger
	e.StdLogger = stdLog.New(e.Logger.Output(), e.Logger.Prefix()+": ", 0)

	e.GET("/status", status)

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status}}` + "\n",
		Output: logger.Output(),
	}))

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Validator: func(key string, c echo.Context) (bool, error) {
			if config.IsTokenAuth() {
				for _, v := range config.Auth.Tokens {
					if key == v {
						return true, nil
					}
				}
			}
			return true, nil
		},
		Skipper: func(c echo.Context) bool {
			return !config.IsTokenAuth()
		},
	}))

	v1 := e.Group("/v1")

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.DB.UserName,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.DBName,
	))
	if err != nil {
		return err
	}
	db.LogMode(os.Getenv("GORM_LOG") == "true")
	db.SetLogger(logger)
	controller.DomainEndpoints(v1, db.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false))
	controller.RecordEndpoints(v1, db.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false))
	controller.VironEndpoints(v1)
	v1.GET("/swagger/*", echoSwagger.WrapHandler)

	docs.SwaggerInfo.Host = config.Listen
	if config.Endpoint != "" {
		u, err := url.Parse(config.Endpoint)
		if err != nil {
			return err
		}
		docs.SwaggerInfo.Schemes = []string{u.Scheme}
		docs.SwaggerInfo.Host = u.Host
		docs.SwaggerInfo.BasePath = u.Path
	}

	go func() {
		if err := e.Start(config.Listen); err != nil {
			e.Logger.Fatalf("shutting down the server: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
func status(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
