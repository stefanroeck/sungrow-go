# Agent instructions

- For all code changes, run the build and tests before considering the work complete.
- From the repository root, build all packages with `go build -o ./bin/`.
- From the repository root, run the full test suite with `go test ./...`.
- Some packages currently have no test files; `[no test files]` is expected, and all test packages must pass.
- Use Go 1.27 for both commands, as required by `go.mod`.
