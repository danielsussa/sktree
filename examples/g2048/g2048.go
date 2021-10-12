package g2048

import (
	"fmt"
	tree "github.com/danielsussa/tmp_tree"
	"math/rand"
	"sort"
)

type g2048 struct {
	board []int
	score int
}

func (g g2048) Copy() tree.State {
	gCopy := make([]int, 16)
	copy(gCopy, g.board)
	return &g2048{
		board: gCopy,
		score: g.score,
	}
}

func (g g2048) PossibleActions() []string {
	iters := make([]string, 0)
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
	return iters
}

func print2048(board []int, score int) {
	fmt.Print("\033[H\033[2J")
	fmt.Println(fmt.Sprintf("------- %v --------", score))
	for i := 0; i < 4; i++ {
		k := i * 4
		fmt.Print(fmt.Sprintf("%-6d %-6d %-6d %-6d", board[0+k], board[1+k], board[2+k], board[3+k]))
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
	//return fmt.Sprintf("%v", convertScalar(g.board))
	return fmt.Sprintf("%v", g.board)
}

func (g *g2048) PlayAction(i string) {
	score := 0
	if i == "D" {
		score += computeDown(g.board)
	} else if i == "U" {
		score += computeUp(g.board)
	} else if i == "R" {
		score += computeRight(g.board)
	} else if i == "L" {
		score += computeLeft(g.board)
	}
	g.score += score
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

func (g g2048) PlaySideEffects() {
	addNumberOnBoard(g.board)
}

func (g g2048) TurnResult(r tree.TurnRequest) tree.TurnResult {
	iters := len(g.PossibleActions())
	return tree.TurnResult{
		EndGame: iters == 0,
	}
}

func (g g2048) GameResult() tree.GameResult {
	return tree.GameResult{
		Score: g.score,
	}
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
				//break
			}
		}
	}
	return score
}

func startNewGame() *g2048 {
	game := &g2048{
		board: make([]int, 16),
	}
	addNumberOnBoard(game.board)
	return game
}

func newArr() []int {
	return make([]int, 16)
}
