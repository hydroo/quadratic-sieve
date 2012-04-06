package main


import (
	"fmt"
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


func TestConvertIndex(t *testing.T) {
	type Test struct {
		columnCount int
		index int
		column int
		bit uint64
		exp uint32
	}

	tests := []Test{
		{155,  0, 2,0x0000000000000001,  0},
		{155,  1, 2,0x0000000000000002,  1},
		{155,  2, 2,0x0000000000000004,  2},
		{155,  3, 2,0x0000000000000008,  3},
		{155,  4, 2,0x0000000000000010,  4},
		{155,  5, 2,0x0000000000000020,  5},
		{155,  6, 2,0x0000000000000040,  6},
		{155,  7, 2,0x0000000000000080,  7},
		{155,  8, 2,0x0000000000000100,  8},

		{155, 63, 2,0x8000000000000000, 63},

		{155, 64, 1,0x0000000000000001,  0},
		{155, 65, 1,0x0000000000000002,  1},
		{155, 66, 1,0x0000000000000004,  2},
		{155, 67, 1,0x0000000000000008,  3},
		{155, 68, 1,0x0000000000000010,  4},
		{155, 69, 1,0x0000000000000020,  5},
		{155, 70, 1,0x0000000000000040,  6},
		{155, 71, 1,0x0000000000000080,  7},
		{155, 72, 1,0x0000000000000100,  8},

		{155,127, 1,0x8000000000000000, 63},

		{155,128, 0,0x0000000000000001,  0},
		{155,129, 0,0x0000000000000002,  1},
		{155,130, 0,0x0000000000000004,  2},
		{155,131, 0,0x0000000000000008,  3},
		{155,132, 0,0x0000000000000010,  4},
		{155,133, 0,0x0000000000000020,  5},
		{155,134, 0,0x0000000000000040,  6},
		{155,135, 0,0x0000000000000080,  7},
		{155,136, 0,0x0000000000000100,  8},
	}

	for _, test := range tests {
		row := NewRow(test.columnCount)
		if c, b, e := row.convertIndex(test.index); c != test.column || b != test.bit || e != test.exp {
			t.Error("index", test.index, "should be column", test.column, ", bit",
					test.bit, "and exp", test.exp, "not", c, ",", b, "and", e)
		}
	}
}


/* *** LinearSystem ******************************************************** */
func TestGaussianElimination(t *testing.T) {

	type Test struct {
		before, expect [][]int
	}

	tests := []Test{

			{
				[][]int{
				{0,0,1,1,0,0},
				{1,1,1,1,1,0},
				{0,0,1,0,1,1},
				{1,0,0,0,1,0},
				{0,0,1,1,0,0},
				},
				[][]int{
				{1,1,1,1,1,0},
				{0,1,1,1,0,0},
				{0,0,1,0,1,1},
				{0,0,0,1,1,1},
				{0,0,0,0,0,0},
				},
			},

			{
				[][]int{
				{1,1,1,1},
				{0,1,1,1},
				{0,0,1,1},
				{0,0,0,1},
				{1,1,0,1},
				{1,1,1,0},
				{0,1,0,1},
				{0,1,1,0},
				},
				[][]int{
				{1,1,1,1},
				{0,1,1,1},
				{0,0,1,1},
				{0,0,0,1},
				{0,0,0,0},
				{0,0,0,0},
				{0,0,0,0},
				{0,0,0,0},
				},
			},
			}

	for _, test := range tests {
		before := linearSystemFromIntMatrix(test.before)
		after := linearSystemFromIntMatrix(test.before)
		expect := linearSystemFromIntMatrix(test.expect)

		after.GaussianElimination(before)

		if after.Equals(expect) == false {
			t.Error("missmatch between result and expectation:\ninput:\n",before,
					"\nresult:\n", after,"\nexpectation:\n", expect)
		}

	}
}


func linearSystemFromIntMatrix(m [][]int) *LinearSystem {

	if len(m) == 0 {
		panic("matrix is not supposed to be empty")
	}

	if len(m[0]) == 0 {
		panic("matrix is not supposed to be empty")
	}


	rows := len(m)
	columns := len(m[0])

	ret := NewLinearSystem(rows, columns)

	for j, v := range m {
		for i, k := range v {
			if k != 0 && k != 1 {
				panic(fmt.Sprint("invalid value k ", k))
			}
			ret.Row(j).SetColumn(columns-1-i, Bit(k))
		}
	}

	return ret
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

