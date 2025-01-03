package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var port int64 = 8000
var timeout = 250 * time.Millisecond

func runContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func newEnv() func(string) string {
	newPort := atomic.AddInt64(&port, 1)
	return func(env string) string {
		switch env {
		case "HOST":
			return "localhost"
		case "PORT":
			return fmt.Sprintf("%d", newPort)
		case "LOG_LEVEL":
			return "DEBUG"
		default:
			return ""
		}
	}
}

func url(getenv func(string) string) string {
	return fmt.Sprintf("http://%s:%s", getenv("HOST"), getenv("PORT"))
}

func assertUp(t *testing.T, urlBase string) {
	res, err := http.Get(fmt.Sprintf("%s/readyz", urlBase))
	if err != nil {
		t.Fatalf("making readyz request: %s", err)
	}

	if status := res.StatusCode; status != http.StatusOK {
		t.Fatalf("server not ready: %d", status)
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

func (r serverRunner) Run(t *testing.T) {
	if err := Run(r.ctx, r.args, r.getenv, r.stdin, &r.stdout, &r.stderr); err != nil && !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("failed to run: %s", err)
	}
	t.Log("stdout:", r.stdout.String())
}

func defaultRunner(override context.Context) serverRunner {
	var ctx context.Context
	if override != nil {
		ctx = override
	} else {
		ctx, _ = runContext()
	}

	runner := serverRunner{
		ctx:    ctx,
		getenv: newEnv(),
		stdin:  strings.NewReader(""),
	}
	return runner
}

func TestInit(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	runner := defaultRunner(nil)
	go func() {
		defer wg.Done()
		runner.Run(t)
	}()

	// hack to make sure the server is up
	time.Sleep(timeout / 100)
	assertUp(t, url(runner.getenv))

}
