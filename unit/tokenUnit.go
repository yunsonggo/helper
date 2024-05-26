package unit

import (
	"fmt"
	"math"
	"math/big"
)

func WeiToEth(b *big.Int) string {
	f := new(big.Float)
	f.SetString(b.String())
	s := new(big.Float).Quo(f, big.NewFloat(math.Pow10(18)))
	return s.String()
}

func EthStrToWei(ethString string) (*big.Int, error) {
	if ethString == "" {
		return new(big.Int), nil
	}
	ethValue, _ := new(big.Float).SetString(ethString)
	weiValue := new(big.Int)
	weiValue, _ = weiValue.SetString(fmt.Sprintf("%.0f", new(big.Float).Mul(ethValue, big.NewFloat(1e18))), 10)
	return weiValue, nil
}

func WeiStrToBig(wei string) *big.Int {
	i := new(big.Int)
	i.SetString(wei, 0)
	return i
}

func WeiToEthFloat(b *big.Int) float64 {
	f := new(big.Float)
	f.SetString(b.String())
	s := new(big.Float).Quo(f, big.NewFloat(math.Pow10(18)))
	v, _ := s.Float64()
	return v
}

func EthFloat64ToWei(eth float64) *big.Int {
	ethValue := new(big.Float).SetFloat64(eth)
	weiValue := new(big.Int)
	weiValue, _ = weiValue.SetString(fmt.Sprintf("%.0f", new(big.Float).Mul(ethValue, big.NewFloat(1e18))), 10)
	return weiValue
}

func GweiToEth(gasFee *big.Int) string {
	f := new(big.Float)
	f.SetString(gasFee.String())
	s := new(big.Float).Quo(f, big.NewFloat(math.Pow10(10)))
	return s.Text('f', 10)
}

func EthToGwei(eth string) (*big.Int, error) {
	if eth == "" {
		return new(big.Int), nil
	}
	ethValue, _ := new(big.Float).SetString(eth)
	gwei := new(big.Int)
	gwei, _ = gwei.SetString(fmt.Sprintf("%.0f", new(big.Float).Mul(ethValue, big.NewFloat(1e10))), 10)
	return gwei, nil
}
