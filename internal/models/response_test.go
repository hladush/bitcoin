package models

import "testing"

func TestErrorResponse(t *testing.T) {
	resp := ErrorResponse("test error")
	if resp.Success {
		t.Error("Expected Success to be false")
	}
	if resp.Error != "test error" {
		t.Errorf("Expected Error to be 'test error', got '%s'", resp.Error)
	}
}

func TestSuccessResponse(t *testing.T) {
	data := map[string]string{"key": "value"}
	resp := SuccessResponse(data)
	if !resp.Success {
		t.Error("Expected Success to be true")
	}
	if resp.Data == nil {
		t.Error("Expected Data to not be nil")
	}
}

func TestMessageResponse(t *testing.T) {
	resp := MessageResponse("test message")
	if !resp.Success {
		t.Error("Expected Success to be true")
	}
	if resp.Message != "test message" {
		t.Errorf("Expected Message to be 'test message', got '%s'", resp.Message)
	}
}
