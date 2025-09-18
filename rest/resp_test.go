package rest

import (
	"encoding/json"
	"testing"
)

// TestResp test that resp marshals to the expected JSON structure.
// This is important because the struct is used in API responses.
// If the JSON structure changes, it could break clients.
// The test ensures that any changes to the struct are intentional.
func TestResp(t *testing.T) {
	t.Parallel()
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
