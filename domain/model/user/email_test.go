package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "空の場合エラー",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.input, string(email))
		})
	}
}

func Test_Reconstruct(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "正常なメールアドレスでインスタンス化できる",
			input:   "test@example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		email, err := Reconstruct(tt.input)

		require.NoError(t, err)
		assert.Equal(t, tt.input, string(*email))
	}
}
