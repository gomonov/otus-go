package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		envName       string
		env           EnvValue
		expectedCode  int
		expectedValue string
	}{
		{
			envName: "BAR",
			env: EnvValue{
				Value:      "bar",
				NeedRemove: false,
			},
			expectedCode:  0,
			expectedValue: "bar\n",
		},
		{
			envName: "EMPTY",
			env: EnvValue{
				Value:      "",
				NeedRemove: false,
			},
			expectedCode:  0,
			expectedValue: "\n",
		},
		{
			envName: "UNSET",
			env: EnvValue{
				Value:      "",
				NeedRemove: true,
			},
			expectedCode:  1,
			expectedValue: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.envName, func(t *testing.T) {
			mapEnv := make(map[string]EnvValue)
			mapEnv[tc.envName] = tc.env

			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			code := RunCmd([]string{"printenv", tc.envName}, mapEnv)

			w.Close()
			os.Stdout = old

			out, _ := io.ReadAll(r)

			require.Equal(t, tc.expectedCode, code)

			if tc.expectedCode == 0 {
				require.Equal(t, tc.expectedValue, string(out))
			}
		})
	}
}
