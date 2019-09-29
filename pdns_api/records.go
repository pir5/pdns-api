package pdns_api

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/pir5/pdns-api/model"
)

// getRecords is getting records.
// @Summary get records
// @Description get records
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id query int false "Record ID"
// @Param domain_id query int false "Domain ID"
// @Success 200 {array} model.Record
// @Failure 404 {array} model.Record
// @Failure 500 {object} pdns_api.HTTPError
// @Router /records [get]
// @Tags records
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

	ids, err := filterDomains(ds.ToIntreface(), c)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	return c.JSON(http.StatusOK, ids)
}

// updateRecord is update record.
// @Summary update record
// @Description update record
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id path int true "Record ID "
// @Param record body model.Record true "Record Object"
// @Success 200 {object} model.Record
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /records/{id} [put]
// @Tags records
func (h *recordHandler) updateRecord(c echo.Context) error {
	if err := h.isAllowRecordID(c); err != nil {
		return err
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

// enableRecord is enable record.
// @Summary enable record
// @Description enable record
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id path int true "Record ID "
// @Param record body model.Record true "Record Object"
// @Success 200 {object} model.Record
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /records/enable/{id} [put]
// @Tags records
func (h *recordHandler) enableRecord(c echo.Context) error {
	return changeState(h, c, false)
}

// disableRecord is disable record.
// @Summary disable record
// @Description disable record
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id path int true "Record ID "
// @Param record body model.Record true "Record Object"
// @Success 200 {object} model.Record
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /records/disable/{id} [put]
// @Tags records
func (h *recordHandler) disableRecord(c echo.Context) error {
	return changeState(h, c, true)
}

func changeState(h *recordHandler, c echo.Context, disabled bool) error {
	if err := h.isAllowRecordID(c); err != nil {
		return err
	}

	nd := &model.Record{
		Disabled: &disabled,
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
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id path int true "Record ID "
// @Success 204 {object} model.Record
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /records/{id} [delete]
// @Tags records
func (h *recordHandler) deleteRecord(c echo.Context) error {
	err := h.isAllowRecordID(c)
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

	return c.NoContent(http.StatusNoContent)
}

// createRecord is create record.
// @Summary create record
// @Description create record
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param record body model.Record true "Record Object"
// @Success 201 {object} model.Record
// @Failure 403 {object} pdns_api.HTTPError
// @Failure 404 {object} pdns_api.HTTPError
// @Failure 500 {object} pdns_api.HTTPError
// @Router /records [post]
// @Tags records
func (h *recordHandler) createRecord(c echo.Context) error {
	d := &model.Record{}
	if err := c.Bind(d); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := h.isAllowDomainID(c, d.DomainID)
	if err != nil {
		return err
	}

	if err := h.recordModel.Create(d); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, d)
}

func (h *recordHandler) isAllowRecordID(c echo.Context) error {
	if !globalConfig.IsHTTPAuth() {
		return nil
	}

	ds, err := h.recordModel.FindBy(map[string]interface{}{
		"id": []string{c.Param("id")},
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if ds == nil || len(ds) == 0 {
		return c.JSON(http.StatusNotFound, "records does not exists")
	}

	return h.isAllowDomainID(c, ds[0].DomainID)
}

func (h *recordHandler) isAllowDomainID(c echo.Context, domainID int) error {
	if !globalConfig.IsHTTPAuth() {
		return nil
	}

	ds, err := h.domainModel.FindBy(map[string]interface{}{
		"id": domainID,
	})

	if err != nil {
		return c.JSON(http.StatusForbidden, nil)
	}

	if ds == nil || len(ds) == 0 {
		return c.JSON(http.StatusNotFound, "domains does not exists")
	}

	return isAllowDomain(c, ds[0].Name)
}

type recordHandler struct {
	recordModel model.RecordModeler
	domainModel model.DomainModeler
}

func NewRecordHandler(r model.RecordModeler, d model.DomainModeler) *recordHandler {
	return &recordHandler{
		recordModel: r,
		domainModel: d,
	}
}
func RecordEndpoints(g *echo.Group, db *gorm.DB) {
	h := NewRecordHandler(
		model.NewRecordModeler(db),
		model.NewDomainModeler(db),
	)
	g.GET("/records", h.getRecords)
	g.PUT("/records/:id", h.updateRecord)
	g.PUT("/records/enable/:id", h.enableRecord)
	g.PUT("/records/disable/:id", h.disableRecord)
	g.DELETE("/records/:id", h.deleteRecord)
	g.POST("/records", h.createRecord)
}
