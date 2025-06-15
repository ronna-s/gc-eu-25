package main

import (
	"github.com/ronna-s/gc-eu-25/pkg/pnp"
	engine "github.com/ronna-s/gc-eu-25/pkg/pnp/engine/tview"
)

func main() {
	var pm ProductManager
	app := pnp.New(&pm)
	app.Run(engine.New())
}

type ProductManager struct {
	Fired bool
}

func (p *ProductManager) PossibleActions(g *pnp.Game) []pnp.Action {
	return []pnp.Action{
		{
			Description: "Pay wages",
			OnSelect: func(g *pnp.Game) pnp.Outcome {
				if g.Coins < len(g.Players) {
					p.Fired = true
					return "Not enough coins to pay wages. Band is bankrupt. PM is fired!"
				}
				g.Coins -= len(g.Players)
				return "Wages paid"
			},
		},
	}
}

func (p *ProductManager) String() string {
	return "Sir Tan Lee Knot"
}

func (p *ProductManager) AsciiArt() string {
	return `
 O
/|\
/ \`
}

func (p *ProductManager) Alive() bool {
	return !p.Fired
}
