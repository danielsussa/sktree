package g2048

import (
	"fmt"
	tree "github.com/danielsussa/tmp_tree"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func Test2048First(t *testing.T) {
	rand.Seed(3)
	game := startNewGame()
	game.board = []int{
		0, 2, 0, 0,
		0, 0, 0, 0,
		2, 4, 0, 0,
		64, 128, 256, 0,
	}
	print2048(game.board, game.score)
	totalPlays := 0
	stateTree := tree.New()
	for i := 0; i < 30; i++ {
		res := stateTree.Train(game, tree.StateTreeConfig{
			MaxDepth: 10,
		})
		fmt.Println("new nodes: ", res.TotalNewNodes)

		result := stateTree.PlayTurn(game)
		if result.EndGame {
			break
		}

		game.PlaySideEffects()

		if game.topFirst() == 512 {
			break
		}

		print2048(game.board, game.score)
		totalPlays++
	}
	print2048(game.board, game.score)
	fmt.Println(totalPlays)
}

func Test2048BestResult(t *testing.T) {
	rand.Seed(3)
	game := startNewGame()
	game.board = []int{
		0, 2, 0, 0,
		0, 16, 0, 0,
		16, 0, 0, 0,
		32, 64, 128, 256,
	}

	stateTree := tree.New()
	res := stateTree.Train(game, tree.StateTreeConfig{
		MaxDepth: 10,
	})
	fmt.Println("new nodes: ", res.TotalNewNodes)

	result := stateTree.PlayTurn(game)
	assert.Equal(t, result.Action.ID, "L")
}

func Test2048GoToRight(t *testing.T) {
	rand.Seed(3)
	game := startNewGame()
	game.board = []int{
		0, 2, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		64, 128, 256, 0,
	}

	stateTree := tree.New()
	res := stateTree.Train(game, tree.StateTreeConfig{
		MaxDepth: 30,
	})
	fmt.Println("new nodes: ", res.TotalNewNodes)

	result := stateTree.PlayTurn(game)
	assert.Equal(t, result.Action.ID, "R")
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
			MaxDepth:   15,
		})

		result := stateTree.PlayTurn(game)
		if result.EndGame {
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
