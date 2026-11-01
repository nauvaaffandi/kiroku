package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testPIN = "654321"

func TestPinVerificationWithTestPin(t *testing.T) {
	hash, err := HashPin(testPIN, 3, 65536, 4)
	require.NoError(t, err)

	ok, err := VerifyPin(testPIN, hash)
	require.NoError(t, err)
	assert.True(t, ok, "test pin should verify correctly")
}

func TestVerifyPin_WrongPin(t *testing.T) {
	hash, err := HashPin(testPIN, 3, 65536, 4)
	require.NoError(t, err)

	ok, err := VerifyPin("wrongpin", hash)
	require.NoError(t, err)
	assert.False(t, ok, "wrong pin should not verify")
}

func TestVerifyPin_InvalidHash(t *testing.T) {
	ok, err := VerifyPin(testPIN, "invalid-hash")
	assert.Error(t, err)
	assert.False(t, ok)
}

func TestValidatePIN_Valid(t *testing.T) {
	require.NoError(t, ValidatePIN(testPIN))
	require.NoError(t, ValidatePIN("000000"))
	require.NoError(t, ValidatePIN("999999"))
}

func TestValidatePIN_InvalidLength(t *testing.T) {
	require.Error(t, ValidatePIN("12345"))
	require.Error(t, ValidatePIN("1234567"))
	require.Error(t, ValidatePIN(""))
}

func TestValidatePIN_NonNumeric(t *testing.T) {
	require.Error(t, ValidatePIN("abcdef"))
	require.Error(t, ValidatePIN("12345a"))
	require.Error(t, ValidatePIN("12 456"))
}
