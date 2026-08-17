package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const opAIBxinleVideoPollInterval = 5 * time.Second

type opAIBxinleVideoRequest struct {
	Model          string      `json:"model"`
	Prompt         string      `json:"prompt"`
	Mode           string      `json:"mode,omitempty"`
	NegativePrompt string      `json:"negative_prompt,omitempty"`
	Duration       *int        `json:"duration,omitempty"`
	Resolution     string      `json:"resolution,omitempty"`
	AspectRatio    string      `json:"aspect_ratio,omitempty"`
	BitrateMode    string      `json:"bitrate_mode,omitempty"`
	Seed           *int        `json:"seed,omitempty"`
	EnableAudio    *bool       `json:"enable_audio,omitempty"`
	GenerateAudio  *bool       `json:"generate_audio,omitempty"`
	Image          interface{} `json:"image,omitempty"`
	FirstFrame     string      `json:"first_frame,omitempty"`
	LastFrame      string      `json:"last_frame,omitempty"`
	ReferenceVideo interface{} `json:"reference_video,omitempty"`
	ReferenceAudio interface{} `json:"reference_audio,omitempty"`
}

func runOpAIBxinleVideoTask(ctx context.Context, input canvasGenerationInput) (map[string]interface{}, error) {
	id := resumedProviderRequestID(ctx)
	var created map[string]interface{}
	if id == "" {
		body, err := opAIBxinleVideoRequestBody(input)
		if err != nil {
			return nil, err
		}
		if err := postJSON(ctx, input.Config, "/videos", body, &created); err != nil {
			return nil, wrapOpAIBxinleHTTPError(err, "提交任务")
		}
		if data, ok := created["data"].(map[string]interface{}); ok {
			created = data
		}
		id = firstNonEmptyString(stringField(created, "id"), stringField(created, "task_id"))
	}
	if id == "" {
		return nil, errors.New("OpenAiBxinle 没有返回任务 ID")
	}
	for deadline := providerPollingDeadline(ctx); time.Now().Before(deadline); {
		var state map[string]interface{}
		if err := getJSON(ctx, input.Config, "/videos/"+id, &state); err != nil {
			return nil, wrapOpAIBxinleHTTPError(err, "查询任务")
		}
		if data, ok := state["data"].(map[string]interface{}); ok {
			state = data
		}
		status := strings.ToLower(strings.TrimSpace(stringField(state, "status")))
		if status == "succeeded" || status == "completed" || status == "success" || status == "done" {
			return downloadOpAIBxinleVideoResult(ctx, input.Config, id, state)
		}
		if status == "failed" || status == "error" {
			return nil, fmt.Errorf("OpenAiBxinle 视频生成失败（任务 %s）：%s", id, opAIBxinleErrorMessage(state))
		}
		if err := sleepContext(ctx, opAIBxinleVideoPollInterval); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("OpenAiBxinle 视频生成超时（任务 %s）", id)
}

func downloadOpAIBxinleVideoResult(ctx context.Context, config providerConfig, id string, state map[string]interface{}) (map[string]interface{}, error) {
	videoURL := opAIBxinleResultURL(state)
	if videoURL != "" {
		data, mimeType, err := getExternalBinary(withProviderRequestKind(ctx, "download"), videoURL)
		if err == nil && len(data) > 0 {
			return map[string]interface{}{"mode": "video", "video": map[string]interface{}{"dataUrl": dataURL(mimeType, data), "mimeType": defaultString(mimeType, "video/mp4"), "sourceUrl": videoURL, "expiresHint": "上游结果地址可能临时有效，请尽快预览或下载"}}, nil
		}
	}
	data, mimeType, err := getBinary(withProviderRequestKind(ctx, "download"), config, "/videos/"+id+"/content")
	if err != nil {
		if videoURL != "" {
			return nil, fmt.Errorf("OpenAiBxinle 任务 %s 结果地址与 content 端点均下载失败，地址可能已失效", id)
		}
		return nil, fmt.Errorf("OpenAiBxinle 任务 %s 已成功但没有返回可下载的视频地址", id)
	}
	return map[string]interface{}{"mode": "video", "video": map[string]interface{}{"dataUrl": dataURL(mimeType, data), "mimeType": defaultString(mimeType, "video/mp4"), "expiresHint": "上游结果地址可能临时有效，请尽快预览或下载"}}, nil
}

func opAIBxinleResultURL(state map[string]interface{}) string {
	if metadata, ok := state["metadata"].(map[string]interface{}); ok {
		if value := stringField(metadata, "url"); value != "" {
			return value
		}
	}
	if output, ok := state["output"].(map[string]interface{}); ok {
		if value := firstNonEmptyString(stringField(output, "url"), stringField(output, "video_url")); value != "" {
			return value
		}
	}
	return firstNonEmptyString(stringField(state, "url"), stringField(state, "video_url"), stringField(state, "object"))
}

func opAIBxinleErrorMessage(state map[string]interface{}) string {
	if errorValue, ok := state["error"].(map[string]interface{}); ok {
		if message := stringField(errorValue, "message"); message != "" {
			return message
		}
	}
	return defaultString(seedanceErrorMessage(state), defaultString(stringField(state, "error"), "上游返回 failed"))
}

func opAIBxinleVideoRequestBody(input canvasGenerationInput) (opAIBxinleVideoRequest, error) {
	if strings.TrimSpace(input.Prompt) == "" {
		return opAIBxinleVideoRequest{}, errors.New("OpenAiBxinle 需要填写 prompt")
	}
	if len(input.ReferenceImages) > 9 {
		return opAIBxinleVideoRequest{}, errors.New("OpenAiBxinle 最多支持 9 张参考图片")
	}
	if len(input.ReferenceVideos) > 3 {
		return opAIBxinleVideoRequest{}, errors.New("OpenAiBxinle 最多支持 3 个参考视频")
	}
	if len(input.ReferenceAudios) > 3 {
		return opAIBxinleVideoRequest{}, errors.New("OpenAiBxinle 最多支持 3 个参考音频")
	}
	if len(input.ReferenceAudios) > 0 && len(input.ReferenceImages) == 0 && len(input.ReferenceVideos) == 0 {
		return opAIBxinleVideoRequest{}, errors.New("OpenAiBxinle 的参考音频必须与至少一个图片或视频参考素材同时使用")
	}

	imageURLs := make([]string, 0, len(input.ReferenceImages))
	for _, image := range input.ReferenceImages {
		url, err := openAIImageInputURL(image)
		if err != nil {
			return opAIBxinleVideoRequest{}, err
		}
		imageURLs = append(imageURLs, url)
	}
	videoURLs := make([]string, 0, len(input.ReferenceVideos))
	for _, video := range input.ReferenceVideos {
		url, err := seedanceVideosMediaURL(video)
		if err != nil {
			return opAIBxinleVideoRequest{}, err
		}
		videoURLs = append(videoURLs, url)
	}
	audioURLs := make([]string, 0, len(input.ReferenceAudios))
	for _, audio := range input.ReferenceAudios {
		url, err := seedanceVideosMediaURL(audio)
		if err != nil {
			return opAIBxinleVideoRequest{}, err
		}
		audioURLs = append(audioURLs, url)
	}

	startFrameID := metadataString(input.Metadata, "videoStartFrameNodeId")
	endFrameID := metadataString(input.Metadata, "videoEndFrameNodeId")
	var firstFrame, lastFrame string
	regularImages := make([]string, 0, len(imageURLs))
	for index, image := range input.ReferenceImages {
		if index >= len(imageURLs) {
			continue
		}
		switch image.ID {
		case startFrameID:
			if startFrameID != "" {
				firstFrame = imageURLs[index]
				continue
			}
		case endFrameID:
			if endFrameID != "" {
				lastFrame = imageURLs[index]
				continue
			}
		}
		regularImages = append(regularImages, imageURLs[index])
	}
	if startFrameID != "" && firstFrame == "" {
		return opAIBxinleVideoRequest{}, errors.New("已配置的首帧参考图未包含在视频请求中")
	}
	if endFrameID != "" && lastFrame == "" {
		return opAIBxinleVideoRequest{}, errors.New("已配置的尾帧参考图未包含在视频请求中")
	}

	body := opAIBxinleVideoRequest{
		Model:  strings.TrimSpace(input.Config.Model),
		Prompt: strings.TrimSpace(input.Prompt),
	}
	if body.Model == "" {
		return opAIBxinleVideoRequest{}, errors.New("OpenAiBxinle 需要配置模型 ID")
	}
	if duration, ok := opAIBxinleDuration(input.Config.VideoSeconds); ok {
		body.Duration = &duration
	}
	if resolution := opAIBxinleResolution(input.Config.VQuality); resolution != "" {
		body.Resolution = resolution
	}
	if ratio := opAIBxinleAspectRatio(input.Config.Size); ratio != "" {
		body.AspectRatio = ratio
	}
	if videoCapabilitySupportsAudio(input) {
		value := parseBool(input.Config.VideoGenerateAudio, true)
		body.EnableAudio = &value
	}

	if firstFrame != "" {
		body.FirstFrame = firstFrame
	}
	if lastFrame != "" {
		body.LastFrame = lastFrame
	}
	if imageValue := opAIBxinleMediaValue(regularImages); imageValue != nil {
		body.Image = imageValue
	}
	if videoValue := opAIBxinleMediaValue(videoURLs); videoValue != nil {
		body.ReferenceVideo = videoValue
	}
	if audioValue := opAIBxinleMediaValue(audioURLs); audioValue != nil {
		body.ReferenceAudio = audioValue
	}
	body.Mode = opAIBxinleMode(body)
	return body, nil
}

func opAIBxinleMode(body opAIBxinleVideoRequest) string {
	if body.FirstFrame != "" && body.LastFrame != "" {
		return "first_last_frame"
	}
	if body.ReferenceVideo != nil {
		return "video2video"
	}
	if images, ok := body.Image.([]string); ok && len(images) > 1 {
		return "reference2video"
	}
	if body.Image != nil || body.FirstFrame != "" {
		return "image2video"
	}
	return "text2video"
}

func opAIBxinleMediaValue(values []string) interface{} {
	if len(values) == 0 {
		return nil
	}
	if len(values) == 1 {
		return values[0]
	}
	return values
}

func opAIBxinleDuration(value string) (int, bool) {
	seconds, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || seconds <= 0 {
		return 0, false
	}
	return seconds, true
}

func opAIBxinleResolution(value string) string {
	resolution := strings.ToUpper(strings.TrimSpace(value))
	resolution = strings.ReplaceAll(resolution, "P", "")
	if resolution == "" {
		return ""
	}
	return resolution + "P"
}

func opAIBxinleAspectRatio(value string) string {
	switch strings.TrimSpace(value) {
	case "1:1", "16:9", "9:16", "4:3", "3:4", "21:9":
		return strings.TrimSpace(value)
	default:
		return ""
	}
}

func wrapOpAIBxinleHTTPError(err error, action string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) {
		return err
	}
	if errors.Is(err, context.DeadlineExceeded) || strings.Contains(strings.ToLower(err.Error()), "timeout") {
		return fmt.Errorf("OpenAiBxinle %s超时", action)
	}
	var httpErr providerHTTPError
	if errors.As(err, &httpErr) {
		switch httpErr.StatusCode {
		case http.StatusUnauthorized:
			return fmt.Errorf("OpenAiBxinle %s失败：认证无效（401）", action)
		case http.StatusPaymentRequired:
			return fmt.Errorf("OpenAiBxinle %s失败：余额不足（402）", action)
		case http.StatusTooManyRequests:
			return fmt.Errorf("OpenAiBxinle %s失败：请求过于频繁（429）", action)
		default:
			if httpErr.StatusCode >= 500 {
				return fmt.Errorf("OpenAiBxinle %s失败：上游服务异常（%d）", action, httpErr.StatusCode)
			}
			return fmt.Errorf("OpenAiBxinle %s失败：上游 HTTP %d", action, httpErr.StatusCode)
		}
	}
	return err
}
