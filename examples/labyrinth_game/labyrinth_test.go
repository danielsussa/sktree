package labyrinth

import (
	tree "github.com/danielsussa/tmp_tree"
	"testing"
)


func TestLabyrinthWithTrain(t *testing.T) {
	game := newGame()
	stateTree := tree.New()
	for {
		stateTree.Train(game, tree.StateTreeConfig{
			MaxDepth: 50,
		})
		res := stateTree.PlayTurn(game)
		if res.EndGame {
			break
		}
		game.Print()
	}
	game.Print()
}
