package labyrinth

import (
	"fmt"
	tree "github.com/danielsussa/tmp_tree"
	"testing"
)

func TestLabyrinth(t *testing.T) {
	labGame := newGame()

	treeGame := tree.New()
	treeGame.DebugState(func(state tree.State, debug tree.Debug) {
		//state.(game).Print()
		//fmt.Println("turn: ", state.(game).TotalMoves)
	})

	treeGame.DebugAction(func(actions []*tree.Action, selected *tree.Action) {
		//actList := make([]string, 0)
		//for _,act := range actions {
		//	actList = append(actList, act.ID.(string))
		//}
		//fmt.Println(fmt.Sprintf("[%d]actions: %v -> %s",selected.GetNVisited(),actList, selected.ID))
	})

	treeGame.Controller(func(req tree.ControllerRequest) tree.ControllerResponse {
		game := req.State.(game)
		if game.HasKey {
			//game.Print()
			//fmt.Println("turn: ", game.TotalMoves)
			if game.TotalMoves < 25 {
				game.Print()
				fmt.Println("turn: ", game.TotalMoves)
				return tree.ControllerResponse{Restart: true}
			}
		}
		return tree.ControllerResponse{Restart: true}
	})

	treeGame.PlayGame(labGame)

}
