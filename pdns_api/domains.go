package pdns_api

import (
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
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
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /domains [get]
func getDomains(c echo.Context) error {
	db, err := getDBConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	d := model.Domains{}
	query := db.Debug().
		Preload("Records")

	for k, v := range c.QueryParams() {
		if k == "name" || k == "id" {
			query = query.Where(k+" = ?", v)
		}
	}

	r := query.Find(&d)
	if r.Error != nil {
		return r.Error
	}

	if len(d) == 0 {
		return c.JSON(http.StatusNotFound, "domains does not exists")
	}

	if conf.IsHTTPAuth() {
		domains := c.Get(AllowDomainsKey).([]string)
		d.FilterOwnerDomains(domains)
	}
	return c.JSON(http.StatusNoContent, d)
}

// updateDomain is update domain.
// @Summary update domain
// @Description update domain
// @Accept  json
// @Produce  json
// @Param name path string true "Domain Name"
// @Param domain body model.Domain true "Domain Object"
// @Success 200
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /domains/{name} [update]
func updateDomain(c echo.Context) error {
	db, err := prepare(c, c.Param("name"))
	if err != nil {
		return err
	}

	d := model.Domain{}
	r := db.Where("name = ?", c.Param("name")).Take(&d)
	if r.Error != nil {
		if r.RecordNotFound() {
			return c.JSON(http.StatusNotFound, "domains does not exists")
		} else {
			return c.JSON(http.StatusInternalServerError, r.Error)
		}
	}

	nd := &model.Domain{}
	if err := c.Bind(nd); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	r = db.Model(&d).Updates(&nd)
	if r.Error != nil {
		return c.JSON(http.StatusInternalServerError, r.Error)
	}
	return c.JSON(http.StatusOK, nil)
}

// deleteDomain is delete domain.
// @Summary delete domain
// @Description delete domain
// @Accept  json
// @Produce  json
// @Param name path string true "Domain Name"
// @Success 204
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /domains/{name} [delete]
func deleteDomain(c echo.Context) error {
	db, err := prepare(c, c.Param("name"))
	if err != nil {
		return err
	}

	d := model.Domain{}
	r := db.Where("name = ?", c.Param("name")).Take(&d)
	if r.Error != nil {
		if r.RecordNotFound() {
			return c.JSON(http.StatusNotFound, "domains does not exists")
		} else {
			return c.JSON(http.StatusInternalServerError, r.Error)
		}
	}

	r = db.Delete(d)
	if r.Error != nil {
		return c.JSON(http.StatusInternalServerError, r.Error)
	}

	return c.JSON(http.StatusOK, nil)
}

// createDomain is create domain.
// @Summary create domain
// @Description create domain
// @Accept  json
// @Produce  json
// @Param domain body model.Domain true "Domain Object"
// @Success 201
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /domains [post]
func createDomain(c echo.Context) error {
	db, err := prepare(c, "")
	if err != nil {
		return err
	}

	d := &model.Domain{}
	if err := c.Bind(d); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := db.Create(&d).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, nil)
}

func prepare(c echo.Context, name string) (*gorm.DB, error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, c.JSON(http.StatusInternalServerError, err)
	}

	if name != "" && !allowDomain(c, name) {
		return nil, c.JSON(http.StatusForbidden, nil)
	}
	return db, nil
}

func allowDomain(c echo.Context, name string) bool {
	if !conf.IsHTTPAuth() {
		return true
	}
	domains := c.Get(AllowDomainsKey).([]string)
	for _, vv := range domains {
		if strings.ToLower(name) == strings.ToLower(vv) {
			return true
		}
	}
	return false
}

func DomainEndpoints(g *echo.Group) {
	g.GET("/domains", getDomains)
	g.PUT("/domains/:name", updateDomain)
	g.DELETE("/domains/:name", deleteDomain)
	g.POST("/domains", createDomain)
}
