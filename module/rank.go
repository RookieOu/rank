package module

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

type PlayerNode struct {
	Uid       string
	Name      string
	RankScore int32
	TimeUnix  int64
	RankNum   int // 排名
}

// Comparator PlayerNode的比较规则
// 1. 先比积分，积分高的排名靠前
// 2. 积分相同，比时间，时间早的排名靠前
// 3. 积分和时间都相同，比uid，uid小的排名靠前
func (d *PlayerNode) Comparator(k1, k2 *PlayerNode) int {
	if k1.RankScore > k2.RankScore {
		return 1
	} else if k1.RankScore < k2.RankScore {
		return -1
	}

	if k1.TimeUnix < k2.TimeUnix {
		return 1
	} else if k1.TimeUnix > k2.TimeUnix {
		return -1
	}

	if k1.Uid < k2.Uid {
		return 1
	} else if k1.Uid > k2.Uid {
		return -1
	}
	return 0
}

type Rank struct {
	rankSlice    []*PlayerNode
	id2RankIndex map[string]int

	cmdChan  chan icmd
	cmdList  []*PlayerNode
	quitChan chan struct{}
}

func NewRank() *Rank {
	return &Rank{
		rankSlice:    make([]*PlayerNode, 0, 2048),
		id2RankIndex: make(map[string]int, 2048),
		cmdList:      make([]*PlayerNode, 0, 512),
		quitChan:     make(chan struct{}),
		cmdChan:      make(chan icmd, 512),
	}
}

func (r *Rank) Start(w *sync.WaitGroup) {
	w.Add(1)
	go func() {
		defer w.Done()
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(fmt.Errorf("rank panic, err: %v", err))
			}
		}()
		// 5秒更新一次
		ticker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-r.quitChan:
				return
			case cmd := <-r.cmdChan:
				func() {
					defer func() {
						if err := recover(); err != nil {
							fmt.Println(fmt.Errorf("rank run panic, err: %v", err))
						}
					}()
					cmd.run(r)
				}()
			case <-ticker.C:
				func() {
					defer func() {
						if err := recover(); err != nil {
							fmt.Println(fmt.Errorf("rank updateRank panic, err: %v", err))
						}
					}()
					r.updateRank()
				}()
				ticker = time.NewTicker(time.Second * 5)
			}
		}
	}()

}

func (r *Rank) Stop() {
	close(r.quitChan)
}

func (r *Rank) sendCmd(cmd icmd) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	select {
	case r.cmdChan <- cmd:
	case <-ctx.Done():
		fmt.Println("sendCmd timeout")
		return fmt.Errorf("sendCmd timeout")
	}
	return nil
}

func (r *Rank) updateRank() {
	for _, c := range r.cmdList {
		// 新上榜玩家，直接append
		if r.id2RankIndex[c.Uid] == 0 {
			r.rankSlice = append(r.rankSlice, c)
			r.id2RankIndex[c.Uid] = len(r.rankSlice)
			continue
		}
		index := r.id2RankIndex[c.Uid] - 1
		// 旧玩家，更新积分
		oldNode := r.rankSlice[index]
		if oldNode.RankScore == c.RankScore {
			continue
		}
		if oldNode.TimeUnix >= c.TimeUnix {
			continue
		}
		oldNode.RankScore = c.RankScore
	}
	r.cmdList = r.cmdList[:0]
	sort.SliceStable(r.rankSlice, func(i, j int) bool {
		return r.rankSlice[i].Comparator(r.rankSlice[i], r.rankSlice[j]) == 1
	})
	r.id2RankIndex = make(map[string]int, len(r.rankSlice))
	rankNum := 1
	for i, v := range r.rankSlice {
		r.id2RankIndex[v.Uid] = i + 1
		// 当前玩家和前一个玩家分数相同
		if i > 0 && v.RankScore == r.rankSlice[i-1].RankScore {
			// 如果是平分玩家，排名不变
			v.RankNum = r.rankSlice[i-1].RankNum
		} else {
			v.RankNum = rankNum
			rankNum++
		}
	}
}

func (r *Rank) GetPlayerRank(id string) *PlayerNode {
	index := r.id2RankIndex[id]
	if index == 0 {
		return nil
	}
	node := r.rankSlice[index-1]
	nodeClone := *node
	return &nodeClone
}

func (r *Rank) GetTopN(n int32) []*PlayerNode {
	if n <= 0 {
		return nil
	}
	if n >= int32(len(r.rankSlice)) {
		n = int32(len(r.rankSlice))
	}
	ret := make([]*PlayerNode, n)
	for i := int32(0); i < n; i++ {
		nodeClone := *r.rankSlice[i]
		ret[i] = &nodeClone
	}
	return ret
}
