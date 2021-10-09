package leaderboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRankListDefined(t *testing.T) {
	rankList := MakeRankList()
	assert.NotNil(t, rankList, "Rank list is nil")
}

func TestShouldInsertItem(t *testing.T) {
	type Test struct {
		id    string
		score float64
	}
	rankList := MakeRankList()
	entries := make([]Test, 6)
	entries[0] = Test{"1", 100}
	entries[1] = Test{"2", 101}
	entries[2] = Test{"3", 110}
	entries[3] = Test{"4", 116}
	entries[4] = Test{"5", 108}
	entries[5] = Test{"6", 104}

	for _, e := range entries {
		r := rankList.Insert(e.id, e.score)
		assert.NotNil(t, r)
		assert.Equal(t, r.id, e.id)
		assert.Equal(t, r.score, e.score)
	}
}

func TestShouldReturnCorrectRank(t *testing.T) {
	type Test struct {
		id            string
		score         float64
		expected_rank uint64
	}
	rankList := MakeRankList()
	entries := make([]Test, 6)
	entries[0] = Test{"1", 100, 0}
	entries[1] = Test{"2", 101, 1}
	entries[2] = Test{"3", 110, 4}
	entries[3] = Test{"4", 116, 5}
	entries[4] = Test{"5", 108, 3}
	entries[5] = Test{"6", 104, 2}

	for _, e := range entries {
		rankList.Insert(e.id, e.score)
	}
	for _, e := range entries {
		r := rankList.GetRank(e.id, e.score)
		assert.Equal(t, r, e.expected_rank)
	}
}
