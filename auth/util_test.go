package auth

import (
	"reflect"
	"testing"

	"gopkg.in/guregu/null.v4"
)

func TestConvertToNullTypeHookFunc(t *testing.T) {
	type args struct {
		f    reflect.Type
		t    reflect.Type
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "ConvertToNullTypeHookFunc",
			args: args{
				f:    reflect.TypeOf(null.String{}),
				t:    reflect.TypeOf(null.Int{}),
				data: "123",
			},
			want:    "123",
			wantErr: false,
		},
		{
			name: "ConvertToNullTypeHookFunc_int",
			args: args{
				f:    reflect.TypeOf(null.String{}),
				t:    reflect.TypeOf(null.Int{}),
				data: int(123),
			},
			want:    null.IntFrom(123),
			wantErr: false,
		},
		{
			name: "ConvertToNullTypeHookFunc_int64",
			args: args{
				f:    reflect.TypeOf(null.String{}),
				t:    reflect.TypeOf(null.Int{}),
				data: int64(123),
			},
			want:    null.IntFrom(123),
			wantErr: false,
		},
		{
			name: "ConvertToNullTypeHookFunc_float",
			args: args{
				f:    reflect.TypeOf(null.String{}),
				t:    reflect.TypeOf(null.Float{}),
				data: float64(123.45),
			},
			want:    null.FloatFrom(123.45),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToNullTypeHookFunc(tt.args.f, tt.args.t, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToNullTypeHookFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				// t.Errorf("ConvertToNullTypeHookFunc() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
				t.Errorf("ConvertToNullTypeHookFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
