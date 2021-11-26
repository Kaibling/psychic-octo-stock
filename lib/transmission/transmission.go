package transmission

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/lucsky/cuid"
)

type Envelope struct {
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
	RequestID string      `json:"request_id"`
}

type Response struct {
	requestID string
	w         http.ResponseWriter
	r         *http.Request
	envelope  *Envelope
}

func NewResponse(w http.ResponseWriter, r *http.Request) *Response {

	var requestID string
	clientRequestID := r.Header.Get("X-REQUEST-ID")
	if clientRequestID != "" {
		requestID = clientRequestID
	} else {
		requestID = cuid.New()
	}
	envelope := &Envelope{RequestID: requestID}
	return &Response{requestID: requestID, envelope: envelope, r: r, w: w}
}

func (s *Response) Send(data interface{}, message string, httpStatus int) {
	s.envelope.Data = data
	s.envelope.Message = message

	var sendEnv *Envelope
	if httpStatus == http.StatusNoContent {
		sendEnv = nil
	} else {
		sendEnv = s.envelope
	}
	render.Status(s.r, httpStatus)
	render.Respond(s.w, s.r, sendEnv)
}

// func SendResponse(w http.ResponseWriter, r *http.Request, data *Envelope, httpStatusCode int) {
// 	render.Status(r, httpStatusCode)

// 	render.Respond(w, r, data)
// }

func GetOrCreateResponse(w http.ResponseWriter, r *http.Request) *Response {
	parameter := r.Context().Value("responseObject")
	if parameter == nil {
		return NewResponse(w, r)
	}
	return parameter.(*Response)
}
