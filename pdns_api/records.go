package pdns_api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/pir5/pdns-api/model"
)

// getRecords is getting records.
// @Summary get records
// @Description get records
// @Accept  json
// @Produce  json
// @Param id query int false "Record ID"
// @Param domain_id query int false "Domain ID"
// @Success 200 {array} model.Record
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /records [get]
func (h *recordHandler) getRecords(c echo.Context) error {
	whereParams := map[string]interface{}{}
	for k, v := range c.QueryParams() {
		if k != "id" && k != "domain_id" {
			return c.JSON(http.StatusForbidden, nil)
		}
		whereParams[k] = v
	}

	ds, err := h.recordModel.FindBy(whereParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if globalConfig.IsHTTPAuth() {
		ret := model.Records{}
		records, err := getAllowDomains(c)
		if err != nil {
			return c.JSON(http.StatusForbidden, err)
		}

		for _, vv := range records {
			for _, v := range ds {
				if strings.ToLower(v.Domain.Name) == strings.ToLower(vv) {
					ret = append(ret, v)
				}
			}
		}
		ds = ret
	}

	if ds == nil || len(ds) == 0 {
		return c.JSON(http.StatusNotFound, "records does not exists")
	}

	return c.JSON(http.StatusOK, ds)
}

// updateRecord is update record.
// @Summary update record
// @Description update record
// @Accept  json
// @Produce  json
// @Param id path int true "Record ID "
// @Param record body model.Record true "Record Object"
// @Success 200
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /records/{id} [update]
func (h *recordHandler) updateRecord(c echo.Context) error {
	if err := h.isAllowRecord(c); err != nil {
		return nil
	}

	nd := &model.Record{}
	if err := c.Bind(nd); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	updated, err := h.recordModel.UpdateByID(c.Param("id"), nd)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if !updated {
		return c.JSON(http.StatusNotFound, "records does not exists")
	}
	return c.JSON(http.StatusOK, nil)
}

// deleteRecord is delete record.
// @Summary delete record
// @Description delete record
// @Accept  json
// @Produce  json
// @Param id path int true "Record ID "
// @Success 204
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /records/{name} [delete]
func (h *recordHandler) deleteRecord(c echo.Context) error {
	err := h.isAllowRecord(c)
	if err != nil {
		return err
	}

	deleted, err := h.recordModel.DeleteByID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if !deleted {
		return c.JSON(http.StatusNotFound, "records does not exists")
	}

	return c.JSON(http.StatusNoContent, nil)
}

// createRecord is create record.
// @Summary create record
// @Description create record
// @Accept  json
// @Produce  json
// @Param record body model.Record true "Record Object"
// @Success 201
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /records [post]
func (h *recordHandler) createRecord(c echo.Context) error {
	d := &model.Record{}
	if err := c.Bind(d); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := h.isAllowRecordByDomainID(c, d.DomainID)
	if err != nil {
		return err
	}

	if err := h.recordModel.Create(d); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, nil)
}

func (h *recordHandler) isAllowRecord(c echo.Context) error {
	if !globalConfig.IsHTTPAuth() {
		return nil
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ds, err := h.recordModel.FindBy(map[string]interface{}{
		"id": id,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if ds == nil || len(ds) == 0 {
		return c.JSON(http.StatusNotFound, "records does not exists")
	}

	err = isAllowDomain(c, ds[0].Name)
	if err != nil {
		return err
	}
	return nil
}

func (h *recordHandler) isAllowRecordByDomainID(c echo.Context, domainID int) error {
	if !globalConfig.IsHTTPAuth() {
		return nil
	}

	ds, err := h.domainModel.FindBy(map[string]interface{}{
		"id": domainID,
	})

	err = isAllowDomain(c, ds[0].Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusForbidden, nil)
}

type recordHandler struct {
	recordModel model.RecordModel
	domainModel model.DomainModel
}

func NewRecordHandler(r model.RecordModel, d model.DomainModel) *recordHandler {
	return &recordHandler{
		recordModel: r,
		domainModel: d,
	}
}
func RecordEndpoints(g *echo.Group, db *gorm.DB) {
	h := NewRecordHandler(
		model.NewRecordModel(db),
		model.NewDomainModel(db),
	)
	g.GET("/records", h.getRecords)
	g.PUT("/records/:id", h.updateRecord)
	g.DELETE("/records/:id", h.deleteRecord)
	g.POST("/records", h.createRecord)
}
