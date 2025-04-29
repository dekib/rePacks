package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculatePacks(t *testing.T) {
	svc := &PackService{}

	tests := []struct {
		name     string
		items    int
		expected []PackResult
	}{
		{
			name:  "1 item with default sizes",
			items: 1,
			expected: []PackResult{
				{Pack: 250, Quantity: 1},
			},
		},
		{
			name:  "250 items with default sizes",
			items: 250,
			expected: []PackResult{
				{Pack: 250, Quantity: 1},
			},
		},
		{
			name:  "251 items with default sizes",
			items: 251,
			expected: []PackResult{
				{Pack: 500, Quantity: 1},
			},
		},
		{
			name:  "501 items with default sizes",
			items: 501,
			expected: []PackResult{
				{Pack: 250, Quantity: 1},
				{Pack: 500, Quantity: 1},
			},
		},
		{
			name:  "12001 items with default sizes",
			items: 12001,
			expected: []PackResult{
				{Pack: 250, Quantity: 1},
				{Pack: 2000, Quantity: 1},
				{Pack: 5000, Quantity: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.CalculatePacks(tt.items)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCustomPackSizes(t *testing.T) {
	svc := &PackService{}
	svc.SetPackSizes([]int{23, 31, 53})

	t.Run("500000 items with custom sizes", func(t *testing.T) {
		expected := []PackResult{
			{Pack: 23, Quantity: 2},
			{Pack: 31, Quantity: 7},
			{Pack: 53, Quantity: 9429},
		}
		result := svc.CalculatePacks(500000)
		assert.Equal(t, expected, result)
	})
}
