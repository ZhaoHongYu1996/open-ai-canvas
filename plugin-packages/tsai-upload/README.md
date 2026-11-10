# TSAI 统一协议插件（上传版）

该 `.yingce-plugin` 使用独立插件 ID `tsai-upload`，Provider ID 均带 `tsai-upload-` 前缀，不会覆盖随应用自动注册的官方系统插件。上传安装后，在渠道配置中选择 Provider 并填写 TSAI API Key。

包使用统一 `yingce.plugin/v1` Manifest 和宿主声明式运行时，不含 Web 入口代码，也不请求 `host:` 执行器。Base URL 填 `https://api.alifenqi.com:8188`，不要重复添加 `/v1`。API Key 只填写在渠道密钥配置中。

插件版本：`1.3.0`。当前覆盖 TSAI `/v1/models` 中的对话、视频、图片和音频模型：

- 对话：`deepseek-v4-flash`、`deepseek-v4-pro`。
- 视频：`seedance-2.0`、`doubao-seedance-2-0-fast-260128`、`doubao-seedance-2-0-mini-260615`、`MiniMax/MiniMax-H3`、`custom/wan3.0-video`。
- 图片：`gpt-image-2`、`banana-2.1`、`doubao-seedream-4-5-251128`。
- 音频：`seed-audio`。

详见 [TSAI 接口映射](docs/interface.md)。Seed Audio 的查询响应没有公开示例，任务 ID、状态与结果映射采用多路径兼容策略，需要用一次真实任务确认。
