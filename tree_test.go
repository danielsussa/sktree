package tree

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMCTS(t *testing.T) {
	assert.Equal(t, 0.79, fSelection(-1, 1, 5))
	assert.Equal(t, 2.79, fSelection(1, 1, 5))
	assert.Equal(t, 1.48, fSelection(0, 2, 9))
	assert.Equal(t, 2.48, fSelection(2, 2, 9))
	assert.Equal(t, 0.89, fSelection(-1, 1, 6))
	assert.Equal(t, 2.52, fSelection(2, 2, 10))
}

func TestMCT2(t *testing.T) {
	fmt.Println(fSelection(0.0, 5, 1))
	fmt.Println(fSelection(0.0, 10, 1))
}

func TestTmp(t *testing.T) {
	nodeMap := make(map[int]Node, 0)
	for i := 0; i < 1000000*20; i++ {
		actionList := make([]*Action, 10)
		for k := 0; k < 10; k++ {
			actionList[k] = &Action{
				ID:       "hellow wkw",
				Score:    313141,
				Turn:     "d",
				deadEnd:  false,
				NVisited: 455,
			}
			//actionList = append(actionList, &Action{
			//	ID:       "hellow wkw",
			//	Score:    313141,
			//	Turn:     "d",
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
				Turn:     "d",
				deadEnd:  false,
				NVisited: 455,
			}
			//actionList = append(actionList, &Action{
			//	ID:       "hellow wkw",
			//	Score:    313141,
			//	Turn:     "d",
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
