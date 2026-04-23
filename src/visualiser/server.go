package visualiser

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/dim-pan/blueprint/src/parser"
)

//go:embed all:web/dist
var distFS embed.FS

const DefaultPort = 8080

type apiModel struct {
	Requirements []parser.Requirement         `json:"requirements"`
	Components   []parser.Component           `json:"components"`
	Interfaces   []parser.InterfaceDefinition `json:"interfaces"`
	TestSpecs    []parser.TestSpec            `json:"test_specs"`
}

func Handler(model parser.SystemModel) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/model", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(apiModel{
			Requirements: model.Requirements,
			Components:   model.Components,
			Interfaces:   model.Interfaces,
			TestSpecs:    model.TestSpecs,
		})
	})

	sub, err := fs.Sub(distFS, "web/dist")
	if err == nil {
		mux.Handle("/", http.FileServer(http.FS(sub)))
	} else {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "UI assets not built (run: cd src/visualiser/web && npm run build)", http.StatusInternalServerError)
		})
	}

	return mux
}

func Serve(model parser.SystemModel, addr string) error {
	server := &http.Server{Addr: addr, Handler: Handler(model)}
	return server.ListenAndServe()
}
