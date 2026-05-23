# xls

[![GoDoc](https://godoc.org/github.com/extrame/xls?status.svg)](https://godoc.org/github.com/extrame/xls)
[![Go Report Card](https://goreportcard.com/badge/github.com/meisbokai/xls)](https://goreportcard.com/report/github.com/meisbokai/xls)

Pure Golang library for reading Microsoft Excel 97-2004 `.xls` files (BIFF8 binary format).

Based on [libxls](https://sourceforge.net/projects/libxls/) and the original [extrame/xls](https://github.com/extrame/xls).

## Install

```bash
go get github.com/meisbokai/xls
```

## Usage

### Open a file

```go
package main

import (
    "fmt"
    "github.com/meisbokai/xls"
)

func main() {
    wb, err := xls.Open("spreadsheet.xls", "utf-8")
    if err != nil {
        panic(err)
    }

    sheet := wb.GetSheet(0) // first sheet
    if sheet == nil {
        return
    }

    for i := 0; i <= int(sheet.MaxRow); i++ {
        row := sheet.Row(i)
        if row == nil {
            continue
        }
        fmt.Println(row.Col(0)) // first column as formatted string
    }
}
```

### Open with closer

```go
wb, closer, err := xls.OpenWithCloser("spreadsheet.xls", "utf-8")
if err != nil {
    panic(err)
}
defer closer.Close()
```

### Open from reader

```go
file, _ := os.Open("spreadsheet.xls")
defer file.Close()
wb, err := xls.OpenReader(file, "utf-8")
```

## API

| Function | Description |
|----------|-------------|
| `Open(file, charset)` | Open .xls file by path |
| `OpenWithCloser(file, charset)` | Open .xls, returns closer for deferred close |
| `OpenReader(reader, charset)` | Open .xls from `io.ReadSeeker` |

### WorkBook fields

| Field | Type | Description |
|-------|------|-------------|
| `Author` | `string` | Document author |
| `Is5ver` | `bool` | BIFF5 format flag |
| `Codepage` | `uint16` | Code page for text decoding |
| `Xfs` | `[]st_xf_data` | Extended format table |
| `Fonts` | `[]Font` | Font table |
| `Formats` | `map[uint16]*Format` | Number format table |

### WorkBook methods

| Method | Description |
|--------|-------------|
| `GetSheet(num)` | Get sheet by index (0-based) |
| `NumSheets()` | Number of sheets |
| `ReadAllCells(max)` | Read all cells, max rows per sheet |

### WorkSheet fields

| Field | Type | Description |
|-------|------|-------------|
| `Name` | `string` | Sheet name |
| `MaxRow` | `uint16` | Highest row number with data |
| `Selected` | `bool` | Whether sheet is selected |
| `Visibility` | `TWorkSheetVisibility` | Sheet visibility state |

### Row methods

| Method | Description |
|--------|-------------|
| `Col(i)` | Formatted cell value as string |
| `ColExact(i)` | Raw cell value (no format) |
| `FirstCol()` | First column index with data |
| `LastCol()` | Last column index with data |

## Supported Data

- String cells (ASCII and Unicode)
- Numeric cells (integer and float, including RK-encoded)
- Date/time cells
- Formula cells (cached results)
- Blank cells with formatting
- Merged cell ranges
- Hyperlinks

## Charset Support

The `charset` parameter controls text decoding. Pass `"utf-8"` for most files, or a specific charset name (e.g. `"windows-1251"`) for legacy encodings.

## Contributing

See [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for full development setup, testing, and PR guidelines.

## License

Apache License 2.0
