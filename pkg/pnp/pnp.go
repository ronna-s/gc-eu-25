package pnp

import (
	_ "embed"
)

type (
	// Game represents a Platforms and Programmers Game
	// where a band of player will attempt to take on production
	Game struct {
		Players        []Player
		Prod           ProductionState
		BandName       string
		Coins          int
		CurrentPlayer  int
		ProductManager string
	}

	Outcome string

	// Player represents a P&P player
	Player interface {
		Options(g *Game) []Option
		AsciiArt() string
		Alive() bool
	}

	// Engine represents the game's user interface rendering engine
	Engine interface {
		Welcome(cb func(bandName string))
		Start()
		GameOver()
		GameWon()
		RenderGame(g *Game)
		SelectOption(g *Game, player Player, cb func(selected Option))
		RenderOutcome(outcome Outcome, cb func())
		PizzaDelivery(cb func())
	}
)

// New returns a new P&P game
func New(players ...Player) *Game {
	g := Game{Players: players, Prod: NewProduction(), Coins: 10}
	return &g
}

// Run starts a new game
func (g *Game) Run(e Engine) {
	g.Welcome(e, func() {
		g.MainLoop(e)
	})
	e.Start()
}
func (g *Game) Welcome(e Engine, fn func()) {
	e.Welcome(func(bandName string) {
		g.BandName = bandName
		fn()
	})
}

// MainLoop kicks off the next players round
func (g *Game) MainLoop(e Engine) {
	e.RenderGame(g)
	e.SelectOption(g, g.Players[g.CurrentPlayer], func(selected Option) {
		outcome := selected.Selected()
		e.RenderOutcome(outcome, func() {
			g.CurrentPlayer++
			if g.CurrentPlayer >= len(g.Players) {
				g.CurrentPlayer = 0
			}
			g.MainLoop(e)
		})
	})
}

type Option struct {
	OnSelect    func() Outcome
	Description string
}

func (o Option) String() string {
	return o.Description
}

func (o Option) Selected() Outcome {
	return o.OnSelect()
}
