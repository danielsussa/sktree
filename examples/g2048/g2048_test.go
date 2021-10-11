package g2048

import (
	tree "github.com/danielsussa/tmp_tree"
	"github.com/danielsussa/tmp_tree/examples/defaultdb"
	"github.com/stretchr/testify/assert"
	"testing"
)

// du -hs g2048
func TestTrain2048(t *testing.T) {
	defaultDb := defaultdb.NewDefaultDiskDB("/media/kanczuk/Seagate Expansion Drive/dataset/game2048_mod")

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

func TestMapConverter(t *testing.T) {
	{
		board := [][]int{
			{0, 0, 0, 4},
			{0, 0, 0, 32},
			{0, 0, 0, 64},
			{0, 0, 0, 64},
		}
		expected := [][]int{
			{0, 0, 0, 1},
			{0, 0, 0, 2},
			{0, 0, 0, 3},
			{0, 0, 0, 3},
		}
		assert.Equal(t, expected, convertScalar(board))
	}
	{
		board := [][]int{
			{0, 0, 0, 0},
			{2, 0, 0, 0},
			{0, 0, 0, 0},
			{4, 8, 128, 128},
		}
		expected := [][]int{
			{0, 0, 0, 0},
			{1, 0, 0, 0},
			{0, 0, 0, 0},
			{2, 3, 4, 4},
		}
		assert.Equal(t, expected, convertScalar(board))
	}
	{
		board := [][]int{
			{0, 0, 0, 0},
			{2, 0, 0, 0},
			{0, 0, 0, 0},
			{8, 16, 128, 128},
		}
		expected := [][]int{
			{0, 0, 0, 0},
			{1, 0, 0, 0},
			{0, 0, 0, 0},
			{2, 3, 4, 4},
		}
		assert.Equal(t, expected, convertScalar(board))
	}
	{
		board := [][]int{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		}
		expected := [][]int{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		}
		assert.Equal(t, expected, convertScalar(board))
	}
}
