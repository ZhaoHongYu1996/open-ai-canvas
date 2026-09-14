# TSAI 网关接口字段

## 协议身份

- 插件 ID：`tsai`。
- Provider：`tsai-chat`（文本）、`tsai-minimax-h3`（视频）、`tsai-seedance-mini`（视频）、`tsai-seedream`（图片）。
- 默认 Base URL：`https://api.alifenqi.com:8188`，不要再重复写 `/v1`。
- 鉴权：渠道密钥，类型 `bearer`，请求头 `Authorization: Bearer <apiKey>`。
- 对话：`POST /v1/chat/completions`，OpenAI Chat Completions 格式。
- MiniMax H3 / Seedance 2.0 mini 创建：`POST /v1/videos/generations`。
- 视频查询：`GET /v1/videos/{{taskId}}`。
- Seedream 4.5 创建：`POST /v1/images/generations`。
- 图片长任务查询：`GET /v1/images/generations/{{taskId}}`。

渠道 Base URL 只填协议服务根。密钥只由后端渠道中转读取，不要放进浏览器 URL、清单或日志。公网素材必须直接可下载，不接受内网 URL、重定向和 Base64。

## 配置字段

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `apiKey` | secret | 是 | TSAI 个人中心生成的 API Key，不要发送 DeepSeek 原始密钥或网页登录令牌。 |

## 统一字段映射

| 统一字段 | 类型 | 必填 | 实际上游映射 | 说明 |
| --- | --- | --- | --- | --- |
| `model` | string | 是 | `model` | 模型列表中的 ID，例如 `deepseek-v4-flash`、`MiniMax/MiniMax-H3`、`doubao-seedance-2-0-mini-260615`、`doubao-seedream-4-5-251128`。 |
| `messages` | message[] | 对话必填 | `messages` | OpenAI `role`/`content` 消息列表。 |
| `instructions` | string | 否 | `system/instructions` | 系统指令。 |
| `prompt` | string | 视频/图片必填 | `prompt` 或 `content[].text` | MiniMax 最多 7000 字；Seedance / Seedream 最多 20000 字。 |
| `images` | media[] | 否 | MiniMax：`first_frame_image` / `last_frame_image` / `image_urls`；Seedance：`content[].image_url`；Seedream：`image` | 按 role 拆到首帧、尾帧或普通参考。首尾帧与普通多模态参考不能混用。最多 9 张。 |
| `videos` | media[] | 否 | MiniMax：`video_urls`；Seedance：`content[].video_url` | 参考视频，最多 3 段。 |
| `audios` | media[] | 否 | MiniMax：`audio_urls`；Seedance：`content[].audio_url` | 参考音频，最多 3 段。MiniMax 音频必须搭配图片或视频。 |
| `duration` | integer | 否 | `duration` | 4–15 秒，默认 5。 |
| `aspectRatio` | string | 否 | MiniMax：`aspect_ratio`；Seedance：`ratio`；Seedream：`size` | 视频画幅或图片尺寸档。 |
| `resolution` | string | 否 | `resolution` | MiniMax 发送前转大写（`2K`/`768P`/`480P`/`1080P`）；Seedance 转小写（`480p`/`720p`）。 |
| `generateAudio` | boolean | 否 | `generate_audio` | Seedance 是否同时生成声音，默认 true。 |
| `watermark` | boolean | 否 | `watermark` | Seedance 默认 false；Seedream 默认 true。MiniMax 当前仅允许 false 或省略。 |
| `temperature` | number | 否 | `temperature` | 对话采样温度。 |
| `top_p` | number | 否 | `top_p` | 对话核采样。 |
| `max_tokens` | integer | 否 | `max_tokens` | 对话最大输出 token。 |
| `tools` | array | 否 | `tools` | 对话工具定义。 |
| `tool_choice` | object\|string | 否 | `tool_choice` | 对话工具选择策略。 |
| `response_format` | object | 否 | `response_format` | 对话结构化输出。 |
| `stream` | boolean | 否 | `stream` | 对话流式开关；后台任务当前以最终响应归一。 |
| `imageCount` | integer | 否 | 固定 1 | Seedream 每次只输出一张图。 |
| `quality` | string | 否 | 不发送 | TSAI 图片不支持独立质量档。 |
| `providerOptions` | object | 否 | `provider-specific fields` | 渠道扩展字段。 |

