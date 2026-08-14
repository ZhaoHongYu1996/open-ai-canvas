package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpAIBxinleVideoRequestBodyText2Video(t *testing.T) {
	enableAudio := true
	input := canvasGenerationInput{
		Prompt: "黄金时刻的电影感跟拍镜头",
		Config: providerConfig{Model: "doubao-seedance-2.0", VideoGenerateAudio: "true"},
	}
	body, err := opAIBxinleVideoRequestBody(input)
	if err != nil {
		t.Fatal(err)
	}
	if body.Model != "doubao-seedance-2.0" || body.Prompt == "" || body.Mode != "text2video" {
		t.Fatalf("unexpected body: %+v", body)
	}
	if body.Duration != nil || body.Resolution != "" || body.AspectRatio != "" {
		t.Fatalf("optional fields should be omitted when unset: %+v", body)
	}
	if body.EnableAudio == nil || *body.EnableAudio != enableAudio {
		t.Fatalf("enable_audio default true, got %+v", body.EnableAudio)
	}
	raw, _ := json.Marshal(body)
	if strings.Contains(string(raw), "negative_prompt") || strings.Contains(string(raw), "bitrate_mode") || strings.Contains(string(raw), "image") {
		t.Fatalf("unexpected omitted fields present: %s", raw)
	}
}

func TestOpAIBxinleVideoRequestBodyModes(t *testing.T) {
	image := providerMedia{ID: "img-1", DataURL: "data:image/png;base64,aaaa", URL: "https://example.com/a.png"}
	image2 := providerMedia{ID: "img-2", DataURL: "data:image/png;base64,bbbb", URL: "https://example.com/b.png"}
	start := providerMedia{ID: "first", DataURL: "data:image/png;base64,cccc", URL: "https://example.com/first.png"}
	end := providerMedia{ID: "last", DataURL: "data:image/png;base64,dddd", URL: "https://example.com/last.png"}
	video := providerMedia{ID: "vid-1", URL: "https://example.com/ref.mp4"}

	t.Run("image2video", func(t *testing.T) {
		body, err := opAIBxinleVideoRequestBody(canvasGenerationInput{
			Prompt:          "镜头",
			Config:          providerConfig{Model: "doubao-seedance-2.0", VideoSeconds: "8", Size: "1:1", VQuality: "720p"},
			ReferenceImages: []providerMedia{image},
		})
		if err != nil {
			t.Fatal(err)
		}
		if body.Mode != "image2video" || body.Image == nil || body.Duration == nil || *body.Duration != 8 || body.Resolution != "720P" || body.AspectRatio != "1:1" {
			t.Fatalf("unexpected image2video body: %+v", body)
		}
	})

	t.Run("reference2video", func(t *testing.T) {
		body, err := opAIBxinleVideoRequestBody(canvasGenerationInput{
			Prompt:          "镜头",
			Config:          providerConfig{Model: "doubao-seedance-2.0"},
			ReferenceImages: []providerMedia{image, image2},
		})
		if err != nil {
			t.Fatal(err)
		}
		if body.Mode != "reference2video" {
			t.Fatalf("mode=%s", body.Mode)
		}
		images, ok := body.Image.([]string)
		if !ok || len(images) != 2 {
			t.Fatalf("image should be array, got %#v", body.Image)
		}
	})

	t.Run("first_last_frame", func(t *testing.T) {
		body, err := opAIBxinleVideoRequestBody(canvasGenerationInput{
			Prompt:          "镜头",
			Config:          providerConfig{Model: "doubao-seedance-2.0"},
			ReferenceImages: []providerMedia{start, end},
			Metadata:        map[string]interface{}{"videoStartFrameNodeId": "first", "videoEndFrameNodeId": "last"},
		})
		if err != nil {
			t.Fatal(err)
		}
		if body.Mode != "first_last_frame" || body.FirstFrame == "" || body.LastFrame == "" {
			t.Fatalf("unexpected first_last_frame body: %+v", body)
		}
	})

	t.Run("video2video", func(t *testing.T) {
		body, err := opAIBxinleVideoRequestBody(canvasGenerationInput{
			Prompt:          "镜头",
			Config:          providerConfig{Model: "doubao-seedance-2.0"},
			ReferenceVideos: []providerMedia{video},
		})
		if err != nil {
			t.Fatal(err)
		}
		if body.Mode != "video2video" || body.ReferenceVideo == nil {
			t.Fatalf("unexpected video2video body: %+v", body)
		}
	})
}

