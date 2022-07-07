package logging

import (
	"9m/consts"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_SetupLogger(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "Test_SetupLogger",
			want: []byte(`{"level":"info","message":"hello world"}` + "\n"),
		},
	}

	fname := consts.ProgramName + "-" + time.Now().Format("20060102") + ".log"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetupLogger()

			Object.Info().Msg("hello world")
			data, _ := ioutil.ReadFile(fname)

			F.Close()
			os.Remove(fname)

			require.Equal(t, tt.want, data, "Not equal")
		})
	}
}

func Test_RenewLogger(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "Test_SetupLogger",
			want: []byte(`{"level":"info","message":"hello world"}` + "\n"),
		},
	}

	fname := consts.ProgramName + "-" + time.Now().Format("20060102") + ".log"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RenewLogger()

			Object.Info().Msg("hello world")
			data, _ := ioutil.ReadFile(fname)

			F.Close()
			os.Remove(fname)

			require.Equal(t, tt.want, data, "Not equal")
		})
	}
}
