package handlers

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLicensePlateValidator(t *testing.T) {
	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{
			name:   "empty number",
			number: "",
			want:   false,
		},
		{
			name:   "too many characters",
			number: "0123456789",
			want:   false,
		},
		{
			name:   "postfix",
			number: "A123AA",
			want:   true,
		},
		{
			name:   "prefix",
			number: "АВ123АВ",
			want:   false,
		},
		{
			name:   "only lowercase letters",
			number: "а123вв",
			want:   true,
		},
		{
			name:   "english letters",
			number: "A123CB",
			want:   true,
		},
		{
			name:   "english lowercase letters",
			number: "b123ca",
			want:   true,
		},
		{
			name:   "with udarenie",
			number: "а́123aa",
			want:   false,
		},
		{
			name:   "english and russian lowercase mix",
			number: "a111аа",
			want:   true,
		},
		{
			name:   "english and russian uppercase mix",
			number: "A123ВС",
			want:   true,
		},
		{
			name:   "english and russian uppercase and lowercase mix",
			number: "a222Вв",
			want:   true,
		},
		{
			name:   "few numbers",
			number: "a22aa",
			want:   false,
		},
		{
			name:   "three zeroes number",
			number: "a000aa",
			want:   false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := LicensePlateValidatorByString(tt.number)

			require.Equal(t, tt.want, res)
		})
	}
}
