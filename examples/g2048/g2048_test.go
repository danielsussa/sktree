package g2048

import (
	"fmt"
	tree "github.com/danielsussa/sktree"
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
	print2048(game)
	totalPlays := 0

	for i := 0; i < 60; i++ {
		stateTree := tree.New()
		trainRes := stateTree.Train(game, tree.StateTreeConfig{
			MaxIterations: 4096,
			ScoreNormalizer: &tree.ScoreNormalizer{
				Min: 0,
				Max: 1024,
			},
		})

		result := stateTree.PlayTurn(game)
		if result.EndGame {
			break
		}

		game.playSideEffects()

		if game.topFirst() == 512 {
			break
		}

		print2048WithRes(game, trainRes, result)
		totalPlays++
	}
	fmt.Println("total plays: ", totalPlays)
	print2048(game)
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
	_ = stateTree.Train(game, tree.StateTreeConfig{
		MaxDepth: 15,
		ScoreNormalizer: &tree.ScoreNormalizer{
			Min: 0,
			Max: 600,
		},
	})

	result := stateTree.PlayTurn(game)
	assert.Equal(t, "L", result.Action.ID)
}

func Test2048GoToRight(t *testing.T) {
	rand.Seed(3)
	game := startNewGame()
	game.board = []int{
		64, 0, 64, 0,
		0, 0, 128, 0,
		0, 0, 1024, 0,
		256, 512, 0, 0,
	}

	stateTree := tree.New()
	stateTree.Controller(func(req tree.ControllerRequest) tree.ControllerResponse {
		if topFirstValue(req.State.(*g2048).board) == 2048 {
			return tree.ControllerResponse{ForceStop: true}
		}
		return tree.ControllerResponse{ForceStop: false}
	})
	res := stateTree.Train(game, tree.StateTreeConfig{
		MaxDepth: 10,
		ScoreNormalizer: &tree.ScoreNormalizer{
			Min: 0,
			Max: 2048,
		},
	})
	fmt.Println(res.TotalNodes)

	result := stateTree.PlayTurn(game)
	assert.Equal(t, "D", result.Action.ID)
}

// /media/kanczuk/146D-1AFD/dataset/game2048
// /home/kanczuk/.tmp/game2048
func TestPlay2048(t *testing.T) {

	rand.Seed(1)
	game := startNewGame()

	print2048(game)
	//stateTree.SetDB(defaultDb)
	for {
		stateTree := tree.New()
		stateTree.Train(game, tree.StateTreeConfig{
			MaxIterations: 4096 * 32,
		})

		result := stateTree.PlayTurn(game)
		if result.EndGame {
			break
		}

		game.playSideEffects()

		if len(game.PossibleActions()) == 0 {
			break
		}

		print2048(game)
	}
	print2048(game)
}
