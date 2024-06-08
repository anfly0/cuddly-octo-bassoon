package storage

import (
	"context"
	"sync"

	"github.com/anfly0/cuddly-octo-bassoon/internal/robot"
)

type RobotStore interface {
	Get(id string, ctx context.Context) *robot.Robot
	Put(id string, r *robot.Robot, ctx context.Context) error
}

type RobotMemeStore struct {
	m map[string]*robot.Robot
	l sync.RWMutex
}

func NewRobotMemStore() *RobotMemeStore {
	m := make(map[string]*robot.Robot)
	return &RobotMemeStore{m: m, l: sync.RWMutex{}}
}

// In this implementation we ignore the context. But in a robot store where a cancellations makes sense it is useful.
func (rs *RobotMemeStore) Get(id string, _ context.Context) *robot.Robot {
	rs.l.RLock()
	defer rs.l.RUnlock()

	return rs.m[id]
}

// In this implementation we ignore the context. But in a robot store where a cancellations makes sense it is useful.
func (rs *RobotMemeStore) Put(id string, r *robot.Robot, _ context.Context) error {
	rs.l.Lock()
	defer rs.l.Unlock()

	rs.m[id] = r
	return nil
}
