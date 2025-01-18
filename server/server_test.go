package server

import (
	"bytes"
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

	"github.com/ohhfishal/alice-rest/database"
	"github.com/stretchr/testify/require"
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

	t.Run("Readyz", Readyz(runner.Url(), 200))
}

func TestPostEvent(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	runner := defaultRunner(nil)
	go runServer(t, &wg, runner)

	tests := []struct {
		Name   string
		User   string
		Body   database.Event
		Status int
	}{
		{
			Name:   "Valid",
			User:   "vaild",
			Body:   database.Event{Description: "foo"},
			Status: 201,
		},
		{
			Name:   "Valid/No Body",
			User:   "vaild",
			Status: 400,
		},
		{
			Name:   "No User/ValidBody",
			User:   "",
			Status: 404,
			Body:   database.Event{Description: "foo"},
		},
		{
			Name:   "No User/No Body",
			User:   "",
			Status: 404,
			Body:   database.Event{Description: "foo"},
		},
	}

	t.Run("Readyz", Readyz(runner.Url(), 200))

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			url := runner.Url() + "/api/v1/event/" + test.User
			_ = PostEvent(t, url, test.Body, test.Status)
		})
	}
}

func TestGetEvent(t *testing.T) {
	return
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	runner := defaultRunner(nil)
	go runServer(t, &wg, runner)

	urlBase := runner.Url()
	tests := []struct {
		Name       string
		IDOverride int64
		PostBody   database.Event
		Status     int
		Expected   database.Event
	}{
		{
			Name:     "EventExists",
			PostBody: database.Event{Description: "foo"},
			Status:   200,
			Expected: database.Event{Description: "foo", State: "in progress"},
		},
		{
			Name:       "EventMissing",
			IDOverride: -1,
			Status:     404,
		},
	}

	var zero database.Event
	t.Run("Readyz", Readyz(runner.Url(), 200))
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			var id int64
			if test.PostBody != zero {
				id = PostEvent(t, urlBase+"/api/v1/event/user", test.PostBody, 201)
			}

			if test.IDOverride != 0 {
				id = test.IDOverride
			}

			url := fmt.Sprintf("%s/api/v1/event/user/%d", urlBase, id)
			GetEvent(t, url, test.Status, test.Expected)
		})
	}
}

func PostEvent(t *testing.T, url string, event database.Event, status int) int64 {
	var id int64
	eventBytes, err := json.Marshal(event)
	require.Nil(t, err)
	reader := bytes.NewReader(eventBytes)

	t.Logf("POST: %s", url)
	res, err := http.Post(url, "application/json", reader)
	body, err := io.ReadAll(res.Body)

	require.Nil(t, err)
	require.Equal(t, status, res.StatusCode, string(body))
	t.Logf("POST: %d", res.StatusCode)

	if res.StatusCode >= 400 {
		return id
	}

	err = json.Unmarshal(body, &id)
	require.Nil(t, err)
	t.Log("Body: ", string(body))
	return id

}

func GetEvent(t *testing.T, url string, status int, expected database.Event) {
	var event database.Event

	res, err := http.Get(url)
	require.Nil(t, err)
	require.Equal(t, status, res.StatusCode)
	defer t.Logf("GET: %d", res.StatusCode)

	if res.StatusCode >= 300 {
		return
	}

	body, err := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &event)
	require.Nil(t, err)
	require.Equal(t, expected.Description, event.Description)
	require.Equal(t, expected.State, event.State)
}

func Readyz(url string, status int) func(*testing.T) {
	return func(t *testing.T) {
		// hack to make sure the server is up
		time.Sleep(timeout / 100)
		res, err := http.Get(url + "/readyz")
		require.Nil(t, err)
		require.Equal(t, res.StatusCode, status)
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
