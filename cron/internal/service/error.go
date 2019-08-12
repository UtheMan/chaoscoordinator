package service

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
)

type InvalidRequestResponse struct {
	Err            error  `json:"-"`               // low-level runtime error
	HTTPStatusCode int    `json:"-"`               // http response status code
	StatusText     string `json:"status"`          // user-level status message
	AppCode        int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText      string `json:"error,omitempty"` // application-level error message, for debugging
}

type ContentNotFoundResponse struct {
	Err            error  `json:"-"`               // low-level runtime error
	HTTPStatusCode int    `json:"-"`               // http response status code
	StatusText     string `json:"status"`          // user-level status message
	AppCode        int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText      string `json:"error,omitempty"` // application-level error message, for debugging
}

func InvalidRequest(e error) render.Renderer {
	return &InvalidRequestResponse{
		Err:            e,
		HTTPStatusCode: 400,
		StatusText:     "Invalid Request",
		ErrorText:      e.Error(),
	}
}

func ContentNotFoundRequest(e error) render.Renderer {
	return &ContentNotFoundResponse{
		Err:            e,
		HTTPStatusCode: 204,
		StatusText:     "Content not found",
		ErrorText:      e.Error(),
	}
}

func InvalidRender(err error) render.Renderer {
	return &InvalidRequestResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func (e *InvalidRequestResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func (e *ContentNotFoundResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func (s *CronJobService) ParamsNotPresent(w http.ResponseWriter, r *http.Request) {
	err := errors.New("name or namespace must be specified")
	render.Render(w, r, InvalidRequest(err))
}
