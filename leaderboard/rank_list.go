package leaderboard

import (
	"math/bits"
	"math/rand"
)

const (
	MAX_LEVEL = 16
	NULL_ID   = "NULL_ID"
)

type listLevel struct {
	span uint64
	next *RankListNode
}

type RankListNode struct {
	id     string
	score  float64
	prev   *RankListNode
	levels []listLevel
}

func (node *RankListNode) Next(i int) *RankListNode {
	if notNil(node) {
		return node.levels[i].next
	}
	return nil
}

func (node *RankListNode) Span(i int) uint64 {
	if notNil(node) {
		return node.levels[i].span
	}
	return 0
}

func makeNode(id string, score float64, level int) *RankListNode {
	return &RankListNode{
		id:     id,
		score:  score,
		prev:   nil,
		levels: make([]listLevel, level),
	}
}

type RankList struct {
	head     *RankListNode
	tail     *RankListNode
	length   uint64
	maxLevel int
}

func (list *RankList) getRandomLevel() int {
	// Use log_2(list.length) * 4 as estimated max level.
	estimated := bits.Len(uint(list.length)) * 4
	// Increase level with probability 1 in kBranching
	var kBranching uint64 = 4
	if estimated > MAX_LEVEL {
		estimated = MAX_LEVEL
	}
	var level int = 1
	for level < MAX_LEVEL && (rand.Uint64()%kBranching == 0) {
		level++
	}
	if level > MAX_LEVEL {
		return MAX_LEVEL
	}
	return level
}

/* utility function used by Insert */
func (list *RankList) findLessThanEqual(
	id string,
	score float64,
	ranks []uint64,
	pathNodes []*RankListNode) *RankListNode {

	x := list.head
	for i := list.maxLevel - 1; i >= 0; i-- {
		if i == list.maxLevel-1 {
			ranks[i] = 0
		} else {
			ranks[i] = ranks[i+1]
		}
		for notNil(x.Next(i)) &&
			(scoreLessThan(x.Next(i), score) ||
				(scoreEqualTo(x.Next(i), score) && x.Next(i).id <= id)) {
			ranks[i] += x.Span(i)
			x = x.Next(i)
		}
		pathNodes[i] = x
	}
	return x
}

func (list *RankList) Insert(id string, score float64) *RankListNode {
	ranks := make([]uint64, MAX_LEVEL)
	pathNodes := make([]*RankListNode, MAX_LEVEL)

	x := list.findLessThanEqual(id, score, ranks, pathNodes)

	if scoreEqualTo(x, score) && x.id == id {
		panic("Duplicate nodes are not allowed")
	}

	level := list.getRandomLevel()
	node := makeNode(id, score, level)

	if list.maxLevel < level {
		for i := list.maxLevel; i < level; i++ {
			ranks[i] = 0
			pathNodes[i] = list.head
			pathNodes[i].levels[i].span = list.length
		}
		list.maxLevel = level
	}

	for i := 0; i < level; i++ {
		node.levels[i].next = pathNodes[i].levels[i].next
		pathNodes[i].levels[i].next = node

		node.levels[i].span = pathNodes[i].levels[i].span - (ranks[0] - ranks[i])
		pathNodes[i].levels[i].span = ranks[0] - ranks[i] + 1
	}

	for i := level; i < list.maxLevel; i++ {
		pathNodes[i].levels[i].span++
	}

	if pathNodes[0] == list.head {
		node.prev = nil
	} else {
		node.prev = pathNodes[0]
	}

	if node.levels[0].next != nil {
		node.levels[0].next.prev = node
	} else {
		list.tail = node
	}

	list.length++
	return node
}

/* utility function used by Delete, UpdateScore */
func (list *RankList) findLessThan(id string, score float64, pathNodes []*RankListNode) *RankListNode {
	x := list.head
	for i := list.maxLevel - 1; i >= 0; i-- {
		for notNil(x.Next(i)) &&
			(scoreLessThan(x.Next(i), score) ||
				(scoreEqualTo(x.Next(i), score) && x.Next(i).id < id)) {
			x = x.Next(i)
		}
		pathNodes[i] = x
	}
	return x
}

/* utility function used by Delete, UpdateScore */
func (list *RankList) deleteNode(x *RankListNode, pathNodes []*RankListNode) {
	for i := 0; i < list.maxLevel; i++ {
		if pathNodes[i].Next(i) == x {
			pathNodes[i].levels[i].next = x.levels[i].next
			pathNodes[i].levels[i].span += x.levels[i].span - 1
		} else {
			pathNodes[i].levels[i].span--
		}
	}

	if notNil(x.Next(0)) {
		x.Next(0).prev = x.prev
	} else {
		list.tail = x.prev
	}

	for list.maxLevel > 1 && list.head.Next(list.maxLevel-1) == nil {
		list.maxLevel--
	}
	list.length--
}

func (list *RankList) UpdateScore(id string, score float64, updatedScore float64) *RankListNode {
	pathNodes := make([]*RankListNode, MAX_LEVEL)
	x := list.findLessThan(id, score, pathNodes)
	if x == nil || !scoreEqualTo(x.Next(0), score) || x.Next(0).id != id {
		return nil
	}
	x = x.Next(0)
	if (x.prev == nil || x.prev.score < updatedScore) && (x.Next(0) == nil || x.Next(0).score > updatedScore) {
		x.score = updatedScore
		return x
	}
	list.deleteNode(x, pathNodes)
	return list.Insert(id, updatedScore)
}

func (list *RankList) Delete(id string, score float64) bool {
	pathNodes := make([]*RankListNode, MAX_LEVEL)
	x := list.findLessThan(id, score, pathNodes)
	if x == nil || !scoreEqualTo(x.Next(0), score) || x.Next(0).id != id {
		return false
	}

	x = x.Next(0)
	list.deleteNode(x, pathNodes)
	return true
}

/* return 1-based rank due to extera head node */
func (list *RankList) GetRank(id string, score float64) uint64 {
	x := list.head
	var traversed uint64 = 0
	for i := list.maxLevel - 1; i >= 0; i-- {
		for notNil(x.Next(i)) &&
			(scoreLessThan(x.Next(i), score) ||
				(scoreEqualTo(x.Next(i), score) && x.Next(i).id <= id)) {
			traversed += x.Span(i)
			x = x.Next(i)
		}
		if scoreEqualTo(x, score) && x.id == id {
			return traversed
		}
	}
	return 0
}

/* need 1-based rank if 0 provided returns default id */
func (list *RankList) GetIdByRank(rank uint64) string {
	if rank > list.length {
		return NULL_ID
	}
	x := list.head
	var traversed uint64 = 0
	for i := list.maxLevel - 1; i >= 0; i-- {
		for notNil(x.Next(i)) && ((traversed + x.Span(i)) <= (rank)) {
			traversed += x.Span(i)
			x = x.Next(i)
		}
		if traversed == rank {
			return x.id
		}
	}
	return NULL_ID
}

func scoreLessThan(x *RankListNode, score float64) bool {
	return x != nil && x.score < score
}

func scoreEqualTo(x *RankListNode, score float64) bool {
	return x != nil && x.score == score
}

func notNil(node *RankListNode) bool {
	return node != nil
}

/* Constructor for RankList */
func MakeRankList() *RankList {
	list := &RankList{}
	list.length = 0
	list.maxLevel = 1
	list.head = makeNode(NULL_ID, 0, MAX_LEVEL)
	for i := 0; i < MAX_LEVEL; i++ {
		list.head.levels[i].span = 0
		list.head.levels[i].next = nil
	}
	list.tail = nil
	return list
}
