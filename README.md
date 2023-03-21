# The zkSNARKOP repo on Golang

Usage example:

Firstly, you should define t(x), h(x) and p(x) there p(x) = t(x) * h (x)
```go
type (
    T func(x *big.Int) *big.Int
    H func(xi []*bn256.G2) *bn256.G2
    P func(xi []*bn256.G2) *bn256.G2
)
```

That zkSNARKOP uses [Cloudflare bn256 Bilinear map implementation](https://github.com/cloudflare/bn256) where G2 is defined.
Explore [Test example](./main_test.go) to see how t,h and p should be defined.

Use `Setup(t T, n uint64) *SetupParams` method to generate SetupParams.

```go
params := Setup(testT, 4)
```

Generate Proof
```go
proof := MakeProof(params, testH, testP)
```

Verify Proof
```go
err := VerifyProof(params, proof)
```