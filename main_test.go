package main

import (
	"bytes"
	"testing"
	"net/http"
)

type mockResponse struct {
	Buffer bytes.Buffer
}

func (m *mockResponse) Write(d []byte) (int, error) {
	return m.Buffer.Write(d)
}

func (m *mockResponse) WriteHeader(d int) {
	panic("Not Implemented")
}

func (m *mockResponse) Header() http.Header {
	panic("Not Implemented")
}

func Test_getReponseWriter(t *testing.T) {
	// Signature: (cfg *config) func (http.ResponseWriter, *http.Request)
	message := "test message"
	cfg := &config{Message: message,}
	expected_output := "<html>\n<head><title>" + message + "</title></head>\n<body>\n<p>" + message + "</p>\n</body>\n</html>"
	handler := getResponseWriter(cfg)

	req, _ := http.NewRequest("GET", "/", nil)
	rsp := &mockResponse{}

	handler(rsp, req)
	if expected_output != rsp.Buffer.String() {
		t.Error("handler output not as expected")
		return
	}
}