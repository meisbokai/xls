# Code Map — xls

> Quick reference for navigating the codebase.

## Entry Points

| Function | File | Description |
|----------|------|-------------|
| `Open(file, charset)` | `xls.go:11` | Open .xls file by path |
| `OpenWithCloser(file, charset)` | `xls.go:20` | Open .xls, return closer for deferred close |
| `OpenReader(reader, charset)` | `xls.go:30` | Open .xls from any `io.ReadSeeker` |

## Core Types

### WorkBook (`workbook.go`)

Top-level parser. Holds all sheets, fonts, formats, XF table, and SST.

```
WorkBook
  +-- Sheets    []WorkSheet    (via GetSheet/NumSheets)
  +-- fonts     []FontInfo
  +-- formats   []Format
  +-- xf        []st_xf_data   (Xf5 or Xf8)
  +-- sst       []string       (shared string table)
```

Key methods:
- `Parse(buf)` -- walks all BIFF records
- `GetSheet(num)` -- returns `*WorkSheet` by index
- `NumSheets()` -- sheet count
- `ReadAllCells(max)` -- read all cells up to max rows per sheet
- `get_string(buf, size)` -- reads SST string (ASCII or UTF-16LE, handles CONTINUE records)

### WorkSheet (`worksheet.go`)

Represents one sheet. Parses sheet-level BIFF records into rows and cells.

```
WorkSheet
  +-- Name      string
  +-- Visible   TWorkSheetVisibility
  +-- rows      map[int]*RowInfo
```

### Row (`row.go`)

A single row of cells.

- `Col(i)` -- formatted cell value as string
- `ColExact(i)` -- raw cell value (no format application)
- `FirstCol()` / `LastCol()` -- column bounds

### Cell Types (`col.go`)

All implement `Coler` interface (`Row`, `FirstCol`, `LastCol`).

| Type | BIFF Record | Description |
|------|------------|-------------|
| `Col` | Base | Row/col range, delegates to `contentHandler` |
| `RkCol` | RK | Encoded integer/float (RK encoding) |
| `MulrkCol` | MULRK | Multiple RK values in one record |
| `NumberCol` | NUMBER | IEEE 754 float64 |
| `FormulaCol` | FORMULA | Formula with cached result |
| `FormulaStringCol` | FORMULA | Formula with string result |
| `LabelsstCol` | LABELSST | Index into SST |
| `labelCol` | LABEL | Inline string (BIFF5/7) |
| `BlankCol` | BLANK | Empty formatted cell |
| `MulBlankCol` | MULBLANK | Multiple blanks in one record |
| `XfRk` | -- | XF index + RK value pair |

## Supporting Types

| Type | File | Purpose |
|------|------|---------|
| `Format` | `format.go` | Number format string (e.g. `#,##0.00`) |
| `FontInfo` / `Font` | `font.go` | Font metadata |
| `Xf5` / `Xf8` | `xf.go` | Extended format for BIFF5/BIFF8 |
| `SstInfo` | `sst.go` | Shared string table header |
| `bof` / `biffHeader` | `bof.go` | BIFF record header |
| `CellRange` | `cell_range.go` | Merged cell ranges |
| `HyperLink` | `cell_range.go` | Hyperlink cells |
| `RK` | `col.go` | RK-encoded number (int or float) |

## Data Flow: Reading a Cell

```
1. Open() -> ole2.Open() -> parse OLE2 container
2. WorkBook.Parse(buf) -> walk BIFF records
   - BOF records -> identify sheet boundaries
   - SST records -> build shared string table
   - XF records -> build format table
   - FONT records -> build font table
3. WorkSheet.parse(buf) -> parse sheet records
   - ROW records -> create Row entries
   - COL/MULRK/NUMBER/FORMULA/etc -> create cell entries
4. Row.Col(i) -> lookup cell by index
   - cell.String(wb) -> apply number format from XF table
   - returns formatted string
```

## Number Format Chain

```
Cell (e.g. NumberCol)
  -> String(wb *WorkBook)
    -> wb.xf[cell.Xf] -> formatNo()
      -> wb.formats[no] -> Format string
        -> applyFormat(value, formatString)
```

Format numbers 0-49 are built-in Excel formats. Custom formats are stored in the `Format` struct.

## Date Handling (`date.go`)

Excel stores dates as float64 days since 1900-01-01 (or 1904-01-01 in 1904 date system).

- `timeFromExcelTime(value, date1904)` -> `time.Time`
- Uses Fliegel-Van Flandern algorithm for Julian date conversion

## Test Files

| File | Tests |
|------|-------|
| `xls_test.go` | Core Open/Parse functionality |
| `example_test.go` | GoDoc examples |
| `bigtable_test.go` | Large table parsing |
| `issue47_test.go` | Regression test for issue #47 |

Test data in `testdata/` includes `.xls`/`.xlsx` pairs for cross-validation via `CompareXlsXlsx`.
