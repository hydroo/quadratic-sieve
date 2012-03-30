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


func (this Row) CheckIndex(index int) {
	if index < 0 || index >= this.columnCount {
		panic(fmt.Sprint("index out of bounds", index, "!! [",0,",",this.columnCount,"]"))
	}
}


func (this Row) Column(index int) Bit {
	this.CheckIndex(index)

	column, bit, exp := ConvertIndex(index)

	ret := Bit((this.columns[column] & bit) >> exp)
	ret.Check()
	return ret
}


func (this *Row) SetColumn(index int, value Bit) {
	value.Check()
	this.CheckIndex(index)

	column, bit, _ := ConvertIndex(index)

	if value == 0 {
		this.columns[column] &= ^bit
	} else {
		this.columns[column] |= bit
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

func ConvertIndex(index int) (column int, bit uint64, exp uint32) {
	return index / 64, 1 << (uint(index) % 64), (uint32(index) % 64)
}

