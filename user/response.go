package user

const (
	msgOK    = "ok"
	msgError = "error"
)

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

func newResponse(message, content string, data interface{}) response {
	var resp response

	switch message {
	case msgOK:
		resp = response{
			messageOK: &messageOK{
				Content: content,
			},
			messageError: nil,
			Data:         data,
		}
	case msgError:
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
