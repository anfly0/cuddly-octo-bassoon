package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anfly0/cuddly-octo-bassoon/internal/robot"
	"github.com/anfly0/cuddly-octo-bassoon/internal/storage"
	"github.com/anfly0/cuddly-octo-bassoon/internal/utils"
)

type rspStatus struct {
	Direction string `json:"direction"`
	X         uint   `json:"x"`
	Y         uint   `json:"y"`
	Id        string `json:"id"`
}

func (rs *rspStatus) fromRobot(r *robot.Robot, id string) {
	d, coo := r.Report()

	rs.Direction = string(d)
	rs.X = coo.X
	rs.Y = coo.Y
	rs.Id = id
}

type reqCreate struct {
	Direction string           `json:"direction"`
	Room      robot.Room       `json:"room"`
	Start     robot.Coordinate `json:"start"`
}

type reqCmd struct {
	Cmd string `json:"cmd"`
}

type RobotHandler struct {
	store storage.RobotStore
}

func (rh *RobotHandler) command(w http.ResponseWriter, r *http.Request) {
	req := reqCmd{}

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	rb := rh.store.Get(id)

	if rb == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = rb.Cmd(req.Cmd)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rsp := &rspStatus{}
	rsp.fromRobot(rb, id)
	j, _ := json.Marshal(rsp)
	fmt.Fprint(w, string(j))
}

func (rh *RobotHandler) create(w http.ResponseWriter, r *http.Request) {

	req := reqCreate{Direction: "N"}

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil || len(req.Direction) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rb := &robot.Robot{Room: req.Room, Compass: *robot.NewCompass(rune(req.Direction[0])), Coordinate: req.Start}
	id, err := utils.RandId(4)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rh.store.Put(id, rb)

	rsp := &rspStatus{}

	rsp.fromRobot(rb, id)

	j, _ := json.Marshal(rsp)
	fmt.Fprint(w, string(j))
}

func (rh *RobotHandler) getStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	rb := rh.store.Get(id)

	if rb == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rsp := &rspStatus{}

	rsp.fromRobot(rb, id)

	j, _ := json.Marshal(rsp)
	fmt.Fprint(w, string(j))
}

func main() {

	rh := RobotHandler{store: storage.NewRobotMemStore()}

	http.Handle("POST /robot", Chain(http.HandlerFunc(rh.create), Logging, ContentHeader))
	http.Handle("GET /robot/{id}", Chain(http.HandlerFunc(rh.getStatus), Logging, ContentHeader))
	http.Handle("POST /robot/{id}", Chain(http.HandlerFunc(rh.command), Logging, ContentHeader))

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
