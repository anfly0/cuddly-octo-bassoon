package robot

import (
	"errors"
	"sync"
	"unicode"
)

var directions = []rune{'N', 'E', 'S', 'W'}

type Compass struct {
	index uint
}

// Creates a new compass set to one of the following directions N, E, S, W
func NewCompass(d rune) *Compass {
	var index uint
	d = unicode.ToUpper(d)

	for i, v := range directions {
		if v == d {
			index = uint(i)
		}
	}

	return &Compass{index: index}
}

func (c *Compass) current() rune {
	return directions[c.index]
}

func (c *Compass) turnR() {
	c.index = (c.index + 1) % uint(len(directions))
}

func (c *Compass) turnL() {
	c.index = (c.index - 1 + uint(len(directions))) % uint(len(directions))
}

type Room struct {
	X uint `json:"x"`
	Y uint `json:"y"`
}

type Coordinate struct {
	X uint `json:"x"`
	Y uint `json:"y"`
}

type Robot struct {
	Compass    Compass
	Room       Room
	Coordinate Coordinate
	l          sync.RWMutex
}

func NewRobot(r Room, d rune, x, y uint) *Robot {
	c := NewCompass(d)
	coo := Coordinate{X: x, Y: y}
	return &Robot{Compass: *c, Coordinate: coo, Room: r}
}

func (r *Robot) Cmd(cs string) error {
	r.l.Lock()
	defer r.l.Unlock()

	for _, c := range cs {
		err := r.doCmd(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Robot) doCmd(c rune) error {
	c = unicode.ToUpper(c)

	switch c {
	case 'L':
		r.Compass.turnL()
	case 'R':
		r.Compass.turnR()
	case 'F':
		switch r.Compass.current() {
		case 'S':
			if r.Coordinate.Y < r.Room.Y-1 {
				r.Coordinate.Y++
			}
		case 'E':
			if r.Coordinate.X < r.Room.X-1 {
				r.Coordinate.X++
			}
		case 'N':
			if r.Coordinate.Y > 0 {
				r.Coordinate.Y--
			}
		case 'W':
			if r.Coordinate.X > 0 {
				r.Coordinate.X--
			}
		}
	default:
		return errors.New("invalid command")
	}

	return nil
}

func (r *Robot) Report() (rune, Coordinate) {
	r.l.RLock()
	defer r.l.RUnlock()

	return r.Compass.current(), r.Coordinate
}
