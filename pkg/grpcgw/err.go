package grpcgw

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
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
	error
	Msg string `json:"message"`
	S   int    `json:"status"`
}

func (gw *gwError) Status() int {
	return gw.S
}

func (gw *gwError) Message() string {
	return gw.Msg
}

// NewBadRequest is the bad request
func NewBadRequest(err error, message string) error {
	return NewBadRequestStatus(err, message, http.StatusBadRequest)
}

// NewBadRequestStatus is the bad request
func NewBadRequestStatus(err error, message string, status int) error {
	return &gwError{
		error: errors.Wrap(err, message),
		Msg:   message,
		S:     status,
	}
}

// defaultHTTPError is my first try to overwrite the default
func defaultHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`

	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType())

	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	re := false
	gw, ok := err.(GWError)
	if ok {
		re = true
		s = status.New(codes.Code(gw.Status()), gw.Message())
	}

	body := struct {
		Error string `protobuf:"bytes,1,name=error" json:"error"`
		// This is to make the error more compatible with users that expect errors to be Status objects:
		// https://github.com/grpc/grpc/blob/master/src/proto/grpc/status/status.proto
		// It should be the exact same message as the Error field.
		Message string     `protobuf:"bytes,1,name=message" json:"message"`
		Code    int32      `protobuf:"varint,2,name=code" json:"code"`
		Details []*any.Any `protobuf:"bytes,3,rep,name=details" json:"details,omitempty"`
	}{
		Error:   s.Message(),
		Message: s.Message(),
		Code:    int32(s.Code()),
		Details: s.Proto().GetDetails(),
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

	st := int(s.Code())
	if !re {
		st = runtime.HTTPStatusFromCode(s.Code())
	}
	w.WriteHeader(st)
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
