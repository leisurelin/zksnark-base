package zksnark_base

import (
	"math/big"
	"testing"

	"github.com/cloudflare/bn256"
)

// test p = (x + 1) * (x + 2) * (x + 3) * (x + 4) = x^4 + 7x3 + 12x2 + 3x3 + 21x2 + 36x + 2x2 + 14x + 24
// = x^4 + 10x^3 + 35x^2 + 50x + 24

// t = (x + 1) * (x + 2) = x^2 + 3x + 2
// h = (x + 3) * (x + 4) = x^2 + 7x + 12

var testT = func(x *big.Int) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Add(x, big.NewInt(1)),
			new(big.Int).Add(x, big.NewInt(2)),
		),
		bn256.Order,
	)
}

var testH = func(x []*bn256.G2) *bn256.G2 {
	var c = []*big.Int{big.NewInt(12), big.NewInt(7), big.NewInt(1), big.NewInt(0), big.NewInt(0)}
	var e = make([]*bn256.G2, 0, 5)

	for i, val := range x {
		e = append(e, new(bn256.G2).ScalarMult(val, c[i]))
	}

	var res = e[0]
	for i := 1; i < 5; i++ {
		res = new(bn256.G2).Add(e[i], res)
	}

	return res
}

var testP = func(x []*bn256.G2) *bn256.G2 {
	var c = []*big.Int{big.NewInt(24), big.NewInt(50), big.NewInt(35), big.NewInt(10), big.NewInt(1)}
	var e = make([]*bn256.G2, 0, 5)

	for i, val := range x {
		e = append(e, new(bn256.G2).ScalarMult(val, c[i]))
	}

	var res = e[0]
	for i := 1; i < 5; i++ {
		res = new(bn256.G2).Add(e[i], res)
	}

	return res
}

func TestProving(_ *testing.T) {
	params := Setup(testT, 4)
	proof := MakeProof(params, testH, testP)

	if err := VerifyProof(params, proof); err != nil {
		panic(err)
	}
}
