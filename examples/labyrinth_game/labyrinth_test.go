package labyrinth

import (
	"encoding/json"
	"fmt"
	tree "github.com/danielsussa/tmp_tree"
	"io/ioutil"
	"testing"
)

type DefaultMemoryDB struct {
	nodeMap map[string]*tree.Node
}

func (dmp DefaultMemoryDB) Find(key string) (*tree.Node, bool) {
	if node, ok := dmp.nodeMap[key]; ok {
		return node, true
	}
	return nil, false
}

func (dmp DefaultMemoryDB) Add(key string, node *tree.Node) error {
	dmp.nodeMap[key] = node
	return nil
}

func (dmp DefaultMemoryDB) purgeToDisk() error {
	file, _ := json.MarshalIndent(dmp.nodeMap, "", " ")
	_ = ioutil.WriteFile("dataset.json", file, 0644)
	return nil
}

func newDefaultDB() DefaultMemoryDB {
	return DefaultMemoryDB{nodeMap: map[string]*tree.Node{}}
}

func TestLabyrinth(t *testing.T) {
	defaultDb := newDefaultDB()
	labGame := newGame()

	treeGame := tree.New().SetDB(defaultDb)
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

	totalWins := 0

	treeGame.Controller(func(req tree.ControllerRequest) tree.ControllerResponse {
		game := req.State.(game)
		if game.WinGame {
			totalWins++
			game.Print()
			fmt.Println("turn: ", game.TotalMoves)
			_ = defaultDb.purgeToDisk()
			if totalWins > 10 {
				return tree.ControllerResponse{Restart: false}
			}
			return tree.ControllerResponse{Restart: true}

		}
		return tree.ControllerResponse{Restart: true}
	})

	treeGame.PlayGame(labGame)

}
