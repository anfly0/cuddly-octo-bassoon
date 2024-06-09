package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anfly0/cuddly-octo-bassoon/internal/robot"
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

			println(string(payload))

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
