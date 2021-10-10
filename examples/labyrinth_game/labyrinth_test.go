package labyrinth

import (
	"fmt"
	tree "github.com/danielsussa/tmp_tree"
	"testing"
)

func TestLabyrinth(t *testing.T) {
	labGame := newGame()

	treeGame := tree.New()
	treeGame.DebugState(func(debug tree.NodeDebug, kind tree.Debug) {
		//game := debug.State.(game)
		//game.Print()
		//
		//fmt.Println("turn: ", game.TotalMoves)
	})

	treeGame.DebugAction(func(actions []*tree.Action, selected *tree.Action) {
		//actList := make([]string, 0)
		//for _, act := range actions {
		//	actList = append(actList, act.ID.(string))
		//}
		//fmt.Println(fmt.Sprintf("[%d]actions: %v -> %s", selected.GetNVisited(), actList, selected.ID))
	})

	totalWins := 0

	treeGame.Controller(func(req tree.ControllerRequest) tree.ControllerResponse {
		game := req.State.(game)
		//game.Print()
		//fmt.Println("turn: ", game.TotalMoves)
		if game.WinGame {
			totalWins++
			game.Print()
			fmt.Println("turn: ", game.TotalMoves)
			if totalWins > 10 {
				return tree.ControllerResponse{Restart: false}
			}
			return tree.ControllerResponse{Restart: true}

		}
		return tree.ControllerResponse{Restart: true}
	})

	treeGame.PlayGame(labGame)

}
