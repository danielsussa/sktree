package g2048

import (
	"fmt"
	tree "github.com/danielsussa/tmp_tree"
	"math/rand"
	"testing"
)

func TestTrain2048(t *testing.T) {
	//defaultDb := defaultdb.NewBadgerDB("/media/kanczuk/146D-1AFD/dataset2/game2048")

	stateTree := tree.New()
	fmt.Println("starting")
	stateTree.DebugState(func(node tree.NodeDebug, debug tree.Debug) {
		//game := node.State.(*g2048)
		//print2048(game.board, game.score)
	})

	stateTree.DebugAction(func(actions []*tree.Action, selected *tree.Action) {
		//actList := make([]string, 0)
		//for _, act := range actions {
		//	actList = append(actList, act.ID)
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

func Test2048First(t *testing.T) {
	rand.Seed(1)
	game := startNewGame()
	game.board = []int{
		0, 0, 0, 0,
		0, 0, 2, 0,
		0, 0, 0, 0,
		16, 32, 64, 128,
	}
	print2048(game.board, game.score)
	totalPlays := 0
	for i := 0; i < 15; i++ {
		stateTree := tree.New()
		stateTree.Train(game, tree.StateTreeConfig{
			MaxIterations: 30000,
		})

		endGame := stateTree.PlayTurn(game)
		if endGame {
			break
		}

		game.PlaySideEffects()

		if game.topFirst() == 256 {
			break
		}

		print2048(game.board, game.score)
		totalPlays++
	}
	print2048(game.board, game.score)
	fmt.Println(totalPlays)
}

// /media/kanczuk/146D-1AFD/dataset/game2048
// /home/kanczuk/.tmp/game2048
func TestPlay2048(t *testing.T) {
	//defaultDb := defaultdb.NewBadgerDB("/media/kanczuk/146D-1AFD/dataset/badger/game2048")
	//defaultDb := defaultdb.NewDefaultDiskDB("/media/kanczuk/mydataset/game2048")
	//defaultDb := defaultdb.NewSqlDB("/media/kanczuk/datasetntfs/data/data.db")

	rand.Seed(1)
	game := startNewGame()

	print2048(game.board, game.score)
	//stateTree.SetDB(defaultDb)
	for {
		stateTree := tree.New()
		stateTree.Train(game, tree.StateTreeConfig{
			MaxIterations: 2048,
		})

		endGame := stateTree.PlayTurn(game)
		if endGame {
			break
		}

		game.PlaySideEffects()

		if len(game.PossibleActions()) == 0 {
			break
		}

		print2048(game.board, game.score)
	}
	print2048(game.board, game.score)
}
