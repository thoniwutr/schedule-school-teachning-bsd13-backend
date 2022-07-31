package util

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	c "github.com/Beam-Data-Company/merchant-config-svc/constant"
)

// HandlerUtil can be used to write generic reponses back to the caller or any errors encountered
type HandlerUtil struct {
	Log *Logger
}

// NewHandlerUtil is a constructor for HandlerUtil
func NewHandlerUtil(log *Logger) *HandlerUtil {
	return &HandlerUtil{log}
}

// APIResponse is the structure of the response to write back
type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// HTTPSuccess writes the `msg` and the `code` to the response
func (hh *HandlerUtil) HTTPSuccess(w http.ResponseWriter, msg string, code int) {
	resp := &APIResponse{
		Message: msg,
		Code:    code,
	}
	hh.writeResponse(w, resp)
}

// HTTPError simply writes the `err.Error()` and the `code` to the response
func (hh *HandlerUtil) HTTPError(w http.ResponseWriter, err error, code int) {
	resp := &APIResponse{
		Message: err.Error(),
		Code:    code,
	}
	hh.logError(resp, err)
	hh.writeResponse(w, resp)
}

// WrappedError determines the http response code is from the error that was wrapped by the caller
// i.hh. if you wrap ErrConflict, it will write http.StatusConflict back to the user
// by default, the reponse code is http.StatusInternalServerError
// Wrapping can be done as so: `fmt.Errorf("wrapped error: %w", err)` where go automatically wraps the
// error under `%w` syntax (ref https://blog.golang.org/go1.13-errors)
func (hh *HandlerUtil) WrappedError(w http.ResponseWriter, err error) {
	code := c.ErrToHTTPCode(err)
	resp := &APIResponse{
		Message: err.Error(),
		Code:    code,
	}
	hh.logError(resp, err)
	hh.writeResponse(w, resp)
}

func (hh *HandlerUtil) logError(resp *APIResponse, err error) {
	hh.Log.Warn().
		Err(err).
		Caller(2).
		Bytes("trace", debug.Stack()).
		Int("status", resp.Code).
		Str("errMsg", resp.Message).
		Msg("unsuccessful operation")
}

// writeResponse is a helper function to write the ErrorResponse to the ResponseWriter
func (hh *HandlerUtil) writeResponse(w http.ResponseWriter, resp *APIResponse) {
	bytes, err := json.Marshal(resp)
	if err != nil {
		// we do not really expect marshalling of the response to fail
		hh.Log.Error().
			Stack().
			Err(err).
			Str("respMsg", resp.Message).
			Int("respCode", resp.Code).
			Msg("failed to marshal response")
	}
	w.WriteHeader(resp.Code)
	w.Write(bytes)
}
