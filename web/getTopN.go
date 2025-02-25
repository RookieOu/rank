package web

import (
	"fmt"
	"net/http"
	"rank/module"
	"strconv"
	"strings"
)

func GetTopNHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(fmt.Errorf("parse form failed, err: %v", err))
	}
	var numStr string
	for k, v := range r.Form {
		if k == "num" {
			numStr = v[0]
		}
	}
	num, err := strconv.Atoi(numStr)
	if err != nil {
		w.Write([]byte("Invalid request"))
		return
	}
	nodes := module.GetTopN(int32(num))
	retStr := strings.Builder{}
	for _, node := range nodes {
		retStr.WriteString(fmt.Sprintf("Player id: %s, name: %s, rank: %d, score: %d\n", node.Uid, node.Name, node.RankNum, node.RankScore))
	}
	w.Write([]byte(retStr.String()))
}
