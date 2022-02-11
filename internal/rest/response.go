package rest

const (
	msgError = "error"
	msgOK    = "ok"
)

type response struct {
	*messageOK    `json:"message_ok,omitempty"`
	*messageError `json:"message_error,omitempty"`
	Data          interface{} `json:"data,omitempty"`
}

type messageOK struct {
	Code    string `json:"code"`
	Content string `json:"content"`
}

type messageError struct {
	Code    string `json:"code"`
	Content string `json:"content"`
}

// newResponse return standar response JSON.
// Usage example: response := newResponse(MsgOK, "OK002", "resource has been updated", data).
func newResponse(message, code, content string, data interface{}) response {
	var resp response

	switch message {
	case msgOK:
		resp = response{
			messageOK: &messageOK{
				Code:    code,
				Content: content,
			},
			messageError: nil,
			Data:         data,
		}
	case msgError:
		resp = response{
			messageOK: nil,
			messageError: &messageError{
				Code:    code,
				Content: content,
			},
			Data: data,
		}
	}

	return resp
}
