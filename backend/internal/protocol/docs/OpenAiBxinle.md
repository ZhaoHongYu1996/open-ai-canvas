# OpenAiBxinle / ZeroFA Seedance 视频

OpenAiBxinle 是影策内置的宿主实现视频协议，面向 ZeroFA / Seedance 兼容渠道。创建与轮询走 JSON `/v1/videos`，成功后优先下载结果 URL，失败再请求 `GET /v1/videos/{id}/content`。执行路径是 `host:opaibxinle`，不走声明式协议 runner。

## 接口与鉴权

{{OPERATIONS}}

```http
POST {channel_base_url}/v1/videos
GET  {channel_base_url}/v1/videos/{id}
GET  {channel_base_url}/v1/videos/{id}/content
Authorization: Bearer <API_KEY>
Content-Type: application/json
```

渠道 Base URL 应指向协议服务根；宿主实现会把相对路径 `/videos` 拼到上游。密钥只由后端渠道中转读取，不要放进浏览器 URL 或日志。

## 模型与能力边界

模型名必须填写渠道实际开放的视频模型 ID。默认时长范围是 4–15 秒整数，默认 5 秒；比例默认 `16:9`，分辨率默认 `720p`，并支持 `1080p`。参考素材上限为 9 张图片、3 个视频、3 个音频；参考音频必须同时带至少一张图片或一个视频。当前实现支持文生视频和图片参考视频，不把任意模型名当作可用。

## 参数与字段映射

{{PARAMETERS}}

宿主请求体使用上游字段：`duration`、`resolution`、`aspect_ratio`、`generate_audio` / `enable_audio`，以及 `image` / `first_frame` / `last_frame` / `reference_video` / `reference_audio`。`images` 会按角色拆到图片、首帧和尾帧；`videos` 与 `audios` 分别映射参考视频和参考音频。空字段不发送。

## 创建任务示例

```bash
curl -X POST "{channel_base_url}/v1/videos" \
  -H "Authorization: Bearer <API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{
    "model":"YOUR_VIDEO_MODEL",
    "prompt":"保持角色一致，人物转身看向镜头",
    "duration":5,
    "resolution":"720p",
    "aspect_ratio":"16:9",
    "generate_audio":true
  }'
```

## 轮询、下载与结果

成功状态接受 `succeeded`、`completed`、`success`、`done`。结果地址可出现在 `metadata.url`、`output.url` / `output.video_url`，或顶层 `url` / `video_url`。若结果 URL 下载失败，再请求 `GET /v1/videos/{id}/content`。上游地址通常有时效，成功后应立即转入影策资源保存。当前没有取消接口。

## 官方与兼容资料

- 该协议是影策针对 ZeroFA / Seedance 兼容 JSON 视频合同的宿主实现，不是 OpenAI Videos multipart 合同。
- 路径相近不等于字段一致：不要把 `newapi` 的 `input_reference` multipart 或 MiniMax `content[]` 套到本协议。
- 具体模型、价格和额度以实际渠道账户为准。

{{CONTRACT}}
