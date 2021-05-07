package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sqlTime/cmd"
)

func main() {
	srv := &http.Server{Addr: ":0"}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
	return
}
