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
	compass    Compass
	room       Room
	coordinate Coordinate
	l          sync.RWMutex
}

func NewRobot(r Room, d rune, c Coordinate) (*Robot, error) {
	comp := NewCompass(d)

	if c.X >= r.X || c.Y >= r.Y {
		// TODO: This should probably be a customer error type.
		return nil, errors.New("the robot coordinates are outside the room")
	}

	return &Robot{compass: *comp, coordinate: c, room: r}, nil
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
		r.compass.turnL()
	case 'R':
		r.compass.turnR()
	case 'F':
		switch r.compass.current() {
		case 'S':
			if r.coordinate.Y < r.room.Y-1 {
				r.coordinate.Y++
			}
		case 'E':
			if r.coordinate.X < r.room.X-1 {
				r.coordinate.X++
			}
		case 'N':
			if r.coordinate.Y > 0 {
				r.coordinate.Y--
			}
		case 'W':
			if r.coordinate.X > 0 {
				r.coordinate.X--
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

	return r.compass.current(), r.coordinate
}
