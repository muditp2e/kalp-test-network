package chaincode_test

import (
	"chaincode-go/chaincode"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	calculator := chaincode.SmartContract{}
	val := calculator.Add(4, 8)
	require.Equal(t, 12, val)
}

func TestSub(t *testing.T) {
	calculator := chaincode.SmartContract{}
	val := calculator.Sub(4, 8)
	require.Equal(t, -4, val)
}

func TestMul(t *testing.T) {
	calculator := chaincode.SmartContract{}
	val := calculator.Mul(4, 8)
	require.Equal(t, 32, val)
}

func TestDiv(t *testing.T) {
	calculator := chaincode.SmartContract{}
	val, err := calculator.Div(8, 3)
	require.Equal(t, 2, val)
	require.NoError(t, err)

	_, err = calculator.Div(8, 0)
	require.EqualError(t, err, "division by zero")
}
