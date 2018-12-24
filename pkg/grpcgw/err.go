package grpcgw

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/go-playground/validator.v9"

	"github.com/pkg/errors"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/grpclog"
)

// GWError is used for the error returned from the grpc implementation
// it can handle custom errors
type GWError interface {
	error
	// Status is the http status code
	Status() int
	// Message to outside user
	Message() string
}

type gwError struct {
	error `json:"-"`
	Msg   string            `json:"message"`
	S     int               `json:"status"`
	F     map[string]string `json:"fields,omitempty"`
}

func (gw *gwError) Status() int {
	return gw.S
}

func (gw *gwError) Message() string {
	return gw.Msg
}

// NewNotFound return not found error
func NewNotFound(err error) error {
	return NewBadRequestStatus(err, "Not found", http.StatusNotFound)
}

// NewBadRequest is the bad request
func NewBadRequest(err error, message string) error {
	return NewBadRequestStatus(err, message, http.StatusBadRequest)
}

// NewBadRequestStatus is the bad request
func NewBadRequestStatus(err error, message string, status int) error {
	ret := &gwError{
		error: errors.Wrap(err, message),
		Msg:   message,
		S:     status,
	}
	if v, ok := err.(validator.ValidationErrors); ok {
		ret.F = make(map[string]string)
		for _, fld := range v {
			ret.F[fld.Field()] = fld.Tag()
		}
	}
	return ret
}

// defaultHTTPError is my first try to overwrite the default
func defaultHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType())

	body, ok := err.(GWError)
	if !ok {
		body = &gwError{
			Msg: "unknown",
			S:   http.StatusInternalServerError,
		}
	}

	buf, merr := marshaler.Marshal(body)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", body, merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Infof("Failed to extract ServerMetadata from context")
	}

	w.WriteHeader(body.Status())
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}

	for k, vs := range md.TrailerMD {
		tKey := fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k)
		for _, v := range vs {
			w.Header().Add(tKey, v)
		}
	}
}

func init() {
	runtime.HTTPError = defaultHTTPError
}
