package grpcgw

import (
	"net/http"

	"github.com/fzerorubigd/engine/pkg/log"
)

type recoverHandler struct {
	original http.Handler
}

func (rh *recoverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer Recover(w)
	rh.original.ServeHTTP(w, r)
}

func newRecover(handler http.Handler) http.Handler {
	return &recoverHandler{original: handler}
}

// Recover is used for recovering from panic in http.ServeHTTP
func Recover(w http.ResponseWriter) {
	e := recover()
	if e == nil {
		return
	}
	log.Error("Recover from panic", log.Any("panic", e))

	w.WriteHeader(http.StatusInternalServerError)
}
