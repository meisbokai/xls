package xls

import "testing"

// newFormatTestWB builds a minimal WorkBook with a single XF at index 0 whose
// format number is fNo, mapped to the given format string.
func newFormatTestWB(fNo uint16, format string) *WorkBook {
	wb := &WorkBook{
		Xfs:     []st_xf_data{&Xf8{Format: fNo}},
		Formats: map[uint16]*Format{},
	}
	f := &Format{str: format}
	f.Head.Index = fNo
	wb.Formats[fNo] = f
	return wb
}

// serial 46090 in the 1900 date system corresponds to 2026-03-09.
const localeDateSerial = 46090

func TestNumberCol_LocaleTaggedDate_ISO(t *testing.T) {
	wb := newFormatTestWB(200, "[$-030009]yyyy mmmm")
	c := &NumberCol{Index: 0, Float: localeDateSerial}
	got := c.String(wb)
	if len(got) != 1 || got[0] != "2026-03-09" {
		t.Fatalf("NumberCol locale date: got %v, want [2026-03-09]", got)
	}
}

func TestXfRk_LocaleTaggedDate_ISO(t *testing.T) {
	// user-defined format (>= 164) triggers the formatter date path
	wb := newFormatTestWB(200, "[$-030009]yyyy mmmm")
	// RK encoding of integer serial: value << 2 with isInt bit set
	xf := &XfRk{Index: 0, Rk: RK(uint32(localeDateSerial)<<2 | 2)}
	got := xf.String(wb)
	if got != "2026-03-09" {
		t.Fatalf("XfRk locale date: got %q, want %q", got, "2026-03-09")
	}
}

func TestNumberCol_NormalDate_Unchanged(t *testing.T) {
	// A goyymmdd-compatible format (uppercase M = month) must still route
	// through goyymmdd untouched and render correctly.
	wb := newFormatTestWB(200, "yyyy-MM-dd")
	c := &NumberCol{Index: 0, Float: localeDateSerial}
	got := c.String(wb)
	if len(got) != 1 || got[0] != "2026-03-09" {
		t.Fatalf("NumberCol normal date: got %v, want [2026-03-09]", got)
	}
}

func TestNumberCol_Numeric_Unaffected(t *testing.T) {
	wb := newFormatTestWB(200, "#,##0.00")
	c := &NumberCol{Index: 0, Float: 1234.5}
	got := c.String(wb)
	if len(got) != 1 || got[0] != "1234.5" {
		t.Fatalf("NumberCol numeric: got %v, want [1234.5]", got)
	}
}

func TestNumberCol_CurrencyLocale_Unaffected(t *testing.T) {
	// currency locale formats start with "[$" + symbol, not "[$-", and also
	// contain "#", so they short-circuit to numeric rendering.
	wb := newFormatTestWB(200, "[$$-409]#,##0.00")
	c := &NumberCol{Index: 0, Float: 1234.5}
	got := c.String(wb)
	if len(got) != 1 || got[0] != "1234.5" {
		t.Fatalf("NumberCol currency locale: got %v, want [1234.5]", got)
	}
}
