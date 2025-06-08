package pnp

import "fmt"

// ProductionState represents a production state
// There are 4 possible states - Calm, Annoyed, Enraged and Legacy
type ProductionState int

const (
	Calm ProductionState = iota
	Annoyed
	Enraged
	Legacy
)

// React returns XP gained, Satisfaction gained and the new PRODUCTION state
func (s ProductionState) React(well bool) ProductionState {
	next := s.nextState(well)
	return next
}

func (s ProductionState) nextState(good bool) ProductionState {
	if good {
		if s == Calm {
			return s
		}
		return s - 1
	}
	if s == Legacy {
		return s
	}
	return s + 1
}

// String ...
func (s ProductionState) String() string {
	switch s {
	case Calm:
		return "Calm"
	case Annoyed:
		return "Annoyed"
	case Enraged:
		return "Enraged"
	case Legacy:
		return "Legacy"
	}
	return "not supported"
}

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Calm-0]
	_ = x[Annoyed-1]
	_ = x[Enraged-2]
	_ = x[Legacy-3]
}

// NewProduction ...
func NewProduction() ProductionState {
	return Calm
}

// Upset upsets production
func (p *ProductionState) Upset() string {
	*p = p.React(false)
	return fmt.Sprintf("PRODUCTION didn't like this move. PRODUCTION is now '%s'", *p)
}

// CalmDown calms production down
func (p *ProductionState) CalmDown() string {
	*p = p.React(true)
	return fmt.Sprintf("PRODUCTION is happy with this move. PRODUCTION is now '%s'", *p)
}

// NoImpact calms production down
func (p *ProductionState) NoImpact() string {
	return fmt.Sprintf("PRODUCTION is indifferent to this move. PRODUCTION is now '%s'", *p)
}
