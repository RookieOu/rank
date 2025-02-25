package web

import (
	"log"
	"net/http"
)

// HandlerFunc 定义处理请求的函数类型
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Middleware 定义中间件类型
type Middleware func(next HandlerFunc) HandlerFunc

// Router 路由结构
type Router struct {
	routes map[string]map[string]HandlerFunc
}

// NewRouter 创建一个新的 Router 实例
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]HandlerFunc),
	}
}

// Handle 注册路由和请求方法的处理函数
func (r *Router) Handle(method, path string, handler HandlerFunc) {
	if _, exists := r.routes[path]; !exists {
		r.routes[path] = make(map[string]HandlerFunc)
	}
	r.routes[path][method] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path

	if methods, exists := r.routes[path]; exists {
		if handler, ok := methods[method]; ok {
			log.Default().Printf("path: %s, method: %s\n", path, method)
			handler(w, req)
			return
		}
	}

	http.NotFound(w, req)
}

func (r *Router) Init() {
	// 注册路由和处理器
	r.Handle("POST", "/updateScore", UpdateScoreHandler)
	r.Handle("POST", "/getPlayerRank", GetPlayerRankHandler)
	r.Handle("POST", "/getTopN", GetTopNHandler)
	r.Handle("POST", "/getPlayerRange", GetPlayerRangeHandler)
}
