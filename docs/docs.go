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
        "/customer": {
            "post": {
                "description": "Set new customer",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Create customer",
                "parameters": [
                    {
                        "description": "application/json",
                        "name": "createCustomerForm",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.createCustomerForm"
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
                                            "$ref": "#/definitions/http.customerProfileDTO"
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
        },
        "/customers": {
            "get": {
                "description": "Paginate customers",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "List customers",
                "parameters": [
                    {
                        "type": "integer",
                        "example": 2,
                        "description": "Limit of pages",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "example": 1,
                        "description": "Current page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "created_at",
                        "description": "Sort results by a value",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "desc",
                        "description": "Order by ascendent o descendent",
                        "name": "direction",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/http.respOkData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/http.customerProfileDTO"
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
        },
        "/customers/{id}": {
            "delete": {
                "description": "Set new customer",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Create customer",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Customer id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/http.respOkData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/http.customerProfileDTO"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "$ref": "#/definitions/http.respError"
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
        },
        "/login": {
            "post": {
                "description": "User authentication",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "application/json",
                        "name": "userLoginForm",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.userLoginForm"
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
                                            "$ref": "#/definitions/http.dataTokenDTO"
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
                    "401": {
                        "description": "Unauthorized",
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
        },
        "/products": {
            "get": {
                "description": "Get a list of products",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "List products",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/http.respData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/http.productCardDTO"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.respError"
                        }
                    }
                }
            },
            "post": {
                "description": "Register product",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Add product",
                "parameters": [
                    {
                        "description": "application/json",
                        "name": "addProductForm",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.addProductForm"
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
                                            "$ref": "#/definitions/http.productCardDTO"
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
                    "401": {
                        "description": "Unauthorized",
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
        },
        "/users": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "List users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/http.respOkData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/http.userProfileDTO"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.respError"
                        }
                    }
                }
            },
            "post": {
                "description": "Register a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "SignUp user",
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
        },
        "/users/{id}": {
            "get": {
                "description": "Find user by its id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Find user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/http.respData"
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
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.respError"
                        }
                    }
                }
            },
            "put": {
                "description": "Update user by its id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Update user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "application/json",
                        "name": "userUpdateForm",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/http.userUpdateForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
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
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.respError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete user by its id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Delete user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/http.respOk"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.respError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.respError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "http.addProductForm": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "observations": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                }
            }
        },
        "http.createCustomerForm": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "http.customerProfileDTO": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "lastName": {
                    "type": "string"
                }
            }
        },
        "http.dataTokenDTO": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "http.productCardDTO": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "observations": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                }
            }
        },
        "http.respData": {
            "type": "object",
            "properties": {
                "data": {}
            }
        },
        "http.respError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "http.respOk": {
            "type": "object",
            "properties": {
                "ok": {
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
        "http.userLoginForm": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
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
        },
        "http.userUpdateForm": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "lorem@ipsum.com"
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
                    "example": "1234567a"
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
