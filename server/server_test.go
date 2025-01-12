package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ohhfishal/alice-rest/lib/event"
	"github.com/stretchr/testify/assert"
)

var port int64 = 8000
var timeout = 250 * time.Millisecond

func runServer(t *testing.T, wg *sync.WaitGroup, runner serverRunner) {
	defer wg.Done()
	runner.Run(t)
	runner.Cleanup()
}

func TestInit(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	runner := defaultRunner(nil)
	go runServer(t, &wg, runner)

	t.Run("Readyz", testUp(runner.Url()))
	t.Run("BadPath", testGet(runner.Url()+"/bad", http.StatusNotFound))
}

func TestPostEvent(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	runner := defaultRunner(nil)
	go runServer(t, &wg, runner)

	url := runner.Url()
	t.Run("Readyz", testUp(url))
	t.Run("UserAuth", testUserAuth("POST", url+"/api/v1/event/test"))

	t.Run("UserExists", func(t *testing.T) {
		url = url + "/api/v1/event/valid"
		// TODO: Register the user
		status, err := testPost(t, url, `{"description":"foo"}`)
		assert.Nil(t, err)
		assert.Equal(t, status, http.StatusCreated)
	})
}

func expectNil(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

func expectStatus(t *testing.T, status, expected int) {
	if status != expected {
		t.Fatalf("status: expected: %d got: %d", expected, status)
	}
}

func TestGetEvent(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	runner := defaultRunner(nil)
	go runServer(t, &wg, runner)

	urlBase := runner.Url()
	t.Run("Readyz", testUp(runner.Url()))
	t.Run("UserAuth", testUserAuth("GET", urlBase+"/api/v1/event/test/0"))

	t.Run("UserExists", func(t *testing.T) {
		// TODO: Register the user
		t.Run("EventMissing", testGet(urlBase+"/api/v1/event/valid/bad", http.StatusNotFound))
		t.Run("EventExists", func(t *testing.T) {
			var expected event.Event
			// TODO: Post the event
			testGetJSON(urlBase+"/api/v1/event/valid/0", 200, &expected)(t)
		})
	})
}

func testPost(t *testing.T, url, object string) (int, error) {
	reader := strings.NewReader(object)
	res, err := http.Post(url, "application/json", reader)
	if err != nil {
		return 0, fmt.Errorf("POST Request: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("Reading response: %w", err)
	}

	t.Log("Response: ", string(body))
	return res.StatusCode, nil
}

func testGetJSON[T comparable](url string, status int, expected T) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		bytes := testGetHelper(t, url, status)

		var result T
		if err := json.Unmarshal(bytes, &result); err != nil {
			t.Fatalf("unmarshaling result: %s", err.Error())
		}

		if result != expected {
			t.Fatalf("expected: %v: got: %v", expected, result)
		}
	}
}

func testGet(url string, expected int) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		_ = testGetHelper(t, url, expected)
	}
}

func testGetHelper(t *testing.T, url string, expected int) []byte {
	t.Helper()
	res, err := http.Get(url)
	if err != nil {
		t.Fatalf("making request: %s", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("reading body response: %s", err)
	}

	if status := res.StatusCode; status != expected {
		t.Fatalf("expected %d: got: %d: %s", expected, status, body)
	}
	return body
}

func testUserAuth(_, url string) func(*testing.T) {
	return testGet(url, http.StatusForbidden)
}

func testUp(urlBase string) func(*testing.T) {
	return func(t *testing.T) {
		// hack to make sure the server is up
		time.Sleep(timeout / 100)
		testGet(urlBase+"/readyz", http.StatusOK)(t)
	}
}

func runContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func addEnv(next func(string) string, key, value string) func(string) string {
	return func(k string) string {
		switch k {
		case key:
			return value
		}
		return next(k)
	}
}

type serverRunner struct {
	ctx    context.Context
	args   []string
	getenv func(string) string
	stdin  io.Reader
	stdout strings.Builder
	stderr strings.Builder
}

func (r serverRunner) Url() string {
	return fmt.Sprintf("http://%s:%s", r.getenv("HOST"), r.getenv("PORT"))
}

func (r serverRunner) Run(t *testing.T) {
	if err := Run(r.ctx, r.args, r.getenv, r.stdin, &r.stdout, &r.stderr); err != nil && !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("failed to run: %s", err)
	}
	// t.Log("stdout:", r.stdout.String())
}

func (r serverRunner) Cleanup() {
	if path := r.getenv("DATABASE_DIRECTORY"); path != "" {
		_ = os.RemoveAll(path)
	}

}

func defaultRunner(override context.Context) serverRunner {
	var ctx context.Context
	if override != nil {
		ctx = override
	} else {
		ctx, _ = runContext()
	}

	newDir, err := os.MkdirTemp("", "testdata-")
	if err != nil {
		panic(err)
	}

	newPort := atomic.AddInt64(&port, 1)
	env := func(env string) string {
		switch env {
		case "HOST":
			return "localhost"
		case "PORT":
			return fmt.Sprintf("%d", newPort)
		case "LOG_LEVEL":
			return "DEBUG"
		case "DATABASE_DIRECTORY":
			return newDir + "/"
		default:
			return ""
		}
	}

	runner := serverRunner{
		ctx:    ctx,
		getenv: env,
		stdin:  strings.NewReader(""),
	}
	return runner
}
