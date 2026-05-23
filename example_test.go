package xls

import (
	"fmt"
	"log"
)

func ExampleOpen() {
	if xlFile, err := Open("Table.xls", "utf-8"); err == nil {
		fmt.Println(xlFile.Author)
	}
}

// Demonstrates handling errors when opening a file.
func ExampleOpen_errorHandling() {
	_, err := Open("nonexistent.xls", "utf-8")
	if err != nil {
		log.Print("failed to open file: ", err)
		return
	}
}

func ExampleWorkBook_NumSheets() {
	if xlFile, err := Open("Table.xls", "utf-8"); err == nil {
		for i := 0; i < xlFile.NumSheets(); i++ {
			sheet := xlFile.GetSheet(i)
			fmt.Println(sheet.Name)
		}
	}
}

// Output: read the content of first two cols in each row
func ExampleWorkBook_GetSheet() {
	if xlFile, err := Open("Table.xls", "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			fmt.Print("Total Lines ", sheet1.MaxRow, sheet1.Name)
			col1 := sheet1.Row(0).Col(0)
			col2 := sheet1.Row(0).Col(0)
			for i := 0; i <= (int(sheet1.MaxRow)); i++ {
				row1 := sheet1.Row(i)
				col1 = row1.Col(0)
				col2 = row1.Col(1)
				fmt.Print("\n", col1, ",", col2)
			}
		}
	}
}

// Demonstrates Col() which searches spanned columns for merged cells,
// and ColExact() which returns a value only if the cell is explicitly present.
func ExampleRow_Col() {
	xlFile, err := Open("Table.xls", "utf-8")
	if err != nil {
		log.Fatal(err)
	}
	sheet := xlFile.GetSheet(0)
	if sheet == nil {
		return
	}
	row := sheet.Row(0)
	if row == nil {
		return
	}
	// Col returns the formatted value, searching across spanned columns for merged cells.
	fmt.Println(row.Col(0))
	// ColExact returns a value only if the cell is explicitly present at this index.
	fmt.Println(row.ColExact(0))
}

// Demonstrates reading all cells from all sheets using the ReadAllCells helper.
func ExampleWorkBook_ReadAllCells() {
	xlFile, err := Open("Table.xls", "utf-8")
	if err != nil {
		log.Fatal(err)
	}
	cells := xlFile.ReadAllCells(100)
	for i, row := range cells {
		fmt.Printf("Row %d: %v\n", i, row)
	}
}
