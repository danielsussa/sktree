package tree

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
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
	controller   func(req ControllerRequest) ControllerResponse
	stats        *rootStats
				// [depth] -> [id]
	nodeMap map[string]*Node
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
	ID       string
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
		if action.NVisited == 0 {
			return action
		}
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
	if len(actionScoreList) == 0 {
		return nil
	}

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
	Expand       Debug = "expand"
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

func (st *StateTree) getOrCreateNode(state State) (*Node, bool) {
	stateId := state.ID()

	if val, ok := st.nodeMap[stateId]; ok {
		return val, false
	}

	actionList := make([]*Action, 0)
	for _, action := range state.PossibleActions() {
		actionList = append(actionList, &Action{
			ID:    action,
			Score: 0,
		})
	}

	rand.Shuffle(len(actionList), func(i, j int) {
		actionList[i], actionList[j] = actionList[j], actionList[i]
	})

	node := &Node{
		Actions: actionList,
		id:      stateId,
	}
	return node, true
}

func newSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func newSHA512(data []byte) string {
	hash := sha512.Sum512(data)
	return hex.EncodeToString(hash[:])
}


type ControllerRequest struct {
	State State
}

type ControllerResponse struct {
	ForceStop bool
}

func (st *StateTree) Controller(f func(req ControllerRequest) ControllerResponse) {
	st.controller = f
}

//func (st *StateTree) SingleSimulation(s State) {
//	for {
//		root := s.Copy()
//	}
//}


type PlayTurnResult struct {
	EndGame bool
	Action  *Action
}

func (st *StateTree) PlayTurn(state State) PlayTurnResult {

	node, _ := st.getOrCreateNode(state)

	currentAction := node.selectAction(st.stats)
	
	if currentAction == nil {
		return PlayTurnResult{}
	}

	state.PlayAction(currentAction.ID)

	return PlayTurnResult{
		EndGame: false,
		Action:  currentAction,
	}
}

func (st *StateTree) flush() {
	st.nodeMap = make(map[string]*Node, 0)
	st.stats = &rootStats{NVisited: 0}
}

type StateTreeConfig struct {
	MaxDepth   int
}

type TrainResult struct {
	TotalNewNodes int
}

func (st *StateTree) Train(s State, config StateTreeConfig) TrainResult{
	st.flush()
	return st.train(s, config)
}

func (st *StateTree) train(s State, config StateTreeConfig)  TrainResult{
	totalGamesWithoutNewNode := 0
	newNodes := 0

	for {
		state := s.Copy()
		actionMap := make(map[string]*Action, 0)
		depth := 0
		gameWithNewNode := false

		// game loop
		for {
			node, newNode := st.getOrCreateNode(state)

			currentAction := node.selectAction(st.stats)

			if currentAction == nil {
				break
			}

			state.PlayAction(currentAction.ID)
			state.PlaySideEffects()


			// new node
			actionMap[state.ID()] = currentAction
			st.nodeMap[node.id] = node

			depth++
			if newNode {
				gameWithNewNode = true
				break
			}
			//if depth > config.MaxDepth {
			//	break
			//}

		}

		if !gameWithNewNode {
			totalGamesWithoutNewNode++
		}

		st.stats.NVisited++
		gameResult := state.GameResult()
		for _, action := range actionMap {
			action.Score += gameResult.Score
			action.NVisited++
		}

		ctrlRes := st.controller(ControllerRequest{State: state})
		if ctrlRes.ForceStop {
			break
		}

		if depth > config.MaxDepth || totalGamesWithoutNewNode > 10 {
			break
		}
	}

	return TrainResult{
		TotalNewNodes: newNodes,
	}
}


type State interface {
	ID() string
	PossibleActions() []string
	Copy() State
	PlayAction(string)
	PlaySideEffects()
	GameResult() GameResult
}

type GameResult struct {
	Score int
}

type TurnRequest struct {
	Depth int
}

type TurnResult struct {
	EndGame bool
}

func (a Action) GetNVisited() int {
	return a.NVisited
}

func New() *StateTree {
	return &StateTree{
		controller: func(req ControllerRequest) ControllerResponse {
			return ControllerResponse{}
		},
		stats: &rootStats{NVisited: 0},
	}
}
