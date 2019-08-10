package grpcgw

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"elbix.dev/engine/pkg/assert"
	"elbix.dev/engine/pkg/log"
)

var (
	errDef    map[string]interface{}
	errFldDef map[string]interface{}
	err400    = map[string]interface{}{
		"description": "Input error",
		"schema": map[string]interface{}{
			"$ref": "#/definitions/ErrorFldResponse",
		},
	}
	err401 = map[string]interface{}{
		"description": "Returned when not authenticated",
		"schema": map[string]interface{}{
			"$ref": "#/definitions/ErrorResponse",
		},
	}
	err403 = map[string]interface{}{
		"description": "Returned when not authorized",
		"schema": map[string]interface{}{
			"$ref": "#/definitions/ErrorResponse",
		},
	}
	err404 = map[string]interface{}{
		"description": "Returned when the route is not correct",
		"schema": map[string]interface{}{
			"$ref": "#/definitions/ErrorResponse",
		},
	}
)

type swaggerFile struct {
	Swagger string `json:"swagger"`
	Info    struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	} `json:"info"`

	Host        string                 `json:"host"`
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
		}{Title: "Engine Swagger", Version: "1.0"},
		Schemes:     []string{"https", "http"},
		Consumes:    []string{"application/json"},
		Produces:    []string{"application/json"},
		Host:        "",
		Paths:       make(map[string]interface{}),
		Definitions: make(map[string]interface{}),
	}
)

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	defer Recover(w)

	swaggerLock.RLock()
	defer swaggerLock.RUnlock()

	fl := strings.TrimPrefix(r.RequestURI, "/v1/swagger/")
	log.Info("Swagger file requested", log.String("file", fl))
	if fl == "index.json" {
		w.Header().Add("Content-Type", "application/json")
		var d = data
		d.Host = r.Host
		if err := json.NewEncoder(w).Encode(d); err != nil {
			log.Error("Failed to serve swagger", log.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if fl == "" {
		fl = "index.html"
	}

	d, err := Asset(fl)
	if err != nil {
		log.Error("Not found", log.Err(err))
		w.WriteHeader(http.StatusNotFound)
	}
	switch strings.ToLower(filepath.Ext(fl)) {
	case ".json":
		w.Header().Add("Content-Type", "application/json")
	case ".css":
		w.Header().Add("Content-Type", "text/css")
	case ".js":
		w.Header().Add("Content-Type", "text/javascript")
	case ".html":
		w.Header().Add("Content-Type", "text/html")
	}
	_, _ = w.Write(d)
}

// RegisterSwagger register a swagger end point
func RegisterSwagger(paths map[string]interface{}, definitions map[string]interface{}) {
	swaggerLock.Lock()
	defer swaggerLock.Unlock()

	for i := range paths {
		_, ok := data.Paths[i]
		// TODO: Currently one path multiple method is not possible, fix it
		assert.False(ok, "Path is already registered", i)
		data.Paths[i] = appendSecurity(paths[i], strings.Contains(i, "{"))
	}

	for i := range definitions {
		_, ok := data.Definitions[i]
		assert.False(ok, "Definition is already registered", i)
		data.Definitions[i] = definitions[i]
	}

	if data.Definitions["ErrorResponse"] == nil {
		data.Definitions["ErrorResponse"] = errDef
		data.Definitions["ErrorFldResponse"] = errFldDef
	}

}

func appendSecurity(d interface{}, has404 bool) map[string]interface{} {
	v := d.(map[string]interface{})
	for i := range v {
		meth, ok := v[i].(map[string]interface{})
		if !ok {
			continue
		}
		if meth["security"] != nil {
			if p, ok := meth["parameters"].([]interface{}); !ok {
				meth["parameters"] = createParameter(nil)
			} else {
				meth["parameters"] = createParameter(p)
			}
		}

		meth["responses"] = create40XResponses(meth["responses"].(map[string]interface{}), meth["security"] != nil, has404)
	}
	return v

}

func createParameter(old []interface{}) []interface{} {
	return append(old, map[string]interface{}{
		"description": "the security token, get it from login route",
		"in":          "header",
		"name":        "authorization",
		"required":    true,
		"type":        "string",
	})
}

func create40XResponses(in map[string]interface{}, forbidden, notFound bool) map[string]interface{} {
	if forbidden {
		in["401"] = err401
		in["403"] = err403
	}
	if notFound {
		in["404"] = err404
	}
	in["400"] = err400
	return in
}

func init() {
	// TODO : better approach
	x1 := `{
			"type": "object",
			"properties": {
				"message": {
					"type": "string"
				},
				"status": {
					"type": "integer",
					"format": "int32"
				},
        		"fields": {
          			"type": "object",
          			"additionalProperties": {
            			"type": "string"
          			}
        		}
			}
		}`
	x2 := `{
			"type": "object",
			"properties": {
				"message": {
					"type": "string"
				},
				"status": {
					"type": "integer",
					"format": "int32"
				}
			}
		}`
	errDef = make(map[string]interface{})
	errFldDef = make(map[string]interface{})
	assert.Nil(json.Unmarshal([]byte(x1), &errFldDef))
	assert.Nil(json.Unmarshal([]byte(x2), &errDef))

}
