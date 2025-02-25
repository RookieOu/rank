package web

import (
	"fmt"
	"net/http"
	"rank/module"
)

func GetPlayerRankHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(fmt.Errorf("parse form failed, err: %v", err))
	}
	var id string
	for k, v := range r.Form {
		if k == "id" {
			id = v[0]
		}
	}
	if id == "" {
		w.Write([]byte("Invalid request"))
		return
	}
	// 更新分数
	node := module.GetPlayerRank(id)
	w.Write([]byte(fmt.Sprintf("Player rank: %d", node.RankNum)))
}
