package pdns_api

import (
	"errors"
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
func (h *domainHandler) getDomains(c echo.Context) error {
	whereParams := map[string]interface{}{}
	for k, v := range c.QueryParams() {
		if k != "id" && k != "name" {
			return c.JSON(http.StatusForbidden, nil)
		}
		whereParams[k] = v
	}

	ds, err := h.domainModel.FindBy(whereParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if globalConfig.IsHTTPAuth() {
		ret := model.Domains{}
		domains, err := getAllowDomains(c)
		if err != nil {
			return c.JSON(http.StatusForbidden, err)
		}
		for _, vv := range domains {
			for _, v := range ds {
				if strings.ToLower(v.Name) == strings.ToLower(vv) {
					ret = append(ret, v)
				}
			}
		}
		ds = ret
	}

	if ds == nil || len(ds) == 0 {
		return c.JSON(http.StatusNotFound, "domains does not exists")
	}

	return c.JSON(http.StatusOK, ds)
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
func (h *domainHandler) updateDomain(c echo.Context) error {
	err := isAllowDomain(c, c.Param("name"))
	if err != nil {
		return err
	}

	nd := &model.Domain{}
	if err := c.Bind(nd); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	updated, err := h.domainModel.UpdateByName(c.Param("name"), nd)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if !updated {
		return c.JSON(http.StatusNotFound, "domains does not exists")
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
func (h *domainHandler) deleteDomain(c echo.Context) error {
	err := isAllowDomain(c, c.Param("name"))
	if err != nil {
		return err
	}

	deleted, err := h.domainModel.DeleteByName(c.Param("name"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if !deleted {
		return c.JSON(http.StatusNotFound, "domains does not exists")
	}

	return c.NoContent(http.StatusNoContent)
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
func (h *domainHandler) createDomain(c echo.Context) error {
	d := &model.Domain{}
	if err := c.Bind(d); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := isAllowDomain(c, d.Name)
	if err != nil {
		return err
	}

	if err := h.domainModel.Create(d); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, nil)
}

func isAllowDomain(c echo.Context, name string) error {
	if !globalConfig.IsHTTPAuth() {
		return nil
	}

	domains, err := getAllowDomains(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, nil)
	}

	for _, vv := range domains {
		if strings.ToLower(name) == strings.ToLower(vv) {
			return nil
		}
	}
	return c.JSON(http.StatusForbidden, nil)
}

func getAllowDomains(c echo.Context) ([]string, error) {
	ds := c.Get(AllowDomainsKey)
	if ds != nil {
		return ds.([]string), nil
	}
	return nil, errors.New("allow domains not exists")
}

type domainHandler struct {
	domainModel model.DomainModel
}

func NewDomainHandler(d model.DomainModel) *domainHandler {
	return &domainHandler{
		domainModel: d,
	}
}
func DomainEndpoints(g *echo.Group, db *gorm.DB) {
	h := NewDomainHandler(model.NewDomainModel(db))
	g.GET("/domains", h.getDomains)
	g.PUT("/domains/:name", h.updateDomain)
	g.DELETE("/domains/:name", h.deleteDomain)
	g.POST("/domains", h.createDomain)
}
