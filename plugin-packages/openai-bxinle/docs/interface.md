# OpenAiBxinle / ZeroFA Seedance 视频接口字段

## 协议身份

- 插件 ID：`openai-bxinle`。
- Provider ID：`openai-bxinle`。
- 能力：`video`。
- 默认 Base URL：渠道服务根地址，不要再重复写 `/v1`。
- 鉴权：渠道密钥，类型 `bearer`。
- 创建：`POST /v1/videos`。
- 查询：`GET /v1/videos/{{taskId}}`。
- 结果下载：`GET /v1/videos/{{taskId}}/content`。

渠道 Base URL 只填协议服务根。本包把相对路径 `/v1/videos` 拼到上游。密钥只由后端渠道中转读取，不要放进浏览器 URL、清单或日志。

## 配置字段

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `apiKey` | secret | 是 | API Key |

## 统一字段映射

| 统一字段 | 类型 | 必填 | 实际上游映射 | 说明 |
| --- | --- | --- | --- | --- |
| `model` | string | 是 | `model` | 视频模型 ID。 |
| `prompt` | string | 是 | `prompt` | 视频提示词。 |
| `images` | media[] | 否 | `image` / `first_frame` / `last_frame` | 按 role 拆到参考图、首帧和尾帧。 |
| `videos` | media[] | 否 | `reference_video` | 参考视频。 |
| `audios` | media[] | 否 | `reference_audio` | 参考音频；必须同时带图片或视频。 |
| `duration` | integer | 否 | `duration` | 时长秒数，默认能力范围 4–15。 |
| `aspectRatio` | string | 否 | `aspect_ratio` | 画幅比例。 |
| `resolution` | string | 否 | `resolution` | 分辨率档位，发送前转成大写。 |
| `generateAudio` | boolean | 否 | `generate_audio` / `enable_audio` | 是否同时生成音频。 |
| `watermark` | boolean | 否 | `watermark` | 水印开关；本协议默认不发送，可走扩展体。 |
| `providerOptions` | object | 否 | `provider-specific fields` | 渠道扩展字段。 |

## 创建请求模型与字段清单

下表由本包清单生成，覆盖 body、query、headers 和结果下载中的每个字段。

| 请求位置 | 值或转换表达式 |
| --- | --- |
| `create.method` | `"POST"` |
| `create.path` | `"/v1/videos"` |
| `create.contentType` | `"application/json"` |
| `create.body.model` | `request.model` |
| `create.body.prompt` | `request.prompt` |
| `create.body.duration` | `request.duration` |
| `create.body.aspect_ratio` | `request.aspectRatio` |
| `create.body.resolution` | `request.resolution` 转大写 |
| `create.body.generate_audio` | `request.generateAudio` |
| `create.body.enable_audio` | `request.generateAudio` |
| `create.body.first_frame` | role=`first_frame` 的图片 URL |
| `create.body.last_frame` | role=`last_frame` 的图片 URL |
| `create.body.image` | 其余参考图；一张发字符串，多张发数组 |
| `create.body.reference_video` | 参考视频 URL；一条发字符串，多条发数组 |
| `create.body.reference_audio` | 参考音频 URL；一条发字符串，多条发数组 |
| `create.body.mode` | 按素材组合推导 `text2video` / `image2video` / `reference2video` / `first_last_frame` / `video2video` |
| `poll.method` | `"GET"` |
| `poll.path` | `"/v1/videos/{{taskId}}"` |
| `result.method` | `"GET"` |
| `result.path` | `"/v1/videos/{{taskId}}/content"` |
| `result.headers.Accept` | `"video/mp4"` |

## Provider 扩展点

- `providerOptions.openai-bxinle.body`
- `providerOptions.openai-bxinle.extra_body`

动态模型或网关需要额外字段时，写入上述对象；该对象不是协议本身的固定 schema，不会被清单校验。

## 响应映射表与字段清单

