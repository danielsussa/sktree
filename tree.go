package tree

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"sort"
)

type StateTree struct {
	nodeMap map[string]*Node
	Root    *Node

	debugState   func(state State, debug Debug)
	debugActions func(actions []*Action, selected *Action)
	stats        *Stats
}

type Stats struct {
	NVisited int
}

type Node struct {
	state   State
	parent  *Node
	actions []*Action
	stats   *Stats
}

type Action struct {
	ID       interface{}
	score    float64
	nVisited int
}

func (a *Action) compute(res TurnResult) {
	a.nVisited++
	a.score += res.Score
}

type actionScore struct {
	action *Action
	score  float64
}

func (n *Node) selectAction() *Action {
	actionScoreList := make([]actionScore, 0)

	for _, action := range n.actions {
		actionScoreList = append(actionScoreList, actionScore{
			action: action,
			score:  fSelection(action.score, action.nVisited, n.stats.NVisited),
		})
	}
	sort.SliceStable(actionScoreList, func(i, j int) bool {
		if actionScoreList[i].score > actionScoreList[j].score {
			return true
		} else if actionScoreList[i].score < actionScoreList[j].score {
			return false
		} else {
			return actionScoreList[i].action.nVisited < actionScoreList[j].action.nVisited
		}
	})
	//n.selected = actionScoreList[0].action
	return actionScoreList[0].action
}

//func (n *Node) backPropagate(res TurnResult) {
//	n.selected.nVisited++
//	n.selected.score += res.Score
//	if n.parent == nil {
//		return
//	}
//	n.parent.backPropagate(res)
//}

func fSelection(total float64, nVisited, NVisited int) float64 {
	exploitation := total / float64(nVisited)
	exploration := math.Sqrt(2 * math.Log(float64(NVisited)) / float64(nVisited))
	sum := exploitation + exploration
	return math.Round(sum*100) / 100
}

// DebugState Mode
type Debug string

const (
	Bootstrap    Debug = "bootstrap"
	CurrentState Debug = "current_action"
)

func (st *StateTree) newNode(state State, parentNode *Node) *Node {
	stateId := newSHA256([]byte(state.ID()))
	if node, ok := st.nodeMap[stateId]; ok {
		node.state = state
		return node
	}

	actionList := make([]*Action, 0)
	for _, action := range state.PossibleActions() {
		actionList = append(actionList, &Action{
			ID:    action,
			score: 0,
		})
	}

	node := &Node{
		parent:  parentNode,
		state:   state,
		actions: actionList,
		stats:   st.stats,
	}

	st.nodeMap[stateId] = node
	return node
}

func newSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (st *StateTree) DebugState(f func(state State, debug Debug)) {
	st.debugState = f
}

func (st *StateTree) DebugAction(f func(actions []*Action, selected *Action)) {
	st.debugActions = f
}

func (st *StateTree) PlayGame(s State) {

	node := st.newNode(s.Copy(), nil)
	st.Root = node
	st.debugState(node.state, Bootstrap)

	actionList := make([]*Action, 0)

	for {
		currentAction := node.selectAction()
		st.debugActions(node.actions, currentAction)

		state := node.state.Copy()
		state = state.PlayAction(currentAction.ID)
		state = state.PlaySideEffects()

		result := state.TurnResult()
		if result.EndGame {
			break
		}

		// new node
		node = st.newNode(result.State.Copy(), node)
		st.debugState(node.state, CurrentState)

		actionList = append(actionList, currentAction)

		for _, action := range actionList {
			action.compute(result)
		}
	}
	// compute

}

type State interface {
	ID() string
	PossibleActions() []interface{}
	Copy() State
	PlayAction(interface{}) State
	PlaySideEffects() State
	TurnResult() TurnResult
}

type TurnResult struct {
	State   State
	Score   float64
	EndGame bool
}

func (a Action) GetNVisited() int {
	return a.nVisited
}

func New() *StateTree {
	return &StateTree{
		nodeMap: map[string]*Node{},
		debugState: func(state State, debug Debug) {
			// default
		},
		debugActions: func(actions []*Action, selected *Action) {

		},
		stats: &Stats{NVisited: 0},
	}
}
