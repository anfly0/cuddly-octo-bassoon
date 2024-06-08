package robot

import (
	"errors"
	"fmt"
	"reflect"
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
		{&Robot{room: Room{X: 3, Y: 3},
			coordinate: Coordinate{X: 1, Y: 1},
			compass:    *NewCompass('N')}, "L",
			'W', Coordinate{X: 1, Y: 1}},
		{&Robot{room: Room{X: 3, Y: 3},
			coordinate: Coordinate{X: 1, Y: 1},
			compass:    *NewCompass('N')}, "R",
			'E', Coordinate{X: 1, Y: 1}},
		{&Robot{room: Room{X: 3, Y: 3},
			coordinate: Coordinate{X: 1, Y: 1},
			compass:    *NewCompass('N')}, "F",
			'N', Coordinate{X: 1, Y: 0}},
		{&Robot{room: Room{X: 5, Y: 5},
			coordinate: Coordinate{X: 1, Y: 2},
			compass:    *NewCompass('N')}, "RFRFFRFRF",
			'N', Coordinate{X: 1, Y: 3}},
		{&Robot{room: Room{X: 5, Y: 5},
			coordinate: Coordinate{X: 0, Y: 0},
			compass:    *NewCompass('E')}, "RFLFFLRF",
			'E', Coordinate{X: 3, Y: 1}},
		{&Robot{room: Room{X: 1, Y: 1},
			coordinate: Coordinate{X: 0, Y: 0},
			compass:    *NewCompass('E')}, "RFLFFLRF",
			'E', Coordinate{X: 0, Y: 0}},
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

func TestNewRobot(t *testing.T) {
	type args struct {
		r Room
		d rune
		c Coordinate
	}
	tests := []struct {
		name    string
		args    args
		want    *Robot
		wantErr error
	}{
		{
			name: "Valid robot",
			args: args{
				r: Room{X: 3, Y: 3},
				d: 'N',
				c: Coordinate{X: 1, Y: 1},
			},
			want: &Robot{room: Room{X: 3, Y: 3}, compass: *NewCompass('N'), coordinate: Coordinate{X: 1, Y: 1}},
		},
		{
			name: "Valid robot",
			args: args{
				r: Room{X: 3, Y: 3},
				d: 'N',
				c: Coordinate{X: 1, Y: 1},
			},
			want:    &Robot{room: Room{X: 3, Y: 3}, compass: *NewCompass('N'), coordinate: Coordinate{X: 1, Y: 1}},
			wantErr: nil,
		},
		{
			name: "Robot created outside the room",
			args: args{
				r: Room{X: 1, Y: 1},
				d: 'N',
				c: Coordinate{X: 1, Y: 1},
			},
			want:    nil,
			wantErr: errors.New("the robot coordinates are outside the room"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := NewRobot(tt.args.r, tt.args.d, tt.args.c); !(reflect.DeepEqual(got, tt.want) && reflect.DeepEqual(err, tt.wantErr)) {
				t.Errorf("NewRobot() = %v, want %v err %v want %v", got, tt.want, err, tt.wantErr)
			}
		})
	}
}
