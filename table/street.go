package table

import "errors"

type Street uint8

const (
	NoStreet Street = iota
	Preflop
	Flop
	Turn
	River
	Finished
)

func NewStreetFromString(str string) (Street, error) {
	switch str {
	case "preflop":
		return Preflop, nil
	case "flop":
		return Flop, nil
	case "turn":
		return Turn, nil
	case "river":
		return River, nil
	case "finished":
		return Finished, nil
	case "":
		return NoStreet, nil
	default:
		return 0, errors.New("unknown street")
	}
}

func (s Street) String() string {
	switch s {
	case Preflop:
		return "preflop"
	case Flop:
		return "flop"
	case Turn:
		return "turn"
	case River:
		return "river"
	case Finished:
		return "finished"
	case NoStreet:
		return "no street"
	default:
		return "unknown"
	}
}

// func (s Street) Add(num int) Street {
// 	next := int(s) + num
// 	if next > int(Finished) {
// 		return Finished
// 	}
// 	return Street(next)
// }
