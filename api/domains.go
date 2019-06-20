package api

import (
	"net/http"

	"github.com/labstack/echo"
	_ "github.com/tredoe/osutil/user/crypt/md5_crypt"
	_ "github.com/tredoe/osutil/user/crypt/sha256_crypt"
)

func getDomains(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
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
