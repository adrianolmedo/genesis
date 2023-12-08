// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
            "name": "Adrián Olmedo",
            "url": "https://twitter.com/adrianolmedo",
            "email": "adrianolmedo.ve@gmail.com"
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
        "/users": {
            "post": {
                "description": "Register a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "SignUp User",
                "parameters": [
                    {
                        "description": "application/json",
                        "name": "userSignUpForm",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.userSignUpForm"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/http.respOkData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/http.userProfileDTO"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.respError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.respError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "http.respError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "http.respOkData": {
            "type": "object",
            "properties": {
                "data": {},
                "ok": {
                    "type": "string"
                }
            }
        },
        "http.userProfileDTO": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "johndoe@aol.com"
                },
                "firstName": {
                    "type": "string",
                    "example": "John"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "lastName": {
                    "type": "string",
                    "example": "Doe"
                }
            }
        },
        "http.userSignUpForm": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "johndoe@aol.com"
                },
                "firstName": {
                    "type": "string",
                    "example": "John"
                },
                "lastName": {
                    "type": "string",
                    "example": "Doe"
                },
                "password": {
                    "type": "string",
                    "example": "1234567b"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/v1/",
	Schemes:          []string{},
	Title:            "Genesis REST API",
	Description:      "This is a sample server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
