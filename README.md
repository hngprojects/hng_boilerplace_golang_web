## hng_boilerplace_golang_web

### Prerequisites

1. **Go 1.17** or **lastest version** already installed on your local machine.
2. 2 Postgresql servers (one serves as normal database server, and another for running tests). You can use a disposable docker instance for tests

### Run Project from Root

1. Ensure your postgres instances are running
2. Create and populate a `app.env` file on the project root with its keys and corresponding values as listed in `app-sample.env`
3. Run from project root directory

```bash
$ go run main.go
```

### Testing

1. Automated unit and integration tests done with golang's builtin [`testing`](https://pkg.go.dev/testing) package.

To run one test file:

```bash
$ go test -v  ./tests/<file name> -timeout 99999s
```

To run all tests:

```bash
$ go test -v  ./tests/<folder name>/<file name> -timeout 99999s
```

```bash
$ go test -v  ./tests/... -timeout 99999s
```

NB: Always add timeout tag to prevent early timeout
