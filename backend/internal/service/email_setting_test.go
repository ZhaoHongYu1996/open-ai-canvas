package service

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEmailSettingRequestAcceptsNumericAndStringPort(t *testing.T) {
	var numeric EmailSettingRequest
	if err := json.Unmarshal([]byte(`{"port":465}`), &numeric); err != nil {
		t.Fatal(err)
	}
	if int(numeric.Port) != 465 {
		t.Fatalf("numeric port = %d, want 465", numeric.Port)
	}

	var text EmailSettingRequest
	if err := json.Unmarshal([]byte(`{"port":"587"}`), &text); err != nil {
		t.Fatal(err)
	}
	if int(text.Port) != 587 {
		t.Fatalf("string port = %d, want 587", text.Port)
	}
}

func TestEmailSettingRequestRejectsInvalidPort(t *testing.T) {
	var req EmailSettingRequest
	err := json.Unmarshal([]byte(`{"port":"smtp"}`), &req)
	if err == nil || !strings.Contains(err.Error(), "SMTP 端口必须是整数") {
		t.Fatalf("error = %v", err)
	}
}
