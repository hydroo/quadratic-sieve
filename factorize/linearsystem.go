package main

/* homogenous system of linear equations with coefficients of GF(2) */

import (
	"fmt"
)


/* *** Bit *** ************************************************************* */
type Bit int

func (this Bit) Check() {
	if this != 0 && this != 1 {
		panic(fmt.Sprint("Invalid Value. Should be 0 or 1, but is", this))
	}
}


func (this Bit) String() string {
	if this == 0 {
		return "0"
	} /* else { */

	return "1"
}


/* *** Row *** ************************************************************* */
type Row struct {
	chunks []uint64
	columnCount int
}


func NewRow(columnCount int) *Row {

	if columnCount < 0 {
		panic("column count has to be >= 0")
	}

	var ret Row
	ret.columnCount = columnCount
	ret.chunks = make([]uint64, ((columnCount-1)/64)+1)

	for i, _ := range ret.chunks {
		/* initialize to zero */
		ret.chunks[i] = 0x0000000000000000
	}

	return &ret
}


func (this Row) Column(index int) Bit {
	this.checkIndex(index)

	column, bit, exp := this.convertIndex(index)

	ret := Bit((this.chunks[column] & bit) >> exp)
	ret.Check()
	return ret
}


func (this *Row) SetColumn(index int, value Bit) {
	value.Check()
	this.checkIndex(index)

	column, bit, _ := this.convertIndex(index)

	if value == 0 {
		this.chunks[column] &= ^bit
	} else {
		this.chunks[column] |= bit
	}
}


func (this *Row) Set(other *Row) *Row {

	this.checkSameSize(other)

	for i, k := range other.chunks {
		this.chunks[i] = k;
	}

	return this
}


func (this *Row) Swap(other *Row) *Row {

	this.checkSameSize(other)

	var tempChunk uint64
	for i, _ := range other.chunks {
		tempChunk = this.chunks[i]
		this.chunks[i] = other.chunks[i];
		other.chunks[i] = tempChunk;
	}

	return this
}


func (this *Row) Xor(a, b *Row) *Row {

	this.checkSameSize(a)
	this.checkSameSize(b)

	for i, _ := range a.chunks {
		this.chunks[i] = a.chunks[i] ^ b.chunks[i]
	}

	return this
}


func (this *Row) Neg(other *Row) *Row {

	this.checkSameSize(other)

	for i, chunk := range other.chunks {
		this.chunks[i] = ^chunk
	}

	return this
}


func (this Row) String() string {
	var ret string
	for i := this.columnCount - 1; i >= 0; i -= 1 {
		ret += fmt.Sprint(this.Column(i)," ")
	}
	return ret
}


func (this Row) Equals(other *Row) bool {

	this.checkSameSize(other)

	for i, chunk := range this.chunks {
		if chunk != other.chunks[i] {
			return false
		}
	}

	return true
}


/* *** private *** */

func (this Row) checkIndex(index int) {
	if index < 0 || index >= this.columnCount {
		panic(fmt.Sprint("index out of bounds", index, "!! [",0,",",this.columnCount,")"))
	}
}


func (this Row) checkSameSize(other *Row) {
	if this.columnCount != other.columnCount {
		panic(fmt.Sprint("cannot perform this operation on two rows of differing size. columnCount",
				this.columnCount, "!=", other.columnCount))
	}
}


func (this Row) convertIndex(index int) (column int, bit uint64, exp uint32) {
	return ((len(this.chunks)*64 - 1 - index) / 64), 1 << (uint(index) % 64), (uint32(index) % 64)
}


/* *** LinearSystem *** **************************************************** */
type LinearSystem struct {
	rows []*Row
	rowCount, columnCount int
}


func NewLinearSystem(rows, columns int) *LinearSystem {

	if rows < 0 || columns < 0 {
		panic(fmt.Sprint("columnCount", columns, "< 0 or rowCount", rows, "< 0"))
	}

	var ret LinearSystem
	ret.rowCount = rows
	ret.columnCount = columns
	ret.rows = make([]*Row, rows)

	for i, _ := range ret.rows {
		ret.rows[i] = NewRow(columns)
	}

	return &ret
}


func (this LinearSystem) Row(index int) *Row {
	this.checkRowIndex(index)
	return this.rows[index]
}


func (this *LinearSystem) SetRow(index int, row *Row) {
	this.checkRowIndex(index)
	this.rows[index].Set(row)
}


func (this *LinearSystem) Set(other *LinearSystem) *LinearSystem {

	this.checkSameSize(other)

	for i, row := range other.rows {
		this.SetRow(i, row)
	}

	return this
}


func (m *LinearSystem) GaussianElimination(other *LinearSystem) *LinearSystem {

	m.checkSameSize(other)

	m.Set(other)

	startingRow := 0

	for column := m.columnCount - 1; column >= 0; column -= 1 {

		var row int
		for row = startingRow; row < m.rowCount; row += 1 {
			if m.Row(row).Column(column) == 1 {
				m.Row(startingRow).Swap(m.Row(row))
				break
			}
		}

		if row == m.rowCount {
			/* no row has been found that has a bit at the wanted column,
			try again using the next column to the left */
			continue
		}

		for row = startingRow + 1; row < m.rowCount; row += 1 {
			if m.Row(row).Column(column) == 1 {
				m.Row(row).Xor(m.Row(row),m.Row(startingRow))
			}
		}

		startingRow += 1
	}

	return m
}


func (this LinearSystem) String() string {
	var ret string
	for _, k := range this.rows {
		ret += fmt.Sprint(k)
		ret += "\n"
	}
	return ret
}


func (this LinearSystem) Equals(other *LinearSystem) bool {

	this.checkSameSize(other)

	for i, row := range this.rows {
		if row.Equals(other.Row(i)) == false {
			return false
		}
	}

	return true
}


/* *** private *** */
func (this LinearSystem) checkRowIndex(i int) {
	if i < 0 || i >= this.rowCount {
		panic(fmt.Sprint("invalid index:", i, " is not element of [0 ,", this.rowCount,")"))
	}
}


func (this LinearSystem) checkSameSize(other *LinearSystem) {
	if this.rowCount != other.rowCount || this.columnCount != other.columnCount {
		panic(fmt.Sprint("cannot perform operation on two linear systems of differing size. columnCount:",
				this.rowCount, "!=", other.rowCount, "or", this.columnCount, "!=", other.columnCount))
	}
}


/* *** helper *** ********************************************************** */