## 创建请求模型与字段清单

| 请求位置 | 值或转换表达式 |
| --- | --- |
| `tsai-chat.create.method` | `"POST"` |
| `tsai-chat.create.path` | `"/v1/chat/completions"` |
| `tsai-chat.create.body.model` | `request.model` |
| `tsai-chat.create.body.messages` | `request.messages` |
| `tsai-chat.agent.path` | `"/v1/chat/completions"` |
| `tsai-minimax-h3.create.path` | `"/v1/videos/generations"` |
| `tsai-minimax-h3.create.body.prompt` | `request.prompt` |
| `tsai-minimax-h3.create.body.first_frame_image` | role=`first_frame` 的图片 URL |
| `tsai-minimax-h3.create.body.last_frame_image` | role=`last_frame` 的图片 URL |
| `tsai-minimax-h3.create.body.image_urls` | 其余参考图 URL 数组 |
| `tsai-minimax-h3.create.body.video_urls` | 参考视频 URL 数组 |
| `tsai-minimax-h3.create.body.audio_urls` | 参考音频 URL 数组 |
| `tsai-minimax-h3.create.body.aspect_ratio` | `request.aspectRatio` |
| `tsai-minimax-h3.create.body.resolution` | `request.resolution` 转大写 |
| `tsai-minimax-h3.poll.path` | `"/v1/videos/{{taskId}}"` |
| `tsai-seedance-mini.create.path` | `"/v1/videos/generations"` |
| `tsai-seedance-mini.create.body.content` | 文本 + 按顺序排列的图片/视频/音频项 |
| `tsai-seedance-mini.create.body.ratio` | `request.aspectRatio` |
| `tsai-seedance-mini.create.body.resolution` | `request.resolution` 转小写 |
| `tsai-seedance-mini.create.body.generate_audio` | `request.generateAudio` |
| `tsai-seedance-mini.poll.path` | `"/v1/videos/{{taskId}}"` |
| `tsai-seedream.create.path` | `"/v1/images/generations"` |
| `tsai-seedream.create.body.size` | `2K`/`4K` 或由比例映射的像素尺寸 |
| `tsai-seedream.create.body.image` | 参考图 URL；一张发字符串，多张发数组 |
| `tsai-seedream.create.body.sequential_image_generation` | `"disabled"` |
| `tsai-seedream.create.body.response_format` | `"url"` |
| `tsai-seedream.create.body.stream` | `false` |
| `tsai-seedream.poll.path` | `"/v1/images/generations/{{taskId}}"` |

## Provider 扩展点

- `providerOptions.tsai-chat.body`
- `providerOptions.tsai-chat.extra_body`
- `providerOptions.tsai-chat.thinking`
- `providerOptions.tsai-minimax-h3.body`
- `providerOptions.tsai-minimax-h3.extra_body`
- `providerOptions.tsai-seedance-mini.body`
- `providerOptions.tsai-seedance-mini.extra_body`
- `providerOptions.tsai-seedream.body`
- `providerOptions.tsai-seedream.extra_body`

动态模型或网关需要额外字段时，写入上述对象；该对象不是协议本身的固定 schema，不会被清单校验。

## 响应映射表与字段清单

| 映射位置 | 解析路径或转换表达式 |
| --- | --- |
| `tsai-chat.response.text` | `choices.0.message.content` / `choices.0.text` |
| `tsai-minimax-h3.response.taskId` | `data.0.task_id` / `data.id` / `task_id` |
| `tsai-minimax-h3.response.status` | `data.0.status` / `data.status` / `status` |
| `tsai-minimax-h3.response.videos` | `data.result.videos` |
| `tsai-seedance-mini.response.taskId` | 与 MiniMax 相同 |
| `tsai-seedance-mini.response.videos` | `data.result.videos` |
| `tsai-seedream.response.taskId` | `id` |
| `tsai-seedream.response.status` | `status` |
| `tsai-seedream.response.images` | `data` |
| `response.errorPaths` | `error.code` / `data.error.code` |

## 响应约定

