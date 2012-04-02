package main


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


/* *** Row *** ************************************************************* */
type Row struct {
	columns []uint64
	columnCount int
}


func NewRow(columnCount int) *Row {

	if columnCount < 0 {
		panic("column count has to be >= 0")
	}

	var ret Row
	ret.columnCount = columnCount
	ret.columns = make([]uint64, (columnCount-1/64)+1)

	for i, _ := range ret.columns {
		/* initialize to zero */
		ret.columns[i] = 0x0000000000000000
	}

	return &ret
}


func (this Row) Column(index int) Bit {
	this.checkIndex(index)

	column, bit, exp := convertIndex(index)

	ret := Bit((this.columns[column] & bit) >> exp)
	ret.Check()
	return ret
}


func (this *Row) SetColumn(index int, value Bit) {
	value.Check()
	this.checkIndex(index)

	column, bit, _ := convertIndex(index)

	if value == 0 {
		this.columns[column] &= ^bit
	} else {
		this.columns[column] |= bit
	}
}


func (this *Row) Set(other *Row) *Row {

	if this.columnCount != other.columnCount {
		panic("cannot assign Rows of different sizes")
	}

	for i, k := range other.columns {
		this.columns[i] = k;
	}

	return this
}


func (this *Row) Xor(a, b *Row) *Row {

	if this.columnCount != a.columnCount || a.columnCount != b.columnCount || this.columnCount != b.columnCount {
		panic(fmt.Sprint("cannot xor/set rows of differing columnCount:",
				this.columnCount, ",", a.columnCount, ",", b.columnCount))
	}

	for i, _ := range a.columns {
		this.columns[i] = a.columns[i] ^ b.columns[i]
	}

	return this
}


func (this *Row) Neg(other *Row) *Row {

	if this.columnCount != other.columnCount {
		panic(fmt.Sprint("cannot set rows of differing columnCount:",
				this.columnCount, " !=", other.columnCount))
	}

	for i, k := range other.columns {
		this.columns[i] = ^k
	}

	return this
}


/* *** private *** */

func (this Row) checkIndex(index int) {
	if index < 0 || index >= this.columnCount {
		panic(fmt.Sprint("index out of bounds", index, "!! [",0,",",this.columnCount,"]"))
	}
}


/* *** LinearSystem *** **************************************************** */
type LinearSystem struct {
	rows map[int]Row
}

func NewLinearSystem() *LinearSystem {
	var ret LinearSystem
	ret.rows = make(map[int]Row)
	return &ret
}


/* *** helper *** ********************************************************** */

func convertIndex(index int) (column int, bit uint64, exp uint32) {
	return index / 64, 1 << (uint(index) % 64), (uint32(index) % 64)
}

