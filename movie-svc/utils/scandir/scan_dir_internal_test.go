package scandir

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasSuffix(t *testing.T) {
	testCases := []struct {
		name      string
		str       string
		arr       []string
		hasSuffix bool
	}{
		{
			name:      "Contains",
			str:       "test.mp4",
			arr:       []string{".avi", ".mp4"},
			hasSuffix: true,
		},
		{
			name:      "Doesn't contain",
			str:       "test.text",
			arr:       []string{".mp4", ".avi"},
			hasSuffix: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			has := hasSuffix(tt.str, tt.arr)
			assert.Equal(t, tt.hasSuffix, has)
		})
	}
}
