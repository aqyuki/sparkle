package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOrErr(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		key     string
		value   string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:  "success to get a value",
			value: "test",
			args: args{
				key: "TEST",
			},
			want:    "test",
			wantErr: false,
		},
		{
			name:    "failed to get a value",
			args:    args{key: "TEST"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				t.Setenv(tt.args.key, tt.value)
			}
			v, err := GetOrErr(tt.args.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, v)
			}
		})
	}
}
