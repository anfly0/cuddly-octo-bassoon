package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/anfly0/cuddly-octo-bassoon/internal/robot"
	"github.com/anfly0/cuddly-octo-bassoon/internal/storage"
)

type robotVoidStore struct{}

func (rs *robotVoidStore) Get(_ string, _ context.Context) *robot.Robot {
	return nil
}

func (rs *robotVoidStore) Put(_ string, _ *robot.Robot, _ context.Context) error {
	return nil
}

func TestRobotHandler_create(t *testing.T) {

	voidStore := &robotVoidStore{}
	robotHandler := RobotHandler{store: voidStore}

	type args struct {
		body reqCreate
	}

	type rsp struct {
		code int
		//TODO: add expected body to make the test cases more complete.
	}
	tests := []struct {
		name string
		args args
		want rsp
	}{
		{
			name: "Create valid robot",
			args: args{body: reqCreate{Direction: "N", Room: robot.Room{X: 1, Y: 1}, Start: robot.Coordinate{X: 0, Y: 0}}},
			want: rsp{code: http.StatusOK},
		},
		{
			name: "Create invalid robot",
			args: args{body: reqCreate{Direction: "N", Room: robot.Room{X: 1, Y: 1}, Start: robot.Coordinate{X: 1, Y: 1}}},
			want: rsp{code: http.StatusBadRequest},
		},
		{
			name: "Create robot with invalid direction",
			args: args{body: reqCreate{Direction: "", Room: robot.Room{X: 1, Y: 1}, Start: robot.Coordinate{X: 0, Y: 0}}},
			want: rsp{code: http.StatusBadRequest},
		},
		{
			name: "Create robot with invalid room",
			args: args{body: reqCreate{Direction: "N", Room: robot.Room{X: 0, Y: 0}, Start: robot.Coordinate{X: 0, Y: 0}}},
			want: rsp{code: http.StatusBadRequest},
			// TODO: This test is passing because all possible robot starting positions are "outside" the room. There should be a check when crating the room to give a better error meg.
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			payload, err := json.Marshal(tt.args.body)
			if err != nil {
				t.Fatal(err)
			}

			b := strings.NewReader(string(payload))

			req, err := http.NewRequest("POST", "/robot", b)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(robotHandler.create)

			handler.ServeHTTP(rr, req)

			// Check the status code
			if rr.Result().StatusCode != tt.want.code {
				t.Errorf("wrong status code: got %v want %v", rr.Code, tt.want.code)
			}
		})
	}
}

func TestRobotHandler_getStatus(t *testing.T) {

	robotStore := storage.NewRobotMemStore()
	robotHandler := RobotHandler{store: robotStore}

	type args struct {
		robotId string
		reqId   string
		room    robot.Room
		coo     robot.Coordinate
		d       rune
	}

	type rsp struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want rsp
	}{
		// Make sure all test cases has a uniq id for the robot. All test cases share the same robot store.
		{
			name: "Get a robot that is in the store",
			args: args{robotId: "abc", reqId: "abc", room: robot.Room{X: 1, Y: 1}, coo: robot.Coordinate{X: 0, Y: 0}, d: 'N'},
			want: rsp{code: http.StatusOK},
		},
		{
			name: "Get a robot that is not in the store",
			args: args{robotId: "abc", reqId: "abcd", room: robot.Room{X: 1, Y: 1}, coo: robot.Coordinate{X: 0, Y: 0}, d: 'N'},
			want: rsp{code: http.StatusNotFound},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest("GET", fmt.Sprintf("/robot/%s", tt.args.reqId), nil)
			if err != nil {
				t.Fatal(err)
			}
			req.SetPathValue("id", tt.args.reqId)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(robotHandler.getStatus)

			r, err := robot.NewRobot(tt.args.room, tt.args.d, tt.args.coo)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()

			robotStore.Put(tt.args.robotId, r, ctx)
			handler.ServeHTTP(rr, req)

			// Check the status code
			if rr.Result().StatusCode != tt.want.code {
				t.Errorf("wrong status code: got %v want %v", rr.Code, tt.want.code)
			}
		})
	}
}

func TestRobotHandler_command(t *testing.T) {

	robotStore := storage.NewRobotMemStore()
	robotHandler := RobotHandler{store: robotStore}

	type args struct {
		robotId string
		reqId   string
		room    robot.Room
		coo     robot.Coordinate
		d       rune
		cmd     string
	}

	type rsp struct {
		code   int
		status rspStatus
	}
	tests := []struct {
		name string
		args args
		want rsp
	}{
		// Make sure all test cases has a uniq id for the robot (robotId). All test cases share the same robot store.
		{
			name: "Command a robot that is in the store",
			args: args{robotId: "abc", reqId: "abc", room: robot.Room{X: 5, Y: 5}, coo: robot.Coordinate{X: 1, Y: 2}, d: 'N', cmd: "RFRFFRFRF"},
			want: rsp{code: http.StatusOK, status: rspStatus{Direction: "N", X: 1, Y: 3, Id: "abc"}},
		},
		{
			name: "Command a robot that is not in the store",
			args: args{robotId: "abc", reqId: "abcd", room: robot.Room{X: 5, Y: 5}, coo: robot.Coordinate{X: 1, Y: 2}, d: 'N', cmd: "RFRFFRFRF"},
			want: rsp{code: http.StatusNotFound, status: rspStatus{}},
		},
		{
			name: "Command a robot with an invalid command string",
			args: args{robotId: "abc", reqId: "abc", room: robot.Room{X: 5, Y: 5}, coo: robot.Coordinate{X: 1, Y: 2}, d: 'N', cmd: "RFRFFRFRFAFFFF"},
			want: rsp{code: http.StatusBadRequest, status: rspStatus{Direction: "N", X: 1, Y: 3, Id: "abc"}},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			c := reqCmd{Cmd: tt.args.cmd}
			payload, err := json.Marshal(c)
			if err != nil {
				t.Fatal(err)
			}
			b := strings.NewReader(string(payload))

			req, err := http.NewRequest("POST", fmt.Sprintf("/robot/%s", tt.args.reqId), b)
			if err != nil {
				t.Fatal(err)
			}
			req.SetPathValue("id", tt.args.reqId)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(robotHandler.command)

			r, err := robot.NewRobot(tt.args.room, tt.args.d, tt.args.coo)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.Background()

			robotStore.Put(tt.args.robotId, r, ctx)
			handler.ServeHTTP(rr, req)

			// Check the status code
			if rr.Result().StatusCode != tt.want.code {
				t.Errorf("wrong status code: got %v want %v", rr.Code, tt.want.code)
			}

			bs := rr.Body.Bytes()

			rsp := rspStatus{}
			json.Unmarshal(bs, &rsp)

			if !reflect.DeepEqual(rsp, tt.want.status) {
				t.Error("Response body did not match the expected value(s)")
			}
		})
	}
}

// This can be used together with pprof as a quick an dirty way to find any general perf issues.
func BenchmarkRobotHandler_create(b *testing.B) {

	voidStore := &robotVoidStore{}
	robotHandler := RobotHandler{store: voidStore}

	payload, err := json.Marshal(reqCreate{Direction: "N", Room: robot.Room{X: 1, Y: 1}, Start: robot.Coordinate{X: 0, Y: 0}})
	if err != nil {
		b.Fatal(err)
	}

	p := strings.NewReader(string(payload))

	req, err := http.NewRequest("POST", "/robot", p)
	if err != nil {
		b.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(robotHandler.create)

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			b.Error("Non 200 response.")
		}
	}
}
