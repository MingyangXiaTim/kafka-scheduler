package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MingyangXiaTim/kafka-scheduler/apiserver/rest"
	hmapcoll "github.com/MingyangXiaTim/kafka-scheduler/internal/collector/hmap"
	"github.com/MingyangXiaTim/kafka-scheduler/internal/store/hmap"
	"github.com/MingyangXiaTim/kafka-scheduler/scheduler"
)

func newServer() (srv rest.Server, closeFunc func()) {
	sch := scheduler.New(hmap.New(), hmapcoll.New(), nil)
	sch.Start(scheduler.StartOfToday())

	srv = rest.New(&sch)

	return srv, func() {
		sch.Close()
	}
}

func executeRequest(router http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Fatalf("unexpected response code %d, expected: %d", actual, expected)
	}
}
