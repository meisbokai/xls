# Contributing to xls

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git

### Getting Started

```bash
git clone https://github.com/meisbokai/xls.git
cd xls
go mod download
```

## Development Commands

<!-- AUTO-GENERATED: commands -->
| Command | Description |
|---------|-------------|
| `go test ./...` | Run all tests |
| `go test -race ./...` | Run tests with race detection |
| `go test -cover ./...` | Run tests with coverage report |
| `go test -run TestOpen ./...` | Run a specific test |
| `go vet ./...` | Static analysis |
| `gofmt -w .` | Format all Go files |
<!-- END AUTO-GENERATED -->

## Writing Tests

- Use standard `go test` with **table-driven tests** where appropriate
- Place test data files in `testdata/`
- For regression tests, include both `.xls` and `.xlsx` versions when possible
- Use the `compareXlsXlsx` helper in `comparexlsxlsx_test.go` to validate output against the reference xlsx library
- Name regression tests after the issue: `TestIssue47`, `TestSstContinue`, etc.

### Test Structure

```go
func TestFeature(t *testing.T) {
    // Arrange
    wb, err := xls.Open("testdata/example.xls", "utf-8")
    if err != nil {
        t.Fatal(err)
    }

    // Act
    sheet := wb.GetSheet(0)
    got := sheet.Row(0).Col(0)

    // Assert
    if got != "expected" {
        t.Errorf("got %q, want %q", got, "expected")
    }
}
```

## Code Style

- Run `gofmt` before committing -- no style debates
- Wrap errors with context: `fmt.Errorf("failed to X: %w", err)`
- Keep functions focused and small
- Add GoDoc comments on all exported types and functions

## Pull Request Checklist

- [ ] `go test -race ./...` passes
- [ ] `go vet ./...` passes
- [ ] `gofmt` shows no diffs
- [ ] New functionality has tests
- [ ] Regression tests include test data in `testdata/`
- [ ] No hardcoded file paths in tests

## Release Process

Releases are automatic. Merging to `master` triggers:

1. Minor version bump (e.g. v0.5.0 -> v0.6.0)
2. Git tag creation
3. GitHub Release with auto-generated notes

Commit messages starting with `bump:` are skipped by the release workflow.

## Architecture Overview

See [CODEMAP.md](./CODEMAP.md) for the full code map and architecture reference.
