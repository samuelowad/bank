package util

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := RandString(8)

	hashedPassword, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
}

func TestComparePassword(t *testing.T) {
	password := RandString(8)

	hashedPassword, _ := HashedPassword(password)
	fmt.Printf("hashedPassword: %v\n", hashedPassword)

	err := ComparePassword(password, hashedPassword)
	require.NoError(t, err)
}
