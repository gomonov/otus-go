package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	tests := []struct {
		env                string
		expectedValue      string
		expectedNeedRemove bool
	}{
		{env: "BAR", expectedValue: "bar", expectedNeedRemove: false},
		{env: "EMPTY", expectedValue: "", expectedNeedRemove: false},
		{env: "FOO", expectedValue: "   foo\nwith new line", expectedNeedRemove: false},
		{env: "UNSET", expectedValue: "", expectedNeedRemove: true},
		{env: "HELLO", expectedValue: "\"hello\"", expectedNeedRemove: false},
	}

	environment, _ := ReadDir("testdata/env")

	for _, tc := range tests {
		t.Run(tc.env, func(t *testing.T) {
			require.Equal(t, tc.expectedValue, environment[tc.env].Value)
			require.Equal(t, tc.expectedNeedRemove, environment[tc.env].NeedRemove)
		})
	}
}

func TestInvalidDir(t *testing.T) {
	t.Run("invalid dir", func(t *testing.T) {
		_, err := ReadDir("testdata/e")
		require.Error(t, err)
		require.ErrorContains(t, err, "open testdata/e: no such file or directory")
	})
}
