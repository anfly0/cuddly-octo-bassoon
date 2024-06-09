package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
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

func RspStatusFromRobot(r *robot.Robot, id string) rspStatus {
	d, coo := r.Report()
	return rspStatus{Direction: string(d), X: coo.X, Y: coo.Y, Id: id}
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

	if err != nil || req.Cmd == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	rb := rh.store.Get(id, r.Context())

	if rb == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	d, coo, err := rb.Cmd(req.Cmd)

	rsp := rspStatus{Direction: string(d), X: coo.X, Y: coo.Y, Id: id}
	j, _ := json.Marshal(rsp)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	io.WriteString(w, string(j))
}

func (rh *RobotHandler) create(w http.ResponseWriter, r *http.Request) {

	req := reqCreate{Direction: "N"}

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil || len(req.Direction) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	rb, err := robot.NewRobot(req.Room, rune(req.Direction[0]), req.Start)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	id, err := utils.RandId(4)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rh.store.Put(id, rb, r.Context())

	rsp := RspStatusFromRobot(rb, id)

	j, _ := json.Marshal(rsp)
	io.WriteString(w, string(j))
}

func (rh *RobotHandler) getStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	rb := rh.store.Get(id, r.Context())

	if rb == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rsp := RspStatusFromRobot(rb, id)

	j, _ := json.Marshal(rsp)
	io.WriteString(w, string(j))
}

func main() {

	addr := flag.String("addr", "", "Ip address the server will listen to")
	port := flag.String("port", "8080", "Port number the server will listen to")
	flag.Parse()
	rh := RobotHandler{store: storage.NewRobotMemStore()}

	http.Handle("POST /robot", Chain(http.HandlerFunc(rh.create), Logging, ContentHeader))
	http.Handle("GET /robot/{id}", Chain(http.HandlerFunc(rh.getStatus), Logging, ContentHeader))
	http.Handle("POST /robot/{id}", Chain(http.HandlerFunc(rh.command), Logging, ContentHeader))

	s_addr := fmt.Sprintf("%s:%s", *addr, *port)
	fmt.Printf("Starting server on %s!", s_addr)
	if err := http.ListenAndServe(s_addr, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
