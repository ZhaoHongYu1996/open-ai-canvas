# OpenAiBxinle / ZeroFA Seedance 视频

本目录是官方系统协议源码。打包后的 `openai-bxinle.yingce-plugin` 可在「插件管理」上传导入，也可随仓库 `plugin-packages/` 在启动时自动安装。本包是声明式协议，不依赖系统宿主 `host:` 执行器。

面向 ZeroFA / Seedance 兼容的 JSON `/v1/videos` 合同：创建与轮询走 JSON，成功后优先使用结果 URL，失败再请求 `GET /v1/videos/{id}/content`。

详细接口见 [docs/interface.md](docs/interface.md)。
