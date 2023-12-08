package http

const (
	msgOK    = "ok"
	msgError = "error"
)

// response struct for JSON standar response.
type response struct {
	*messageOK    `json:"messageOk,omitempty"`
	*messageError `json:"messageError,omitempty"`
	Data          interface{} `json:"data,omitempty"`
	Links         interface{} `json:"links,omitempty"`
	Meta          interface{} `json:"meta,omitempty"`
}

type messageOK struct {
	Content string `json:"content"`
}

type messageError struct {
	Content string `json:"content"`
}

// links set Links.
func (r response) setLinks(i interface{}) response {
	r.Links = i
	return r
}

// meta set Meta.
func (r response) setMeta(i interface{}) response {
	r.Meta = i
	return r
}

// respJSON return standar response JSON.
// Usage example: resp := respJSON(msgOK, "resource has been updated", data).
func respJSON(msgType, content string, data interface{}) response {
	var resp response

	switch msgType {
	case msgOK:
		mOK := &messageOK{
			Content: content,
		}

		if content == "" {
			mOK = nil
		}

		resp = response{
			messageOK:    mOK,
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

// Response represents standar JSON response.
type Response struct {
	Ok    string      `json:"ok,omitempty" example:"resource created"`
	Data  interface{} `json:"data,omitempty"`
	Meta  interface{} `json:"meta,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

type respOk struct {
	Msg string `json:"ok,omitempty"`
}

type respOkData struct {
	Msg  string      `json:"ok,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type respError struct {
	Msg string `json:"error,omitempty"`
}

type respErrorData struct {
	Msg  string      `json:"error,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type respMetaData struct {
	Meta interface{} `json:"meta,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
