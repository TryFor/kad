package dht

import (
	"math/big"
)



func binew(x int64) *big.Int {
	return big.NewInt(x)
}

func id2bi(id Identifier) *big.Int {
	return big.NewInt(0).SetBytes(id)
}

func birsh(x int64, n uint) *big.Int {
	return big.NewInt(0).Rsh(binew(x), n)
}

func bilsh(x int64, n uint) *big.Int {
	return big.NewInt(0).Lsh(binew(x), n)
}

func bisub(x, y *big.Int) *big.Int {
	return big.NewInt(0).Sub(x, y)
}

func bidiv(x, y *big.Int) *big.Int {
	return big.NewInt(0).Div(x, y)
}

func bimid(max, min *big.Int) *big.Int {
	d := bisub(max, min)
	return d.Rsh(d, 1).Add(d, min)
}

func biadd(x, y *big.Int) *big.Int {
	return big.NewInt(0).Add(x, y)
}
