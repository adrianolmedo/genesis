package rest

import (
	"encoding/json"
	"testing"
)

func TestResp(t *testing.T) {
	b, err := json.Marshal(resp{
		Status: "Success",
		detailsResp: detailsResp{
			Message: "Hello world",
		},
	})
	if err != nil {
		t.Fatalf("failed to marshal resp: %v", err)
	}
	expected := `{"status":"Success","message":"Hello world"}`
	if string(b) != expected {
		t.Errorf("unexpected resp JSON: got %s, want %s", string(b), expected)
	}
}
