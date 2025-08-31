// Package docs Code generated manually. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "version": "{{.Version}}",
        "contact": {}
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "schemes": {{ marshal .Schemes }},
    "paths": {
        "/users": {
            "get": {
                "tags": ["users"],
                "summary": "Get users",
                "description": "Get paginated list of users",
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "page",
                        "in": "query",
                        "description": "Page number",
                        "required": false,
                        "type": "integer"
                    },
                    {
                        "name": "limit",
                        "in": "query",
                        "description": "Page size",
                        "required": false,
                        "type": "integer"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": { "$ref": "#/definitions/UsersResponse" }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": { "$ref": "#/definitions/MessageResponse" }
                    }
                }
            }
        },
        "/users/{mobile}": {
            "get": {
                "tags": ["users"],
                "summary": "Get user by mobile",
                "description": "Get user info by mobile number",
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "mobile",
                        "in": "path",
                        "description": "Mobile number",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": { "$ref": "#/definitions/UserResponse" }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": { "$ref": "#/definitions/MessageResponse" }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "tags": ["auth"],
                "summary": "Register user",
                "description": "Register a user with mobile number",
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "description": "Register payload",
                        "required": true,
                        "schema": { "$ref": "#/definitions/RegisterRequest" }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User registered successfully",
                        "schema": { "$ref": "#/definitions/RegisterResponse" }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": { "$ref": "#/definitions/MessageResponse" }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "tags": ["auth"],
                "summary": "Login user",
                "description": "Login user using mobile + OTP",
                "produces": ["application/json"],
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "description": "Login payload",
                        "required": true,
                        "schema": { "$ref": "#/definitions/LoginRequest" }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User logged in successfully",
                        "schema": { "$ref": "#/definitions/LoginResponse" }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": { "$ref": "#/definitions/MessageResponse" }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": { "$ref": "#/definitions/MessageResponse" }
                    }
                }
            }
        }
    },
    "definitions": {
        "User": {
            "type": "object",
            "properties": {
                "id": {"type": "string","example": "64c89d2a23f8"},
                "mobile": {"type": "string","example": "09123456789"},
                "created_at": {"type": "string","format": "date-time","example": "2023-07-30T12:34:56Z"}
            }
        },
        "UsersResponse": {
            "type": "object",
            "additionalProperties": {
                "type": "array",
                "users": { "$ref": "#/definitions/User" }
            }
        },
        "UserResponse": {
            "type": "object",
            "properties": {
                "user": { "$ref": "#/definitions/User" }
            }
        },
        "MessageResponse": {
            "type": "object",
            "properties": {
                "message": {"type": "string"}
            }
        },
        "RegisterRequest": {
            "type": "object",
            "properties": {
                "mobileNumber": {"type": "string","example": "09123456789"}
            },
            "required": ["mobileNumber"]
        },
        "RegisterResponse": {
            "type": "object",
            "properties": {
                "message": {"type": "string"}
            }
        },
        "LoginRequest": {
            "type": "object",
            "properties": {
                "mobileNumber": {"type": "string","example": "09123456789"},
                "OTP": {"type": "integer","example": 123456}
            },
            "required": ["mobileNumber","OTP"]
        },
        "LoginResponse": {
            "type": "object",
            "properties": {
                "message": {"type": "string"},
                "token": {"type": "string","example": "jwt.token.here"}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Dekamond API",
	Description:      "Dekamond task API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
