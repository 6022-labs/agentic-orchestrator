# HTTP Client Unit Test Templates

Use these templates to keep HTTP client unit tests structurally consistent. They exercise a thin client wrapper against a real `httptest.Server` so serialization, method, path, headers, status handling, and response decoding are all verified at the transport boundary.

Placeholders to replace:
- `<pkg>` — package of the client under test
- `<Client>` — exported client type (e.g. `FooClient`)
- `<NewClient>` — constructor (e.g. `NewFooClient`)
- `<NewClientParams>` — constructor params struct
- `<Settings>` — settings struct exposing `BaseUrl`
- `<Behavior>` — PascalCase description of the method under test (e.g. `CallingCompletions`)
- `<Method>` — method on the client being exercised
- `<Request>` / `<Response>` — request/response model types

## 1) Request + decoded response template

Use when the client sends a request and parses a JSON response body.

```go
package clients_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"<module>/src/<domain>/<pkg>/clients"
	"<module>/src/<domain>/<pkg>/models"
	"<module>/src/<domain>/<pkg>/settings"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type When<Behavior>TestingSuite struct {
	sut    *clients.<Client>
	server *httptest.Server

	handler http.HandlerFunc
}

func When<Behavior>BeforeEach() *When<Behavior>TestingSuite {
	suite := &When<Behavior>TestingSuite{}

	suite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if suite.handler == nil {
			http.NotFound(w, r)
			return
		}
		suite.handler(w, r)
	}))

	suite.sut = clients.<NewClient>(clients.<NewClientParams>{
		Logger:     zap.NewNop(),
		HttpClient: http.DefaultClient,
		<Settings>: &settings.<Settings>{BaseUrl: suite.server.URL},
	})

	return suite
}

func When<Behavior>AfterEach(suite *When<Behavior>TestingSuite) {
	suite.server.Close()
}

func TestWhen<Behavior>(t *testing.T) {
	t.Parallel()

	t.Run("Given the upstream returns 200", func(t *testing.T) {
		t.Parallel()

		var capturedMethod string
		var capturedPath string
		var capturedBody []byte

		initSuite := func(suite *When<Behavior>TestingSuite) {
			suite.handler = func(w http.ResponseWriter, r *http.Request) {
				capturedMethod = r.Method
				capturedPath = r.URL.Path
				capturedBody, _ = io.ReadAll(r.Body)
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(&models.<Response>{ /* fields */ })
			}
		}

		t.Run("Should <send the expected request and decode the response>", func(t *testing.T) {
			t.Parallel()

			suite := When<Behavior>BeforeEach()
			defer When<Behavior>AfterEach(suite)
			initSuite(suite)

			result, err := suite.sut.<Method>(/* args */, &models.<Request>{ /* fields */ })

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, http.MethodPost, capturedMethod)
			assert.Equal(t, "/expected/path", capturedPath)
			assert.Contains(t, string(capturedBody), `"field":"value"`)
		})
	})

	t.Run("Given the upstream returns non-200", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *When<Behavior>TestingSuite) {
			suite.handler = func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := When<Behavior>BeforeEach()
			defer When<Behavior>AfterEach(suite)
			initSuite(suite)

			_, err := suite.sut.<Method>(/* args */, &models.<Request>{})

			assert.Error(t, err)
		})
	})
}
```

## 2) Fire-and-forget request template

Use when the client sends a request but returns only an error (no decoded response) — e.g. callback/webhook posts. The client typically receives a full URL rather than a `BaseUrl`.

```go
package clients_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"<module>/src/<domain>/<pkg>/clients"
	"<module>/src/<domain>/<pkg>/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type When<Behavior>TestingSuite struct {
	sut    *clients.<Client>
	server *httptest.Server

	handler http.HandlerFunc
}

func When<Behavior>BeforeEach() *When<Behavior>TestingSuite {
	suite := &When<Behavior>TestingSuite{}

	suite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if suite.handler == nil {
			http.NotFound(w, r)
			return
		}
		suite.handler(w, r)
	}))

	suite.sut = clients.<NewClient>(clients.<NewClientParams>{
		Logger:     zap.NewNop(),
		HttpClient: http.DefaultClient,
	})

	return suite
}

func When<Behavior>AfterEach(suite *When<Behavior>TestingSuite) {
	suite.server.Close()
}

func TestWhen<Behavior>(t *testing.T) {
	t.Parallel()

	t.Run("Given the upstream returns 2xx", func(t *testing.T) {
		t.Parallel()

		var capturedMethod string
		var capturedPath string
		var capturedBody []byte

		initSuite := func(suite *When<Behavior>TestingSuite) {
			suite.handler = func(w http.ResponseWriter, r *http.Request) {
				capturedMethod = r.Method
				capturedPath = r.URL.Path
				capturedBody, _ = io.ReadAll(r.Body)
				w.WriteHeader(http.StatusNoContent)
			}
		}

		t.Run("Should POST the payload to the callback URL", func(t *testing.T) {
			t.Parallel()

			suite := When<Behavior>BeforeEach()
			defer When<Behavior>AfterEach(suite)
			initSuite(suite)

			body := &models.<Request>{ /* fields */ }

			err := suite.sut.<Method>(suite.server.URL+"/expected/path", body)

			assert.NoError(t, err)
			assert.Equal(t, http.MethodPost, capturedMethod)
			assert.Equal(t, "/expected/path", capturedPath)

			var decoded models.<Request>
			assert.NoError(t, json.Unmarshal(capturedBody, &decoded))
			assert.Equal(t, body, &decoded)
		})
	})

	t.Run("Given the upstream returns non-2xx", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *When<Behavior>TestingSuite) {
			suite.handler = func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := When<Behavior>BeforeEach()
			defer When<Behavior>AfterEach(suite)
			initSuite(suite)

			err := suite.sut.<Method>(suite.server.URL, &models.<Request>{})

			assert.Error(t, err)
		})
	})
}
```

## Capture patterns

- `capturedMethod` — assert verb matches expectation (`http.MethodPost`, `http.MethodGet`, etc.).
- `capturedPath` — assert URL path, including any path-param interpolation (e.g. `/agent/0xabc/7/compute`).
- `capturedQuery` (`r.URL.RawQuery`) — assert query string when the client builds one.
- `capturedHeader` (`r.Header.Get("Authorization")`) — assert auth/custom headers when the client sets them.
- `capturedBody` — `io.ReadAll(r.Body)`; decode with `json.Unmarshal` for structural asserts, or `assert.Contains(string(body), ...)` for a quick substring check.

## Scenario checklist

Cover these at minimum; add more when the client does more.

- `Given the upstream returns 2xx` → `Should send expected method/path/body` (+ decode response if applicable).
- `Given the upstream returns non-2xx` → `Should return an error`.
- `Given the upstream returns a malformed body` → `Should return an error` (only when the client decodes a response).
- `Given the upstream is unreachable` → `Should return an error` (optional; close the server before calling the SUT to simulate).
