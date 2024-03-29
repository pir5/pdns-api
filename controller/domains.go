package controller

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/pir5/pdns-api/model"
)

// getDomains is getting domains.
// @Summary get domains
// @Description get domains
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id query integer false "Domain ID"
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

	ds, total, err := h.domainModel.FindBy(c.Request(), whereParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	setPaginationHeader(c.Response().Writer, c.Request(), total)
	return c.JSON(http.StatusOK, ds)
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
	nd := &model.Domain{}
	if err := c.Bind(nd); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := validate.Struct(nd); err != nil {
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
// @Param id path int true "Dorain ID"
// @Param domain body model.Domain true "Domain Object"
// @Success 200 {object} model.Domain
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /domains/{id} [put]
// @Tags domains
func (h *domainHandler) updateDomainByID(c echo.Context) error {
	nd := &model.Domain{}
	if err := c.Bind(nd); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := validate.Struct(nd); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
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
// @Param id path integer true "Domain ID"
// @Success 204 {object} model.Domain
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /domains/{id} [delete]
// @Tags domains
func (h *domainHandler) deleteDomainByID(c echo.Context) error {
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

	if err := validate.Struct(d); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.domainModel.Create(d); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, nil)
}

type domainHandler struct {
	domainModel model.DomainModeler
}

func NewDomainHandler(d model.DomainModeler) *domainHandler {
	return &domainHandler{
		domainModel: d,
	}
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