func TestOpAIBxinleVideoRequestBodyRejectsAudioOnly(t *testing.T) {
	_, err := opAIBxinleVideoRequestBody(canvasGenerationInput{
		Prompt:          "镜头",
		Config:          providerConfig{Model: "doubao-seedance-2.0"},
		ReferenceAudios: []providerMedia{{ID: "a1", URL: "https://example.com/a.mp3"}},
	})
	if err == nil || !strings.Contains(err.Error(), "参考音频") {
		t.Fatalf("expected audio constraint, got %v", err)
	}
}

func TestOpAIBxinleVideoTaskPollAndDownload(t *testing.T) {
	t.Setenv("CANVAS_ALLOW_PRIVATE_UPSTREAMS", "true")
	var createdBody opAIBxinleVideoRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/v1/videos" {
			if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if err := json.NewDecoder(r.Body).Decode(&createdBody); err != nil {
				t.Errorf("decode create: %v", err)
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"id": "task-1", "status": "queued"})
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/v1/videos/task-1" {
			_ = json.NewEncoder(w).Encode(map[string]string{"id": "task-1", "status": "succeeded"})
			return
		}
		if r.Method == http.MethodGet && r.URL.Path == "/v1/videos/task-1/content" {
			w.Header().Set("Content-Type", "video/mp4")
			_, _ = w.Write([]byte("mp4-bytes"))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	result, err := runOpAIBxinleVideoTask(context.Background(), canvasGenerationInput{
		Prompt: "黄金时刻的电影感跟拍镜头",
		Config: providerConfig{BaseURL: server.URL, APIKey: "sk-zerofa-test", Model: "doubao-seedance-2.0", VideoGenerateAudio: "true"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if createdBody.Model != "doubao-seedance-2.0" || createdBody.Mode != "text2video" {
		t.Fatalf("create body: %+v", createdBody)
	}
	video, _ := result["video"].(map[string]interface{})
	if video == nil || !strings.Contains(stringField(video, "dataUrl"), "video/mp4") {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestWrapOpAIBxinleHTTPError(t *testing.T) {
	if err := wrapOpAIBxinleHTTPError(providerHTTPError{StatusCode: 401, Status: "401"}, "提交任务"); err == nil || !strings.Contains(err.Error(), "401") {
		t.Fatalf("401: %v", err)
	}
	if err := wrapOpAIBxinleHTTPError(providerHTTPError{StatusCode: 402, Status: "402"}, "提交任务"); err == nil || !strings.Contains(err.Error(), "402") {
		t.Fatalf("402: %v", err)
	}
	if err := wrapOpAIBxinleHTTPError(providerHTTPError{StatusCode: 429, Status: "429"}, "查询任务"); err == nil || !strings.Contains(err.Error(), "429") {
		t.Fatalf("429: %v", err)
	}
	if err := wrapOpAIBxinleHTTPError(providerHTTPError{StatusCode: 503, Status: "503"}, "查询任务"); err == nil || !strings.Contains(err.Error(), "503") {
		t.Fatalf("503: %v", err)
	}
	if err := wrapOpAIBxinleHTTPError(context.DeadlineExceeded, "查询任务"); err == nil || !strings.Contains(err.Error(), "超时") {
		t.Fatalf("timeout: %v", err)
	}
}

func TestValidateGenerationInterfaceOpAIBxinle(t *testing.T) {
	if err := validateGenerationInterface("video", "opAIBxinle"); err != nil {
		t.Fatal(err)
	}
	if err := validateGenerationInterface("image", "opAIBxinle"); err == nil {
		t.Fatal("image should reject video protocol")
	}
}

