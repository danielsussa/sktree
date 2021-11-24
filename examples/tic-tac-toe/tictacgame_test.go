package tictactoe

import (
	tree "github.com/danielsussa/tmp_tree"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFirstSecondMoveMove(t *testing.T) {
	game := ticTacGame{
		board: []player{
			E, E, E,
			E, E, E,
			E, E, E,
		},
	}

	expected := []player{
		E, E, E,
		E, X, E,
		E, E, E,
	}

	stateTree := tree.New()
	stateTree.Train(game, tree.StateTreeConfig{
		MaxDepth: 10,
	})
	stateTree.PlayTurn(game)

	assert.Equal(t, expected, game.board)
}

func TestBestSecondMoveMove(t *testing.T) {
	game := ticTacGame{
		board: []player{
			O, E, E,
			E, E, E,
			E, E, E,
		},
	}

	expected := []player{
		O, E, E,
		E, X, E,
		E, E, E,
	}

	stateTree := tree.New()
	stateTree.Train(game, tree.StateTreeConfig{
		MaxDepth: 10,
	})
	stateTree.PlayTurn(game)

	assert.Equal(t, expected, game.board)
}

func TestDontLoseMoveMove(t *testing.T) {
	game := ticTacGame{
		board: []player{
			O, E, O,
			E, X, E,
			E, E, E,
		},
	}

	expected := []player{
		O, X, O,
		E, X, E,
		E, E, E,
	}

	stateTree := tree.New()
	stateTree.Train(game, tree.StateTreeConfig{
		MaxDepth: 10,
	})
	stateTree.PlayTurn(game)

	assert.Equal(t, expected, game.board)
}
