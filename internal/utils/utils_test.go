package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMD5Hash(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "empty string",
			text: "",
			want: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name: "alphanumeric string",
			text: "abc123",
			want: "6ca13d52ca70c883e0f0bb101e425a89e8624de51db2d2392593af6a84118090",
		},
		{
			name: "special characters",
			text: "#$!",
			want: "ba62056c842406d4ec2c34b43127806fbb8a5ebc0fcc04b6c43db75926e92053",
		},
		{
			name: "long string",
			text: "ThisIsALongStringWithSomeDifferentCharacters!!",
			want: "d876841623daa8467309f4d4f101febce67f25f92d146704f15eb5e3da729323",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetMD5Hash(tt.text))
		})
	}
}
