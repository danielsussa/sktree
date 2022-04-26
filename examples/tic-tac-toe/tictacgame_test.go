package tictactoe

import (
	tree "github.com/danielsussa/sktree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFirstSecondMoveMove(t *testing.T) {
	game := &ticTacGame{
		currentPlayer: H,
		board: []player{
			E, E, E,
			E, E, E,
			E, E, E,
		},
	}

	expected := []player{
		E, E, E,
		E, H, E,
		E, E, E,
	}

	stateTree := tree.New()
	stateTree.Train(game, tree.StateTreeConfig{
		MaxIterations: 4096,
	})
	stateTree.PlayTurn(game)

	assert.Equal(t, expected, game.board)
}

func TestBestSecondMoveMove(t *testing.T) {
	game := &ticTacGame{
		currentPlayer: H,
		board: []player{
			M, E, E,
			E, E, E,
			E, E, E,
		},
	}

	expected := []player{
		M, E, E,
		E, H, E,
		E, E, E,
	}

	stateTree := tree.New()
	stateTree.Train(game, tree.StateTreeConfig{
		MaxIterations: 4096,
	})
	stateTree.PlayTurn(game)

	assert.Equal(t, expected, game.board)
}

func TestDontLoseMoveMove(t *testing.T) {
	game := &ticTacGame{
		currentPlayer: H,
		board: []player{
			M, E, M,
			E, H, E,
			E, E, E,
		},
	}

	expected := []player{
		M, H, M,
		E, H, E,
		E, E, E,
	}

	stateTree := tree.New()
	stateTree.Train(game, tree.StateTreeConfig{
		MaxIterations: 4096,
	})
	stateTree.PlayTurn(game)

	assert.Equal(t, expected, game.board)
}

func TestTotalIterations(t *testing.T) {
	game := &ticTacGame{
		currentPlayer: M,
		board: []player{
			H, M, E,
			M, M, H,
			H, H, E,
		},
	}

	stateTree := tree.New()
	res := stateTree.Train(game, tree.StateTreeConfig{
		MaxIterations: 4096,
	})
	stateTree.PlayTurn(game)

	assert.Equal(t, 3, res.TotalNodes)
}
