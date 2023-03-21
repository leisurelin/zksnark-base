package zksnark_base

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
)

func GetRandomWithMax(max *big.Int) *big.Int {
	val, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}

	return val
}

func log(i interface{}) {
	data, err := json.MarshalIndent(i, "", "")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
