// Package swagger_assets_http_api Code generated by swaggo/swag. DO NOT EDIT
package swagger_assets_http_api

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/assets/{id}": {
            "put": {
                "description": "Modify an asset",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "assets"
                ],
                "summary": "Modify an asset",
                "parameters": [
                    {
                        "type": "string",
                        "default": "00000000-0000-0000-0000-000000000000",
                        "description": "Asset ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Asset data",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/assetshttp.ModifyAssetRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new asset",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "assets"
                ],
                "summary": "Create a new asset",
                "parameters": [
                    {
                        "type": "string",
                        "default": "00000000-0000-0000-0000-000000000000",
                        "description": "Asset ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Asset data",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/assetshttp.CreateAssetRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete an asset",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "assets"
                ],
                "summary": "Delete an asset",
                "parameters": [
                    {
                        "type": "string",
                        "default": "00000000-0000-0000-0000-000000000000",
                        "description": "Asset ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "assetshttp.CreateAssetRequest": {
            "type": "object",
            "properties": {
                "assetMoneyAmount": {
                    "type": "number",
                    "example": 1000
                },
                "assetMoneyCurrency": {
                    "type": "string",
                    "example": "USD"
                },
                "assetName": {
                    "type": "string",
                    "example": "My Asset"
                },
                "assetType": {
                    "type": "string",
                    "example": "cash"
                }
            }
        },
        "assetshttp.ModifyAssetRequest": {
            "type": "object",
            "required": [
                "asset_money_amount",
                "asset_money_currency",
                "asset_name",
                "asset_type"
            ],
            "properties": {
                "asset_money_amount": {
                    "type": "number"
                },
                "asset_money_currency": {
                    "type": "string"
                },
                "asset_name": {
                    "type": "string"
                },
                "asset_type": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:6000",
	BasePath:         "/api/v1",
	Schemes:          []string{"http"},
	Title:            "Asset Management APIs",
	Description:      "This is the API for managing assets in the Finantrack application.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
