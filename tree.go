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
	nodeMap map[string]*Node
	logs    bytes.Buffer
}

type rootStats struct {
	NVisited int
}

type Node struct {
	Actions []*Action
	id      string
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

func (st *StateTree) getOrCreateNode(state State) (*Node, bool) {
	stateId := state.ID()

	if val, ok := st.nodeMap[stateId]; ok {
		return val, false
	}

	actionList := make([]*Action, 0)
	for _, action := range state.PossibleActions() {
		actionList = append(actionList, &Action{
			ID:    action,
			Turn:  state.Turn(),
			Score: 0,
		})
	}

	//rand.Shuffle(len(actionList), func(i, j int) {
	//	actionList[i], actionList[j] = actionList[j], actionList[i]
	//})

	node := &Node{
		Actions: actionList,
		id:      stateId,
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

func (st *StateTree) flush(s State, c StateTreeConfig) {
	if st.nodeMap == nil {
		st.nodeMap = make(map[string]*Node, 0)
	}

	if c.Flush {
		st.nodeMap = make(map[string]*Node, 0)
		st.stats = &rootStats{NVisited: 0}
	} else {
		node, _ := st.getOrCreateNode(s)
		st.stats = &rootStats{NVisited: node.totalNVisited()}
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

		// game loop
	gameLoop:
		for {
			node, newNode := st.getOrCreateNode(state)

			currentAction := node.selectAction(st.stats)

			if currentAction == nil {
				deadEnd = true
				statePreSim = state.Copy()
				break gameLoop
			}

			state.PlayAction(currentAction.ID)

			// new node
			actionList = append(actionList, currentAction)
			st.nodeMap[node.id] = node

			depth++
			currIteration++
			oneLoopGame = false
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
		TotalNodes:      len(st.nodeMap),
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
	ID() string
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
