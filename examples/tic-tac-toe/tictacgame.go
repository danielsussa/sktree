package tictactoe

import (
	"fmt"
	tree "github.com/danielsussa/sktree"
	"math/rand"
)

type player string

const (
	E player = "E"
	H player = "H"
	M player = "M"
)

type ticTacGame struct {
	board         []player
	lastMove      int
	currentPlayer player
}

func (t ticTacGame) Turn() tree.TurnKind {
	if t.currentPlayer == H {
		return tree.Human
	}
	return tree.Machine
}

func (t ticTacGame) Copy() tree.State {
	newBoard := make([]player, len(t.board))
	copy(newBoard, t.board)
	return &ticTacGame{
		board:         newBoard,
		lastMove:      t.lastMove,
		currentPlayer: t.currentPlayer,
	}
}

func (p player) toScore() float64 {
	switch p {
	case H:
		return 1
	case E:
		return 0
	case M:
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

func (t ticTacGame) PossibleActions() []any {
	iters := make([]any, 0)
	for idx, place := range t.board {
		if place == E {
			iters = append(iters, idx)
		}
	}
	return iters
}

func (t *ticTacGame) PlayAction(id any) {
	t.move(id.(int), t.currentPlayer)
	t.changePlayer()
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

func (t ticTacGame) GameResult() float64 {
	return t.winner().toScore()
}

func (t *ticTacGame) changePlayer() {
	if t.currentPlayer == H {
		t.currentPlayer = M
	} else {
		t.currentPlayer = H
	}
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
