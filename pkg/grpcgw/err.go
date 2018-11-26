package grpcgw

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
	Msg string `json:"message"`
}

func (gw gwError) Error() string {
	return gw.Msg
}

func (gw gwError) Status() int {
	panic("implement me")
}

func (gw gwError) Message() string {
	panic("implement me")
}
