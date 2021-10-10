package leaderboard

type LeaderboardI interface {
	Insert(id string, score float64) bool
	GetRank(id string) uint64
	GetByRank(rank uint64) (string, float64)
	Delete(id string) bool
	UpdateScore(id string, updatedScore float64) bool
}

type Leaderboard struct {
	rlist   *RankList
	lookupt *LookupTable
}

func CreateLDB() LeaderboardI {
	return &Leaderboard{
		rlist:   CreateRankList(),
		lookupt: CreateLookupTable(),
	}
}

func (ldb Leaderboard) Insert(id string, score float64) bool {
	presentScore, ok := ldb.lookupt.Get(id)

	if ok && presentScore == score {
		return false
	}

	if ok {
		ldb.lookupt.Update(id, score)
		ldb.rlist.UpdateScore(id, presentScore, score)
		return true
	}

	ldb.lookupt.Set(id, score)
	ldb.rlist.Insert(id, score)
	return true
}

func (ldb Leaderboard) GetRank(id string) uint64 {
	score, ok := ldb.lookupt.Get(id)

	if !ok {
		return 0
	}

	return ldb.rlist.GetRank(id, score)
}

func (ldb Leaderboard) GetByRank(rank uint64) (string, float64) {
	id := ldb.rlist.GetIdByRank(rank)

	if id == NULL_ID {
		return NULL_ID, 0
	}

	score, ok := ldb.lookupt.Get(id)

	if !ok {
		return NULL_ID, 0
	}

	return id, score
}

func (ldb Leaderboard) Delete(id string) bool {
	score, ok := ldb.lookupt.Get(id)
	if !ok {
		return false
	}
	return ldb.rlist.Delete(id, score)
}

func (ldb Leaderboard) UpdateScore(id string, updatedScore float64) bool {
	score, ok := ldb.lookupt.Get(id)
	if !ok {
		return false
	}
	return ldb.rlist.UpdateScore(id, score, updatedScore) != nil
}
