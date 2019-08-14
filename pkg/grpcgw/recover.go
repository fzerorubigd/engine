package grpcgw

import (
	"net/http"

	"elbix.dev/engine/pkg/log"
	"elbix.dev/engine/pkg/sentry"
)

// Recover is used for recovering from panic in http.ServeHTTP
func Recover(w http.ResponseWriter) {
	e := recover()
	if e == nil {
		return
	}
	_ = sentry.Recover(e)
	log.Error("Recover from panic", log.Any("panic", e))

	w.WriteHeader(http.StatusInternalServerError)
}
