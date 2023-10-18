package main

import (
	"fmt"
)

const blockLength = 64
const hexsInBlock = blockLength / 4

type MyBigInt struct {
	value      []uint64
	blockCount int
}

func (this *MyBigInt) SetHex(s string) {
	this.blockCount = len(s) / hexsInBlock
	if len(s)%hexsInBlock > 0 {
		this.blockCount++
	}
	this.value = make([]uint64, this.blockCount)
	for i, end := 0, len(s); end > 0; end -= hexsInBlock {
		start := max(end-hexsInBlock, 0)
		block := s[start:end]
		fmt.Sscan("0x"+block, &this.value[i])
		i++
	}
}

func (this *MyBigInt) GetHex() string {
	s := ""
	for i := this.blockCount - 1; i >= 0; i-- {
		s += fmt.Sprintf("%x", this.value[i])
	}
	return s
}

func (this *MyBigInt) String() string {
	return this.GetHex()
}

func INV(a *MyBigInt) *MyBigInt {
	len := a.blockCount
	res := MyBigInt{make([]uint64, len), len}
	for i := 0; i < len; i++ {
		res.value[i] = ^a.value[i]
	}
	return &res
}

func AND(a, b *MyBigInt) *MyBigInt {
	min := min(a.blockCount, b.blockCount)
	res := MyBigInt{make([]uint64, min), min}
	for i := 0; i < min; i++ {
		res.value[i] = a.value[i] & b.value[i]
	}
	return &res
}

func OR(a, b *MyBigInt) *MyBigInt {
	len := max(a.blockCount, b.blockCount)
	min := a.blockCount + b.blockCount - len
	res := MyBigInt{make([]uint64, len), len}
	for i := 0; i < min; i++ {
		res.value[i] = a.value[i] | b.value[i]
	}
	if a.blockCount > min {
		copy(res.value[min+1:], a.value[min+1:])
	} else if b.blockCount > min {
		copy(res.value[min+1:], b.value[min+1:])
	}
	return &res
}

func XOR(a, b *MyBigInt) *MyBigInt {
	len := max(a.blockCount, b.blockCount)
	min := a.blockCount + b.blockCount - len
	res := MyBigInt{make([]uint64, len), len}
	for i := 0; i < min; i++ {
		res.value[i] = a.value[i] ^ b.value[i]
	}
	if a.blockCount > min {
		copy(res.value[min+1:], a.value[min+1:])
	} else if b.blockCount > min {
		copy(res.value[min+1:], b.value[min+1:])
	}
	return &res
}

func (this *MyBigInt) trunc() *MyBigInt {
	i := this.blockCount - 1
	for ; i > 0; i-- {
		if this.value[i] != 0 {
			break
		}
	}
	this.blockCount = i + 1
	this.value = this.value[:this.blockCount]
	return this
}

func ShiftR(a *MyBigInt, shift int) *MyBigInt {
	res := MyBigInt{make([]uint64, a.blockCount), a.blockCount}
	blockShift := shift / 64
	shift %= 64
	copy(res.value, a.value[blockShift:])
	res.trunc()
	if shift == 0 {
		return &res
	}
	buf, t := uint64(0), uint64(0)
	for i := res.blockCount - 1; i >= 0; i-- {
		t = res.value[i] << (blockLength - shift)
		res.value[i] >>= shift
		res.value[i] += buf
		buf = t
	}
	return res.trunc()
}

func ShiftL(a *MyBigInt, shift int) *MyBigInt {
	blockShift := shift / 64
	shift %= 64
	len := a.blockCount + blockShift + 1
	res := MyBigInt{make([]uint64, len), len}
	copy(res.value[blockShift:], a.value)
	if shift != 0 {
		buf, t := uint64(0), uint64(0)
		for i := blockShift; i < len; i++ {
			t = res.value[i] >> (blockLength - shift)
			res.value[i] <<= shift
			res.value[i] += buf
			buf = t
		}
	}
	return res.trunc()
}

func main() {
	a := &MyBigInt{}
	b := &MyBigInt{}
	a.SetHex("a6378bc8e8ada6345678e8ada6378bc8e8ada637abeebaad1")
	b.SetHex("5a6378bc8e8ada6345678e8ada6378bc8e8ada637abeebaad6")
	fmt.Println("A:\n" + a.String())
	fmt.Println("B:\n" + b.String())

	c := OR(a, b)
	fmt.Println("A or B:\n" + c.String())

	d := AND(a, b)
	fmt.Println("A and B:\n" + d.String())

	e := XOR(a, b)
	fmt.Println("A xor B:\n" + e.String())

	f := ShiftR(a, 80)
	fmt.Println("A >> 80:\n" + f.String())

	g := ShiftL(a, 24)
	fmt.Println("A << 24:\n" + g.String())

	fmt.Println()
	fmt.Println()

	/*f.SetHex("000000000000000000000000000000c8e8ada637abeebaad1")
	fmt.Println(f.trunc().value)

	/*e := &MyBigInt{}
	e.SetHex("6378bc8e8ada6345678e8ada6378bc8e8ada637abeebaad1")
	fmt.Println("E:\n" + e.GetHex())
	f := INV(e)
	fmt.Println("inverted E:\n" + f.GetHex())
	g := INV(f)
	fmt.Println("E inverted twice:\n" + g.GetHex())*/
}
