package main

import (
	"go-server/handler"
	"log"
	"net/http"
)

const (
	DefaultQeueSize = 100
	DefaultGoNum    = 1000
)

func main() {
	cmdHandler := handler.NewWebHandler(DefaultGoNum, CmdHandleFunc)
	cmdHandler.Serve()

	http.HandleFunc("/test", cmdHandler.HandleRequest)
	http.HandleFunc("/hello", cmdHandler.PrintHello)

	log.Println("Server listening on port 8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

// CmdHandleFunc handle function to define how to handle request(cmd)
func CmdHandleFunc(cmd *handler.WebCmd) {
	// if cmd.InputArgs[0] == "0" {
	// 	log.Println("Hello ", cmd.AppName, " from handler func")
	// } else if cmd.InputArgs[0] == "1" {
	// 	log.Println("Already said Hi ", cmd.AppName, " from handler func")
	// }

	for {
		select {
		case r := <-cmd.Fanout.Read:
			go ReadFunc(r)
		case w := <-cmd.Fanout.Write:
			go WriteFunc(w)
		}
	}
}

func ReadFunc(r string) {

	log.Println("reading ...", r)
}

func WriteFunc(w string) {
	log.Println("writing ...", w)
}
