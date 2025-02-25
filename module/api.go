package module

import (
	"context"
	"fmt"
	"time"
)

// 玩家更新积分
func UpdateScore(id, name string, score int32, timeUnix int64) {
	cmd := &UpdateScoreCmd{
		id:       id,
		score:    score,
		name:     name,
		timeUnix: timeUnix,
	}
	err := GetInstance().sendCmd(cmd)
	if err != nil {
		fmt.Println(err)
	}
}

// 获取玩家当前排名
func GetPlayerRank(id string) *PlayerNode {
	cmd := &GetPlayerRankCmd{
		id:  id,
		ret: make(chan *PlayerNode),
	}
	err := GetInstance().sendCmd(cmd)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var ret *PlayerNode
	select {
	case ret = <-cmd.ret:
	case <-ctx.Done():
		fmt.Println("GetPlayerRank timeout")
		return nil
	}

	return ret
}

// 获取排行榜前N名
func GetTopN(n int32) []*PlayerNode {
	cmd := &GetTopNCmd{
		n:   n,
		ret: make(chan []*PlayerNode),
	}
	err := GetInstance().sendCmd(cmd)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var ret []*PlayerNode
	select {
	case ret = <-cmd.ret:
	case <-ctx.Done():
		fmt.Println("GetTopN timeout")
		return nil
	}

	return ret
}

// 获取玩家周边排名
func GetPlayerRankRange(id string, number int32) []*PlayerNode {
	cmd := &GetPlayerRankRangeCmd{
		id:  id,
		num: number,
		ret: make(chan []*PlayerNode),
	}

	err := GetInstance().sendCmd(cmd)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var ret []*PlayerNode
	select {
	case ret = <-cmd.ret:
	case <-ctx.Done():
		fmt.Println("GetPlayerRankRange timeout")
	}
	return ret
}
