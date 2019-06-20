package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pir5/pdns-api/model"
)

// getDomains is getting domains.
// @Summary get domains
// @Description get domains
// @Accept  json
// @Produce  json
// @Param id query int false "Domain ID"
// @Param name query string false "Name"
// @Success 200 {array} model.Domain
// @Failure 500 {object} pdns_api.HTTPError
// @Router /domains [get]
func getDomains(c echo.Context) error {
	return c.JSON(http.StatusOK, model.Domains{model.Domain{
		ID: 1,
	}})
}
func updateDomain(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
func deleteDomain(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
func createDomain(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func DomainEndpoints(g *echo.Group) {
	g.GET("/domains", getDomains)
	g.PUT("/domains/:name", updateDomain)
	g.DELETE("/domains/:name", deleteDomain)
	g.POST("/domains", createDomain)
}
