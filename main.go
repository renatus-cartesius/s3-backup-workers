package main

import (
	"backup-workers/cmd"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		http.ListenAndServe("localhost:8080", nil)
	}()
	cmd.Execute()
}
