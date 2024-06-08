package robot

import (
	"fmt"
	"testing"
)

func TestNewCompass(t *testing.T) {
	tests := []struct{ direction, want rune }{
		{'N', 'N'},
		{'E', 'E'},
		{'S', 'S'},
		{'W', 'W'},
		{'A', 'N'},
		{'n', 'N'},
		{'e', 'E'},
		{'s', 'S'},
		{'w', 'W'},
	}

	for _, tt := range tests {
		tname := fmt.Sprintf("NewCompass: %s", string(tt.direction))

		t.Run(tname, func(t *testing.T) {
			r := NewCompass(tt.direction)
			if r.current() != tt.want {
				t.Errorf("Got %s, want %s", string(r.current()), string(tt.want))
			}
		})
	}
}

func TestCompassTurnR(t *testing.T) {
	tests := []struct{ direction, want rune }{
		{'N', 'E'},
		{'E', 'S'},
		{'S', 'W'},
		{'W', 'N'},
		{'A', 'E'},
		{'n', 'E'},
		{'e', 'S'},
		{'s', 'W'},
		{'w', 'N'},
	}

	for _, tt := range tests {
		tname := fmt.Sprintf("NewCompass: %s", string(tt.direction))

		t.Run(tname, func(t *testing.T) {
			r := NewCompass(tt.direction)
			r.turnR()
			if r.current() != tt.want {
				t.Errorf("Got %s, want %s", string(r.current()), string(tt.want))
			}
		})
	}
}

func TestCompassTurnL(t *testing.T) {
	tests := []struct{ direction, want rune }{
		{'N', 'W'},
		{'E', 'N'},
		{'S', 'E'},
		{'W', 'S'},
		{'A', 'W'},
		{'n', 'W'},
		{'e', 'N'},
		{'s', 'E'},
		{'w', 'S'},
	}

	for _, tt := range tests {
		tname := fmt.Sprintf("NewCompass: %s", string(tt.direction))

		t.Run(tname, func(t *testing.T) {
			r := NewCompass(tt.direction)
			r.turnL()
			if r.current() != tt.want {
				t.Errorf("Got %s, want %s", string(r.current()), string(tt.want))
			}
		})
	}
}

func TestRobotCmd(t *testing.T) {
	tests := []struct {
		robot  *Robot
		cmd    string
		want_d rune
		want_c Coordinate
	}{
		{&Robot{Room: Room{X: 3, Y: 3},
			Coordinate: Coordinate{X: 1, Y: 1},
			Compass:    *NewCompass('N')}, "L",
			'W', Coordinate{X: 1, Y: 1}},
		{&Robot{Room: Room{X: 3, Y: 3},
			Coordinate: Coordinate{X: 1, Y: 1},
			Compass:    *NewCompass('N')}, "R",
			'E', Coordinate{X: 1, Y: 1}},
		{&Robot{Room: Room{X: 3, Y: 3},
			Coordinate: Coordinate{X: 1, Y: 1},
			Compass:    *NewCompass('N')}, "F",
			'N', Coordinate{X: 1, Y: 0}},
		{&Robot{Room: Room{X: 5, Y: 5},
			Coordinate: Coordinate{X: 1, Y: 2},
			Compass:    *NewCompass('N')}, "RFRFFRFRF",
			'N', Coordinate{X: 1, Y: 3}},
		{&Robot{Room: Room{X: 5, Y: 5},
			Coordinate: Coordinate{X: 0, Y: 0},
			Compass:    *NewCompass('E')}, "RFLFFLRF",
			'E', Coordinate{X: 3, Y: 1}},
	}

	for _, tt := range tests {
		tname := fmt.Sprintf("Test cmd: %s", tt.cmd)

		t.Run(tname, func(t *testing.T) {
			r := tt.robot

			r.Cmd(tt.cmd)

			d, c := r.Report()

			if d != tt.want_d || c.X != tt.want_c.X || c.Y != tt.want_c.Y {
				t.Errorf("failed to process cmd %s. Got: %d %d %s want: %d %d %s", string(tt.cmd), c.X, c.Y, string(d), tt.want_c.X, tt.want_c.Y, string(tt.want_d))
			}
		})
	}
}
