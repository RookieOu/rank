package module

type icmd interface {
	run(mgr *Rank)
}

type UpdateScoreCmd struct {
	id       string
	score    int32
	name     string
	timeUnix int64
}

func (c *UpdateScoreCmd) run(mgr *Rank) {
	newNode := &PlayerNode{
		Uid:       c.id,
		RankScore: c.score,
		Name:      c.name,
		TimeUnix:  c.timeUnix,
	}
	mgr.cmdList = append(mgr.cmdList, newNode)
}

type GetPlayerRankCmd struct {
	id  string
	ret chan *PlayerNode
}

func (c *GetPlayerRankCmd) run(mgr *Rank) {
	ret := mgr.GetPlayerRank(c.id)
	c.ret <- ret
}

type GetTopNCmd struct {
	n   int32
	ret chan []*PlayerNode
}

func (c *GetTopNCmd) run(mgr *Rank) {
	ret := mgr.GetTopN(c.n)
	c.ret <- ret
}

type GetPlayerRankRangeCmd struct {
	id  string
	num int32
	ret chan []*PlayerNode
}

func (c *GetPlayerRankRangeCmd) run(mgr *Rank) {
	var playerRank = mgr.id2RankIndex[c.id]
	if playerRank == 0 {
		c.ret <- nil
		return
	}

	// 计算范围，避免越界
	start := playerRank - 1 - int(c.num)
	if start < 0 {
		start = 0
	}

	end := playerRank - 1 + int(c.num)
	if end >= len(mgr.rankSlice) {
		end = len(mgr.rankSlice) - 1
	}

	// 获取玩家周边排名
	ret := make([]*PlayerNode, 0, end-start+1)
	for i := start; i <= end; i++ {
		nodeClone := *mgr.rankSlice[i]
		ret = append(ret, &nodeClone)
	}
	c.ret <- ret
}
