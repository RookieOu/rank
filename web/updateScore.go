package web

import (
	"fmt"
	"net/http"
	"rank/module"
	"strconv"
	"time"
)

func UpdateScoreHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(fmt.Errorf("parse form failed, err: %v", err))
	}
	var id, name string
	var scoreStr string
	for k, v := range r.Form {
		if k == "id" {
			id = v[0]
		} else if k == "score" {
			scoreStr = v[0]
		} else if k == "name" {
			name = v[0]
		}
	}
	score, err := strconv.Atoi(scoreStr)
	if id == "" || err != nil {
		w.Write([]byte("Invalid request"))
		return
	}
	// 更新分数
	module.UpdateScore(id, name, int32(score), time.Now().Unix())
	w.Write([]byte("Success"))
}
