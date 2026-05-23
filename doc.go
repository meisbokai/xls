// Package xls provides a pure Go reader for Microsoft Excel 97-2004 binary
// (.xls) files (BIFF8 format). It does not support the newer .xlsx format.
//
// # Opening a file
//
// Use Open, OpenWithCloser, or OpenReader to obtain a WorkBook:
//
//	wb, err := xls.Open("spreadsheet.xls", "utf-8")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// The charset parameter controls text decoding. Pass "utf-8" for most files.
//
// # Reading data
//
// Access sheets by index with WorkBook.GetSheet, then iterate rows and columns:
//
//	sheet := wb.GetSheet(0)
//	for i := 0; i <= int(sheet.MaxRow); i++ {
//	    row := sheet.Row(i)
//	    if row == nil {
//	        continue
//	    }
//	    fmt.Println(row.Col(0))
//	}
//
// # Key types
//
//   - WorkBook: top-level container holding sheets, fonts, and formats
//   - WorkSheet: a single sheet with rows indexed by number
//   - Row: a row of cells; use Col(i) to search spanned columns for merged cells, or ColExact(i) for explicit cell values only
//
// # Supported features
//
// String cells (ASCII and Unicode), numeric cells (integer, float, RK-encoded),
// date/time cells, formula cells (cached results), blank cells, merged cell
// ranges, and hyperlinks.
package xls
