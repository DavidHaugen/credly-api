package httpservice

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	testClient := NewTestClient(t)
	handler := Handler{}

	// build the request using the test context
	req := testClient.Get(nil)

	// pass the context into the handler.
	handler.Ping(req.Context)

	// assert success
	req.AssertStatus(http.StatusOK)
}
