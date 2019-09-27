package pdns_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// vironAuthType
// @Summary get auth type
// @Description get auth type
// @ID viron_authtype#get
// @Accept  json
// @Produce  json
// @Router /viron_authtype [get]
// @Tags viron
func vironAuthType(c echo.Context) error {
	encodedJSON := []byte(`
[
    {
      "type": "email",
      "provider": "viron-demo",
      "url": "/signin",
      "method": "POST",
    },
    {
      "type": "signout",
      "provider": "",
      "url": "",
      "method": "POST",
    },
]
`)
	return c.JSONBlob(http.StatusOK, encodedJSON)

}

//vironGlobalMenu
// @Summary get global menu
// @Description get global menu
// @ID viron#get
// @Accept json
// @Produce json
// @Router /viron [get]
// @Tags viron
func vironGlobalMenu(c echo.Context) error {
	encodedJSON := []byte(`{
  "theme": "standard",
  "color": "white",
  "name": "Viron example - local",
  "tags": [
    "domains",
    "records"
  ],
  "pages": [
    {
      "section": "manage",
      "id": "domains",
      "name": "Domains",
      "components": [
        {
          "api": {
            "method": "get",
            "path": "/domains"
          },
	  "query": [
	    { key: "id", type: "integer" },
	    { key: "name", type: "string" },
          ],
	  "primary": "id",
          "name": "Domain",
	  "style": "table",
          "pagination": true,
	  "table_labels": [
	    "id",
            "name",
            "type",
            "notified_serial",
            "last_check"
	  ]
        }
      ]
    },
    {
      "section": "manage",
      "id": "records",
      "name": "Records",
      "components": [
        {
          "api": {
            "method": "get",
            "path": "/records"
          },
	  "query": [
	    { key: "id", type: "integer" },
	    { key: "domain_id", type: "integer" },
	    { key: "name", type: "string" },
          ],
          "name": "Record",
	  "style": "table",
	  "primary": "id",
	  "table_labels": [
	    "id",
	    "domain_id",
            "name",
            "type",
            "content",
            "ttl",
            "prio",
            "disabled",
            "ordername",
            "auth"
	  ]
        }
      ]
    }
  ]
}`)
	return c.JSONBlob(http.StatusOK, encodedJSON)

}

func VironEndpoints(g *echo.Group) {
	g.GET("/viron", vironGlobalMenu)
	g.GET("/viron_authtype", vironAuthType)
}
