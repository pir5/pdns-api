package middleware

import (
	"github.com/labstack/echo/v4"
	em "github.com/labstack/echo/v4/middleware"
)

type (
	HeaderAuthConfig struct {
		Skipper em.Skipper

		Validator HeaderAuthValidator
	}

	HeaderAuthValidator func(string, string, echo.Context) (bool, error)
)

const (
	basic        = "basic"
	defaultRealm = "Restricted"
)

var (
	DefaultHeaderAuthConfig = HeaderAuthConfig{
		Skipper: em.DefaultSkipper,
	}
)

func HeaderAuth(fn HeaderAuthValidator) echo.MiddlewareFunc {
	c := DefaultHeaderAuthConfig
	c.Validator = fn
	return HeaderAuthWithConfig(c)
}

const PIR5HeaderID = "PIR5-ID"
const PIR5HeaderSecret = "PIR5-SECRET"

func HeaderAuthWithConfig(config HeaderAuthConfig) echo.MiddlewareFunc {
	if config.Validator == nil {
		panic("echo: header-auth middleware requires a validator function")
	}
	if config.Skipper == nil {
		config.Skipper = DefaultHeaderAuthConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			id := c.Request().Header.Get(PIR5HeaderID)
			secret := c.Request().Header.Get(PIR5HeaderSecret)

			valid, err := config.Validator(id, secret, c)
			if err != nil {
				return err
			} else if valid {
				return next(c)
			}

			return echo.ErrUnauthorized
		}
	}
}
