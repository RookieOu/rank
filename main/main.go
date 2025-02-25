package main

import (
	"log"
	"net/http"
	"rank/module"
	"rank/web"
	"sync"
)

func main() {
	// http
	router := web.NewRouter()
	router.Init()
	log.Println("Server started on http://localhost:7777")
	w := &sync.WaitGroup{}
	go func() {
		w.Add(1)
		defer w.Done()
		if err := http.ListenAndServe(":7777", router); err != nil {
			log.Fatal(err)
		}
	}()

	rankInstance := module.GetInstance()
	rankInstance.Start(w)
	defer rankInstance.Stop()
	w.Wait()
}
