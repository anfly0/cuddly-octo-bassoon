package storage

import (
	"sync"

	"github.com/anfly0/cuddly-octo-bassoon/internal/robot"
)

type RobotStore interface {
	Get(id string) *robot.Robot
	Put(id string, r *robot.Robot) error
}

type RobotMemeStore struct {
	m map[string]*robot.Robot
	l sync.RWMutex
}

func NewRobotMemStore() *RobotMemeStore {
	m := make(map[string]*robot.Robot)
	return &RobotMemeStore{m: m, l: sync.RWMutex{}}
}

func (rs *RobotMemeStore) Get(id string) *robot.Robot {
	rs.l.RLock()
	defer rs.l.RUnlock()

	return rs.m[id]
}

func (rs *RobotMemeStore) Put(id string, r *robot.Robot) error {
	rs.l.Lock()
	defer rs.l.Unlock()

	rs.m[id] = r
	return nil
}
