# TSAI 统一协议插件

该源码目录构建为一个版本化 `.yingce-plugin` 包，使用统一的 `yingce.plugin/v1` Manifest 和宿主声明式 Provider runtime。插件本身不执行 Web 代码；上游请求、鉴权、媒体公网 URL、任务轮询和结果保存由宿主处理。

随应用自动注册的官方包使用保留插件 ID `tsai`。如果要在「插件管理」中手动上传，请使用 [uploadable/tsai-upload.yingce-plugin](../uploadable/tsai-upload.yingce-plugin)，它使用独立的插件和 Provider ID，不会覆盖内置记录。

当前协议以 [TSAI API 文档](https://api.alifenqi.com:8188/app/#docs) 为准。渠道 Base URL 填 `https://api.alifenqi.com:8188`，不要重复添加 `/v1`。API Key 只填写在渠道密钥配置，不放入 Manifest、URL 或任务正文。

插件版本：`1.3.0`。当前覆盖 `/v1/models` 中的对话、视频、图片和音频模型：

- 对话：`deepseek-v4-flash`、`deepseek-v4-pro`，使用 `/v1/chat/completions`。
- 视频：`seedance-2.0`、`doubao-seedance-2-0-fast-260128`、`doubao-seedance-2-0-mini-260615`、`MiniMax/MiniMax-H3`、`custom/wan3.0-video`，使用 `/v1/videos/generations` 和 `/v1/videos/{task_id}`。
- 图片：`gpt-image-2`、`banana-2.1`、`doubao-seedream-4-5-251128`，使用 `/v1/images/generations` 和 `/v1/images/generations/{task_id}`。
- 音频：`seed-audio`，使用 `/v1/audio/generations` 和 `/v1/audio/generations/{task_id}`。

接口字段、枚举、响应与错误说明见 [docs/interface.md](docs/interface.md)。Seed Audio 的请求字段和查询路径按当前文档映射；查询响应没有公开示例，音频任务 ID、状态及结果字段采用多路径兼容映射，需用真实任务确认。
