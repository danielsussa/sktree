package labyrinth

import (
	tree "github.com/danielsussa/tmp_tree"
	"testing"
)


func TestLabyrinthWithTrain(t *testing.T) {
	game := newGame()
	stateTree := tree.New()
	stateTree.Controller(func(req tree.ControllerRequest) tree.ControllerResponse {
		game := req.State.(*labyrinth)

		return tree.ControllerResponse{
			ForceStop: game.WinGame,
		}
	})
	for {
		stateTree.Train(game, tree.StateTreeConfig{
			MaxDepth: 500,
		})
		res := stateTree.PlayTurn(game)
		if res.EndGame {
			break
		}
		game.Print()
	}
	game.Print()
}
