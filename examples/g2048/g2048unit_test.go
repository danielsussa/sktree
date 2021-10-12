package g2048

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// du -hs g2048

//
func TestMapConverter(t *testing.T) {
	{
		board := []int{
			0, 0, 0, 4,
			0, 0, 0, 32,
			0, 0, 0, 64,
			0, 0, 0, 64,
		}
		expected := []int{
			0, 0, 0, 1,
			0, 0, 0, 2,
			0, 0, 0, 3,
			0, 0, 0, 3,
		}
		assert.Equal(t, expected, convertScalar(board))
	}
	{
		board := []int{
			0, 0, 0, 0,
			2, 0, 0, 0,
			0, 0, 0, 0,
			4, 8, 128, 128,
		}
		expected := []int{
			0, 0, 0, 0,
			1, 0, 0, 0,
			0, 0, 0, 0,
			2, 3, 4, 4,
		}
		assert.Equal(t, expected, convertScalar(board))
	}
	{
		board := []int{
			0, 0, 0, 0,
			2, 0, 0, 0,
			0, 0, 0, 0,
			8, 16, 128, 128,
		}
		expected := []int{
			0, 0, 0, 0,
			1, 0, 0, 0,
			0, 0, 0, 0,
			2, 3, 4, 4,
		}
		assert.Equal(t, expected, convertScalar(board))
	}
	{
		board := []int{
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		expected := []int{
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, expected, convertScalar(board))
	}
}

func TestCanMove(t *testing.T) {
	{
		board := []int{
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, false, canMoveUp(board))
	}
	{
		board := []int{
			4, 0, 4, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, false, canMoveUp(board))
	}
	{
		board := []int{
			4, 0, 4, 0,
			0, 0, 8, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, false, canMoveUp(board))
	}
	{
		board := []int{
			4, 0, 4, 0,
			4, 0, 8, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, true, canMoveUp(board))
	}
	{
		board := []int{
			4, 0, 4, 0,
			2, 0, 8, 0,
			4, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, false, canMoveUp(board))
	}
	// move down
	{
		board := []int{
			4, 0, 4, 0,
			2, 0, 8, 0,
			4, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, true, canMoveDown(board))
	}
	{
		board := []int{
			4, 0, 4, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, true, canMoveDown(board))
	}
	{
		board := []int{
			4, 0, 0, 0,
			16, 0, 0, 0,
			32, 0, 0, 0,
			64, 0, 0, 0,
		}
		assert.Equal(t, false, canMoveDown(board))
	}
	{
		board := []int{
			0, 0, 0, 4,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, false, canMoveRight(board))
	}
	{
		board := []int{
			0, 0, 2, 4,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, false, canMoveRight(board))
	}
	{
		board := []int{
			0, 4, 2, 4,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, true, canMoveLeft(board))
		assert.Equal(t, false, canMoveRight(board))
	}
	{
		board := []int{
			4, 4, 2, 4,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, true, canMoveLeft(board))
		assert.Equal(t, true, canMoveRight(board))
	}
	{
		board := []int{
			2, 4, 2, 4,
			4, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, false, canMoveLeft(board))
		assert.Equal(t, true, canMoveRight(board))
	}
	{
		board := []int{
			2, 4, 2, 4,
			4, 8, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, false, canMoveLeft(board))
		assert.Equal(t, true, canMoveRight(board))
	}
	{
		board := []int{
			2, 4, 2, 4,
			4, 8, 16, 32,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, false, canMoveLeft(board))
		assert.Equal(t, false, canMoveRight(board))
	}
}

func TestComputeMoves(t *testing.T) {
	{
		board := []int{
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		expected := []int{
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, 0, computeUp(board))
		assert.Equal(t, expected, board)
	}
	{
		board := []int{
			0, 0, 0, 0,
			0, 4, 4, 8,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		expected := []int{
			0, 4, 4, 8,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, 0, computeUp(board))
		assert.Equal(t, expected, board)
	}
	{
		board := []int{
			0, 0, 0, 0,
			0, 4, 4, 8,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		expected := []int{
			0, 0, 0, 0,
			0, 0, 8, 8,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
		assert.Equal(t, 8, computeRight(board))
		assert.Equal(t, expected, board)
	}
	{
		board := []int{
			0, 0, 0, 0,
			4, 4, 4, 8,
			0, 8, 0, 8,
			0, 16, 32, 16,
		}
		expected := []int{
			0, 0, 0, 0,
			0, 4, 8, 8,
			0, 0, 0, 16,
			0, 16, 32, 16,
		}
		assert.Equal(t, 24, computeRight(board))
		assert.Equal(t, expected, board)
	}
}

func TestMerge(t *testing.T) {
	{
		curr := []int{2, 0, 2, 4}
		assert.Equal(t, true, canMerge(curr))
		expected := []int{4, 4, 0, 0}
		score := merge(curr)
		assert.Equal(t, expected, curr)
		assert.Equal(t, score, 4)
	}
	{
		curr := []int{2, 2, 2, 4}
		assert.Equal(t, true, canMerge(curr))
		expected := []int{4, 2, 4, 0}
		score := merge(curr)
		assert.Equal(t, expected, curr)
		assert.Equal(t, score, 4)
	}
	{
		curr := []int{2, 2, 2, 2}
		assert.Equal(t, true, canMerge(curr))
		expected := []int{4, 4, 0, 0}
		score := merge(curr)
		assert.Equal(t, expected, curr)
		assert.Equal(t, score, 8)
	}
	{
		curr := []int{2, 0, 2, 0}
		assert.Equal(t, true, canMerge(curr))
		expected := []int{4, 0, 0, 0}
		score := merge(curr)
		assert.Equal(t, expected, curr)
		assert.Equal(t, score, 4)
	}
	{
		curr := []int{2, 2, 2, 0}
		assert.Equal(t, true, canMerge(curr))
		expected := []int{4, 2, 0, 0}
		score := merge(curr)
		assert.Equal(t, expected, curr)
		assert.Equal(t, score, 4)
	}
	{
		curr := []int{4, 0, 0, 4}
		assert.Equal(t, true, canMerge(curr))
		expected := []int{8, 0, 0, 0}
		score := merge(curr)
		assert.Equal(t, expected, curr)
		assert.Equal(t, score, 8)
	}
	{
		curr := []int{2, 4, 8, 16}
		assert.Equal(t, false, canMerge(curr))
		expected := []int{2, 4, 8, 16}
		score := merge(curr)
		assert.Equal(t, expected, curr)
		assert.Equal(t, score, 0)
	}
	{
		curr := []int{4, 2, 4, 0}
		assert.Equal(t, false, canMerge(curr))
		expected := []int{4, 2, 4, 0}
		score := merge(curr)
		assert.Equal(t, expected, curr)
		assert.Equal(t, score, 0)
	}
}
