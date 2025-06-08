package main

import (
	"github.com/ronna-s/gc-eu-25/pkg/pnp"
	"github.com/ronna-s/gc-eu-25/pkg/pnp/engine/tview"
)

func main() {
	var pm ProductManager
	game := pnp.New(pm)
	game.Run(tview.New())
}

type ProductManager struct {
	Fired bool
}

func (p ProductManager) Options(g *pnp.Game) []pnp.Option {
	return []pnp.Option{
		{
			Description: "Pay wages",
			OnSelect: func() pnp.Outcome {
				if g.Coins < len(g.Players) {
					p.Fired = true
					return "Not enough coins to pay wages. PM was fired!"
				}
				g.Coins -= len(g.Players)
				return "Wages paid"
			},
		},
	}
}

func (p ProductManager) String() string {
	return "Sir Tan Lee Knot"
}

func (p ProductManager) AsciiArt() string {
	return `
 O
/|\
/ \`
}

func (p ProductManager) Alive() bool {
	return !p.Fired
}
