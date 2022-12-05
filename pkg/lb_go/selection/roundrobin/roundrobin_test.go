package roundrobin

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ggrangia/lb_go/test"
	"github.com/google/go-cmp/cmp"
)

func TestNextServerAlive(t *testing.T) {
	teardown, backends := test.SetupBackends(t, 3)
	defer teardown(t)
	rr := New()

	for i := 0; i < len(backends); i++ {
		// Set all backends alive
		backends[i].Alive = true
		rr.AddBackend(backends[i])
	}

	tests := []struct {
		want    string
		counter int
	}{
		{want: backends[0].Url.Host, counter: 1},
		{want: backends[1].Url.Host, counter: 2},
		{want: backends[2].Url.Host, counter: 0},
		{want: backends[0].Url.Host, counter: 1},
	}

	for _, tc := range tests {
		got, err := rr.nextServer()
		if err != nil {
			t.Fatalf("Got error %v", err)
		}
		if !cmp.Equal(got.Url.Host, tc.want) {
			t.Errorf("Expected output %v got %v", tc.want, got.Url.Host)
		}
		if !cmp.Equal(rr.Counter, tc.counter) {
			t.Errorf("Expected counter %d got %d", tc.counter, rr.Counter)
		}
	}
}

func TestNextServerNotAlive(t *testing.T) {
	teardown, backends := test.SetupBackends(t, 3)
	defer teardown(t)
	rr := New()

	// Backends alive is false by default
	for i := 0; i < len(backends); i++ {
		rr.AddBackend(backends[i])
	}

	tests := []struct {
		want    string
		counter int
	}{
		{want: backends[0].Url.Host, counter: 1},
		{want: backends[1].Url.Host, counter: 2},
		{want: backends[2].Url.Host, counter: 0},
		{want: backends[0].Url.Host, counter: 1},
	}

	for _, tc := range tests {
		got, err := rr.nextServer()
		if err == nil {
			t.Fatalf("No error connecting to %v. Got: %v", tc.want, got)
		}
		if !cmp.Equal(err.Error(), ErrNoServer.Error()) {
			t.Errorf("Expected output %v got %v", ErrNoServer.Error(), err.Error())
		}
	}
}

func TestNextServerOneAlive(t *testing.T) {
	teardown, backends := test.SetupBackends(t, 3)
	defer teardown(t)
	rr := New()
	// Backends alive is false by default
	for i := 0; i < len(backends); i++ {
		rr.AddBackend(backends[i])
	}

	// Set only one backend alive
	rr.Backends[1].Alive = true

	tests := []struct {
		want    string
		counter int
	}{
		{want: backends[1].Url.Host, counter: 2},
		{want: backends[1].Url.Host, counter: 2},
		{want: backends[1].Url.Host, counter: 2},
		{want: backends[1].Url.Host, counter: 2},
	}

	for _, tc := range tests {
		got, err := rr.nextServer()
		if err != nil {
			t.Fatalf("Got error %v", err)
		}
		if !cmp.Equal(got.Url.Host, tc.want) {
			t.Errorf("Expected output %v got %v", tc.want, got.Url.Host)
		}
		if !cmp.Equal(rr.Counter, tc.counter) {
			t.Errorf("Expected counter %d got %d", tc.counter, rr.Counter)
		}
	}
}

func TestServeHttpAlive(t *testing.T) {
	teardown, backends := test.SetupBackends(t, 3)
	defer teardown(t)
	rr := New()
	// Backends alive is false by default
	for i := 0; i < len(backends); i++ {
		rr.AddBackend(backends[i])
		rr.Backends[i].Alive = true
	}

	tests := []struct {
		want   string
		status int
	}{
		{want: "0", status: 200},
		{want: "1", status: 200},
		{want: "2", status: 200},
		{want: "0", status: 200},
	}

	for _, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		rr.ServeHTTP(w, r)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		raw := string(data)
		// remove trailing newlines
		got := strings.TrimRight(raw, "\r\n")
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}
		if res.StatusCode != tc.status {
			t.Errorf("expected status code %v got %v", tc.status, res.StatusCode)
		}
		if diff := cmp.Diff(got, tc.want); diff != "" {
			t.Errorf("expected %v got %v, Diff: %v", tc.want, got, diff)
		}
	}
}

func TestServeHttpNotAlive(t *testing.T) {
	teardown, backends := test.SetupBackends(t, 3)
	defer teardown(t)
	rr := New()
	// Backends alive is false by default
	for i := 0; i < len(backends); i++ {
		rr.AddBackend(backends[i])
	}

	tests := []struct {
		status int
	}{
		{status: 500},
		{status: 500},
		{status: 500},
		{status: 500},
	}

	for _, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		rr.ServeHTTP(w, r)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		raw := string(data)
		// remove trailing newlines
		got := strings.TrimRight(raw, "\r\n")
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}
		if res.StatusCode != tc.status {
			t.Errorf("expected status code %v got %v", tc.status, res.StatusCode)
		}

		if diff := cmp.Diff(got, ErrNoServer.Error()); diff != "" {
			t.Errorf("expected %v got %v, Diff: %v", ErrNoServer.Error(), got, diff)
		}
	}
}
