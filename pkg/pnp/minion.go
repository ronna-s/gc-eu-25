package pnp

import _ "embed"

func NewMinion(name string) Minion {
	return Minion{}
}

type Minion struct {
	Name string
}

//go:embed resources/minion.txt
var minionArt string

func (m Minion) AsciiArt() string {
	return minionArt
}

func (m Minion) PossibleActions(g *Game) []Action {
	var actions []Action
	if g.Coins > 0 {
		actions = append(actions, Action{
			Description: "Buy a banana and eat it (costs 1 gold coin)",
			OnSelect: func(g *Game) Outcome {
				g.Coins--
				return "You ate a banana"
			},
		})
	}
	actions = append(actions, Action{
		Description: "Add a bug to the code",
		OnSelect: func(g *Game) Outcome {
			return Outcome(g.Prod.Upset())
		},
	})
	return actions
}

func (m Minion) String() string {
	return "Minion"
}

func (m Minion) IsMinion() bool {
	return true
}

type minionPlayer interface {
	isMinion() bool
}

func isMinion(p Player) bool {
	if mp, ok := p.(minionPlayer); ok {
		return mp.isMinion()
	}
	return false
}
