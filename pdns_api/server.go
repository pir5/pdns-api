package pdns_api

import (
	"context"
	"errors"
	"fmt"
	stdLog "log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/facebookgo/pidfile"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pir5/pdns-api/api"
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
	conf, err := NewConfig(*cmdFlags.ConfPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	logger := log.New("pdns-api")
	pidfile.SetPidfilePath(*cmdFlags.PidPath)
	if *cmdFlags.LogPath != "" {
		f, err := os.OpenFile(*cmdFlags.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return errors.New("error opening file :" + err.Error())
		}
		logger.SetOutput(f)
	} else {
		logger.SetLevel(log.DEBUG)
	}

	e := echo.New()
	e.Logger = logger
	e.StdLogger = stdLog.New(e.Logger.Output(), e.Logger.Prefix()+": ", 0)

	e.GET("/status", status)

	if err := pidfile.Write(); err != nil {
		return err
	}

	defer func() {
		if err := os.Remove(pidfile.GetPidfilePath()); err != nil {
			e.Logger.Fatalf("Error removing %s: %s", pidfile.GetPidfilePath(), err)
		}
	}()

	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status}}` + "\n",
		Output: logger.Output(),
	}))

	go func() {
		if err := e.Start(conf.Listen); err != nil {
			e.Logger.Fatalf("shutting down the server: %s", err)
		}
	}()

	v1 := e.Group("/v1")
	api.DomainEndpoints(v1)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello! STNS!!1")
	})

	quit := make(chan os.Signal)
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
