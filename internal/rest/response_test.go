package rest

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewResponse(t *testing.T) {
	tt := []struct {
		name         string
		inputMsgType string
		inputCode    string
		inputContent string
		inputData    interface{}
		want         string
	}{
		{
			name:         "test-1",
			inputMsgType: msgOK,
			inputCode:    "OK200",
			inputContent: "resource created",
			inputData:    nil,
			want:         `{"message_ok":{"code":"OK200","content":"resource created"}}`,
		},
		{
			name:         "test-2",
			inputMsgType: msgError,
			inputCode:    "ER404",
			inputContent: "resource not found",
			inputData:    nil,
			want:         `{"message_error":{"code":"ER404","content":"resource not found"}}`,
		},
	}

	for _, tc := range tt {
		resp := newResponse(tc.inputMsgType, tc.inputCode, tc.inputContent, tc.inputData)
		got, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}

		if ok := strings.EqualFold(string(got), tc.want); ok == false {
			t.Errorf("%s: it's not the same string \n want: %v \n got: %v", tc.name, tc.want, string(got))
		}
	}
}
