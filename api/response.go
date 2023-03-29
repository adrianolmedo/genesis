package api

const (
	MsgOK    = "ok"
	MsgError = "error"
)

// response it's a struct for JSON response.
type response struct {
	*messageOK    `json:"message_ok,omitempty"`
	*messageError `json:"message_error,omitempty"`
	Data          interface{} `json:"data,omitempty"`
}

type messageOK struct {
	Content string `json:"content"`
}

type messageError struct {
	Content string `json:"content"`
}

// RespJSON return standar response JSON.
// Usage example: response := RespJSON(MsgOK, "resource has been updated", data).
func RespJSON(message, content string, data interface{}) response {
	var resp response

	switch message {
	case MsgOK:
		resp = response{
			messageOK: &messageOK{
				Content: content,
			},
			messageError: nil,
			Data:         data,
		}
	case MsgError:
		resp = response{
			messageOK: nil,
			messageError: &messageError{
				Content: content,
			},
			Data: data,
		}
	}

	return resp
}
