package tree

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

type Database interface {
	Find(string) (string, bool)
	Add(string, string) error
}

type DefaultMemoryDB struct {
	nodeMap map[string]string
}

func (dmp DefaultMemoryDB) Find(key string) (string, bool) {
	if node, ok := dmp.nodeMap[key]; ok {
		return node, true
	}
	return "", false
}

func (dmp DefaultMemoryDB) Add(key, val string) error {
	dmp.nodeMap[key] = val
	return nil
}

type StateTree struct {
	debugState   func(node NodeDebug, debug Debug)
	debugActions func(actions []*Action, selected *Action)
	controller   func(req ControllerRequest) ControllerResponse
	stats        *rootStats
	db           Database
}

type rootStats struct {
	NVisited int
}

type Node struct {
	Actions []*Action
	id      string
}

func (n Node) toDB() string {
	var b strings.Builder
	for _, act := range n.Actions {
		b.WriteString(fmt.Sprintf("%s;%d;%d;", act.ID, act.NVisited, act.Score))
	}
	return strings.TrimRight(b.String(), ";")
}

type NodeDebug struct {
	Id    string
	State State
}

type Action struct {
	ID       interface{}
	Score    int
	NVisited int
}

type actionScore struct {
	action *Action
	score  float64
}

func (n *Node) selectAction(stats *rootStats) *Action {
	actionScoreList := make([]actionScore, 0)

	for _, action := range n.Actions {
		actionScoreList = append(actionScoreList, actionScore{
			action: action,
			score:  fSelection(float64(action.Score), action.NVisited, stats.NVisited),
		})
	}
	sort.SliceStable(actionScoreList, func(i, j int) bool {
		if actionScoreList[i].score > actionScoreList[j].score {
			return true
		} else if actionScoreList[i].score < actionScoreList[j].score {
			return false
		} else {
			return actionScoreList[i].action.NVisited < actionScoreList[j].action.NVisited
		}
	})
	//n.selected = actionScoreList[0].action
	return actionScoreList[0].action
}

func fSelection(total float64, nVisited, NVisited int) float64 {
	exploitation := total / float64(nVisited)
	exploration := math.Sqrt(2 * math.Log(float64(NVisited)) / float64(nVisited))
	sum := exploitation + exploration
	return (sum * 100) / 100
}

// DebugState Mode
type Debug string

const (
	Bootstrap    Debug = "bootstrap"
	CurrentState Debug = "current_action"
)

func parseToNode(key, val string) *Node {
	valSpl := strings.Split(val, ";")

	actions := make([]*Action, 0)
	for i := 0; i < len(valSpl); i += 3 {
		id := valSpl[i]
		nVisited, _ := strconv.Atoi(valSpl[i+1])
		score, _ := strconv.Atoi(valSpl[i+2])
		actions = append(actions, &Action{
			ID:       id,
			Score:    score,
			NVisited: nVisited,
		})
	}

	return &Node{
		Actions: actions,
		id:      key,
	}
}

func (st *StateTree) newNode(state State, nodeMap map[string]*Node) *Node {
	stateId := state.ID()
	if val, ok := nodeMap[stateId]; ok {
		return val
	}
	if val, ok := st.db.Find(stateId); ok {
		node := parseToNode(stateId, val)
		node.id = stateId
		return node
	}

	actionList := make([]*Action, 0)
	for _, action := range state.PossibleActions() {
		actionList = append(actionList, &Action{
			ID:    action,
			Score: 0,
		})
	}

	node := &Node{
		Actions: actionList,
		id:      stateId,
	}
	return node
}

func newSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func newSHA512(data []byte) string {
	hash := sha512.Sum512(data)
	return hex.EncodeToString(hash[:])
}

func (st *StateTree) SetDB(db Database) *StateTree {
	st.db = db
	return st
}

func (st *StateTree) DebugState(f func(n NodeDebug, debug Debug)) {
	st.debugState = f
}

func (st *StateTree) DebugAction(f func(actions []*Action, selected *Action)) {
	st.debugActions = f
}

type ControllerRequest struct {
	State State
}

type ControllerResponse struct {
	Restart bool
}

func (st *StateTree) Controller(f func(req ControllerRequest) ControllerResponse) {
	st.controller = f
}

func (st *StateTree) PlayGame(s State) {
	for {
		res := st.controller(st.playGame(s))
		if !res.Restart {
			break
		}
	}
}

func (st *StateTree) playGame(s State) ControllerRequest {
	state := s.Copy()

	nodeMap := make(map[string]*Node, 0)
	actionList := make([]*Action, 0)

	for {
		node := st.newNode(state, nodeMap)
		currentAction := node.selectAction(st.stats)
		st.debugActions(node.Actions, currentAction)

		state = state.Copy()
		state = state.PlayAction(currentAction.ID)
		state = state.PlaySideEffects()

		result := state.TurnResult(TurnRequest{Depth: len(actionList)})
		state = result.State

		// new node
		st.debugState(NodeDebug{State: state, Id: node.id}, CurrentState)

		actionList = append(actionList, currentAction)

		for _, action := range actionList {
			action.NVisited++
		}
		st.stats.NVisited++
		nodeMap[node.id] = node
		if result.EndGame {
			break
		}
	}
	gameResult := state.GameResult()
	for _, action := range actionList {
		action.Score += gameResult.Score
	}
	for key, val := range nodeMap {
		_ = st.db.Add(key, val.toDB())
	}

	return ControllerRequest{
		State: state,
	}
}

type State interface {
	ID() string
	PossibleActions() []interface{}
	Copy() State
	PlayAction(interface{}) State
	PlaySideEffects() State
	TurnResult(TurnRequest) TurnResult
	GameResult() GameResult
}

type GameResult struct {
	State State
	Score int
}

type TurnRequest struct {
	Depth int
}

type TurnResult struct {
	State   State
	EndGame bool
}

func (a Action) GetNVisited() int {
	return a.NVisited
}

func New() *StateTree {
	return &StateTree{
		debugState: func(n NodeDebug, debug Debug) {
			// default
		},
		debugActions: func(actions []*Action, selected *Action) {

		},
		controller: func(req ControllerRequest) ControllerResponse {
			return ControllerResponse{}
		},
		db:    DefaultMemoryDB{nodeMap: map[string]string{}},
		stats: &rootStats{NVisited: 0},
	}
}
