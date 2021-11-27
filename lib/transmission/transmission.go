package transmission

import (
	"net/http"

	"github.com/go-chi/render"
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

func NewResponse(w http.ResponseWriter, r *http.Request, requestID string) *Response {
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

func (s *Response) GetRequestId() string {
	return s.requestID
}

func GetResponse(r *http.Request) *Response {
	parameter := r.Context().Value("responseObject")
	return parameter.(*Response)
}
