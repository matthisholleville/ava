// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/chat": {
            "post": {
                "description": "used to chat with Ava",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chat"
                ],
                "summary": "Chat with Ava",
                "parameters": [
                    {
                        "description": "Create a new chat",
                        "name": "_",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.CreateNewChat"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/api.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/chat/webhook": {
            "post": {
                "description": "used to chat with Ava from AlertManager Webhook",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chat"
                ],
                "summary": "Chat with Ava from AlertManager Webhook",
                "parameters": [
                    {
                        "description": "Webhook payload",
                        "name": "_",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.WebhookPayload"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/api.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/chat/{id}": {
            "get": {
                "description": "used to chat with Ava",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chat"
                ],
                "summary": "Chat with Ava",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.FetchMessagesResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "used to respond to Ava",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chat"
                ],
                "summary": "Chat with Ava",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Create a new chat",
                        "name": "_",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.CreateNewChat"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/api.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/event/slack": {
            "post": {
                "description": "used to chat with Ava when a slack event is received",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Event"
                ],
                "summary": "Receive a Slack event and chat with Ava",
                "parameters": [
                    {
                        "description": "ReceiveSlackEvent payload",
                        "name": "_",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/slack.ReceiveSlackEvent"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/api.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/knowledge": {
            "post": {
                "description": "used to add knowledge to Ava",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Knowledge"
                ],
                "summary": "Add knowledge to Ava",
                "parameters": [
                    {
                        "description": "CreateNewKnowledge payload",
                        "name": "_",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.CreateNewKnowledge"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/api.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "used to purge Ava's knowledge base",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Knowledge"
                ],
                "summary": "Purge Ava's knowledge base",
                "parameters": [
                    {
                        "description": "PurgeKnowledge payload",
                        "name": "_",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.PurgeKnowledge"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/api.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/live": {
            "get": {
                "description": "used by Kubernetes liveness probe",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Kubernetes"
                ],
                "summary": "Liveness check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "503": {
                        "description": "KO",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/readyz": {
            "get": {
                "description": "used by Kubernetes readiness probe",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Kubernetes"
                ],
                "summary": "Readiness check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "503": {
                        "description": "KO",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.Alert": {
            "type": "object",
            "properties": {
                "annotations": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "endsAt": {
                    "type": "string"
                },
                "fingerprint": {
                    "type": "string"
                },
                "generatorURL": {
                    "type": "string"
                },
                "labels": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "startsAt": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "api.CreateNewChat": {
            "type": "object",
            "properties": {
                "language": {
                    "type": "string",
                    "example": "en"
                },
                "message": {
                    "type": "string",
                    "example": "Pod web-server-5b866987d8-sxmtj in namespace default Crashlooping."
                }
            }
        },
        "api.CreateNewKnowledge": {
            "type": "object",
            "properties": {
                "gitAuthToken": {
                    "type": "string",
                    "example": ""
                },
                "gitBranch": {
                    "type": "string",
                    "example": ""
                },
                "gitRepositoryURL": {
                    "type": "string",
                    "example": ""
                },
                "path": {
                    "type": "string",
                    "example": "./docs/runbooks"
                },
                "source": {
                    "type": "string",
                    "example": "local"
                }
            }
        },
        "api.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "api.FetchMessagesResponse": {
            "type": "object",
            "properties": {
                "chat": {
                    "type": "string"
                },
                "response": {
                    "type": "string"
                }
            }
        },
        "api.PurgeKnowledge": {
            "type": "object"
        },
        "api.SuccessResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "api.WebhookPayload": {
            "type": "object",
            "properties": {
                "alerts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.Alert"
                    }
                },
                "commonAnnotations": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "commonLabels": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "externalURL": {
                    "type": "string"
                },
                "groupKey": {
                    "type": "string"
                },
                "groupLabels": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "receiver": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "truncatedAlerts": {
                    "type": "integer"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "slack.ReceiveSlackEvent": {
            "type": "object",
            "properties": {
                "api_app_id": {
                    "type": "string"
                },
                "authorizations": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/slack.SlackAuthorization"
                    }
                },
                "challenge": {
                    "type": "string"
                },
                "context_enterprise_id": {
                    "type": "string"
                },
                "context_team_id": {
                    "type": "string"
                },
                "event": {
                    "$ref": "#/definitions/slack.SlackEvent"
                },
                "event_context": {
                    "type": "string"
                },
                "event_id": {
                    "type": "string"
                },
                "event_time": {
                    "type": "integer"
                },
                "is_ext_shared_channel": {
                    "type": "boolean"
                },
                "team_id": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "slack.SlackAttachment": {
            "type": "object",
            "properties": {
                "color": {
                    "type": "string"
                },
                "fallback": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "mrkdwn_in": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "text": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "title_link": {
                    "type": "string"
                }
            }
        },
        "slack.SlackAuthorization": {
            "type": "object",
            "properties": {
                "enterprise_id": {
                    "type": "string"
                },
                "is_bot": {
                    "type": "boolean"
                },
                "is_enterprise_install": {
                    "type": "boolean"
                },
                "team_id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "slack.SlackBlock": {
            "type": "object",
            "properties": {
                "block_id": {
                    "type": "string"
                },
                "elements": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/slack.SlackBlockSection"
                    }
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "slack.SlackBlockElement": {
            "type": "object",
            "properties": {
                "text": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "slack.SlackBlockSection": {
            "type": "object",
            "properties": {
                "elements": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/slack.SlackBlockElement"
                    }
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "slack.SlackEvent": {
            "type": "object",
            "properties": {
                "attachments": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/slack.SlackAttachment"
                    }
                },
                "blocks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/slack.SlackBlock"
                    }
                },
                "bot_id": {
                    "type": "string"
                },
                "channel": {
                    "type": "string"
                },
                "channel_type": {
                    "type": "string"
                },
                "client_msg_id": {
                    "type": "string"
                },
                "event_ts": {
                    "type": "string"
                },
                "subtype": {
                    "type": "string"
                },
                "team": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                },
                "thread_ts": {
                    "type": "string"
                },
                "ts": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "user": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
