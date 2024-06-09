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

// Returns the current direction the compass is pointing towards.
func (c *Compass) current() rune {
	return directions[c.index]
}

// Updates the compass to point the next direction that is right of the current direction.
// N->E->S->W->back to (N)orth
func (c *Compass) turnR() {
	c.index = (c.index + 1) % uint(len(directions))
}

// Updates the compass to point to the next direction that is left of the current direction.
// back to (N)orth<-E<-S<-W<-N
func (c *Compass) turnL() {
	c.index = (c.index - 1 + uint(len(directions))) % uint(len(directions))
}

// The Room struct is a record of the dimensions of the room that the robot is navigating in.
type Room struct {
	X uint `json:"x"`
	Y uint `json:"y"`
}

/*
*
The Coordinate struct is a record of the robots location in the room.
Note that a valid coordinate must always have X and Y values that are less that ditto values in the Room.
*
*/
type Coordinate struct {
	X uint `json:"x"`
	Y uint `json:"y"`
}

// The Robot struct contains only unexported fields. Use the exported methods to create/manipulate robots.
type Robot struct {
	compass    Compass
	room       Room
	coordinate Coordinate
	l          sync.RWMutex
}

// Creates a new robot with an initial state according to the argument. If the state is invalid the return value will be nil and an error
func NewRobot(r Room, d rune, c Coordinate) (*Robot, error) {
	comp := NewCompass(d)

	if c.X >= r.X || c.Y >= r.Y {
		// TODO: This should probably be a customer error type.
		return nil, errors.New("the robot coordinates are outside the room")
	}

	return &Robot{compass: *comp, coordinate: c, room: r}, nil
}

/*
Cmd executes a series of commands on the robot and returns the new state of the robot.
The commands will be executed one by one and the robots internal state will be updated.
If an invalid command is encountered processing is stopped and the latest state of the robot is returned.
*/
func (r *Robot) Cmd(cs string) (rune, Coordinate, error) {
	r.l.Lock()
	defer r.l.Unlock()

	for _, c := range cs {
		err := r.doCmd(c)
		if err != nil {
			d, coo := r.report()
			return d, coo, err
		}
	}
	d, coo := r.report()
	return d, coo, nil
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

// This is not thread safe version of the Report function and is used to avoid deadlocks when functions that has taken a exclusive lock needs the report data.
func (r *Robot) report() (rune, Coordinate) {
	return r.compass.current(), r.coordinate
}

// This is a thread safe and exported wrapper for the un-exported report function.
func (r *Robot) Report() (rune, Coordinate) {
	r.l.RLock()
	defer r.l.RUnlock()

	return r.report()
}
