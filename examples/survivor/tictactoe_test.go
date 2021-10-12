package tictactoe

import (
	tree "github.com/danielsussa/tmp_tree"
	"testing"
)

func TestTicTac(t *testing.T) {
	game := ticTacGame{}

	gameTree := tree.New(game)

	for {
		result := gameTree.PlayGame(game)
		result.(ticTacGame).print()
	}
}
