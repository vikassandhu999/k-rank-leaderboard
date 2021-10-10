package main

import (
	ldb "github.com/vikassandhu999/k-rank-leaderboard/leaderboard"
)

var leaderboards map[string]ldb.LeaderboardI = make(map[string]ldb.LeaderboardI)

type CommanderI interface {
	CreateLDB(name string) bool
	DeleteLDB(name string)
	Set(name string, id string, score float64) bool
	GetRank(name string, id string) uint64
	Delete(name string, id string) bool
	GetByRank(name string, rank uint64) (string, float64)
}

type Commander struct{}

func (cmd Commander) CreateLDB(name string) bool {
	_, present := leaderboards[name]
	if present {
		return false
	}
	leaderboards[name] = ldb.CreateLDB()
	return true
}

func (cmd Commander) DeleteLDB(name string) {
	delete(leaderboards, name)
}

func (cmd Commander) Set(name string, id string, score float64) bool {
	fLDB, ok := leaderboards[name]
	if !ok {
		return false
	}
	return fLDB.Insert(id, score)
}

func (cmd Commander) GetRank(name string, id string) uint64 {
	fLDB, ok := leaderboards[name]
	if !ok {
		return 0
	}
	return fLDB.GetRank(id)
}

func (cmd Commander) Delete(name string, id string) bool {
	fLDB, ok := leaderboards[name]
	if !ok {
		return false
	}
	return fLDB.Delete(id)
}

func (cmd Commander) GetByRank(name string, rank uint64) (string, float64) {
	fLDB, ok := leaderboards[name]
	if !ok {
		return ldb.NULL_ID, 0
	}
	return fLDB.GetByRank(rank)
}
