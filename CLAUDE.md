# CLAUDE.md - Project Context for Claude Code

## Project Overview

Go library for parsing Microsoft Excel 97-2004 `.xls` (BIFF8 binary format). **Not** `.xlsx`.

- **Module**: `github.com/meisbokai/xls`
- **Go version**: 1.21+
- **Fork of**: `github.com/extrame/xls`
- **License**: Apache 2.0 (check LICENSE)

## Build & Test Commands

```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run with coverage
go test -cover ./...

# Run a specific test
go test -run TestFunctionName ./...

# Format
gofmt -w .

# Vet
go vet ./...
```

## Architecture

### Data Flow

```text
.xls file -> OLE2 container (ole2 package)
         -> WorkBook (workbook.go)
           -> WorkSheet (worksheet.go)
             -> Row (row.go)
               -> Col types (col.go) -- cell data in various formats
```

### Key Files

| File | Responsibility |
|------|---------------|
| `xls.go` | Public API: `Open`, `OpenWithCloser`, `OpenReader` |
| `workbook.go` | WorkBook struct, BIFF record parsing, SST string handling, `get_string` |
| `worksheet.go` | WorkSheet struct, sheet-level BOF parsing, cell/row management |
| `row.go` | Row struct with `Col(i)` / `ColExact(i)` / `LastCol()` / `FirstCol()` |
| `col.go` | All cell type structs (Col, XfRk, NumberCol, FormulaCol, LabelSstCol, etc.) |
| `format.go` | Format struct for number format definitions |
| `date.go` | Excel date/time to `time.Time` conversion |
| `sst.go` | Shared String Table info struct |
| `xf.go` | Extended Format records (Xf5 for BIFF5, Xf8 for BIFF8) |
| `bof.go` | BIFF record header parsing, `bof` and `biffHeader` structs |
| `font.go` | Font info and font record parsing |
| `cell_range.go` | Cell ranges and hyperlinks |
| `comparexlsxlsx_test.go` | Test utility to compare .xls vs .xlsx output (unexported) |
| `format_parity_test.go` | Cross-format parity tests (xls vs xlsx comparison) |
| `large_dataset_test.go` | Large dataset parsing tests (date formatting, scale) |
| `reading_test.go` | Basic file opening and row/column iteration tests |
| `example_test.go` | Godoc examples for public API |

### Key Types

- **`WorkBook`** -- Top-level container; holds sheets, fonts, formats, XF table, SST
- **`WorkSheet`** -- Single sheet; rows indexed by row number
- **`Row`** -- Row of cells; `Col(i)` returns formatted string (searches spanned columns), `ColExact(i)` returns value only for explicit cells
- **`Col`** (interface `Coler`) -- Base cell type with row/col range and `String(wb) []string`
- **Cell variants**: `NumberCol`, `RkCol`, `MulrkCol`, `FormulaCol`, `FormulaStringCol`, `LabelsstCol`, `labelCol`, `BlankCol`, `MulBlankCol`

### Number Format System

Cell values go through format application in `col.go` (`NumberCol.String`, etc.) which uses `WorkBook` format strings from `xf.go`/`format.go` to produce formatted output. The `general` format is a special case.

## Conventions

- **Public API** is in `xls.go` (Open functions) and exported methods on `WorkBook`/`WorkSheet`/`Row`
- **BIFF parsing** uses binary reads from `io.ReadSeeker` with `encoding/binary`
- **Strings** can be ASCII or UTF-16LE; `get_string` in `workbook.go` handles both
- **SST (Shared String Table)** strings can span CONTINUE records -- handled in `workbook.go`
- **Test data** lives in `testdata/` with `.xls` and `.xlsx` pairs for comparison testing

## CI

- **Auto-release**: Pushes to `master` trigger `.github/workflows/auto-release.yml` which bumps minor version, creates a tag, and publishes a GitHub Release with auto-generated notes.
- Commit messages starting with `bump:` are skipped by the release workflow.

## Dependencies

- `github.com/extrame/ole2` -- OLE2 compound document parsing
- `github.com/extrame/goyymmdd` -- Date format helpers
- `github.com/tealeg/xlsx` -- Test-only dependency, used in `comparexlsxlsx_test.go`
- `golang.org/x/text` -- Charset encoding support

## Important Notes

- This is a **binary format parser** -- changes to `workbook.go` `get_string` or SST handling are high-risk for data corruption
- The `comparexlsxlsx_test.go` test utility validates output against the `tealeg/xlsx` library
- Format number handling (0-49 + custom) is in `col.go` -- changes affect how all numeric cells render
