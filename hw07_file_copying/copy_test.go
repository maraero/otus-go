package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	cases := []struct {
		offset int64
		limit  int64
	}{
		{offset: 0, limit: 0},
		{offset: 0, limit: 10},
		{offset: 0, limit: 1000},
		{offset: 0, limit: 10000},
		{offset: 100, limit: 1000},
		{offset: 6000, limit: 1000},
	}
	fromPath := "testdata/input.txt"

	for _, c := range cases {
		c := c
		expected := fmt.Sprintf("testdata/out_offset%d_limit%d.txt", c.offset, c.limit)
		toPath := fmt.Sprintf("testdata/output_offset%d_limit%d.txt", c.offset, c.limit)
		t.Run(expected, func(t *testing.T) {
			err := Copy(fromPath, toPath, c.offset, c.limit)
			defer os.Remove(toPath)
			require.NoError(t, err)

			toContent, err := os.ReadFile(toPath)
			require.NoError(t, err)
			expectedContent, err := os.ReadFile(expected)
			require.NoError(t, err)
			require.Equal(t, string(toContent), string(expectedContent))
		})
	}
}
