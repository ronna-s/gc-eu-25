package pnp

import (
	_ "embed"
	"math/rand"

	"github.com/ronna-s/gc-eu-25/pkg/repo"
)

type (
	// Game represents a Platforms and Programmers Game
	// where a band of player will attempt to take on production
	Game struct {
		Score          int
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
		PossibleActions(g *Game) []Action
		AsciiArt() string
		Alive() bool
	}

	// Engine represents the game's user interface rendering engine
	Engine interface {
		Welcome(leaderboard []repo.ScoreEntry, fn func(bandName string))
		GameOver()
		GameWon()
		RenderGame(g *Game)
		SelectAction(g *Game, player Player, cb func(selected Action))
		RenderOutcome(outcome Outcome, cb func())
		PizzaDelivery(cb func())
		WithOnExit(cb func()) Engine
	}
)

// New returns a new P&P game
func New(players ...Player) *Game {
	g := Game{Players: append(players), Prod: NewProduction(), Coins: 10}
	return &g
}

// Run starts a new game
func (g *Game) Run(e Engine) {
	leaderboard, _ := repo.GetTop(10)
	e.Welcome(leaderboard, func(bandName string) {
		g.BandName = bandName
		e = e.WithOnExit(func() {
			repo.Persist(repo.ScoreEntry{Score: g.Score, BandName: g.BandName})
		})
		g.MainLoop(e)
	})
}

// MainLoop kicks off the next players round
func (g *Game) MainLoop(e Engine) {
	g.Score = rand.Intn(10000)
	e.RenderGame(g)
	e.SelectAction(g, g.Players[g.CurrentPlayer], func(selected Action) {
		outcome := selected.Selected(g)
		e.RenderOutcome(outcome, func() {
			g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)
			g.MainLoop(e)
		})
	})
}

func allPlayersDead(players []Player) bool {
	for _, player := range players {
		if player.Alive() {
			return false
		}
	}
	return true
}

type Action struct {
	OnSelect    func(g *Game) Outcome
	Description string
}

func (o Action) String() string {
	return o.Description
}

func (o Action) Selected(g *Game) Outcome {
	return o.OnSelect(g)
}
