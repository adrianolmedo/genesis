package http

type respOk struct {
	Msg string `json:"ok"`
}

type respError struct {
	Msg string `json:"error"`
}

type respData struct {
	Data any `json:"data"`
}

type respOkData struct {
	Msg  string `json:"ok"`
	Data any    `json:"data"`
}

type respErrorData struct {
	Msg  string `json:"error"`
	Data any    `json:"data"`
}

type respMetaData struct {
	Meta  any `json:"meta"`
	Data  any `json:"data"`
	Links any `json:"links,omitempty"`
}
