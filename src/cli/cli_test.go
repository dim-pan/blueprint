package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func cleanModel(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "requirements/app.req"), `req REQ-001 "X"
  priority: must-have
  The system shall X.
`)
	writeFile(t, filepath.Join(dir, "components/svc/svc.component"), `component COMP-001 "Svc"
  responsibility: Do X.
  satisfies: REQ-001
`)
	return dir
}

// TC-CLI-01-01
func TestTC_CLI_01_01_ValidModel(t *testing.T) {
	dir := cleanModel(t)
	var out, errBuf bytes.Buffer
	code := Run([]string{"validate", dir}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d, stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "valid") {
		t.Fatalf("stdout missing 'valid': %s", out.String())
	}
}

// TC-CLI-01-02
func TestTC_CLI_01_02_InvalidModel(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "components/svc/svc.component"), `component COMP-001 "Svc"
  responsibility: X.
  satisfies: REQ-MISSING
`)
	var out, errBuf bytes.Buffer
	code := Run([]string{"validate", dir}, &out, &errBuf)
	if code == 0 {
		t.Fatalf("expected non-zero exit, stdout=%s", out.String())
	}
	combined := out.String() + errBuf.String()
	if !strings.Contains(combined, "REQ-MISSING") {
		t.Fatalf("output missing REQ-MISSING: %s", combined)
	}
	if !strings.Contains(combined, "svc.component") {
		t.Fatalf("output missing source file: %s", combined)
	}
}

// TC-CLI-01-03
func TestTC_CLI_01_03_ParseError(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "requirements/broken.req"), `req "No id"
  priority: must-have
  The system shall X.
`)
	var out, errBuf bytes.Buffer
	code := Run([]string{"validate", dir}, &out, &errBuf)
	if code == 0 {
		t.Fatalf("expected non-zero exit")
	}
	combined := out.String() + errBuf.String()
	if !strings.Contains(combined, "broken.req") {
		t.Fatalf("output missing filename: %s", combined)
	}
}

// TC-CLI-06-01
func TestTC_CLI_06_01_PositionalDir(t *testing.T) {
	dir := cleanModel(t)
	var out, errBuf bytes.Buffer
	code := Run([]string{"validate", dir}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d, stderr=%s", code, errBuf.String())
	}
}

// TC-CLI-06-02
func TestTC_CLI_06_02_DefaultDir(t *testing.T) {
	dir := t.TempDir()
	// Create a sys/ subdirectory inside a scratch dir, then cd there.
	sysDir := filepath.Join(dir, defaultSysDir)
	writeFile(t, filepath.Join(sysDir, "requirements/app.req"), `req REQ-001 "X"
  priority: must-have
  The system shall X.
`)

	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	var out, errBuf bytes.Buffer
	code := Run([]string{"validate"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d, stderr=%s", code, errBuf.String())
	}
}

// TC-CLI-06-03
func TestTC_CLI_06_03_UnknownCommand(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := Run([]string{"bogus"}, &out, &errBuf)
	if code == 0 {
		t.Fatalf("expected non-zero exit")
	}
	if !strings.Contains(strings.ToLower(errBuf.String()), "unknown command") {
		t.Fatalf("stderr missing 'unknown command': %s", errBuf.String())
	}
}

func traceableModel(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "requirements/app.req"), `req REQ-001 "X"
  priority: must-have
  The system shall X.
`)
	writeFile(t, filepath.Join(dir, "components/a/a.component"), `component COMP-A "A"
  responsibility: Do X.
  satisfies: REQ-001
`)
	writeFile(t, filepath.Join(dir, "tests/app.testspec"), `test TC-001 "Does X"
  verifies: REQ-001
  given: x
  expect: y
`)
	return dir
}

// TC-CLI-04-01
func TestTC_CLI_04_01_ForwardTrace(t *testing.T) {
	dir := traceableModel(t)
	var out, errBuf bytes.Buffer
	code := Run([]string{"trace", "REQ-001", dir}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d, stderr=%s", code, errBuf.String())
	}
	combined := out.String()
	if !strings.Contains(combined, "COMP-A") || !strings.Contains(combined, "TC-001") {
		t.Fatalf("output missing links: %s", combined)
	}
}

// TC-CLI-04-02
func TestTC_CLI_04_02_BackwardTrace(t *testing.T) {
	dir := traceableModel(t)
	var out, errBuf bytes.Buffer
	code := Run([]string{"trace", "TC-001", "backward", dir}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d, stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "REQ-001") {
		t.Fatalf("output missing REQ-001: %s", out.String())
	}
}

// TC-CLI-04-03
func TestTC_CLI_04_03_UnknownIDNonzero(t *testing.T) {
	dir := traceableModel(t)
	var out, errBuf bytes.Buffer
	code := Run([]string{"trace", "REQ-NONE", dir}, &out, &errBuf)
	if code == 0 {
		t.Fatalf("expected non-zero exit")
	}
}

// TC-CLI-05-01
func TestTC_CLI_05_01_VerifyPrintsPercentage(t *testing.T) {
	dir := traceableModel(t)
	var out, errBuf bytes.Buffer
	_ = Run([]string{"verify", dir}, &out, &errBuf)
	if !strings.Contains(out.String(), "coverage") {
		t.Fatalf("output missing 'coverage': %s", out.String())
	}
	if !strings.Contains(out.String(), "%") {
		t.Fatalf("output missing percentage: %s", out.String())
	}
}

// TC-CLI-05-02
func TestTC_CLI_05_02_VerifyListsUncovered(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "requirements/app.req"), `req REQ-GAP "X"
  priority: must-have
  The system shall X.
`)
	var out, errBuf bytes.Buffer
	code := Run([]string{"verify", dir}, &out, &errBuf)
	if code == 0 {
		t.Fatalf("expected non-zero exit")
	}
	if !strings.Contains(out.String(), "REQ-GAP") {
		t.Fatalf("output missing uncovered id: %s", out.String())
	}
}

// TC-CLI-05-03
func TestTC_CLI_05_03_VerifyZeroWhenClean(t *testing.T) {
	dir := traceableModel(t)
	var out, errBuf bytes.Buffer
	code := Run([]string{"verify", dir}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("expected 0, got %d; stdout=%s stderr=%s", code, out.String(), errBuf.String())
	}
}
