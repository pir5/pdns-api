// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-06-20 21:26:11.908522 +0900 JST m=+0.046083334

package docs

import (
	"bytes"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "This is PDNS RESTful API Server.",
        "title": "PDNS-API",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "host": "127.0.0.1:8080",
    "basePath": "/v1",
    "paths": {
        "/domains": {
            "get": {
                "description": "get domains",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "get domains",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Domain ID",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Name",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Domain"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/pdns_api.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/pdns_api.HTTPError"
                        }
                    }
                }
            },
            "post": {
                "description": "create domain",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "create domain",
                "parameters": [
                    {
                        "description": "Domain Object",
                        "name": "domain",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/model.Domain"
                        }
                    }
                ],
                "responses": {
                    "201": {},
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/pdns_api.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/pdns_api.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/pdns_api.HTTPError"
                        }
                    }
                }
            }
        },
        "/domains/{name}": {
            "delete": {
                "description": "delete domain",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "delete domain",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Domain Name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {},
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/pdns_api.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/pdns_api.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/pdns_api.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Domain": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "last_check": {
                    "type": "integer"
                },
                "master": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "notified_serial": {
                    "type": "integer"
                },
                "records": {
                    "type": "object",
                    "$ref": "#/definitions/model.Records"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "model.Record": {
            "type": "object",
            "properties": {
                "auth": {
                    "type": "boolean"
                },
                "content": {
                    "type": "string"
                },
                "disabled": {
                    "type": "boolean"
                },
                "domain_id": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "older_name": {
                    "type": "string"
                },
                "prio": {
                    "type": "integer"
                },
                "ttl": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "model.Records": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "auth": {
                        "type": "boolean"
                    },
                    "content": {
                        "type": "string"
                    },
                    "disabled": {
                        "type": "boolean"
                    },
                    "domain_id": {
                        "type": "integer"
                    },
                    "id": {
                        "type": "integer"
                    },
                    "name": {
                        "type": "string"
                    },
                    "older_name": {
                        "type": "string"
                    },
                    "prio": {
                        "type": "integer"
                    },
                    "ttl": {
                        "type": "integer"
                    },
                    "type": {
                        "type": "string"
                    }
                }
            }
        },
        "pdns_api.HTTPError": {
            "type": "object"
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
