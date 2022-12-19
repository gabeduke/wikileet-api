// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Wikileet Support",
            "url": "http://www.leetserve.com/support",
            "email": "support@leetserve.com"
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
        "/item/:id": {
            "get": {
                "description": "Get item by id.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "Get item",
                "parameters": [
                    {
                        "type": "string",
                        "default": "api_docs@leetserve.com",
                        "description": "user email",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "default",
                        "description": "user workspace",
                        "name": "X-Workspace",
                        "in": "header"
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
                "description": "Delete item by id.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "Delete item",
                "parameters": [
                    {
                        "type": "string",
                        "default": "api_docs@leetserve.com",
                        "description": "user email",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "default",
                        "description": "user workspace",
                        "name": "X-Workspace",
                        "in": "header"
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
            "patch": {
                "description": "Get item by id.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "Update item",
                "parameters": [
                    {
                        "type": "string",
                        "default": "api_docs@leetserve.com",
                        "description": "user email",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "default",
                        "description": "user workspace",
                        "name": "X-Workspace",
                        "in": "header"
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
        },
        "/items": {
            "get": {
                "description": "List items by user email.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "Get items",
                "parameters": [
                    {
                        "type": "string",
                        "format": "email",
                        "description": "search by email address",
                        "name": "user_email",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "api_docs@leetserve.com",
                        "description": "user email",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "default",
                        "description": "user workspace",
                        "name": "X-Workspace",
                        "in": "header",
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
            },
            "post": {
                "description": "Create a new item.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "Creat item",
                "parameters": [
                    {
                        "type": "string",
                        "format": "email",
                        "description": "associate item with user",
                        "name": "user_email",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "format": "email",
                        "description": "associate item with workspace",
                        "name": "workspace",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "api_docs@leetserve.com",
                        "description": "user email",
                        "name": "X-User",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "default",
                        "description": "user workspace",
                        "name": "X-Workspace",
                        "in": "header",
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
    "securityDefinitions": {
        "OAuth2Application": {
            "type": "oauth2",
            "flow": "application",
            "tokenUrl": "https://oauth2.googleapis.com/token",
            "scopes": {
                "https://www.googleapis.com/auth/userinfo.email": "\tSee your primary Google Account email address"
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "http://localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Wikileet API",
	Description:      "Wikileet gift lists API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
