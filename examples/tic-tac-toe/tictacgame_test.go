package tictactoe

import (
	"fmt"
	tree "github.com/danielsussa/sktree"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestGameScore(t *testing.T) {
	{
		game := &ticTacGame{
			CurrentPlayer: H,
			MainPlayer:    H,
			Board: []player{
				M, M, H,
				H, M, M,
				H, E, H,
			},
		}
		assert.Equal(t, E, game.winner())
	}
}

func TestFirstMove(t *testing.T) {
	game := &ticTacGame{
		CurrentPlayer: H,
		MainPlayer:    H,
		Board: []player{
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
		MaxIterations: 1000,
	})
	stateTree.PlayTurn(game)

	assert.Equal(t, expected, game.Board)
}

func TestFirstMachineMove(t *testing.T) {
	rand.Seed(11)
	game := &ticTacGame{
		MainPlayer:    M,
		CurrentPlayer: M,
		Board: []player{
			E, E, E,
			E, E, E,
			E, E, E,
		},
	}

	expected := []player{
		E, E, E,
		E, M, E,
		E, E, E,
	}

	stateTree := tree.New()
	stateTree.Train(game, tree.StateTreeConfig{
		MaxIterations: 1000,
	})
	stateTree.PlayTurn(game)

	assert.Equal(t, expected, game.Board)
}

func TestBestSecondMove(t *testing.T) {
	rand.Seed(2)
	game := &ticTacGame{
		CurrentPlayer: M,
		MainPlayer:    M,
		Board: []player{
			H, E, E,
			E, E, E,
			E, E, E,
		},
	}

	expected := []player{
		H, E, E,
		E, M, E,
		E, E, E,
	}

	stateTree := tree.New()
	stateTree.Train(game, tree.StateTreeConfig{
		MaxIterations:    1000,
		TotalSimulations: 1,
	})
	stateTree.PlayTurn(game)

	assert.Equal(t, expected, game.Board)
}

func TestDontLoseMoveHuman(t *testing.T) {
	game := &ticTacGame{
		CurrentPlayer: H,
		MainPlayer:    H,
		Board: []player{
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

	assert.Equal(t, expected, game.Board)
}

func TestDontLoseMoveHuman2(t *testing.T) {
	rand.Seed(12)
	game := &ticTacGame{
		CurrentPlayer: H,
		MainPlayer:    H,
		Board: []player{
			M, E, H,
			E, M, M,
			E, E, H,
		},
	}

	expected := []player{
		M, E, H,
		H, M, M,
		E, E, H,
	}

	stateTree := tree.New()
	res := stateTree.Train(game, tree.StateTreeConfig{
		MaxIterations: 1000,
	})
	fmt.Println(res.TotalNodes)
	stateTree.PlayTurn(game)

	assert.Equal(t, expected, game.Board)
}

func TestMachineWonMovement(t *testing.T) {
	rand.Seed(12)
	game := &ticTacGame{
		CurrentPlayer: M,
		MainPlayer:    M,
		Board: []player{
			M, E, H,
			E, M, M,
			E, H, H,
		},
	}

	expected := []player{
		M, E, H,
		M, M, M,
		E, H, H,
	}

	stateTree := tree.New()
	res := stateTree.Train(game, tree.StateTreeConfig{
		MaxIterations: 1000,
	})
	fmt.Println(res.TotalNodes)
	stateTree.PlayTurn(game)

	assert.Equal(t, expected, game.Board)
}

func TestDontLoseMachineMove(t *testing.T) {
	game := &ticTacGame{
		CurrentPlayer: M,
		MainPlayer:    M,
		Board: []player{
			H, E, M,
			E, H, H,
			E, E, M,
		},
	}

	expected := []player{
		H, E, M,
		M, H, H,
		E, E, M,
	}

	stateTree := tree.New()
	res := stateTree.Train(game, tree.StateTreeConfig{
		MaxIterations: 1000,
	})
	fmt.Println(res.TotalNodes)
	stateTree.PlayTurn(game)

	assert.Equal(t, expected, game.Board)
}

func TestTotalIterations(t *testing.T) {
	game := &ticTacGame{
		CurrentPlayer: M,
		Board: []player{
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

	assert.Equal(t, 5, res.TotalNodes)
}