| 映射位置 | 解析路径或转换表达式 |
| --- | --- |
| `response.taskId` | `id` / `task_id` / `data.id` |
| `response.status` | `status` / `state` / `data.status` |
| `response.message` | `error.message` / `message` / `fail_reason` |
| `response.videos` | `metadata.url` / `output.url` / `output.video_url` / `url` / `video_url` / `object` |
| `response.errorPaths` | `error.code` |
| `response.resultEphemeral` | `true` |

## 响应约定

成功状态接受 `succeeded`、`completed`、`success`、`done`。任务 ID 可能包在 `data` 里。结果 URL 通常短期有效；没有可下载 URL 时，宿主会走 `GET /v1/videos/{id}/content`。HTTP 非 2xx、业务 code 或 error object 记为失败。当前没有取消接口。

<!-- YINGCE_MANIFEST_CONTRACT_START -->
## Manifest 完整接口定义

以下 JSON 与插件包内实际 `manifest.json` 逐字段一致，覆盖插件身份、权限、配置、鉴权、参数、校验、创建、Agent、查询、取消、结果下载、响应和 Agent 响应映射。`documentation` 字段的值就是当前完整文档；为避免文档在自身内部无限递归，JSON 中仅用等义占位文本表示正文。

```json
{
  "apiVersion": "yingce.plugin/v2",
  "id": "openai-bxinle",
  "name": "OpenAiBxinle / ZeroFA Seedance 视频",
  "version": "1.0.0",
  "author": "ZeroFA / 影策",
  "description": "ZeroFA / Seedance 兼容 JSON 视频协议，可导入为系统协议插件。",
  "enabled": true,
  "installable": true,
  "permissions": [
    "generation.run",
    "media.read"
  ],
  "configuration": {
    "fields": [
      {
        "name": "apiKey",
        "type": "secret",
        "label": "API Key",
        "required": true
      }
    ]
  },
  "contributes": {
    "providers": [
      {
        "id": "openai-bxinle",
        "label": "OpenAiBxinle / ZeroFA Seedance 视频",
        "capabilities": [
          "video"
        ],
        "scopes": [
          "admin.system-channel",
          "user.custom-channel",
          "canvas",
          "creation",
          "agent"
        ],
        "requiresPublicMediaUrls": false,
        "auth": {
          "type": "bearer",
          "field": "apiKey"
        },
        "parameters": [
          {
            "name": "model",
            "type": "string",
            "required": true,
            "mapping": "model",
            "description": "视频模型 ID。"
          },
          {
            "name": "prompt",
            "type": "string",
            "required": true,
            "mapping": "prompt",
            "description": "视频提示词。"
          },
          {
            "name": "images",
            "type": "media[]",
            "required": false,
            "mapping": "image/first_frame/last_frame",
            "description": "按 role 拆到参考图、首帧和尾帧。"
          },
          {
            "name": "videos",
            "type": "media[]",
            "required": false,
            "mapping": "reference_video",
            "description": "参考视频。"
          },
          {
            "name": "audios",
            "type": "media[]",
            "required": false,
            "mapping": "reference_audio",
            "description": "参考音频；必须同时带图片或视频。"
          },
          {
            "name": "duration",
            "type": "integer",
            "required": false,
            "mapping": "duration",
            "description": "时长秒数，默认能力范围 4–15。"
          },
          {
            "name": "aspectRatio",
            "type": "string",
            "required": false,
            "mapping": "aspect_ratio",
            "description": "画幅比例。"
          },
          {
            "name": "resolution",
            "type": "string",
            "required": false,
            "mapping": "resolution",
            "description": "分辨率档位，发送前转成大写。"
          },
          {
            "name": "generateAudio",
            "type": "boolean",
            "required": false,
            "mapping": "generate_audio/enable_audio",
            "description": "是否同时生成音频。"
          },
          {
            "name": "watermark",
            "type": "boolean",
            "required": false,
            "mapping": "watermark",
            "description": "水印开关；本协议默认不发送，可走扩展体。"
          },
          {
            "name": "providerOptions",
            "type": "object",
            "required": false,
            "mapping": "provider-specific fields",
            "description": "渠道扩展字段。"
          }
        ],
        "validations": [
          {
            "assert": {
              "$gt": [
                {
                  "$len": {
                    "$trim": {
                      "$ref": "request.prompt"
                    }
                  }
                },
                0
              ]
            },
            "message": "OpenAiBxinle 需要填写 prompt"
          },
          {
            "assert": {
              "$lte": [
                {
                  "$len": {
                    "$ref": "request.images"
                  }
                },
                9
              ]
            },
            "message": "OpenAiBxinle 最多支持 9 张参考图片"
          },
          {
            "assert": {
              "$lte": [
                {
                  "$len": {
                    "$ref": "request.videos"
                  }
                },
                3
              ]
            },
            "message": "OpenAiBxinle 最多支持 3 个参考视频"
          },
          {
            "assert": {
              "$lte": [
                {
                  "$len": {
                    "$ref": "request.audios"
                  }
                },
                3
              ]
            },
            "message": "OpenAiBxinle 最多支持 3 个参考音频"
          },
          {
            "assert": {
              "$or": [
                {
                  "$eq": [
                    {
                      "$len": {
                        "$ref": "request.audios"
                      }
                    },
                    0
                  ]
                },
                {
                  "$gt": [
                    {
                      "$len": {
                        "$ref": "request.images"
                      }
                    },
                    0
                  ]
                },
                {
                  "$gt": [
                    {
                      "$len": {
                        "$ref": "request.videos"
                      }
                    },
                    0
                  ]
                }
              ]
            },
            "message": "OpenAiBxinle 的参考音频必须与至少一个图片或视频参考素材同时使用"
          }
        ],
        "create": {
          "method": "POST",
          "path": "/v1/videos",
          "contentType": "application/json",
          "body": {
            "$merge": [
              {
                "model": {
                  "$ref": "request.model"
                },
                "prompt": {
                  "$ref": "request.prompt"
                },
                "duration": {
                  "$omitEmpty": {
                    "$ref": "request.duration"
                  }
                },
                "aspect_ratio": {
                  "$omitEmpty": {
                    "$ref": "request.aspectRatio"
                  }
                },
                "resolution": {
                  "$omitEmpty": {
                    "$upper": {
                      "$ref": "request.resolution"
                    }
                  }
                },
                "generate_audio": {
                  "$omitEmpty": {
                    "$ref": "request.generateAudio"
                  }
                },
                "enable_audio": {
                  "$omitEmpty": {
                    "$ref": "request.generateAudio"
                  }
                },
                "first_frame": {
                  "$omitEmpty": {
                    "$first": {
                      "$map": {
                        "from": {
                          "$filter": {
                            "from": {
                              "$sortByOrder": {
                                "$ref": "request.images"
                              }
                            },
                            "as": "media",
                            "where": {
                              "$eq": [
                                {
                                  "$ref": "media.role"
                                },
                                "first_frame"
                              ]
                            }
                          }
                        },
                        "as": "media",
                        "in": {
                          "$ref": "media.value"
                        }
                      }
                    }
                  }
                },
                "last_frame": {
                  "$omitEmpty": {
                    "$first": {
                      "$map": {
                        "from": {
                          "$filter": {
                            "from": {
                              "$sortByOrder": {
                                "$ref": "request.images"
                              }
                            },
                            "as": "media",
                            "where": {
                              "$eq": [
                                {
                                  "$ref": "media.role"
                                },
                                "last_frame"
                              ]
                            }
                          }
                        },
                        "as": "media",
                        "in": {
                          "$ref": "media.value"
                        }
                      }
                    }
                  }
                },
                "image": {
                  "$omitEmpty": {
                    "$if": {
                      "condition": {
                        "$eq": [
                          {
                            "$len": {
                              "$filter": {
                                "from": {
                                  "$sortByOrder": {
                                    "$ref": "request.images"
                                  }
                                },
                                "as": "media",
                                "where": {
                                  "$and": [
                                    {
                                      "$ne": [
                                        {
                                          "$ref": "media.role"
                                        },
                                        "first_frame"
                                      ]
                                    },
                                    {
                                      "$ne": [
                                        {
                                          "$ref": "media.role"
                                        },
                                        "last_frame"
                                      ]
                                    }
                                  ]
                                }
                              }
                            }
                          },
                          1
                        ]
                      },
                      "then": {
                        "$first": {
                          "$map": {
                            "from": {
                              "$filter": {
                                "from": {
                                  "$sortByOrder": {
                                    "$ref": "request.images"
                                  }
                                },
                                "as": "media",
                                "where": {
                                  "$and": [
                                    {
                                      "$ne": [
                                        {
                                          "$ref": "media.role"
                                        },
                                        "first_frame"
                                      ]
                                    },
                                    {
                                      "$ne": [
                                        {
                                          "$ref": "media.role"
                                        },
                                        "last_frame"
                                      ]
                                    }
                                  ]
                                }
                              }
                            },
                            "as": "media",
                            "in": {
                              "$ref": "media.value"
                            }
                          }
                        }
                      },
                      "else": {
                        "$map": {
                          "from": {
                            "$filter": {
                              "from": {
                                "$sortByOrder": {
                                  "$ref": "request.images"
                                }
                              },
                              "as": "media",
                              "where": {
                                "$and": [
                                  {
                                    "$ne": [
                                      {
                                        "$ref": "media.role"
                                      },
                                      "first_frame"
                                    ]
                                  },
                                  {
                                    "$ne": [
                                      {
                                        "$ref": "media.role"
                                      },
                                      "last_frame"
                                    ]
                                  }
                                ]
                              }
                            }
                          },
                          "as": "media",
                          "in": {
                            "$ref": "media.value"
                          }
                        }
                      }
                    }
                  }
                },
                "reference_video": {
                  "$omitEmpty": {
                    "$if": {
                      "condition": {
                        "$eq": [
                          {
                            "$len": {
                              "$ref": "request.videos"
                            }
                          },
                          1
                        ]
                      },
                      "then": {
                        "$first": {
                          "$map": {
                            "from": {
                              "$sortByOrder": {
                                "$ref": "request.videos"
                              }
                            },
                            "as": "media",
                            "in": {
                              "$ref": "media.value"
                            }
                          }
                        }
                      },
                      "else": {
                        "$map": {
                          "from": {
                            "$sortByOrder": {
                              "$ref": "request.videos"
                            }
                          },
                          "as": "media",
                          "in": {
                            "$ref": "media.value"
                          }
                        }
                      }
                    }
                  }
                },
                "reference_audio": {
                  "$omitEmpty": {
                    "$if": {
                      "condition": {
                        "$eq": [
                          {
                            "$len": {
                              "$ref": "request.audios"
                            }
                          },
                          1
                        ]
                      },
                      "then": {
                        "$first": {
                          "$map": {
                            "from": {
                              "$sortByOrder": {
                                "$ref": "request.audios"
                              }
                            },
                            "as": "media",
                            "in": {
                              "$ref": "media.value"
                            }
                          }
                        }
                      },
                      "else": {
                        "$map": {
                          "from": {
                            "$sortByOrder": {
                              "$ref": "request.audios"
                            }
                          },
                          "as": "media",
                          "in": {
                            "$ref": "media.value"
                          }
                        }
                      }
                    }
                  }
                },
                "mode": {
                  "$switch": {
                    "cases": [
                      {
                        "when": {
                          "$and": [
                            {
                              "$gt": [
                                {
                                  "$len": {
                                    "$filter": {
                                      "from": {
                                        "$ref": "request.images"
                                      },
                                      "as": "media",
                                      "where": {
                                        "$eq": [
                                          {
                                            "$ref": "media.role"
                                          },
                                          "first_frame"
                                        ]
                                      }
                                    }
                                  }
                                },
                                0
                              ]
                            },
                            {
                              "$gt": [
                                {
                                  "$len": {
                                    "$filter": {
                                      "from": {
                                        "$ref": "request.images"
                                      },
                                      "as": "media",
                                      "where": {
                                        "$eq": [
                                          {
                                            "$ref": "media.role"
                                          },
                                          "last_frame"
                                        ]
                                      }
                                    }
                                  }
                                },
                                0
                              ]
                            }
                          ]
                        },
                        "then": "first_last_frame"
                      },
                      {
                        "when": {
                          "$gt": [
                            {
                              "$len": {
                                "$ref": "request.videos"
                              }
                            },
                            0
                          ]
                        },
                        "then": "video2video"
                      },
                      {
                        "when": {
                          "$gt": [
                            {
                              "$len": {
                                "$filter": {
                                  "from": {
                                    "$ref": "request.images"
                                  },
                                  "as": "media",
                                  "where": {
                                    "$and": [
                                      {
                                        "$ne": [
                                          {
                                            "$ref": "media.role"
                                          },
                                          "first_frame"
                                        ]
                                      },
                                      {
                                        "$ne": [
                                          {
                                            "$ref": "media.role"
                                          },
                                          "last_frame"
                                        ]
                                      }
                                    ]
                                  }
                                }
                              }
                            },
                            1
                          ]
                        },
                        "then": "reference2video"
                      },
                      {
                        "when": {
                          "$gt": [
                            {
                              "$len": {
                                "$ref": "request.images"
                              }
                            },
                            0
                          ]
                        },
                        "then": "image2video"
                      }
                    ],
                    "default": "text2video"
                  }
                }
              },
              {
                "$coalesce": [
                  {
                    "$ref": "request.providerOptions.openai-bxinle.body"
                  },
                  {
                    "$ref": "request.providerOptions.openai-bxinle.extra_body"
                  },
                  {}
                ]
              }
            ]
          }
        },
        "poll": {
          "method": "GET",
          "path": "/v1/videos/{{taskId}}",
          "contentType": "application/json"
        },
        "result": {
          "method": "GET",
          "path": "/v1/videos/{{taskId}}/content",
          "contentType": "application/json",
          "headers": {
            "Accept": "video/mp4"
          }
        },
        "response": {
          "taskId": {
            "$coalesce": [
              {
                "$ref": "response.id"
              },
              {
                "$ref": "response.task_id"
              },
              {
                "$ref": "response.taskId"
              },
              {
                "$ref": "response.data.id"
              },
              {
                "$ref": "taskId"
              }
            ]
          },
          "status": {
            "$coalesce": [
              {
                "$ref": "response.status"
              },
              {
                "$ref": "response.state"
              },
              {
                "$ref": "response.data.status"
              },
              "pending"
            ]
          },
          "message": {
            "$coalesce": [
              {
                "$ref": "response.error.message"
              },
              {
                "$ref": "response.message"
              },
              {
                "$ref": "response.fail_reason"
              }
            ]
          },
          "videos": {
            "$coalesce": [
              {
                "$ref": "response.metadata.url"
              },
              {
                "$ref": "response.data.metadata.url"
              },
              {
                "$ref": "response.output.url"
              },
              {
                "$ref": "response.output.video_url"
              },
              {
                "$ref": "response.video_url"
              },
              {
                "$ref": "response.url"
              },
              {
                "$ref": "response.object"
              }
            ]
          },
          "errorPaths": [
            "error.code"
          ],
          "resultEphemeral": true
        }
      }
    ]
  },
  "documentation": "<当前插件的完整 documentation，由 README.md 与 docs/interface.md 拼接而成；为避免 JSON 递归，此处不重复展开正文。>"
}
```
<!-- YINGCE_MANIFEST_CONTRACT_END -->
