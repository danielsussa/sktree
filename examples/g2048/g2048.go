package g2048

import (
	"fmt"
	tree "github.com/danielsussa/sktree"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

type g2048 struct {
	board      []int
	turnsCount int
	score      int
	turn       tree.TurnKind
}

func (g g2048) Turn() tree.TurnKind {
	return g.turn
}

func (g g2048) Copy() tree.State {
	gCopy := make([]int, 16)
	copy(gCopy, g.board)
	return &g2048{
		board: gCopy,
		score: g.score,
		turn:  g.turn,
	}
}

func (g g2048) PossibleActions() []any {
	iters := make([]any, 0)
	if g.turn == tree.Human {
		if canMoveUp(g.board) {
			iters = append(iters, "U")
		}
		if canMoveDown(g.board) {
			iters = append(iters, "D")
		}
		if canMoveRight(g.board) {
			iters = append(iters, "R")
		}
		if canMoveLeft(g.board) {
			iters = append(iters, "L")
		}
	} else {
		freePlaces := getFreePlaces(g.board)
		for _, freePlace := range freePlaces {
			iters = append(iters, fmt.Sprintf("2-%d", freePlace))
			iters = append(iters, fmt.Sprintf("4-%d", freePlace))
		}
	}

	return iters
}

func print2048WithRes(g *g2048, trainRes tree.TrainResult, res tree.PlayTurnResult) {
	fmt.Println()
	fmt.Println(fmt.Sprintf("------- %s --------", res.Action.ID))
	for i := 0; i < 4; i++ {
		k := i * 4
		fmt.Print(fmt.Sprintf("%-6d %-6d %-6d %-6d", g.board[0+k], g.board[1+k], g.board[2+k], g.board[3+k]))
		fmt.Println()
	}
	fmt.Println(fmt.Sprintf("nodes: %d", trainRes.TotalNodes))
}

func print2048(g *g2048) {
	fmt.Print("\033[H\033[2J")
	fmt.Println(fmt.Sprintf("------- %v --------", g.score))
	for i := 0; i < 4; i++ {
		k := i * 4
		fmt.Print(fmt.Sprintf("%-6d %-6d %-6d %-6d", g.board[0+k], g.board[1+k], g.board[2+k], g.board[3+k]))
		fmt.Println()
	}
}

func convertScalar(board []int) []int {
	mapConverter := make(map[int]int, 0)
	for _, val := range board {
		mapConverter[val] = 0
	}
	keys := make([]int, 0)
	for k := range mapConverter {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	i := 0
	for _, k := range keys {
		mapConverter[k] = i
		i++
	}

	newBoard := newArr()
	for idx, _ := range newBoard {
		newBoard[idx] = mapConverter[board[idx]]
	}

	return newBoard
}

func (g g2048) ID() string {
	return fmt.Sprintf("%s-%v", g.turn, g.board)
}

func (g *g2048) PlayAction(i any) {
	score := 0
	if i == "D" {
		score += computeDown(g.board)
		g.turn = tree.Machine
	} else if i == "U" {
		score += computeUp(g.board)
		g.turn = tree.Machine
	} else if i == "R" {
		score += computeRight(g.board)
		g.turn = tree.Machine
	} else if i == "L" {
		score += computeLeft(g.board)
		g.turn = tree.Machine
	} else if strings.Contains(i.(string), "-") {
		ispl := strings.Split(i.(string), "-")
		number, _ := strconv.Atoi(ispl[0])
		place, _ := strconv.Atoi(ispl[1])
		g.board[place] = number
		g.turn = tree.Human
	}
	g.score += score
	g.turnsCount++
}

func canMoveUp(board []int) bool {
	for i := 0; i < 4; i++ {
		lane := []int{board[0+i], board[4+i], board[8+i], board[12+i]}
		if canMerge(lane) {
			return true
		}
	}
	return false
}

func computeUp(board []int) int {
	score := 0
	for i := 0; i < 4; i++ {
		lane := []int{board[0+i], board[4+i], board[8+i], board[12+i]}
		score += merge(lane)
		board[0+i] = lane[0]
		board[4+i] = lane[1]
		board[8+i] = lane[2]
		board[12+i] = lane[3]
	}
	return score
}

func topFirstValue(board []int) int {
	maxVal := 2
	for _, val := range board {
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}

func canMoveDown(board []int) bool {
	for i := 0; i < 4; i++ {
		lane := []int{board[12+i], board[8+i], board[4+i], board[0+i]}
		if canMerge(lane) {
			return true
		}
	}
	return false
}

func computeDown(board []int) int {
	score := 0
	for i := 0; i < 4; i++ {
		lane := []int{board[12+i], board[8+i], board[4+i], board[0+i]}
		score += merge(lane)
		board[12+i] = lane[0]
		board[8+i] = lane[1]
		board[4+i] = lane[2]
		board[0+i] = lane[3]
	}
	return score
}

func canMoveRight(board []int) bool {
	for i := 0; i < 4; i++ {
		k := 4 * i
		lane := []int{board[3+k], board[2+k], board[1+k], board[0+k]}
		if canMerge(lane) {
			return true
		}
	}
	return false
}

func computeRight(board []int) int {
	score := 0
	for i := 0; i < 4; i++ {
		k := 4 * i
		lane := []int{board[3+k], board[2+k], board[1+k], board[0+k]}
		score += merge(lane)
		board[3+k] = lane[0]
		board[2+k] = lane[1]
		board[1+k] = lane[2]
		board[0+k] = lane[3]
	}
	return score
}

func canMoveLeft(board []int) bool {
	for i := 0; i < 4; i++ {
		k := 4 * i
		lane := []int{board[0+k], board[1+k], board[2+k], board[3+k]}
		if canMerge(lane) {
			return true
		}
	}
	return false
}

func computeLeft(board []int) int {
	score := 0
	for i := 0; i < 4; i++ {
		k := 4 * i
		lane := []int{board[0+k], board[1+k], board[2+k], board[3+k]}
		score += merge(lane)
		board[0+k] = lane[0]
		board[1+k] = lane[1]
		board[2+k] = lane[2]
		board[3+k] = lane[3]
	}
	return score
}

func (g *g2048) playSideEffects() {
	addNumberOnBoard(g.board)
	g.turn = tree.Human
}

func (g g2048) simpleScore() float64 {
	return float64(g.score)
}

func (g g2048) freePlacesAndActions() float64 {
	return float64(len(getFreePlaces(g.board)) + len(g.PossibleActions()))
}

func (g g2048) freePlaces() float64 {
	return float64(len(getFreePlaces(g.board)))
}

func (g g2048) topFirst() int {
	mapConverter := make(map[int]int, 0)
	for _, val := range g.board {
		mapConverter[val] = 0
	}
	keys := make([]int, 0)
	for k := range mapConverter {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys[len(keys)-1]
}

func (g g2048) top3Score() float64 {
	mapConverter := make(map[int]int, 0)
	for _, val := range g.board {
		mapConverter[val] = 0
	}
	keys := make([]int, 0)
	for k := range mapConverter {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	total := 0
	for i := 0; i < len(keys); i++ {
		total += keys[i]
	}
	return float64(total)
}

func (g g2048) GameResult() float64 {
	//freePlace := len(getFreePlaces(g.board))
	//possibles :=  len(g.PossibleActions())
	//return tree.GameResult{Score: g.topFirst()}
	//return g.freePlacesAndActions()
	return g.simpleScore()
}

func getFreePlaces(board []int) []int {
	freePlaces := make([]int, 0)
	for idx, val := range board {
		if val == 0 {
			freePlaces = append(freePlaces, idx)
		}
	}
	return freePlaces
}

func addNumberOnBoard(board []int) {
	freePlaces := getFreePlaces(board)
	if len(freePlaces) == 0 {
		return
	}

	freePlace := freePlaces[rand.Intn(len(freePlaces))]
	fRand := rand.Float64()
	val := 2
	if fRand >= 0.9 {
		val = 4
	}
	board[freePlace] = val
}

func canMerge(c []int) bool {
	for i := 0; i < 4; i++ {
		val := c[i]
		for j := i + 1; j < 4; j++ {
			iterV := c[j]
			if val != iterV && iterV != 0 && val != 0 {
				break
			}
			if val == iterV && val != 0 {
				return true
			}
			if val == 0 && iterV != 0 {
				return true
			}
		}
	}
	return false
}
func merge(c []int) int {
	score := 0
	for i := 0; i < 4; i++ {
		val := c[i]
		for j := i + 1; j < 4; j++ {
			iterV := c[j]
			if val != iterV && iterV != 0 && val != 0 {
				break
			}
			if val == iterV && val != 0 {
				score += 2 * val
				c[i] = val * 2
				c[j] = 0
				break
			}
			if val == 0 && iterV != 0 {
				c[i] = iterV
				c[j] = 0
				val = iterV
			}
		}
	}
	return score
}

func startNewGame() *g2048 {
	game := &g2048{
		board: make([]int, 16),
		turn:  tree.Human,
	}
	addNumberOnBoard(game.board)
	return game
}

func newArr() []int {
	return make([]int, 16)
}
