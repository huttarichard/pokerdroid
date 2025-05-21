package fsutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClean(t *testing.T) {
	testCases := []struct {
		desc      string
		dir, path string
		expected  string
		err       bool
	}{
		{
			desc:     "1",
			dir:      "base",
			path:     "ptx",
			expected: "base/ptx",
		},
		{
			desc:     "2",
			dir:      "base",
			path:     "/ptx",
			expected: "base/ptx",
		},
		{
			desc:     "3",
			dir:      "/base",
			path:     "/ptx",
			expected: "/base/ptx",
		},
		{
			desc: "4",
			dir:  "/base",
			path: "../ptx",
			err:  true,
		},
		{
			desc:     "5",
			dir:      "x",
			path:     "../x/ok",
			expected: "x/ok",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			x, err := clean(tC.dir, tC.path)
			if tC.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, x, tC.expected)
			}
		})
	}
}
