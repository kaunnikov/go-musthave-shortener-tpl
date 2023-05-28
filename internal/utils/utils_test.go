package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)
import "golang.org/x/exp/slices"

// Тестируем создание 10000 уникальных алиасов сокращения URL
func TestRandSeq(t *testing.T) {
	var alias string
	var aliases = make([]string, 10000)
	for i := 0; i <= 10000; i++ {
		alias = RandSeq(5)
		assert.False(t, slices.Contains(aliases, alias))
		aliases = append(aliases, alias)
	}
}
