package module

var instance *Rank

func GetInstance() *Rank {
	if instance == nil {
		instance = NewRank()
	}
	return instance
}
