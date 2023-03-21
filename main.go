package zksnark_base

import (
	"bytes"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/cloudflare/bn256"
)

type (
	T func(x *big.Int) *big.Int
	H func(xi []*bn256.G2) *bn256.G2
	P func(xi []*bn256.G2) *bn256.G2

	SetupParams struct {
		G1       *bn256.G1
		G2       *bn256.G2
		G1_ts    *bn256.G1
		G1_alpha *bn256.G1

		G1_si []*bn256.G1
		G2_si []*bn256.G2

		G1_alphasi []*bn256.G1
		G2_alphasi []*bn256.G2

		N uint64
	}

	Proof struct {
		G_delta_ps       *bn256.G2
		G_delta_hs       *bn256.G2
		G_delta_alpha_ps *bn256.G2
	}
)

func Setup(t T, n uint64) *SetupParams {
	s, alpha := GetRandomWithMax(bn256.Order), GetRandomWithMax(bn256.Order)
	defer func() {
		s.SetUint64(0)
		alpha.SetUint64(0)
	}()

	_, g1, err := bn256.RandomG1(rand.Reader)
	if err != nil {
		panic(err)
	}

	_, g2, err := bn256.RandomG2(rand.Reader)
	if err != nil {
		panic(err)
	}

	g1_si := make([]*bn256.G1, 0, n+1)
	g2_si := make([]*bn256.G2, 0, n+1)

	g1_alphasi := make([]*bn256.G1, 0, n+1)
	g2_alphasi := make([]*bn256.G2, 0, n+1)

	for i := uint64(0); i <= n; i++ {
		si := new(big.Int).Exp(s, new(big.Int).SetUint64(i), bn256.Order)
		g1_si = append(g1_si, new(bn256.G1).ScalarMult(g1, si))
		g2_si = append(g2_si, new(bn256.G2).ScalarMult(g2, si))
	}

	for i := uint64(0); i <= n; i++ {
		si := new(big.Int).Exp(s, new(big.Int).SetUint64(i), bn256.Order)
		alphasi := new(big.Int).Mod(new(big.Int).Mul(alpha, si), bn256.Order)
		g1_alphasi = append(g1_alphasi, new(bn256.G1).ScalarMult(g1, alphasi))
		g2_alphasi = append(g2_alphasi, new(bn256.G2).ScalarMult(g2, alphasi))
	}

	return &SetupParams{
		G1:         g1,
		G2:         g2,
		G1_ts:      new(bn256.G1).ScalarMult(g1, t(s)),
		G1_alpha:   new(bn256.G1).ScalarMult(g1, alpha),
		G1_si:      g1_si,
		G2_si:      g2_si,
		G1_alphasi: g1_alphasi,
		G2_alphasi: g2_alphasi,
		N:          n,
	}
}

func MakeProof(params *SetupParams, h H, p P) *Proof {
	delta := GetRandomWithMax(bn256.Order)
	defer func() {
		delta.SetUint64(0)
	}()

	return &Proof{
		G_delta_ps:       new(bn256.G2).ScalarMult(p(params.G2_si), delta),
		G_delta_hs:       new(bn256.G2).ScalarMult(h(params.G2_si), delta),
		G_delta_alpha_ps: new(bn256.G2).ScalarMult(p(params.G2_alphasi), delta),
	}
}

func VerifyProof(params *SetupParams, proof *Proof) error {
	if !bytes.Equal(bn256.Pair(params.G1, proof.G_delta_ps).Marshal(), bn256.Pair(params.G1_ts, proof.G_delta_hs).Marshal()) {
		return errors.New("check #1 : e(g^delta*p(s), g) == e(g^t(s), g^delta*h(s)) failed")
	}

	if !bytes.Equal(bn256.Pair(params.G1, proof.G_delta_alpha_ps).Marshal(), bn256.Pair(params.G1_alpha, proof.G_delta_ps).Marshal()) {
		return errors.New("check #2 : e(g^delta*alpha*p(s), g) == e(g^alpha, g^delta*p(s)) failed")
	}

	return nil
}
