package zksnark_base

import (
	"crypto/rand"
	"math/big"
)

func GetRandomWithMax(max *big.Int) *big.Int {
	val, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}

	return val
}
