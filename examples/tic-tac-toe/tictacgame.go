package tictactoe

import (
	"fmt"
	tree "github.com/danielsussa/tmp_tree"
	"math/rand"
	"strconv"
)

type player string

const (
	E player = "E"
	X player = "X"
	O player = "O"
)

type ticTacGame struct {
	board    []player
	lastMove int
}

func (t ticTacGame) Copy() tree.State {
	newBoard := make([]player, len(t.board))
	copy(newBoard, t.board)
	return ticTacGame{
		board:    newBoard,
		lastMove: t.lastMove,
	}
}

func (p player) toScore() int {
	switch p {
	case X:
		return 1
	case E:
		return 0
	case O:
		return -1
	}
	return 0
}

func (t ticTacGame) winner() player {
	b := t.board
	if b[0] == b[1] && b[0] == b[2] && b[0] != E {
		return b[0]
	}
	if b[3] == b[4] && b[3] == b[5] && b[3] != E {
		return b[3]
	}
	if b[6] == b[7] && b[6] == b[8] && b[6] != E {
		return b[6]
	}

	if b[0] == b[3] && b[0] == b[6] && b[0] != E {
		return b[0]
	}
	if b[1] == b[4] && b[1] == b[7] && b[1] != E {
		return b[1]
	}
	if b[2] == b[5] && b[2] == b[8] && b[2] != E {
		return b[2]
	}

	if b[0] == b[4] && b[0] == b[8] && b[0] != E {
		return b[0]
	}
	if b[2] == b[4] && b[2] == b[6] && b[2] != E {
		return b[2]
	}
	return E
}

func (t ticTacGame) PossibleActions() []string {
	iters := make([]string, 0)
	for idx, place := range t.board {
		if place == E {
			iters = append(iters, fmt.Sprintf("%d", idx))
		}
	}
	return iters
}

func (t ticTacGame) PlayAction(a string) {
	id, _ := strconv.Atoi(a)
	t.move(id, X)
}

func (t ticTacGame) ID() string {
	return fmt.Sprintf("%v", t.board)
}

func (t ticTacGame) PlaySideEffects() {
	t.randomMove(O)
}

func (t ticTacGame) TurnResult(req tree.TurnRequest) tree.TurnResult {
	if t.winner() != E {
		return tree.TurnResult{
			EndGame: true,
		}
	}
	if len(t.PossibleActions()) == 0 {
		return tree.TurnResult{
			EndGame: true,
		}
	}
	return tree.TurnResult{
		EndGame: false,
	}
}

func (t ticTacGame) GameResult() tree.GameResult {
	return tree.GameResult{Score: t.winner().toScore()}
}

func (t ticTacGame) randomMove(p player) bool {
	free := make([]int, 0)

	for idx, place := range t.board {
		if place == E {
			free = append(free, idx)
		}
	}
	if len(free) == 0 {
		return false
	}
	place := free[rand.Intn(len(free))]
	t.board[place] = p
	return true
}

func (t ticTacGame) print() {
	fmt.Println("----------------------")
	fmt.Println(fmt.Sprintf("%s|%s|%s", t.board[0], t.board[1], t.board[2]))
	fmt.Println(fmt.Sprintf("%s|%s|%s", t.board[3], t.board[4], t.board[5]))
	fmt.Println(fmt.Sprintf("%s|%s|%s", t.board[6], t.board[7], t.board[8]))
}

func (t ticTacGame) move(idx int, p player) {
	t.board[idx] = p
}
