package leaderboard

import (
	"math"
	"math/bits"
	"math/rand"
)

const (
	MAX_LEVEL = 16
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
	rand     *rand.Rand
}

func (list *RankList) getRandomLevel() int {
	length := list.length
	maxLevel := list.maxLevel

	if maxLevel <= 1 {
		return 1
	}

	// Use log_2(list.length) * 4 as estimated max level.
	estimated := bits.Len(uint(length)) * 4

	if estimated > maxLevel {
		estimated = maxLevel
	}

	// prob should be a bit more than 1/2 chance to make level balanced.
	// This value is selected by some random Get tests on a random generated list.
	// The magic number 3/8 works best when list size < 10,000,000.
	const prob = math.MaxInt32 * 3 / 8

	for i := 1; i < estimated; i++ {
		if list.rand.Int31() < prob {
			return i
		}
	}

	return maxLevel - 1
}

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
				(scoreEqualTo(x.Next(i), score) && x.Next(i).id < id)) {
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

	for i := 0; i < list.maxLevel; i++ {
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

func (list *RankList) GetRank(id string, score float64) uint64 {
	x := list.head
	var traversed uint64 = 0
	for i := list.maxLevel - 1; i >= 0; i-- {
		for notNil(x.Next(i)) &&
			(scoreLessThan(x.Next(i), score) ||
				(scoreEqualTo(x.Next(i), score) && x.Next(i).id < id)) {
			traversed += x.Span(i)
			x = x.Next(i)
		}
		if scoreEqualTo(x, score) && x.id == id {
			return traversed
		}
	}
	return traversed
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

func MakeRankList() *RankList {
	list := &RankList{}
	list.length = 0
	list.maxLevel = 1
	list.head = makeNode("_", 0, MAX_LEVEL)
	for i := 0; i < MAX_LEVEL; i++ {
		list.head.levels[i].span = 0
		list.head.levels[i].next = nil
	}
	list.tail = nil
	return list
}
