package tree

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMCTS(t *testing.T) {
	assert.Equal(t, 0.79, fSelection(-1, 1, 5))
	assert.Equal(t, 2.79, fSelection(1, 1, 5))
	assert.Equal(t, 1.48, fSelection(0, 2, 9))
	assert.Equal(t, 2.48, fSelection(2, 2, 9))
	assert.Equal(t, 0.89, fSelection(-1, 1, 6))
	assert.Equal(t, 2.52, fSelection(2, 2, 10))
}

func TestMCT2(t *testing.T) {
	fmt.Println(fSelection(2, 997440, 1500000))
	fmt.Println(fSelection(0.0, 997440, 1500000))
}
