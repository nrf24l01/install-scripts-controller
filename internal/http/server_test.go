package httpsrv

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"install-scripts-controller/internal/config"
	"install-scripts-controller/internal/database"
)

func newTestConfig(t *testing.T, ttl string) *config.Config {
	t.Helper()
	dir := t.TempDir()
	content := fmt.Sprintf(`
site:
  password: "testpass"
  install_key_ttl: "%s"
server:
  addr: ":0"
database:
  path: %q
`, ttl, filepath.Join(dir, "test.db"))

	path := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CONFIG_PATH", path)

	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func newTestServer(t *testing.T, ttl string) *httptest.Server {
	t.Helper()
	cfg := newTestConfig(t, ttl)

	db, err := database.Open(cfg.Database.Path)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}

	srv := New(cfg, db, "testdata/web")
	srv.e.Logger.SetOutput(io.Discard)
	t.Cleanup(func() { _ = db.Close() })

	ts := httptest.NewServer(srv.Handler())
	t.Cleanup(ts.Close)
	return ts
}

func request(t *testing.T, method, url, token string, body io.Reader) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func login(t *testing.T, base, password string) string {
	t.Helper()
	res := request(t, http.MethodPost, base+"/api/login", "",
		strings.NewReader(fmt.Sprintf(`{"password":%q}`, password)))
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("login: status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	return body.Token
}

func createScript(t *testing.T, base, token string) Script {
	t.Helper()
	payload := `{"name":"Test script","description":"desc","script":"#!/bin/bash\necho hi\n"}`
	res := request(t, http.MethodPost, base+"/api/scripts", token, strings.NewReader(payload))
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create: status = %d, want 201", res.StatusCode)
	}
	var sc Script
	if err := json.NewDecoder(res.Body).Decode(&sc); err != nil {
		t.Fatal(err)
	}
	return sc
}

func TestAuthAndScriptsFlow(t *testing.T) {
	ts := newTestServer(t, "1h")
	base := ts.URL

	// Wrong password is rejected.
	res := request(t, http.MethodPost, base+"/api/login", "", strings.NewReader(`{"password":"nope"}`))
	res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("wrong password: status = %d, want 401", res.StatusCode)
	}

	// Scripts require auth.
	res = request(t, http.MethodGet, base+"/api/scripts", "", nil)
	res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("list without auth: status = %d, want 401", res.StatusCode)
	}

	token := login(t, base, "testpass")

	// Empty list initially.
	res = request(t, http.MethodGet, base+"/api/scripts", token, nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("list: status = %d, want 200", res.StatusCode)
	}
	raw, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if strings.TrimSpace(string(raw)) != "[]" {
		t.Fatalf("initial list = %s, want []", raw)
	}

	// Create.
	sc := createScript(t, base, token)
	if sc.ID != 1 {
		t.Errorf("created id = %d, want 1", sc.ID)
	}
	if !strings.Contains(sc.InstallURL, "/install?id=1&key=") {
		t.Errorf("install_url = %q, want /install?id=1&key=...", sc.InstallURL)
	}

	// List shows the script.
	res = request(t, http.MethodGet, base+"/api/scripts", token, nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("list: status = %d, want 200", res.StatusCode)
	}
	var list []Script
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if len(list) != 1 || list[0].Name != "Test script" {
		t.Fatalf("list = %+v, want 1 script named Test script", list)
	}
	if list[0].Script != "" {
		t.Error("list response must not include the script body")
	}

	// Get full script.
	res = request(t, http.MethodGet, fmt.Sprintf("%s/api/scripts/%d", base, sc.ID), token, nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("get: status = %d, want 200", res.StatusCode)
	}
	var full Script
	if err := json.NewDecoder(res.Body).Decode(&full); err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if !strings.Contains(full.Script, "echo hi") {
		t.Errorf("script body = %q, want to contain echo hi", full.Script)
	}

	// Install endpoint: valid key.
	res = request(t, http.MethodGet, list[0].InstallURL, "", nil)
	bodyBytes, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != http.StatusOK || !strings.Contains(string(bodyBytes), "echo hi") {
		t.Errorf("install: status = %d body = %q", res.StatusCode, bodyBytes)
	}

	// Install endpoint: wrong key.
	bad := strings.ReplaceAll(list[0].InstallURL, "&key=", "&key=wrong-")
	res = request(t, http.MethodGet, bad, "", nil)
	res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("install wrong key: status = %d, want 401", res.StatusCode)
	}

	// Install endpoint: valid key but missing id.
	key := strings.SplitN(list[0].InstallURL, "key=", 2)[1]
	res = request(t, http.MethodGet, base+"/install?key="+key, "", nil)
	res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("install missing id: status = %d, want 400", res.StatusCode)
	}

	// Delete.
	res = request(t, http.MethodDelete, fmt.Sprintf("%s/api/scripts/%d", base, sc.ID), token, nil)
	res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("delete: status = %d, want 204", res.StatusCode)
	}

	// Empty again.
	res = request(t, http.MethodGet, base+"/api/scripts", token, nil)
	res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("list after delete: status = %d, want 200", res.StatusCode)
	}
}

func TestInstallKeyRotation(t *testing.T) {
	ts := newTestServer(t, "100ms")
	base := ts.URL
	token := login(t, base, "testpass")
	sc := createScript(t, base, token)

	// Works now.
	res := request(t, http.MethodGet, sc.InstallURL, "", nil)
	res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("install: status = %d, want 200", res.StatusCode)
	}

	// Wait past the TTL.
	time.Sleep(300 * time.Millisecond)

	// The old URL (with the rotated-out key) is rejected.
	res = request(t, http.MethodGet, sc.InstallURL, "", nil)
	res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("install with expired key: status = %d, want 401", res.StatusCode)
	}

	// A fresh list returns a new key that works.
	res = request(t, http.MethodGet, base+"/api/scripts", token, nil)
	var list []Script
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if len(list) != 1 {
		t.Fatalf("list = %d scripts, want 1", len(list))
	}
	if list[0].InstallURL == sc.InstallURL {
		t.Error("install key did not rotate")
	}
	res = request(t, http.MethodGet, list[0].InstallURL, "", nil)
	res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("install with fresh key: status = %d, want 200", res.StatusCode)
	}
}

func TestDeleteScriptNotFoundIsOk(t *testing.T) {
	ts := newTestServer(t, "1h")
	base := ts.URL
	token := login(t, base, "testpass")

	res := request(t, http.MethodDelete, base+"/api/scripts/999", token, nil)
	res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("delete missing: status = %d, want 204", res.StatusCode)
	}
}

func TestSPAServesIndex(t *testing.T) {
	ts := newTestServer(t, "1h")
	for _, path := range []string{"/", "/scripts", "/unknown-route"} {
		res := request(t, http.MethodGet, ts.URL+path, "", nil)
		body, _ := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode != http.StatusOK || !strings.Contains(string(body), "test") {
			t.Errorf("GET %s: status = %d, want 200 serving index.html", path, res.StatusCode)
		}
	}
}
