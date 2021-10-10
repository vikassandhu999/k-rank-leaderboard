package leaderboard

type LookupTable struct {
	ht map[string]float64
}

func (lookupt *LookupTable) Get(id string) (float64, bool) {
	if lookupt != nil && lookupt.ht != nil {
		val, ok := lookupt.ht[id]
		return val, ok
	}
	return 0, false
}

func (lookupt *LookupTable) Set(id string, score float64) bool {
	if lookupt != nil && lookupt.ht != nil {
		lookupt.ht[id] = score
		return true
	}
	return false
}

func (lookupt *LookupTable) Update(id string, score float64) bool {
	if lookupt != nil && lookupt.ht != nil {
		lookupt.ht[id] = score
		return true
	}
	return false
}

func (lookupt *LookupTable) Delete(id string, score float64) bool {
	if lookupt != nil && lookupt.ht != nil {
		delete(lookupt.ht, id)
		return true
	}
	return false
}

func CreateLookupTable() *LookupTable {
	return &LookupTable{ht: make(map[string]float64)}
}
