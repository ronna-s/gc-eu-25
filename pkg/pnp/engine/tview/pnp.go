// Package tview provides a simple P&P engine based on tview
package tview

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ronna-s/gc-eu-25/pkg/pnp"
	"github.com/ronna-s/gc-eu-25/pkg/pnp/engine"
)

var _ pnp.Engine = &Engine{}

type Engine struct {
	App       *tview.Application
	Pages     *tview.Pages
	Menu      *tview.List
	Inventory *tview.TextView
	Prod      *tview.TextView
	ProdState pnp.ProductionState
	mu        sync.Mutex
}

func New() *Engine {
	return &Engine{
		App:       tview.NewApplication(),
		Pages:     tview.NewPages(),
		Menu:      tview.NewList(),
		Inventory: tview.NewTextView(),
		Prod:      tview.NewTextView(),
	}
}
func (e *Engine) Start() {
	e.Prod.SetText(strings.Repeat("A", 2000)).
		SetTextColor(tcell.ColorGreen).
		SetBorder(true).
		SetTitle(fmt.Sprintf("PRODUCTION is `%s`", e.ProdState))
	go func() {
		e.Prod.SetChangedFunc(func() {
			e.App.Draw()
		})
		time.Sleep(time.Second)
		for {
			time.Sleep(time.Millisecond * 30)
			e.App.QueueUpdate(e.RenderProd)
		}
	}()
	if err := e.App.SetRoot(e.Pages, true).SetFocus(e.Pages).EnableMouse(true).Run(); err != nil {
		log.Fatal(err)
	}
}
func (e *Engine) Stop() {
	e.App.Stop()
}

func (e *Engine) RenderGame(g *pnp.Game) {
	players := g.Players
	currentPlayer := g.CurrentPlayer
	e.ProdState = g.Prod
	const pageName = "main"
	e.Pages.RemovePage(pageName)
	view := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(e.RenderPlayers(g.BandName, players, currentPlayer), 0, 2, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(e.Menu, 0, 1, true).
			AddItem(e.Inventory, 0, 1, false).
			AddItem(e.Prod, 0, 1, false), 0, 1, true)
	e.Inventory.Clear()
	e.Inventory.SetTitle("Inventory").SetBorder(true)
	e.Inventory.SetText(fmt.Sprintf("Coins: %d", g.Coins))
	e.Pages.AddAndSwitchToPage(pageName, view, true)
}

func (e *Engine) SelectOption(g *pnp.Game, player pnp.Player, cb func(selected pnp.Option)) {
	e.Menu.Clear()
	for i, o := range player.Options(g) {
		e.Menu.AddItem(o.String(), "", rune(49+i), nil)
	}
	e.Menu.SetCurrentItem(0)
	e.Menu.SetBorder(true).SetTitle("Select move...")
	e.Menu.SetSelectedFunc(func(choice int, s string, s2 string, r rune) {
		option := player.Options(g)[choice]
		cb(option)
	})
}

func (e *Engine) RenderOutcome(outcome pnp.Outcome, cb func()) {
	m := tview.NewModal()
	style := tcell.StyleDefault.Background(tcell.ColorBlack)
	m.SetText(string(outcome)).SetBackgroundColor(tcell.ColorBlack).SetBorderColor(tcell.ColorWhite).SetBorderStyle(style)

	m.AddButtons([]string{"ok"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "ok" {
			e.Pages.RemovePage("modal")
			cb()
		}
	})
	e.Pages.AddPage("modal", m, true, true)

}

type Livable interface {
	Alive() bool
}

func alive(p pnp.Player) bool {
	if livable, ok := p.(Livable); ok {
		return livable.Alive()
	}
	return true
}

func (e *Engine) RenderPlayers(bandName string, players []pnp.Player, current int) *tview.Flex {
	playersView := tview.NewFlex().SetDirection(tview.FlexRow)
	playersView.SetTitle(bandName).
		SetTitleColor(tcell.ColorWhite).
		SetBorder(true).SetBorderPadding(0, 0, 1, 0).SetBorderColor(tcell.ColorDeepPink)
	for i, p := range players {
		var color tcell.Color
		color = tcell.ColorWhite
		art := tview.NewTextView()
		art.SetBorderColor(tcell.ColorWhite)

		if alive(p) {
			art.SetTextColor(tcell.ColorWhite)
			art.SetText(p.AsciiArt())
		} else {
			art.SetText(engine.Gravestone).SetTextColor(tcell.ColorPurple)
		}
		if i == current {
			art.SetTitle(fmt.Sprintf("It's %s's turn", p)).
				SetBorderColor(tcell.ColorYellow)
		}

		art.SetTextColor(color).SetBorder(true).SetBorderPadding(0, 0, 1, 0)
		playersView.AddItem(art, 0, 1, false)
	}
	return playersView
}

var Rand = rand.Intn

func (e *Engine) RenderProd() {
	var color tcell.Color
	switch e.ProdState {
	case pnp.Calm:
		color = tcell.ColorGreen
	case pnp.Annoyed:
		color = tcell.ColorYellow
	case pnp.Enraged:
		color = tcell.ColorRed
	case pnp.Legacy:
		color = tcell.ColorPurple
	}

	text := e.Prod.GetText(false)
	for i := 0; i < 10; i++ {
		c := string(rune(Rand(128-48) + 48))
		r := Rand(len(text))
		text = text[:r] + c + text[r+1:]
	}
	e.Prod.SetText(text).SetTextColor(color)
	e.Prod.SetTitle(fmt.Sprintf("PRODUCTION is `%s`", e.ProdState))
	e.Prod.ScrollToBeginning()
}

