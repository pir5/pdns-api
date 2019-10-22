package pdns_api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/pir5/pdns-api/model"
	"gopkg.in/go-playground/validator.v9"
)

// getDomains is getting domains.
// @Summary get domains
// @Description get domains
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id query int false "Domain ID"
// @Param name query string false "Name"
// @Success 200 {array} model.Domain
// @Failure 404 {array} model.Domain
// @Failure 500 {object} pdns_api.HTTPError
// @Router /domains [get]
// @Tags domains
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

	ids, err := filterDomains(ds.ToIntreface(), c)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	return c.JSON(http.StatusOK, ids)
}

// updateDomainByName is update domain.
// @Summary update domain
// @Description update domain
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param name path string true "Dorain Name"
// @Param domain body model.Domain true "Domain Object"
// @Success 200 "OK"
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /domains/{name} [put]
// @Tags domains
func (h *domainHandler) updateDomainByName(c echo.Context) error {
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

// updateDomainByID is update domain.
// @Summary update domain
// @Description update domain
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id path integer true "Dorain ID"
// @Param domain body model.Domain true "Domain Object"
// @Success 200 {object} model.Domain
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /domains/{id} [put]
// @Tags domains
func (h *domainHandler) updateDomainByID(c echo.Context) error {
	if err := h.isAllowDomainByID(c); err != nil {
		return err
	}

	nd := &model.Domain{}
	if err := c.Bind(nd); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	updated, err := h.domainModel.UpdateByID(c.Param("id"), nd)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if !updated {
		return c.JSON(http.StatusNotFound, "domains does not exists")
	}
	return c.JSON(http.StatusOK, nil)
}

// deleteDomainByName is delete domain.
// @Summary delete domain
// @Description delete domain
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param name path string true "Domain Name"
// @Success 204 "No Content"
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /domains/{name} [delete]
// @Tags domains
func (h *domainHandler) deleteDomainByName(c echo.Context) error {
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

// deleteDomainByID is delete domain.
// @Summary delete domain
// @Description delete domain
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id path interger true "Domain ID"
// @Success 204 {object} model.Domain
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /domains/{id} [delete]
// @Tags domains
func (h *domainHandler) deleteDomainByID(c echo.Context) error {
	if err := h.isAllowDomainByID(c); err != nil {
		return err
	}

	deleted, err := h.domainModel.DeleteByID(c.Param("id"))
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
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param domain body model.Domain true "Domain Object"
// @Success 201 "Created"
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /domains [post]
// @Tags domains
func (h *domainHandler) createDomain(c echo.Context) error {
	d := &model.Domain{}
	if err := c.Bind(d); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := validator.New().Struct(d); err != nil {
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
func (h *domainHandler) isAllowDomainByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	whereParams := map[string]interface{}{
		"id": id,
	}

	ds, err := h.domainModel.FindBy(whereParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if ds == nil || len(ds) == 0 {
		return c.JSON(http.StatusNotFound, "domains does not exists")
	}

	err = isAllowDomain(c, ds[0].Name)
	if err != nil {
		return err
	}
	return nil
}

func getAllowDomains(c echo.Context) ([]string, error) {
	ds := c.Get(AllowDomainsKey)
	if ds != nil {
		return ds.([]string), nil
	}
	return nil, errors.New("allow domains not exists")
}

type domainHandler struct {
	domainModel model.DomainModeler
}

func NewDomainHandler(d model.DomainModeler) *domainHandler {
	return &domainHandler{
		domainModel: d,
	}
}

func filterDomains(ds []interface{}, c echo.Context) ([]interface{}, error) {
	ret := []interface{}{}
	if globalConfig.IsHTTPAuth() {
		domains, err := getAllowDomains(c)
		if err != nil {
			return nil, c.JSON(http.StatusForbidden, err)
		}
		for _, vv := range domains {
			for _, v := range ds {
				var name string
				switch v := v.(type) {
				case model.Domain:
					name = v.Name
				case model.Record:
					name = v.Domain.Name
				default:
					return nil, fmt.Errorf("unmatch type %s", v)
				}

				if strings.ToLower(name) == strings.ToLower(vv) {
					ret = append(ret, v)
				}
			}
		}
	} else {
		return ds, nil
	}
	return ret, nil
}
func DomainEndpoints(g *echo.Group, db *gorm.DB) {
	h := NewDomainHandler(model.NewDomainModeler(db))
	g.GET("/domains", h.getDomains)
	g.PUT("/domains/name/:name", h.updateDomainByName)
	g.DELETE("/domains/name/:name", h.deleteDomainByName)
	g.PUT("/domains/:id", h.updateDomainByID)
	g.DELETE("/domains/:id", h.deleteDomainByID)
	g.POST("/domains", h.createDomain)
}
