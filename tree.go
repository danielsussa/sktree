package tree

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
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
	MaxNodes  int
	AllocNode int
	NumNodes  int
}

type Node struct {
	OpponentTurn bool
	Actions      []*Action
	Idx          int
}

type NodeDebug struct {
	Id    string
	State State
}

type Action struct {
	ID           any
	Score        float64
	OpponentTurn bool
	NVisited     int
	NodeIdx      int
}

func (a Action) GetNVisited() int {
	return a.NVisited
}

func (a Action) scoreAvg() float64 {
	return a.Score / float64(a.NVisited)
}

type actionScore struct {
	action *Action
	score  float64
}

func (n *Node) selectBestAction() *Action {
	var bestAction *Action
	for _, action := range n.Actions {
		if bestAction == nil {
			bestAction = action
		}
		if bestAction.scoreAvg() < action.scoreAvg() {
			bestAction = action
		}

	}
	return bestAction
}

func (n *Node) selectAction(stats *rootStats) *Action {
	actionScoreList := make([]actionScore, 0)

	for _, action := range n.Actions {
		if action.NVisited == 0 {
			return action
		}

		preScore := action.Score
		if n.OpponentTurn {
			preScore = preScore * -1
		}
		score := fSelection(preScore, action.NVisited, stats.NVisited)

		actionScoreList = append(actionScoreList, actionScore{
			action: action,
			score:  score,
		})
	}

	sort.Slice(actionScoreList, func(i, j int) bool {
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
	nVisited := 0
	for _, action := range n.Actions {
		nVisited += action.NVisited
	}
	return nVisited
}

func fSelection(value float64, nVisited, NVisited int) float64 {
	exploitation := value / float64(nVisited)
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

func (st *StateTree) allocArr(idx int) {
	if idx >= st.stats.MaxNodes {
		st.stats.MaxNodes = idx * 2
		st.nodeList = st.nodeList[:st.stats.MaxNodes]
	}
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
	opponentTurn := state.OpponentTurn()
	node := &Node{
		Idx:          idx,
		OpponentTurn: opponentTurn,
		Actions:      actionList,
	}
	st.nodeList[idx] = node
	st.stats.addNumNodes()
	for actIdx, action := range actions {
		actionList[actIdx] = &Action{
			ID:           action,
			NodeIdx:      st.stats.setIdx(),
			OpponentTurn: opponentTurn,
			Score:        0,
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

	currentAction := node.selectBestAction()

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
	st.stats.NVisited = -1
	st.stats.MaxNodes = 1000
	if st.nodeList == nil {
		st.nodeList = make([]*Node, st.stats.MaxNodes, 1000000*100)
	}

	if c.Flush {
		st.nodeList = make([]*Node, st.stats.MaxNodes, 1000000*100)
		st.stats = &rootStats{NVisited: 0}
	}
}

type StateTreeConfig struct {
	MaxDepth         int
	MaxIterations    int
	ScoreNormalizer  *ScoreNormalizer
	Flush            bool
	TotalSimulations int
}

type ScoreNormalizer struct {
	Min float64
	Max float64
}

type TrainResult struct {
	TotalNodes      int
	TotalIterations int
}

func (st *StateTree) Train(s State, c StateTreeConfig) TrainResult {
	if c.TotalSimulations == 0 {
		c.TotalSimulations = 1
	}
	st.flush(s, c)
	return st.train(s, c)
}

func (st *StateTree) train(s State, config StateTreeConfig) TrainResult {
	if st.copyError(s) {
		panic("object copy with different values")
	}

	for {
		state := s.Copy()
		var statePreSim State
		actionList := make([]*Action, 0)
		depth := 0
		oneLoopGame := true
		currSelectedNodeIdx := 0
		score := 0.0

		// game loop
	gameLoop:
		for {
			st.allocArr(currSelectedNodeIdx)
			node, newNode := st.getOrCreateNode(state, currSelectedNodeIdx)

			currentAction := node.selectAction(st.stats)

			if currentAction == nil {
				score = state.GameResult()
				statePreSim = state.Copy()
				break gameLoop
			}

			depth++
			oneLoopGame = false
			currSelectedNodeIdx = currentAction.NodeIdx

			if newNode {
				statePreSim = state.Copy()
				score = st.simulate(state, config)
				break gameLoop
			} else {
				state.PlayAction(currentAction.ID)
				actionList = append(actionList, currentAction)
			}
		}

		st.stats.NVisited++
		for _, action := range actionList {
			action.Score += score
			action.NVisited++
		}

		if st.finishGame(config, finishGameRequest{
			currDepth:     depth,
			currIteration: st.stats.NVisited,
			state:         statePreSim,
			oneLoopGame:   oneLoopGame,
		}) {
			break
		}
	}

	totalNodes := 0
	for _, n := range st.nodeList {
		if n != nil {
			totalNodes++
		}
	}

	return TrainResult{
		TotalIterations: st.stats.NVisited,
		TotalNodes:      totalNodes,
	}
}

func (st *StateTree) simulate(s State, config StateTreeConfig) float64 {
	score := 0.0
	for i := 0; i < config.TotalSimulations; i++ {
		state := s.Copy()
		for {
			actions := state.PossibleActions()
			if actions == nil || len(actions) == 0 {
				break
			}
			rndIdx := rand.Intn(len(actions))
			state.PlayAction(actions[rndIdx])
		}
		score += normalize(state.GameResult(), config.ScoreNormalizer)
	}
	return score / float64(config.TotalSimulations)
}

func (st *StateTree) copyError(s State) bool {
	sCopy := s.Copy()
	equal := cmp.Equal(s, sCopy)
	return !equal
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
	Copy() State
	PlayAction(any)
	PossibleActions() []any
	GameResult() float64
	OpponentTurn() bool
}

type TurnRequest struct {
	Depth int
}

type TurnResult struct {
	EndGame bool
}

func New() *StateTree {
	return &StateTree{
		controller: func(req ControllerRequest) ControllerResponse {
			return ControllerResponse{}
		},
		stats: &rootStats{NVisited: 0},
	}
}
