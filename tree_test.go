package tree

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestMCTS(t *testing.T) {
	round := func(val float64) float64 {
		return math.Round(val*100) / 100
	}
	assert.Equal(t, 0.79, round(fSelection(-1, 1, 5)))
	assert.Equal(t, 2.79, round(fSelection(1, 1, 5)))
	assert.Equal(t, 1.48, round(fSelection(0, 2, 9)))
	assert.Equal(t, 2.48, round(fSelection(2, 2, 9)))
	assert.Equal(t, 0.89, round(fSelection(-1, 1, 6)))
	assert.Equal(t, 2.52, round(fSelection(2, 2, 10)))
	assert.Equal(t, 0.06, round(fSelection(-56, 543, 1000)))

	assert.Equal(t, 0.09, round(fSelection(-36, 281, 800)))
	assert.Equal(t, 0.08, round(fSelection(-25, 69, 800)))

	assert.Equal(t, 1.46, round(fSelection(15, 22, 800)))
	assert.Equal(t, 0.17, round(fSelection(15, 650, 800)))

	assert.Equal(t, 0.1, round(fSelection(-15, 22, 800)))
	assert.Equal(t, 0.12, round(fSelection(-15, 650, 800)))
	assert.Equal(t, 0.14, round(fSelection(-13, 404, 500)))

}

func TestMCT2(t *testing.T) {
	fmt.Println(fSelection(0.0, 5, 1))
	fmt.Println(fSelection(0.0, 10, 1))
}

func TestSelectAction(t *testing.T) {
	{
		n := Node{
			OpponentTurn: true,
			Actions: []*Action{
				{Score: -4, NVisited: 4, ID: "correct"},
				{Score: -3, NVisited: 7},
				{Score: 648, NVisited: 687},
				{Score: 2, NVisited: 15},
				{Score: 19, NVisited: 37},
			},
			Idx: 0,
		}
		action := n.selectAction(&rootStats{NVisited: 800})
		assert.Equal(t, "correct", action.ID)
	}
	{
		n := Node{
			OpponentTurn: true,
			Actions: []*Action{
				{Score: -1, NVisited: 1, ID: "correct"},
				{Score: 0, NVisited: 1},
				{Score: 1, NVisited: 1},
				{Score: 735, NVisited: 759, ID: "incorrect"},
				{Score: 0, NVisited: 1},
			},
			Idx: 0,
		}
		action := n.selectAction(&rootStats{NVisited: 800})
		assert.Equal(t, "correct", action.ID)
	}
}

func TestTmp(t *testing.T) {
	nodeMap := make(map[int]Node, 0)
	for i := 0; i < 1000000*20; i++ {
		actionList := make([]*Action, 10)
		for k := 0; k < 10; k++ {
			actionList[k] = &Action{
				ID:       "hellow wkw",
				Score:    313141,
				NVisited: 455,
			}
			//actionList = append(actionList, &Action{
			//	ID:       "hellow wkw",
			//	Score:    313141,
			//	IsOpponentTurn:     "d",
			//	deadEnd:  false,
			//	NVisited: 455,
			//})
		}
		nodeMap[i] = Node{
			Actions: actionList,
		}
	}
	fmt.Println(len(nodeMap))
}

func TestTmp2(t *testing.T) {
	nodeMap := make([]*Node, 1000000*1000)
	for i := 0; i < 1000000*10; i++ {
		actionList := make([]*Action, 10)
		for k := 0; k < 10; k++ {
			actionList[k] = &Action{
				ID:       "hellow wkw",
				Score:    313141,
				NVisited: 455,
			}
			//actionList = append(actionList, &Action{
			//	ID:       "hellow wkw",
			//	Score:    313141,
			//	IsOpponentTurn:     "d",
			//	deadEnd:  false,
			//	NVisited: 455,
			//})
		}
		nodeMap[i] = &Node{
			Actions: actionList,
		}
		//nodeMap = append(nodeMap, &Node{
		//	Actions: actionList,
		//})
	}
	fmt.Println(len(nodeMap))
}
