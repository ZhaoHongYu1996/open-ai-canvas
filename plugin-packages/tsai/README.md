# TSAI 网关协议

本目录是随应用启动自动注册的官方系统协议源码。`tsai.yingce-plugin` 使用保留的官方插件 ID `tsai`，不能再通过「插件管理」上传覆盖；需要手动上传时，请使用 `plugin-packages/uploadable/tsai-upload.yingce-plugin`，该副本使用独立的插件和 Provider ID。本包是声明式协议，不依赖系统宿主 `host:` 执行器。

面向 [TSAI 开发者文档](https://api.alifenqi.com:8188/app/#docs)：统一 v1 Manifest 与一个 `.yingce-plugin` 包覆盖对话、MiniMax H3、Seedance 2.0 / 2.0 mini 视频，以及 Seedream 4.5 图片。渠道 Base URL 填 `https://api.alifenqi.com:8188`，不要再重复写 `/v1`。模型 ID 以 `/v1/models` 返回结果为准。


详细接口见 [docs/interface.md](docs/interface.md)。
