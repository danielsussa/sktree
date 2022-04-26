package tree

import (
	"bytes"
	"encoding/json"
	"math"
	"math/rand"
	"sort"
)

type StateTree struct {
	controller func(req ControllerRequest) ControllerResponse
	stats      *rootStats
	// [depth] -> [id]
	nodeList []*Node
	logs     bytes.Buffer
}

type rootStats struct {
	NVisited  int
	AllocNode int
	NumNodes  int
}

type Node struct {
	Actions []*Action
}

type NodeDebug struct {
	Id    string
	State State
}

type Action struct {
	ID       any
	Score    float64
	Turn     TurnKind
	deadEnd  bool
	NVisited int
	NodeIdx  int
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
		if action.deadEnd {
			continue
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

func (n *Node) totalNVisited() int {
	nvisited := 0
	for _, action := range n.Actions {
		nvisited += action.NVisited
	}
	return nvisited
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

func (rs *rootStats) setIdx() int {
	rs.AllocNode++
	return rs.AllocNode
}

func (rs *rootStats) addNumNodes() {
	rs.NumNodes++
}

func (st *StateTree) getOrCreateNode(state State, idx int) (*Node, bool) {
	if idx < st.stats.NumNodes {
		node := st.nodeList[idx]
		if node != nil {
			return st.nodeList[idx], false
		}
	}
	actions := state.PossibleActions()
	actionList := make([]*Action, len(actions))
	node := &Node{
		Actions: actionList,
	}
	st.nodeList[idx] = node
	st.stats.addNumNodes()
	for actIdx, action := range actions {
		actionList[actIdx] = &Action{
			ID:      action,
			NodeIdx: st.stats.setIdx(),
			Turn:    state.Turn(),
			Score:   0,
		}
	}

	return node, true
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

type PlayTurnResult struct {
	EndGame bool
	Action  *Action
}

func (st *StateTree) PlayTurn(state State) PlayTurnResult {

	node := st.nodeList[0]

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

func (st *StateTree) flush(s State, c StateTreeConfig) {
	if st.nodeList == nil {
		st.nodeList = make([]*Node, 1000*1000000)
	}

	if c.Flush {
		st.nodeList = make([]*Node, 1000*1000000)
		st.stats = &rootStats{NVisited: 0}
	}
}

type StateTreeConfig struct {
	MaxDepth        int
	MaxIterations   int
	ScoreNormalizer *ScoreNormalizer
	Flush           bool
}

type ScoreNormalizer struct {
	Min float64
	Max float64
}

type TrainResult struct {
	TotalNodes      int
	TotalIterations int
}

func (st *StateTree) Train(s State, config StateTreeConfig) TrainResult {
	st.flush(s, config)
	return st.train(s, config)
}

func (st *StateTree) train(s State, config StateTreeConfig) TrainResult {
	if st.copyError(s) {
		return TrainResult{}
	}

	currIteration := 0
	for {
		state := s.Copy()
		var statePreSim State
		actionList := make([]*Action, 0)
		depth := 0
		deadEnd := false
		oneLoopGame := true
		crrSelectedIdx := 0

		// game loop
	gameLoop:
		for {
			node, newNode := st.getOrCreateNode(state, crrSelectedIdx)

			currentAction := node.selectAction(st.stats)

			if currentAction == nil {
				deadEnd = true
				statePreSim = state.Copy()
				break gameLoop
			}

			state.PlayAction(currentAction.ID)

			depth++
			currIteration++
			oneLoopGame = false
			crrSelectedIdx = currentAction.NodeIdx
			if newNode {
				statePreSim = state.Copy()
				for {
					actions := state.PossibleActions()
					if actions == nil || len(actions) == 0 {
						break gameLoop
					}
					rndIdx := rand.Intn(len(actions))
					state.PlayAction(actions[rndIdx])
				}
			} else {
				actionList = append(actionList, currentAction)
			}

		}

		st.stats.NVisited++
		score := normalize(state.GameResult(), config.ScoreNormalizer)
		for idx, action := range actionList {
			if deadEnd && len(actionList)-1 == idx {
				action.deadEnd = true
			}
			if action.Turn == Human {
				action.Score += score
			} else {
				action.Score -= score
			}

			action.NVisited++
		}

		if st.finishGame(config, finishGameRequest{
			currDepth:     depth,
			currIteration: currIteration,
			state:         statePreSim,
			oneLoopGame:   oneLoopGame,
		}) {
			break
		}
	}

	return TrainResult{
		TotalIterations: currIteration,
		TotalNodes:      len(st.nodeList),
	}
}

func (st *StateTree) copyError(s State) bool {
	sCopy := s.Copy()
	b1, err := json.Marshal(s)
	if err != nil {
		return true
	}

	b2, err := json.Marshal(sCopy)
	if err != nil {
		return true
	}

	if bytes.Compare(b1, b2) != 0 {
		return true
	}
	return false
}

type finishGameRequest struct {
	state         State
	currDepth     int
	currIteration int
	oneLoopGame   bool
}

func normalize(val float64, n *ScoreNormalizer) float64 {
	if n == nil {
		return val
	}
	return (val - n.Min) / (n.Max - n.Min)
}

func (st *StateTree) finishGame(config StateTreeConfig, req finishGameRequest) bool {
	ctrlRes := st.controller(ControllerRequest{State: req.state})
	if ctrlRes.ForceStop {
		return true
	}
	if config.MaxDepth > 0 && req.currDepth >= config.MaxDepth {
		return true
	}
	if config.MaxIterations > 0 && req.currIteration >= config.MaxIterations {
		return true
	}
	if req.oneLoopGame {
		return true
	}
	return false
}

type State interface {
	PossibleActions() []any
	Copy() State
	PlayAction(any)
	GameResult() float64
	Turn() TurnKind
}

type TurnKind string

const (
	Human   TurnKind = "human"
	Machine TurnKind = "machine"
)

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
