package xls

import (
	"testing"
	"unicode"
)

// Regression test for SST CONTINUE handling on files produced by tools that
// emit long strings crossing record boundaries (e.g. JasperReports Library).
// Pre-fix symptom: from row ~285 onwards, string cells came back either
// garbled (ASCII bytes mis-decoded as UTF-16, producing CJK glyphs) or empty,
// because get_string treated CONTINUE-resumption bytes as fresh-string flags
// and consumed phantom richtext/phonetic length fields.
func TestSSTContinueDoesNotDesync(t *testing.T) {
	wb, err := Open("testdata/sst-error.xls", "utf-8")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sheet := wb.GetSheet(0)
	if sheet == nil {
		t.Fatal("no sheet 0")
	}

	nonEmpty := 0
	garbled := 0
	for i := 0; i <= int(sheet.MaxRow); i++ {
		r := sheet.Row(i)
		if r == nil {
			continue
		}
		c0 := r.Col(0)
		c3 := r.Col(3)
		if c0 != "" && c3 != "" {
			nonEmpty++
		}
		// Heuristic: an ASCII-source description should not contain CJK
		// codepoints. Pre-fix we saw e.g. 吠慲獮敦ੲ䥓䝎偁 instead of
		// "Funds Transfer-IB\nFITNESS...".
		for _, ch := range c3 {
			if unicode.Is(unicode.Han, ch) {
				garbled++
				break
			}
		}
	}

	t.Logf("non-empty rows: %d, rows with CJK in col 3: %d", nonEmpty, garbled)

	if garbled > 0 {
		t.Errorf("found %d rows where col 3 contains CJK characters — SST desync re-introduced", garbled)
	}
	// Pre-fix this was ~283. The file has ~437 populated data rows.
	if nonEmpty < 400 {
		t.Errorf("only %d rows have both col 0 and col 3 populated; expected >= 400", nonEmpty)
	}

	// Spot-check: row 285's cells should be non-garbled (no CJK from byte-
	// pair misinterpretation). Pre-fix, this row was the first one with
	// gibberish; post-fix all cells should be empty-or-clean.
	if r := sheet.Row(285); r != nil {
		for c := uint16(0); c < 10; c++ {
			s := r.Col(int(c))
			for _, ch := range s {
				if unicode.Is(unicode.Han, ch) {
					t.Errorf("row 285 col %d contains CJK: %q", c, s)
					break
				}
			}
		}
	}
}

