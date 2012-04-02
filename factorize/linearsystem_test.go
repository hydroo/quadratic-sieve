package main


import (
	//"fmt"
	"math/rand"
	"testing"
)


/* *** Row ***************************************************************** */

/* *** Get/Set Column *** */
func TestGetSetColumn(t *testing.T) {

	rand.Seed(1234) /* use a fixed seed to make this text non-flacky but still complex */

	row := NewRow(1678)

	copyOfRow := NewRow(1678)

	for i := 0; i < 2000; i += 1 {
		copyOfRow.Set(row)

		/* set random columns to random value and watch for side effects on other columns randomly */
		value := Bit(rand.Int() % 2)
		index := rand.Int() % 1678

		row.SetColumn(index, value)

		if row.Column(index) != value {
			t.Error("Set index", index, "to value", value, "but it's still", row.Column(index))
		}

		testAllButOneColumnAreEqual(t, index, value, row, copyOfRow)
	}
}


/* *** LinearSystem ******************************************************** */


/* *** ConvertIndex ******************************************************** */
func TestConvertIndex(t *testing.T) {
	type Test struct {
		index int
		column int
		bit uint64
		exp uint32
	}

	tests := []Test{
		{  0, 0,0x0000000000000001,  0},
		{  1, 0,0x0000000000000002,  1},
		{  2, 0,0x0000000000000004,  2},
		{  3, 0,0x0000000000000008,  3},
		{  4, 0,0x0000000000000010,  4},
		{  5, 0,0x0000000000000020,  5},
		{  6, 0,0x0000000000000040,  6},
		{  7, 0,0x0000000000000080,  7},
		{  8, 0,0x0000000000000100,  8},

		{ 63, 0,0x8000000000000000, 63},

		{ 64, 1,0x0000000000000001,  0},
		{ 65, 1,0x0000000000000002,  1},
		{ 66, 1,0x0000000000000004,  2},
		{ 67, 1,0x0000000000000008,  3},
		{ 68, 1,0x0000000000000010,  4},
		{ 69, 1,0x0000000000000020,  5},
		{ 70, 1,0x0000000000000040,  6},
		{ 71, 1,0x0000000000000080,  7},
		{ 72, 1,0x0000000000000100,  8},

		{127, 1,0x8000000000000000, 63},

		{128, 2,0x0000000000000001,  0},
		{129, 2,0x0000000000000002,  1},
		{130, 2,0x0000000000000004,  2},
		{131, 2,0x0000000000000008,  3},
		{132, 2,0x0000000000000010,  4},
		{133, 2,0x0000000000000020,  5},
		{134, 2,0x0000000000000040,  6},
		{135, 2,0x0000000000000080,  7},
		{136, 2,0x0000000000000100,  8},

		{767,11,0x8000000000000000, 63},

		{768,12,0x0000000000000001,  0},
		{769,12,0x0000000000000002,  1},
		{770,12,0x0000000000000004,  2},
		{771,12,0x0000000000000008,  3},
		{772,12,0x0000000000000010,  4},
		{773,12,0x0000000000000020,  5},
		{774,12,0x0000000000000040,  6},
		{775,12,0x0000000000000080,  7},
		{776,12,0x0000000000000100,  8},
	}

	for _, test := range tests {
		if c, b, e := convertIndex(test.index); c != test.column || b != test.bit || e != test.exp {
			t.Error("index", test.index, "should be column", test.column, ", bit",
					test.bit, "and exp", test.exp, "not", c, ",", b, "and", e)
		}
	}
}


/* *** Helper *** ********************************************************** */
func testAllButOneColumnAreEqual(t *testing.T, indexDoNotCheck int, value Bit, newRow *Row, oldRow *Row) {

	for i := 0; i < newRow.columnCount; i += 1 {

		if i == indexDoNotCheck {
			continue
		}

		if newRow.Column(i) != oldRow.Column(i) {
			t.Error("Index", indexDoNotCheck, "has been set to value", value, ". But index", i, "'s value changed from",
				oldRow.Column(i), "to", newRow.Column(i))
		}

	}
}

