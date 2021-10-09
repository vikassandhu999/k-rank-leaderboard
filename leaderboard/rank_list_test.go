package leaderboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRankListDefined(t *testing.T) {
	rankList := MakeRankList()
	assert.NotNil(t, rankList)
}

func TestInsert(t *testing.T) {
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

func TestGetRank(t *testing.T) {
	type Test struct {
		id            string
		score         float64
		expected_rank uint64
	}
	rankList := MakeRankList()
	entries := make([]Test, 6)
	entries[0] = Test{"1", 100, 1}
	entries[1] = Test{"2", 101, 2}
	entries[2] = Test{"3", 110, 5}
	entries[3] = Test{"4", 116, 6}
	entries[4] = Test{"5", 108, 4}
	entries[5] = Test{"6", 104, 3}

	for _, e := range entries {
		rankList.Insert(e.id, e.score)
	}
	for _, e := range entries {
		r := rankList.GetRank(e.id, e.score)
		assert.Equal(t, r, e.expected_rank)
	}
}

func TestShouldGetRankById(t *testing.T) {
	type Data struct {
		id    string
		score float64
	}
	rankList := MakeRankList()
	entries := make([]Data, 6)
	entries[0] = Data{"1", 100}
	entries[1] = Data{"2", 101}
	entries[2] = Data{"3", 110}
	entries[3] = Data{"4", 116}
	entries[4] = Data{"5", 108}
	entries[5] = Data{"6", 104}

	for _, e := range entries {
		rankList.Insert(e.id, e.score)
	}

	type Test struct {
		rank        uint64
		expected_id string
	}
	tests := make([]Test, 6)
	tests[0] = Test{1, "1"}
	tests[1] = Test{2, "2"}
	tests[2] = Test{3, "6"}
	tests[3] = Test{4, "5"}
	tests[4] = Test{5, "3"}
	tests[5] = Test{6, "4"}

	for _, e := range tests {
		r := rankList.GetIdByRank(e.rank)
		assert.Equal(t, r, e.expected_id)
	}
}

func TestDelete(t *testing.T) {
	type Data struct {
		id    string
		score float64
	}
	rankList := MakeRankList()
	entries := make([]Data, 6)
	entries[0] = Data{"1", 100}
	entries[1] = Data{"2", 101}
	entries[2] = Data{"3", 110}
	entries[3] = Data{"4", 116}
	entries[4] = Data{"5", 108}
	entries[5] = Data{"6", 104}

	for _, e := range entries {
		rankList.Insert(e.id, e.score)
	}

	assert.Equal(t, rankList.Delete("1", 100), true)
	assert.Equal(t, rankList.Delete("1000", 100), false)
	assert.Equal(t, rankList.Delete("5", 5454), false)

	assert.Equal(t, rankList.GetRank("1", 1000), uint64(0))
	assert.Equal(t, rankList.length, uint64(5))

	assert.Equal(t, rankList.Delete("2", 101), true)

	assert.Equal(t, rankList.GetRank("2", 101), uint64(0))
	assert.Equal(t, rankList.length, uint64(4))
}

func TestUpdateScore(t *testing.T) {
	type Data struct {
		id    string
		score float64
	}
	rankList := MakeRankList()
	entries := make([]Data, 6)
	entries[0] = Data{"1", 100}
	entries[1] = Data{"2", 101}
	entries[2] = Data{"3", 110}
	entries[3] = Data{"4", 116}
	entries[4] = Data{"5", 108}
	entries[5] = Data{"6", 104}

	for _, e := range entries {
		rankList.Insert(e.id, e.score)
	}

	assert.NotNil(t, rankList.UpdateScore("1", 100, 1666))
	assert.NotNil(t, rankList.UpdateScore("4", 116, 99))

	assert.Equal(t, uint64(6), rankList.GetRank("1", 1666))
	assert.Equal(t, uint64(1), rankList.GetRank("4", 99))
	assert.Equal(t, uint64(5), rankList.GetRank("3", 110))

}
