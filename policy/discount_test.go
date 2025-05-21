package policy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMSBEven(t *testing.T) {
	iter := uint64(256)
	msbEven := MSBEven(iter)
	require.Equal(t, uint64(iter), msbEven)
}
