package rest

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewResponse(t *testing.T) {
	tt := []struct {
		name    string
		msgType string
		code    string
		content string
		data    interface{}
		want    string
	}{
		{
			name:    "test-1",
			msgType: MsgOK,
			code:    "OK200",
			content: "resource created",
			data:    nil,
			want:    `{"message_ok":{"code":"OK200","content":"resource created"}}`,
		},
		{
			name:    "test-2",
			msgType: MsgError,
			code:    "ER404",
			content: "resource not found",
			data:    nil,
			want:    `{"message_error":{"code":"ER404","content":"resource not found"}}`,
		},
	}

	for _, tc := range tt {
		resp := newResponse(tc.msgType, tc.code, tc.content, tc.data)
		got, err := json.Marshal(resp)
		if err != nil {
			t.Logf("error marshal: %v", err)
		}

		gotS := string(got)
		if ok := strings.EqualFold(gotS, tc.want); ok == false {
			t.Errorf("%s: it's not the same string \n want: %v \n got: %v", tc.name, tc.want, gotS)
		}
	}
}
