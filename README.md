# go-hello-world

Provides a simple example of a production ready RPC service in Go. Instead of attempting to abstract away the details of bootstrapping a new service, this is meant to be used as a blueprint of how to piece together various libraries.

# Libraries

* [conf](github.com/segmentio/conf): Loading program configuration from multiple sources.
* [events](https://github.com/segmentio/events): Handles routing, formatting and publishing logs.
* [stats](https://github.com/segmentio/stats): Handles routing, formatting and publishing stats.
* [rpc](https://github.com/segmentio/rpc): jsonrpc v2 server (and client) that integrates with `events` and `stats`.
* [alice](https://github.com/justinas/alice): HTTP middleware chaining used for HTTP stats and logs.
* [assert](https://github.com/stretchr/testify#assert-package): Assertion helper methods for testing.

# Tools

* [govendor](https://github.com/kardianos/govendor): Dependency management tool.
* [Golang Build Image](https://github.com/segmentio/golang-private-image): Containerized build process for Go.

# Get started

* make deps
* make test