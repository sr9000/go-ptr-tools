# 2. Getting Started

This library is a zero-dependency utility package designed to be embedded directly into your Go projects. It contains simple building blocks and helper functions around safe value handling using pointers, references, and optional types, with no runtime configuration or special environment required.

There are **no binaries to build**, no frameworks to install, and no background processes involved. The repository provides a set of utility packages and a **Makefile** to run tests, benchmarks, and lint checks easily.

**Installation**:

To start using the utilities, simply import the relevant packages into your Go project.

```bash
go get https://github.com/sr9000/go-ptr-tools
```

**Get Sources Repo**:

```bash
git clone https://github.com/sr9000/go-ptr-tools.git
```

**Makefile Commands**:

A Makefile is included to facilitate common development tasks like linting, running tests, generating benchmarks, or clearing the test cache. This helps keep your workflow clean and consistent.

Here are the available **Makefile targets**:

| Command      | Description                                                  |
| ------------ | ------------------------------------------------------------ |
| `make lint`  | Runs `golangci-lint` to lint the code and optionally auto-fix style issues. |
| `make test`  | Runs all Go tests in the project using `go test`.            |
| `make cover` | Runs tests and generates a test coverage report.             |
| `make bench` | Runs benchmarks using the `go test -bench` command.          |
| `make clean` | Clears the Go test cache using `go clean -testcache`.        |
| `make all`   | Performs `clean`, `lint`, `test`, and `bench` in sequence.   |

