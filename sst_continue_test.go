package xls

import (
	"regexp"
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

// Regression test for a second SST CONTINUE edge case: a record may end
// after a string's character data is complete but mid-way through its
// trailing rgRun (rich-text format runs) or ExtRst (phonetic) bytes. The
// loop used to break without advancing the SST index, so the next CONTINUE
// record then wrote the *next* string onto the already-complete entry via
// `sst[i] = sst[i] + str`, gluing two strings into one slot and shifting
// every subsequent SST entry by one.
//
// In testdata/sst-error.xls this manifests at row 285: column 2 (timestamp)
// previously held "02:38:35 PMFunds Transfer-IB\n..." (timestamp + the
// description that should have been in column 3), with the remainder of the
// row shifted left by one column. Post-fix, column 2 holds only the
// timestamp and column 3 holds the description as authored.
func TestSSTContinueTrailerOverflowDoesNotMerge(t *testing.T) {
	wb, err := Open("testdata/sst-error.xls", "utf-8")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	sheet := wb.GetSheet(0)
	if sheet == nil {
		t.Fatal("no sheet 0")
	}

	r := sheet.Row(285)
	if r == nil {
		t.Fatal("row 285 missing")
	}

	timestampRe := regexp.MustCompile(`^\d{2}:\d{2}:\d{2} (AM|PM)$`)

	c2 := r.Col(2)
	if !timestampRe.MatchString(c2) {
		t.Errorf("row 285 col 2 should be a bare HH:MM:SS AM/PM timestamp, got %q "+
			"— a trailing description glued on indicates the trailer-overflow "+
			"SST merge bug has regressed", c2)
	}

	c3 := r.Col(3)
	if c3 == "" {
		t.Errorf("row 285 col 3 is empty — the description shifted out, "+
			"likely re-merged into col 2 (col 2 was %q)", c2)
	}
	// "Funds Transfer-IB" is the literal first line of row 285's description
	// in the source file; if the SST is shifted by one, col 3 will instead
	// contain the next row's value ("PAYNOW-FAST\n...") or a numeric token.
	if c3 != "" && !startsWithAny(c3, "Funds Transfer-IB", "Inward", "PAYNOW-FAST", "Bulk -", "SVC Chg", "CR Retail", "Outward") {
		t.Errorf("row 285 col 3 = %q does not look like a transaction description; "+
			"likely an SST shift", c3)
	}

	// Walk rows 285..320 and assert column 2 is timestamp-shaped on every row
	// that has data. If the SST is off by one anywhere in this window, col 2
	// will hold description text or a numeric string instead.
	checked := 0
	for i := 285; i <= 320; i++ {
		rr := sheet.Row(i)
		if rr == nil {
			continue
		}
		c0 := rr.Col(0)
		if c0 == "" {
			// Skip header/blank rows that don't carry a transaction.
			continue
		}
		c2 := rr.Col(2)
		if c2 == "" {
			continue
		}
		if !timestampRe.MatchString(c2) {
			t.Errorf("row %d col 2 = %q is not a HH:MM:SS AM/PM timestamp; "+
				"SST is shifted in this window", i, c2)
		}
		checked++
	}
	if checked < 30 {
		t.Errorf("only %d rows in [285,320] had a non-empty timestamp column; "+
			"expected ~36, suggests rows are missing", checked)
	}
}

func startsWithAny(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if len(s) >= len(p) && s[:len(p)] == p {
			return true
		}
	}
	return false
}

