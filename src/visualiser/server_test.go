package visualiser

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dim-pan/blueprint/src/parser"
)

func model(reqs []parser.Requirement, comps []parser.Component) parser.SystemModel {
	return parser.SystemModel{
		Requirements: reqs,
		Components:   comps,
		Interfaces:   []parser.InterfaceDefinition{},
		TestSpecs:    []parser.TestSpec{},
	}
}

// TC-VIS-01-01
func TestTC_VIS_01_01_APIReturnsModel(t *testing.T) {
	m := model(
		[]parser.Requirement{{ID: "R1", Title: "T"}},
		[]parser.Component{{ID: "C1", Name: "C"}},
	)
	srv := httptest.NewServer(Handler(m))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/model")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("status %d", resp.StatusCode)
	}
	var decoded map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"components", "requirements", "interfaces", "test_specs"} {
		if _, ok := decoded[k]; !ok {
			t.Fatalf("missing key %q in response: %+v", k, decoded)
		}
	}
}

// TC-VIS-01-02
func TestTC_VIS_01_02_CountsMatch(t *testing.T) {
	m := model(
		[]parser.Requirement{{ID: "R1"}, {ID: "R2"}, {ID: "R3"}},
		[]parser.Component{{ID: "C1"}, {ID: "C2"}},
	)
	srv := httptest.NewServer(Handler(m))
	defer srv.Close()

	resp, _ := http.Get(srv.URL + "/api/model")
	var decoded struct {
		Components   []any `json:"components"`
		Requirements []any `json:"requirements"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&decoded)
	if len(decoded.Components) != 2 {
		t.Fatalf("components count = %d", len(decoded.Components))
	}
	if len(decoded.Requirements) != 3 {
		t.Fatalf("requirements count = %d", len(decoded.Requirements))
	}
}

// TC-VIS-02-01
func TestTC_VIS_02_01_RootServesCanvas(t *testing.T) {
	srv := httptest.NewServer(Handler(parser.SystemModel{}))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(strings.ToLower(string(body)), "<html") {
		t.Fatalf("expected an HTML response, got: %s", string(body))
	}
}

// TC-VIS-08-01
func TestTC_VIS_08_01_ServesOnConfiguredPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()

	m := model(nil, []parser.Component{{ID: "COMP-X", Name: "X"}})
	server := &http.Server{Addr: addr, Handler: Handler(m)}
	go func() { _ = server.ListenAndServe() }()
	t.Cleanup(func() { _ = server.Close() })

	// Wait until the server accepts a connection, up to 2s.
	deadline := time.Now().Add(2 * time.Second)
	var resp *http.Response
	for time.Now().Before(deadline) {
		if r, err := http.Get("http://" + addr + "/api/model"); err == nil {
			resp = r
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if resp == nil {
		t.Fatal("server never accepted a connection")
	}
	if resp.StatusCode != 200 {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

// TC-VIS-08-02
func TestTC_VIS_08_02_DefaultPort(t *testing.T) {
	if DefaultPort != 8080 {
		t.Fatalf("DefaultPort = %d, want 8080", DefaultPort)
	}
}
