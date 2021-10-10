package g2048

import (
	tree "github.com/danielsussa/tmp_tree"
	"github.com/danielsussa/tmp_tree/examples/defaultdb"
	"testing"
)

// du -hs g2048
func TestTrain2048(t *testing.T) {
	defaultDb := defaultdb.NewDefaultDiskDB("/media/kanczuk/total64/dataset/game2048")

	stateTree := tree.New().SetDB(defaultDb)
	stateTree.DebugState(func(node tree.NodeDebug, debug tree.Debug) {
		//game := node.State.(g2048)
		//print2048(game.board, game.score)
	})

	stateTree.DebugAction(func(actions []*tree.Action, selected *tree.Action) {
		//actList := make([]string, 0)
		//for _, act := range actions {
		//	actList = append(actList, act.ID.(string))
		//}
		//fmt.Println(fmt.Sprintf("[%d]actions: %v", selected.GetNVisited(), selected.ID))
	})

	maxScore := 0

	stateTree.Controller(func(req tree.ControllerRequest) tree.ControllerResponse {
		game := req.State.(*g2048)
		if game.score > maxScore {
			maxScore = game.score
			print2048(game.board, maxScore)
		}
		return tree.ControllerResponse{Restart: true}
	})

	stateTree.PlayGame(startNewGame())
}