func (e *Engine) GameWon() {
	m := NewModal().AddButtons("Yay!").
		SetButtonsAlign(tview.AlignCenter).
		SetText(engine.GameWon).
		SetTextColor(tcell.ColorLime).
		SetBorder(true).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			e.App.Stop()
		})
	m.ResizeItem(m.innerFlex, 0, 5)
	m.innerFlex.ResizeItem(m.modalFlex, 0, 5)

	e.Pages.AddPage("game won", m, true, true)
}

func (e *Engine) GameOver() {
	m := NewModal().AddButtons("Oh well...").
		SetButtonsAlign(tview.AlignCenter).
		SetText(engine.GameOver).
		SetTextColor(tcell.ColorLime).
		SetBorder(true).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			e.App.Stop()
		})
	m.ResizeItem(m.innerFlex, 0, 3)
	m.innerFlex.ResizeItem(m.modalFlex, 0, 3)

	e.Pages.AddPage("game over", m, true, true)
}

func (e *Engine) PizzaDelivery(fn func()) {
	const pageName = "pizza"
	m := NewModal().
		SetText(engine.Pizza).
		SetTextAlign(tview.AlignLeft).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			fn()
			e.Pages.RemovePage(pageName)
		}).
		AddButtons("Thanks, Boss!").
		SetButtonsAlign(tview.AlignCenter)
	m.SetBorder(true)
	m.SetBackgroundColor(tcell.ColorBlack).
		SetTextColor(tcell.ColorGreen)
	m.ResizeItem(m.innerFlex, 0, 3)
	m.innerFlex.ResizeItem(m.modalFlex, 0, 3)
	e.Pages.AddPage(pageName, m, true, true)
}

func (e *Engine) Welcome(fn func(bandName string)) {
	const modalName = "welcome modal"
	newGameText := tview.NewTextView()
	newGameText.SetText("A band of developers will attempt to survive against PRODUCTION!")
	gameArt := tview.NewTextView()
	gameArt.SetText(engine.Gamestarted).SetTextColor(tcell.ColorAqua)
	nameInput := tview.NewInputField().SetLabel("What is the name of your band?  ").SetText("Cool Band").SetFieldTextColor(tcell.ColorBlack).SetFieldBackgroundColor(tcell.ColorDarkCyan).SetFieldWidth(32)
	nameInput.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}
		bandName := nameInput.GetText()
		welcomeModal := tview.NewModal()
		welcomeModal.SetText("Hello, " + bandName + "! Are you ready?").SetBackgroundColor(tcell.ColorBlack)
		welcomeModal.SetTextColor(tcell.ColorDarkCyan)
		welcomeModal.AddButtons([]string{"Let's do this!"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			e.Pages.RemovePage(modalName)
			e.Pages.RemovePage("load")
			fn(bandName)
		})
		e.Pages.AddAndSwitchToPage(modalName, welcomeModal, true)
	})

	form := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(gameArt, 20, 20, false).
		AddItem(newGameText, 1, 1, false).
		AddItem(nameInput, 1, 1, true)

	form.SetBorderPadding(0, 0, 20, 0)
	form.SetBorder(true).SetTitle("New game started!").SetTitleAlign(tview.AlignLeft)
	e.Pages.AddAndSwitchToPage("load", tview.NewFlex().AddItem(form, 0, 1, true), true)
}

type Modal struct {
	*tview.Flex
	text      *tview.TextView
	form      *tview.Form
	innerFlex *tview.Flex
	modalFlex *tview.Flex
	done      func(idx int, label string)
}

func NewModal() *Modal {
	m := &Modal{form: tview.NewForm(), text: tview.NewTextView()}
	m.modalFlex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(m.text, 0, 4, false).
		AddItem(m.form, 0, 1, true)
	m.innerFlex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(m.modalFlex, 0, 2, true).
		AddItem(nil, 0, 1, false)
	m.Flex = tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(m.innerFlex, 0, 2, true).
		AddItem(nil, 0, 1, false)
	return m
}

func (m *Modal) AddButtons(labels ...string) *Modal {
	for i, label := range labels {
		m.form.AddButton(label, func() {
			m.done(i, label)
		})
	}
	return m
}

func (m *Modal) SetText(text string) *Modal {
	m.text.SetText(text)
	return m
}

func (m *Modal) SetTextAlign(align int) *Modal {
	m.text.SetTextAlign(align)
	return m

}
func (m *Modal) SetButtonsAlign(align int) *Modal {
	m.form.SetButtonsAlign(align)
	return m
}

func (m *Modal) SetBackgroundColor(color tcell.Color) *Modal {
	m.modalFlex.SetBackgroundColor(color)
	m.form.SetBackgroundColor(color)
	m.text.SetBackgroundColor(color)
	return m
}

func (m *Modal) SetBorder(show bool) *Modal {
	m.modalFlex.SetBorder(show)
	return m
}

func (m *Modal) SetDoneFunc(done func(buttonIndex int, buttonLabel string)) *Modal {
	m.done = done
	return m
}

func (m *Modal) SetTextColor(color tcell.Color) *Modal {
	m.text.SetTextColor(color)
	return m
}