视频提交返回 `{ code: 200, data: [{ status: "submitted", task_id }] }`，必须用 `data[0].task_id` 查询，不能使用上游内部任务 ID。轮询对象在 `data` 里，成功状态为 `succeeded`，结果取 `data.result.videos[0].url`。`submitted`/`queued` 视为等待，`processing` 视为处理中。图片通常同步返回 `status=succeeded` 与平台文件 URL；超过约 120 秒则返回 HTTP 202 和同一 `id`，再用创建任务的同一把 Key 查询，不要重新提交。对话成功后直接取 `choices[0].message.content`。HTTP 非 2xx、业务 error object 记为失败。当前没有取消接口。

<!-- YINGCE_MANIFEST_CONTRACT_START -->
## Manifest 完整接口定义

以下 JSON 与插件包内实际 `manifest.json` 逐字段一致，覆盖插件身份、权限、配置、鉴权、参数、校验、创建、Agent、查询、取消、结果下载、响应和 Agent 响应映射。`documentation` 字段的值就是当前完整文档；为避免文档在自身内部无限递归，JSON 中仅用等义占位文本表示正文。

```json
{
  "apiVersion": "yingce.plugin/v2",
  "id": "tsai",
  "name": "TSAI 网关",
  "version": "1.0.0",
  "author": "TSAI / 影策",
  "description": "TSAI 对话、MiniMax H3 / Seedance 2.0 mini 视频与 Seedream 4.5 图片协议，可导入为系统协议插件。",
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
        "id": "tsai-chat",
        "label": "TSAI 通用对话",
        "capabilities": [
          "text"
        ],
        "scopes": [
          "admin.system-channel",
          "user.custom-channel",
          "canvas",
          "creation",
          "agent"
        ],
        "baseUrl": "https://api.alifenqi.com:8188",
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
            "description": "对话模型 ID。"
          },
          {
            "name": "messages",
            "type": "message[]",
            "required": true,
            "mapping": "provider message container",
            "description": "包含历史消息和当前用户输入。"
          },
          {
            "name": "instructions",
            "type": "string",
            "required": false,
            "mapping": "system/instructions",
            "description": "系统指令。"
          },
          {
            "name": "temperature",
            "type": "number",
            "required": false,
            "mapping": "temperature",
            "description": "采样温度。"
          },
          {
            "name": "top_p",
            "type": "number",
            "required": false,
            "mapping": "top_p",
            "description": "核采样参数。"
          },
          {
            "name": "max_tokens",
            "type": "integer",
            "required": false,
            "mapping": "max_tokens",
            "description": "最大输出 token。"
          },
          {
            "name": "tools",
            "type": "array",
            "required": false,
            "mapping": "tools",
            "description": "工具定义。"
          },
          {
            "name": "tool_choice",
            "type": "object|string",
            "required": false,
            "mapping": "tool_choice",
            "description": "工具选择策略。"
          },
          {
            "name": "response_format",
            "type": "object",
            "required": false,
            "mapping": "response_format",
            "description": "结构化输出配置。"
          },
          {
            "name": "stream",
            "type": "boolean",
            "required": false,
            "mapping": "stream",
            "description": "流式开关；后台任务当前以最终响应归一。"
          },
          {
            "name": "providerOptions",
            "type": "object",
            "required": false,
            "mapping": "provider-specific fields",
            "description": "渠道扩展字段。"
          }
        ],
        "create": {
          "method": "POST",
          "path": "/v1/chat/completions",
          "contentType": "application/json",
          "body": {
            "$merge": [
              {
                "model": {
                  "$ref": "request.model"
                },
                "messages": {
                  "$ref": "request.messages"
                },
                "temperature": {
                  "$omitEmpty": {
                    "$ref": "request.providerOptions.tsai-chat.temperature"
                  }
                },
                "top_p": {
                  "$omitEmpty": {
                    "$ref": "request.providerOptions.tsai-chat.top_p"
                  }
                },
                "max_tokens": {
                  "$omitEmpty": {
                    "$coalesce": [
                      {
                        "$ref": "request.extra.max_tokens"
                      },
                      {
                        "$ref": "request.providerOptions.tsai-chat.max_tokens"
                      }
                    ]
                  }
                },
                "thinking": {
                  "$omitEmpty": {
                    "$ref": "request.providerOptions.tsai-chat.thinking"
                  }
                },
                "tools": {
                  "$omitEmpty": {
                    "$ref": "request.providerOptions.tsai-chat.tools"
                  }
                },
                "tool_choice": {
                  "$omitEmpty": {
                    "$ref": "request.providerOptions.tsai-chat.tool_choice"
                  }
                },
                "response_format": {
                  "$omitEmpty": {
                    "$ref": "request.providerOptions.tsai-chat.response_format"
                  }
                },
                "stream": {
                  "$omitEmpty": {
                    "$ref": "request.providerOptions.tsai-chat.stream"
                  }
                }
              },
              {
                "$coalesce": [
                  {
                    "$ref": "request.providerOptions.tsai-chat.body"
                  },
                  {
                    "$ref": "request.providerOptions.tsai-chat.extra_body"
                  },
                  {}
                ]
              }
            ]
          }
        },
        "agent": {
          "method": "POST",
          "path": "/v1/chat/completions",
          "contentType": "application/json",
          "body": {
            "$merge": [
              {
                "$ref": "request.extra.agent.chatCompletion"
              },
              {
                "model": {
                  "$ref": "request.model"
                }
              }
            ]
          }
        },
        "response": {
          "status": "succeeded",
          "textPaths": [
            "choices.0.message.content",
            "choices.0.text"
          ],
          "reasoningPaths": [
            "choices.0.message.reasoning_content"
          ],
          "usage": {
            "$ref": "response.usage"
          },
          "errorPaths": [
            "error.code"
          ],
          "messagePaths": [
            "error.message"
          ]
        },
        "agentResponse": {
          "textPaths": [
            "choices.0.message.content",
            "choices.0.text"
          ],
          "reasoningPaths": [
            "choices.0.message.reasoning_content"
          ],
          "toolCallsPath": "choices.0.message.tool_calls",
          "toolCallIdPaths": [
            "id"
          ],
          "toolCallNamePaths": [
            "function.name"
          ],
          "toolCallArgumentsPaths": [
            "function.arguments"
          ]
        }
      },
      {
        "id": "tsai-minimax-h3",
        "label": "TSAI MiniMax H3 视频",
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
        "baseUrl": "https://api.alifenqi.com:8188",
        "requiresPublicMediaUrls": true,
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
            "mapping": "first_frame_image/last_frame_image/image_urls",
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
            "description": "参考音频。"
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
            "description": "分辨率档位。"
          },
          {
            "name": "generateAudio",
            "type": "boolean",
            "required": false,
            "mapping": "generate_audio",
            "description": "是否同时生成音频。"
          },
          {
            "name": "watermark",
            "type": "boolean",
            "required": false,
            "mapping": "watermark",
            "description": "水印开关。"
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
            "message": "TSAI MiniMax H3 需要填写 prompt"
          },
          {
            "assert": {
              "$lte": [
                {
                  "$len": {
                    "$ref": "request.prompt"
                  }
                },
                7000
              ]
            },
            "message": "TSAI MiniMax H3 的 prompt 最多 7000 个字符"
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
            "message": "TSAI MiniMax H3 最多支持 9 张参考图片"
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
            "message": "TSAI MiniMax H3 最多支持 3 个参考视频"
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
            "message": "TSAI MiniMax H3 最多支持 3 个参考音频"
          },
          {
            "assert": {
              "$or": [
                {
                  "$and": [
                    {
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
                {
                  "$and": [
                    {
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
                        0
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "request.videos"
                          }
                        },
                        0
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "request.audios"
                          }
                        },
                        0
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "TSAI MiniMax H3 的首尾帧不能与普通参考图、参考视频或参考音频混用"
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
            "message": "TSAI MiniMax H3 的参考音频必须与至少一个图片或视频参考素材同时使用"
          }
        ],
        "create": {
          "method": "POST",
          "path": "/v1/videos/generations",
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
                "first_frame_image": {
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
                "last_frame_image": {
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
                "image_urls": {
                  "$omitEmpty": {
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
                "video_urls": {
                  "$omitEmpty": {
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
                "audio_urls": {
                  "$omitEmpty": {
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
                "watermark": {
                  "$omitEmpty": {
                    "$ref": "request.watermark"
                  }
                }
              },
              {
                "$coalesce": [
                  {
                    "$ref": "request.providerOptions.tsai-minimax-h3.body"
                  },
                  {
                    "$ref": "request.providerOptions.tsai-minimax-h3.extra_body"
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
        "response": {
          "taskId": {
            "$coalesce": [
              {
                "$ref": "response.data.0.task_id"
              },
              {
                "$ref": "response.data.task_id"
              },
              {
                "$ref": "response.data.id"
              },
              {
                "$ref": "response.task_id"
              },
              {
                "$ref": "response.id"
              },
              {
                "$ref": "taskId"
              }
            ]
          },
          "status": {
            "$coalesce": [
              {
                "$ref": "response.data.0.status"
              },
              {
                "$ref": "response.data.status"
              },
              {
                "$ref": "response.status"
              },
              "pending"
            ]
          },
          "message": {
            "$coalesce": [
              {
                "$ref": "response.data.error.message"
              },
              {
                "$ref": "response.error.message"
              },
              {
                "$ref": "response.message"
              },
              {
                "$ref": "response.detail"
              }
            ]
          },
          "videos": {
            "$coalesce": [
              {
                "$ref": "response.data.result.videos"
              },
              {
                "$ref": "response.result.videos"
              }
            ]
          },
          "errorPaths": [
            "error.code",
            "data.error.code"
          ],
          "resultEphemeral": false
        }
      },
      {
        "id": "tsai-seedance-mini",
        "label": "TSAI Seedance 2.0 mini 视频",
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
        "baseUrl": "https://api.alifenqi.com:8188",
        "requiresPublicMediaUrls": true,
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
            "mapping": "content[].image_url",
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
            "description": "参考音频。"
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
            "mapping": "ratio",
            "description": "画幅比例。"
          },
          {
            "name": "resolution",
            "type": "string",
            "required": false,
            "mapping": "resolution",
            "description": "分辨率档位。"
          },
          {
            "name": "generateAudio",
            "type": "boolean",
            "required": false,
            "mapping": "generate_audio",
            "description": "是否同时生成音频。"
          },
          {
            "name": "watermark",
            "type": "boolean",
            "required": false,
            "mapping": "watermark",
            "description": "水印开关。"
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
            "message": "TSAI Seedance 2.0 mini 需要填写 prompt"
          },
          {
            "assert": {
              "$lte": [
                {
                  "$len": {
                    "$ref": "request.prompt"
                  }
                },
                20000
              ]
            },
            "message": "TSAI Seedance 2.0 mini 的 prompt 最多 20000 个字符"
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
            "message": "TSAI Seedance 2.0 mini 最多支持 9 张参考图片"
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
            "message": "TSAI Seedance 2.0 mini 最多支持 3 个参考视频"
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
            "message": "TSAI Seedance 2.0 mini 最多支持 3 个参考音频"
          },
          {
            "assert": {
              "$or": [
                {
                  "$and": [
                    {
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
                {
                  "$and": [
                    {
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
                        0
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "request.videos"
                          }
                        },
                        0
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "request.audios"
                          }
                        },
                        0
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "TSAI Seedance 2.0 mini 的首尾帧不能与普通参考图、参考视频或参考音频混用"
          }
        ],
        "create": {
          "method": "POST",
          "path": "/v1/videos/generations",
          "contentType": "application/json",
          "body": {
            "$merge": [
              {
                "model": {
                  "$ref": "request.model"
                },
                "content": {
                  "$concatArrays": [
                    [
                      {
                        "type": "text",
                        "text": {
                          "$ref": "request.prompt"
                        }
                      }
                    ],
                    {
                      "$map": {
                        "from": {
                          "$sortByOrder": {
                            "$ref": "request.images"
                          }
                        },
                        "as": "media",
                        "in": {
                          "type": "image_url",
                          "image_url": {
                            "url": {
                              "$ref": "media.value"
                            }
                          },
                          "role": {
                            "$switch": {
                              "cases": [
                                {
                                  "when": {
                                    "$eq": [
                                      {
                                        "$ref": "media.role"
                                      },
                                      "first_frame"
                                    ]
                                  },
                                  "then": "first_frame"
                                },
                                {
                                  "when": {
                                    "$eq": [
                                      {
                                        "$ref": "media.role"
                                      },
                                      "last_frame"
                                    ]
                                  },
                                  "then": "last_frame"
                                }
                              ],
                              "default": "reference_image"
                            }
                          }
                        }
                      }
                    },
                    {
                      "$map": {
                        "from": {
                          "$sortByOrder": {
                            "$ref": "request.videos"
                          }
                        },
                        "as": "media",
                        "in": {
                          "type": "video_url",
                          "video_url": {
                            "url": {
                              "$ref": "media.value"
                            }
                          },
                          "role": "reference_video"
                        }
                      }
                    },
                    {
                      "$map": {
                        "from": {
                          "$sortByOrder": {
                            "$ref": "request.audios"
                          }
                        },
                        "as": "media",
                        "in": {
                          "type": "audio_url",
                          "audio_url": {
                            "url": {
                              "$ref": "media.value"
                            }
                          },
                          "role": "reference_audio"
                        }
                      }
                    }
                  ]
                },
                "duration": {
                  "$omitEmpty": {
                    "$ref": "request.duration"
                  }
                },
                "ratio": {
                  "$omitEmpty": {
                    "$ref": "request.aspectRatio"
                  }
                },
                "resolution": {
                  "$omitEmpty": {
                    "$lower": {
                      "$ref": "request.resolution"
                    }
                  }
                },
                "generate_audio": {
                  "$omitEmpty": {
                    "$ref": "request.generateAudio"
                  }
                },
                "watermark": {
                  "$omitEmpty": {
                    "$ref": "request.watermark"
                  }
                },
                "mode": "reference"
              },
              {
                "$coalesce": [
                  {
                    "$ref": "request.providerOptions.tsai-seedance-mini.body"
                  },
                  {
                    "$ref": "request.providerOptions.tsai-seedance-mini.extra_body"
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
        "response": {
          "taskId": {
            "$coalesce": [
              {
                "$ref": "response.data.0.task_id"
              },
              {
                "$ref": "response.data.task_id"
              },
              {
                "$ref": "response.data.id"
              },
              {
                "$ref": "response.task_id"
              },
              {
                "$ref": "response.id"
              },
              {
                "$ref": "taskId"
              }
            ]
          },
          "status": {
            "$coalesce": [
              {
                "$ref": "response.data.0.status"
              },
              {
                "$ref": "response.data.status"
              },
              {
                "$ref": "response.status"
              },
              "pending"
            ]
          },
          "message": {
            "$coalesce": [
              {
                "$ref": "response.data.error.message"
              },
              {
                "$ref": "response.error.message"
              },
              {
                "$ref": "response.message"
              },
              {
                "$ref": "response.detail"
              }
            ]
          },
          "videos": {
            "$coalesce": [
              {
                "$ref": "response.data.result.videos"
              },
              {
                "$ref": "response.result.videos"
              }
            ]
          },
          "errorPaths": [
            "error.code",
            "data.error.code"
          ],
          "resultEphemeral": false
        }
      },
      {
        "id": "tsai-seedream",
        "label": "TSAI Seedream 4.5 图片",
        "capabilities": [
          "image"
        ],
        "scopes": [
          "admin.system-channel",
          "user.custom-channel",
          "canvas",
          "creation",
          "agent"
        ],
        "baseUrl": "https://api.alifenqi.com:8188",
        "requiresPublicMediaUrls": true,
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
            "description": "图片模型 ID。"
          },
          {
            "name": "prompt",
            "type": "string",
            "required": true,
            "mapping": "prompt",
            "description": "图片提示词。"
          },
          {
            "name": "images",
            "type": "media[]",
            "required": false,
            "mapping": "image",
            "description": "单图或多参考图 URL，按提示词中的图1、图2顺序。"
          },
          {
            "name": "imageCount",
            "type": "integer",
            "required": false,
            "mapping": "n",
            "description": "输出数量；本协议固定一张。"
          },
          {
            "name": "aspectRatio",
            "type": "string",
            "required": false,
            "mapping": "size",
            "description": "2K、4K 或宽高像素。"
          },
          {
            "name": "resolution",
            "type": "string",
            "required": false,
            "mapping": "size",
            "description": "分辨率档位，可映射到 size。"
          },
          {
            "name": "quality",
            "type": "string",
            "required": false,
            "mapping": "unused",
            "description": "TSAI 图片不支持独立质量档。"
          },
          {
            "name": "watermark",
            "type": "boolean",
            "required": false,
            "mapping": "watermark",
            "description": "是否添加水印，默认 true。"
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
            "message": "TSAI Seedream 需要填写 prompt"
          },
          {
            "assert": {
              "$lte": [
                {
                  "$len": {
                    "$ref": "request.prompt"
                  }
                },
                20000
              ]
            },
            "message": "TSAI Seedream 的 prompt 最多 20000 个字符"
          }
        ],
        "create": {
          "method": "POST",
          "path": "/v1/images/generations",
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
                "size": {
                  "$omitEmpty": {
                    "$switch": {
                      "cases": [
                        {
                          "when": {
                            "$eq": [
                              {
                                "$lower": {
                                  "$ref": "request.aspectRatio"
                                }
                              },
                              "2k"
                            ]
                          },
                          "then": "2K"
                        },
                        {
                          "when": {
                            "$eq": [
                              {
                                "$lower": {
                                  "$ref": "request.aspectRatio"
                                }
                              },
                              "4k"
                            ]
                          },
                          "then": "4K"
                        },
                        {
                          "when": {
                            "$eq": [
                              {
                                "$lower": {
                                  "$ref": "request.aspectRatio"
                                }
                              },
                              "auto"
                            ]
                          },
                          "then": "2K"
                        },
                        {
                          "when": {
                            "$eq": [
                              {
                                "$ref": "request.aspectRatio"
                              },
                              "1:1"
                            ]
                          },
                          "then": "2048x2048"
                        },
                        {
                          "when": {
                            "$eq": [
                              {
                                "$ref": "request.aspectRatio"
                              },
                              "16:9"
                            ]
                          },
                          "then": "2560x1440"
                        },
                        {
                          "when": {
                            "$eq": [
                              {
                                "$ref": "request.aspectRatio"
                              },
                              "9:16"
                            ]
                          },
                          "then": "1440x2560"
                        },
                        {
                          "when": {
                            "$eq": [
                              {
                                "$ref": "request.aspectRatio"
                              },
                              "4:3"
                            ]
                          },
                          "then": "2304x1728"
                        },
                        {
                          "when": {
                            "$eq": [
                              {
                                "$ref": "request.aspectRatio"
                              },
                              "3:4"
                            ]
                          },
                          "then": "1728x2304"
                        },
                        {
                          "when": {
                            "$eq": [
                              {
                                "$ref": "request.aspectRatio"
                              },
                              "3:2"
                            ]
                          },
                          "then": "2496x1664"
                        },
                        {
                          "when": {
                            "$eq": [
                              {
                                "$ref": "request.aspectRatio"
                              },
                              "2:3"
                            ]
                          },
                          "then": "1664x2496"
                        },
                        {
                          "when": {
                            "$eq": [
                              {
                                "$ref": "request.aspectRatio"
                              },
                              "21:9"
                            ]
                          },
                          "then": "3024x1296"
                        }
                      ],
                      "default": {
                        "$ref": "request.aspectRatio"
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
                              "$ref": "request.images"
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
                                "$ref": "request.images"
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
                              "$ref": "request.images"
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
                "sequential_image_generation": "disabled",
                "response_format": "url",
                "stream": false,
                "watermark": {
                  "$omitEmpty": {
                    "$ref": "request.watermark"
                  }
                }
              },
              {
                "$coalesce": [
                  {
                    "$ref": "request.providerOptions.tsai-seedream.body"
                  },
                  {
                    "$ref": "request.providerOptions.tsai-seedream.extra_body"
                  },
                  {}
                ]
              }
            ]
          }
        },
        "poll": {
          "method": "GET",
          "path": "/v1/images/generations/{{taskId}}",
          "contentType": "application/json"
        },
        "response": {
          "taskId": {
            "$coalesce": [
              {
                "$ref": "response.id"
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
              "succeeded"
            ]
          },
          "message": {
            "$coalesce": [
              {
                "$ref": "response.error.message"
              },
              {
                "$ref": "response.message"
              }
            ]
          },
          "images": {
            "$ref": "response.data"
          },
          "errorPaths": [
            "error.code"
          ],
          "resultEphemeral": false
        }
      }
    ]
  },
  "documentation": "<当前插件的完整 documentation，由 README.md 与 docs/interface.md 拼接而成；为避免 JSON 递归，此处不重复展开正文。>"
}
```
<!-- YINGCE_MANIFEST_CONTRACT_END -->
