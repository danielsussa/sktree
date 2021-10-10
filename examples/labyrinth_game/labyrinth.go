package labyrinth

import (
	"encoding/json"
	"fmt"
	tree "github.com/danielsussa/tmp_tree"
)

func labyrinthMap() [][]string {
	return [][]string{
		{"X", "X", "X", "X", "X", "X", "X", "X"},
		{"X", "X", "E", "X", "X", "X", "X", "X"},
		{"X", "X", " ", "X", "X", "X", "X", "X"},
		{"X", "X", "D", "X", "X", "X", "X", "X"},
		{"X", "X", " ", "X", "X", " ", "X", "X"},
		{"X", "X", " ", "X", "X", " ", "X", "X"},
		{"X", "X", " ", " ", " ", " ", " ", "X"},
		{"X", "X", "X", "X", " ", "X", " ", "X"},
		{"X", "K", " ", "X", " ", "X", "X", "X"},
		{"X", " ", " ", " ", " ", "X", "X", "X"},
		{"X", " ", " ", "X", "P", "X", "X", "X"},
		{"X", "X", "X", "X", "X", "X", "X", "X"},
	}
}

func playerMap() [][]string {
	labMap := labyrinthMap()
	for j, _ := range labMap {
		for i, _ := range labMap[j] {
			if labMap[j][i] == "P" {
				continue
			}
			labMap[j][i] = "X"
		}
	}
	return labMap
}

func findPlace(labMap [][]string, p string) (j, i int) {
	for j, _ := range labMap {
		for i, _ := range labMap[j] {
			if labMap[j][i] == p {
				return j, i
			}
		}
	}
	panic("error cannot find start place")
}

type game struct {
	PlaceMap  [][]string
	PlayerMap [][]string
	HasKey    bool
	// limit 50 to max
	TotalMoves int
	MaxMoves   int
	Score      float64
	WinGame    bool
}

func (g game) ID() string {
	return fmt.Sprintf("%v-%v", g.PlaceMap, g.HasKey)
}

func (g game) PossibleActions() []interface{} {
	j, i := findPlace(g.PlaceMap, "P")

	up := g.PlaceMap[j-1][i]
	down := g.PlaceMap[j+1][i]
	left := g.PlaceMap[j][i-1]
	right := g.PlaceMap[j][i+1]

	actions := make([]interface{}, 0)

	if up == " " {
		actions = append(actions, "move_up")
	}
	if down == " " {
		actions = append(actions, "move_down")
	}
	if right == " " {
		actions = append(actions, "move_right")
	}
	if left == " " {
		actions = append(actions, "move_left")
	}

	if up == "K" || down == "K" || right == "K" || left == "K" {
		actions = append(actions, "get_key")
	}
	if (up == "D" || down == "D" || right == "D" || left == "D") && g.HasKey {
		actions = append(actions, "open_door")
	}
	if up == "E" || down == "E" || right == "E" || left == "E" {
		actions = append(actions, "go_to_end")
	}

	return actions
}

func (g game) PlayAction(action interface{}) tree.State {
	j, i := findPlace(g.PlaceMap, "P")
	switch action.(string) {
	case "move_up":
		g.PlaceMap[j][i] = " "
		g.PlaceMap[j-1][i] = "P"

		g.PlayerMap[j][i] = " "
		g.PlayerMap[j-1][i] = "P"
	case "move_down":
		g.PlaceMap[j][i] = " "
		g.PlaceMap[j+1][i] = "P"

		g.PlayerMap[j][i] = " "
		g.PlayerMap[j+1][i] = "P"
	case "move_left":
		g.PlaceMap[j][i] = " "
		g.PlaceMap[j][i-1] = "P"

		g.PlayerMap[j][i] = " "
		g.PlayerMap[j][i-1] = "P"
	case "move_right":
		g.PlaceMap[j][i] = " "
		g.PlaceMap[j][i+1] = "P"

		g.PlayerMap[j][i] = " "
		g.PlayerMap[j][i+1] = "P"
	case "get_key":
		g.HasKey = true
		jk, ik := findPlace(g.PlaceMap, "K")
		g.PlaceMap[j][i] = " "
		g.PlaceMap[jk][ik] = "P"

		g.PlayerMap[j][i] = " "
		g.PlayerMap[jk][ik] = "P"
	case "open_door":
		jd, id := findPlace(g.PlaceMap, "D")
		g.PlaceMap[j][i] = " "
		g.PlaceMap[jd][id] = "P"

		g.PlayerMap[j][i] = " "
		g.PlayerMap[jd][id] = "P"
	case "go_to_end":
		g.WinGame = true
	}
	g.TotalMoves++
	return g
}

func (g game) Print() {
	for j, _ := range g.PlaceMap {
		for i, _ := range g.PlaceMap[j] {
			fmt.Print(g.PlaceMap[j][i] + "   ")
		}
		fmt.Print("    ")
		for i, _ := range g.PlaceMap[j] {
			fmt.Print(g.PlayerMap[j][i] + "   ")
		}
		fmt.Println()
	}
	fmt.Println()
}

func (g game) PlaySideEffects() tree.State {
	return g
}

func (g game) TurnResult(request tree.TurnRequest) tree.TurnResult {
	endGame := false
	if g.TotalMoves >= g.MaxMoves {
		endGame = true
	}
	if g.WinGame {
		endGame = true
	}
	return tree.TurnResult{
		State:   g,
		EndGame: endGame,
	}
}

func (g game) GameResult() tree.GameResult {
	return tree.GameResult{
		State: g,
		Score: g.MaxMoves - g.TotalMoves,
	}
}

func (g game) Copy() tree.State {
	b, _ := json.Marshal(g)
	gCopy := game{}
	_ = json.Unmarshal(b, &gCopy)
	return gCopy
}

func newGame() game {
	return game{
		PlaceMap:  labyrinthMap(),
		PlayerMap: playerMap(),
		MaxMoves:  25,
	}
}
