package main

import (
	"fmt"
	"math"
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
	s += fmt.Sprintf("%x", this.value[this.blockCount-1])
	for i := this.blockCount - 2; i >= 0; i-- {
		s += fmt.Sprintf("%016x", this.value[i])
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

func Add(a, b *MyBigInt) *MyBigInt {
	len := max(a.blockCount, b.blockCount)
	min := a.blockCount + b.blockCount - len
	res := MyBigInt{make([]uint64, len+1), len + 1}
	carry := uint64(0)
	aTopBit, bTopBit, resTopBit := uint64(0), uint64(0), uint64(0)
	for i := 0; i < min; i++ {
		res.value[i] = a.value[i] + b.value[i] + carry
		aTopBit = a.value[i] >> 63
		bTopBit = b.value[i] >> 63
		if aTopBit+bTopBit == 2 {
			carry = 1
			continue
		}
		if aTopBit+bTopBit == 0 {
			carry = 0
			continue
		}
		resTopBit = res.value[i] >> 63
		if resTopBit == 0 {
			carry = 1
		} else {
			carry = 0
		}
	}
	if a.blockCount > min {
		for i := min; carry != 0 && i < len; i++ {
			res.value[i] = a.value[i] + carry
			if res.value[i] != 0 {
				carry = 0
			} else {
				carry = 1
			}
		}
	} else if b.blockCount > min {
		for i := min; carry != 0 && i < len; i++ {
			res.value[i] = b.value[i] + carry
			if res.value[i] != 0 {
				carry = 0
			} else {
				carry = 1
			}
		}
	}
	if carry == 1 {
		res.value[len] = 1
	}
	return res.trunc()
}

func (a *MyBigInt) Comp(b *MyBigInt) int {
	a.trunc()
	b.trunc()
	aLen := a.blockCount
	bLen := b.blockCount
	if aLen > bLen {
		return 1
	}
	if aLen < bLen {
		return -1
	}
	for i := aLen - 1; i >= 0; i-- {
		if a.value[i] > b.value[i] {
			return 1
		}
		if a.value[i] < b.value[i] {
			return -1
		}
	}
	return 0
}

func Sub(a, b *MyBigInt) (*MyBigInt, error) {
	switch a.Comp(b) {
	case -1:
		return nil, fmt.Errorf("The second number is bigger")
	case 0:
		return &MyBigInt{value: []uint64{uint64(0)}, blockCount: 1}, nil
	}

	len := a.blockCount
	min := b.blockCount
	res := MyBigInt{make([]uint64, len), len}

	borrow := uint64(0)
	for i := 0; i < min; i++ {
		if a.value[i]-borrow >= b.value[i] {
			res.value[i] = a.value[i] - borrow - b.value[i]
			borrow = 0
			continue
		}
		res.value[i] = math.MaxUint64 - b.value[i] + 1 - borrow + a.value[i]
		borrow = 1
	}
	for i := min; i < len; i++ {
		if a.value[i]-borrow >= 0 {
			res.value[i] = a.value[i] - borrow
			borrow = 0
			continue
		}
		res.value[i] = math.MaxUint64
		borrow = 1
	}
	return res.trunc(), nil
}

func main() {
	a := &MyBigInt{}
	a.SetHex("51bf608414ad5726a3c1bec098f77b1b54ffb2787f8d528a74c1d7fde6470ea4")
	fmt.Println("A:\n" + a.GetHex())

	b := &MyBigInt{}
	b.SetHex("403db8ad88a3932a0b7e8189aed9eeffb8121dfac05c3512fdb396dd73f6331c")
	fmt.Println("B:\n" + b.GetHex())

	c := XOR(a, b)
	fmt.Println("A xor B:\n" + c.GetHex())
	if c.GetHex() == "1182d8299c0ec40ca8bf3f49362e95e4ecedaf82bfd167988972412095b13db8" {
		fmt.Println("XOR is correct")
	} else {
		fmt.Println("XOR is wrong")
	}

	c = OR(a, b)
	fmt.Println("A or B:\n" + c.GetHex())
	if c.GetHex() == "51bff8ad9cafd72eabffbfc9befffffffcffbffaffdd779afdf3d7fdf7f73fbc" {
		fmt.Println("OR is correct")
	} else {
		fmt.Println("OR is wrong")
	}

	c = AND(a, b)
	fmt.Println("A and B:\n" + c.GetHex())
	if c.GetHex() == "403d208400a113220340808088d16a1b10121078400c1002748196dd62460204" {
		fmt.Println("AND is correct")
	} else {
		fmt.Println("AND is wrong")
	}

	e := &MyBigInt{}
	e.SetHex("6378bc8e8ada6345678e8ada6378bc8e8ada637abeebaad1")
	fmt.Println("\nE:\n" + e.GetHex())
	f := INV(e)
	fmt.Println("inverted E:\n" + f.GetHex())
	g := INV(f)
	fmt.Println("E inverted twice:\n" + g.GetHex() + "\n")

	a.SetHex("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80")
	fmt.Println("A:\n" + a.GetHex())
	b.SetHex("70983d692f648185febe6d6fa607630ae68649f7e6fc45b94680096c06e4fadb")
	fmt.Println("B:\n" + b.GetHex())

	c = Add(a, b)
	fmt.Println("A + B:\n" + c.GetHex())
	if c.GetHex() == "a78865c13b14ae4e25e90771b54963ee2d68c0a64d4a8ba7c6f45ee0e9daa65b" {
		fmt.Println("Addition is correct")
	} else {
		fmt.Println("Addition is wrong")
	}

	a.SetHex("33ced2c76b26cae94e162c4c0d2c0ff7c13094b0185a3c122e732d5ba77efebc")
	fmt.Println("A:\n" + a.GetHex())
	b.SetHex("22e962951cb6cd2ce279ab0e2095825c141d48ef3ca9dabf253e38760b57fe03")
	fmt.Println("B:\n" + b.GetHex())

	c, er := Sub(a, b)
	if er == nil {
		fmt.Println("A - B:\n" + c.GetHex())
		if c.GetHex() == "10e570324e6ffdbc6b9c813dec968d9bad134bc0dbb061530934f4e59c2700b9" {
			fmt.Println("Subtraction is correct")
		} else {
			fmt.Println("Subtraction is wrong")
		}
	} else {
		fmt.Println(er)
	}

	fmt.Println("A >> 80:\n" + ShiftR(a, 80).String())
	fmt.Println("A << 36:\n" + ShiftL(a, 36).String())

}
