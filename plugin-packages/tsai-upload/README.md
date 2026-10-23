# TSAI 网关（上传版）

这是供「插件管理 → 上传插件」使用的 `.yingce-plugin` 包，插件 ID 为 `tsai-upload`，Provider ID 均带 `tsai-upload-` 前缀，不会覆盖随应用自动注册的官方 `tsai` 插件。上传后在系统渠道中选择相应 Provider，并填写 TSAI API Key。

本包是统一 `yingce.plugin/v1` 声明式协议，包含 TSAI 对话、MiniMax H3、Seedance 2.0 / 2.0 mini 视频及 Seedream 4.5 图片能力；不含脚本入口，不使用 `host:` 执行器。Base URL 填 `https://api.alifenqi.com:8188`，不要重复填写 `/v1`。模型 ID 以 TSAI `/v1/models` 返回值为准。

测试请求必须包含非空 prompt。



详细接口见 [docs/interface.md](docs/interface.md)。
