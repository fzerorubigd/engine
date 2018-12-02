package grpcgw

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/fzerorubigd/balloon/pkg/assert"
	"github.com/fzerorubigd/balloon/pkg/log"
)

type swaggerFile struct {
	Swagger string `json:"swagger"`
	Info    struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	} `json:"info"`
	Schemes     []string               `json:"schemes"`
	Consumes    []string               `json:"consumes"`
	Produces    []string               `json:"produces"`
	Paths       map[string]interface{} `json:"paths"`
	Definitions map[string]interface{} `json:"definitions"`
}

var (
	swaggerLock sync.RWMutex
	data        = swaggerFile{
		Swagger: "2.0",
		Info: struct {
			Title   string `json:"title"`
			Version string `json:"version"`
		}{Title: "Balloon Swagger", Version: "1.0"},
		Schemes:     []string{"http", "https"},
		Consumes:    []string{"application/json"},
		Produces:    []string{"application/json"},
		Paths:       make(map[string]interface{}),
		Definitions: make(map[string]interface{}),
	}
)

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	defer Recover(w)

	swaggerLock.RLock()
	defer swaggerLock.RUnlock()

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error("Failed to serve swagger", log.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// RegisterSwagger register a swagger end point
func RegisterSwagger(paths map[string]interface{}, definitions map[string]interface{}) {
	swaggerLock.Lock()
	defer swaggerLock.Unlock()

	for i := range paths {
		_, ok := data.Paths[i]
		assert.False(ok, "Path is already registered", i)
		data.Paths[i] = paths[i]
	}

	for i := range definitions {
		_, ok := data.Definitions[i]
		assert.False(ok, "Definition is already registered", i)
		data.Definitions[i] = definitions[i]
	}

}
