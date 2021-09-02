// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package swagger

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/network": {
            "get": {
                "description": "获取网络信息",
                "summary": "网络查询",
                "parameters": [
                    {
                        "type": "string",
                        "description": "网络名称",
                        "name": "network",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "返回网络信息",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "创建新的网络基础信息",
                "produces": [
                    "application/json"
                ],
                "summary": "网络创建",
                "parameters": [
                    {
                        "description": "网络创建",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.NetCreate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "创建成功",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/network/init": {
            "post": {
                "description": "根据请求参数对网络进行初始化，生成 fabric 网络",
                "produces": [
                    "application/json"
                ],
                "summary": "网络初始化",
                "parameters": [
                    {
                        "description": "网络初始化请求参数",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.NetInit"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "code:0 - 网络成功初始化",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "500": {
                        "description": "初始化出错",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "request.CouchDB": {
            "type": "object",
            "properties": {
                "admin": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "port": {
                    "type": "integer"
                }
            }
        },
        "request.NetCreate": {
            "type": "object",
            "properties": {
                "desc": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "request.NetInfo": {
            "type": "object",
            "properties": {
                "desc": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "request.NetInit": {
            "type": "object",
            "properties": {
                "crypto_type": {
                    "description": "SW、GM",
                    "type": "string"
                },
                "genesis_config": {
                    "type": "object",
                    "$ref": "#/definitions/request.OrdererConfig"
                },
                "hosts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/request.RuntimeHost"
                    }
                },
                "image_peer": {
                    "type": "string"
                },
                "images_orderer": {
                    "type": "string"
                },
                "members": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/request.Org"
                    }
                },
                "net_signs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/request.NetSign"
                    }
                },
                "network": {
                    "type": "object",
                    "$ref": "#/definitions/request.NetInfo"
                },
                "node_orderers": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "array",
                        "items": {
                            "$ref": "#/definitions/request.NodeOrderer"
                        }
                    }
                },
                "node_peers": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "array",
                        "items": {
                            "$ref": "#/definitions/request.NodePeer"
                        }
                    }
                },
                "runtime": {
                    "description": "Docker 、 k8s",
                    "type": "string"
                },
                "state_db": {
                    "description": "levelDB、couchDB",
                    "type": "string"
                },
                "version": {
                    "description": "fabric 版本（镜像版本）",
                    "type": "string"
                }
            }
        },
        "request.NetSign": {
            "type": "object",
            "properties": {
                "addr": {
                    "type": "string"
                },
                "data_center": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "port": {
                    "type": "integer"
                }
            }
        },
        "request.NodeOrderer": {
            "type": "object",
            "properties": {
                "cpu": {
                    "description": "核",
                    "type": "number"
                },
                "data_center": {
                    "type": "string"
                },
                "host_id": {
                    "type": "integer"
                },
                "memory": {
                    "description": "G",
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "port": {
                    "type": "integer"
                },
                "stroge": {
                    "description": "G",
                    "type": "integer"
                }
            }
        },
        "request.NodePeer": {
            "type": "object",
            "properties": {
                "cpu": {
                    "description": "核",
                    "type": "number"
                },
                "data_center": {
                    "type": "string"
                },
                "host_id": {
                    "type": "integer"
                },
                "is_anchor": {
                    "type": "boolean"
                },
                "is_bootstrap": {
                    "type": "boolean"
                },
                "memory": {
                    "description": "G",
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "port": {
                    "type": "integer"
                },
                "state_db": {
                    "type": "object",
                    "$ref": "#/definitions/request.CouchDB"
                },
                "stroge": {
                    "description": "G",
                    "type": "integer"
                }
            }
        },
        "request.OrdererConfig": {
            "type": "object",
            "properties": {
                "absoluteMaxBytes": {
                    "description": "M",
                    "type": "integer"
                },
                "batchTimeout": {
                    "description": "s",
                    "type": "integer"
                },
                "maxMessageCount": {
                    "type": "integer"
                },
                "preferredMaxBytes": {
                    "description": "KB",
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "request.Org": {
            "type": "object",
            "properties": {
                "domain": {
                    "type": "string"
                },
                "mspid": {
                    "type": "string"
                }
            }
        },
        "request.RuntimeHost": {
            "type": "object",
            "properties": {
                "desc": {
                    "type": "string"
                },
                "helmConfig": {
                    "type": "object",
                    "properties": {
                        "kubefile": {
                            "type": "string"
                        },
                        "repoConfig": {
                            "type": "object",
                            "properties": {
                                "password": {
                                    "type": "string"
                                },
                                "private": {
                                    "type": "boolean"
                                },
                                "repo": {
                                    "type": "string"
                                },
                                "username": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                },
                "host": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "scheme": {
                    "type": "string"
                },
                "tls": {
                    "type": "boolean"
                },
                "tlsconfig": {
                    "type": "object",
                    "properties": {
                        "tlsca": {
                            "type": "string"
                        },
                        "tlscert": {
                            "type": "string"
                        },
                        "tlskey": {
                            "type": "string"
                        }
                    }
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "response.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "object"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
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
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8080",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "baas 平台后端 API",
	Description: "fabric 区块链管控台 - baas 后端 API",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
