package utils_test

import (
	"testing"

	"github.com/informatik-mannheim/cmg-ss2025/services/job-scheduler/utils"
)

func TestIsUrlValid(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"ftp://example.com", true},
		{"invalid-url", false},
		{"", false},
	}

	for _, test := range tests {
		result := utils.IsUrlValid(test.url)
		if result != test.expected {
			t.Errorf("IsUrlValid(%s) = %v; expected %v", test.url, result, test.expected)
		}
	}
}
