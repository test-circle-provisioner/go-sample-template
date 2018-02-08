package helloWorld

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/segmentio/rpc"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	client, done := setup(t)
	defer done()

	var res string
	err := client.Call(context.Background(), "HelloWorld.HelloWorld", nil, &res)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", res)
}

func TestHello(t *testing.T) {
	client, done := setup(t)
	defer done()

	var res string
	err := client.Call(context.Background(), "HelloWorld.Hello", "prateek", &res)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, prateek!", res)
}

func setup(t *testing.T) (*rpc.Client, func()) {
	srv := rpc.NewServer()
	srv.Register("HelloWorld", rpc.NewService(New()))

	server := httptest.NewServer(srv)
	client := rpc.NewClient(rpc.ClientConfig{
		UserAgent: "hello-world-test",
		URL:       server.URL + "/rpc",
	})

	return client, func() {
		server.Close()
	}
}
