package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	DefaultQeueSize = 100
	DefaultGoNum    = 10
)

type WebHandler struct {
	q          chan *WebCmd
	handler    HandlerFunc
	handlerNum int
}
type HandlerFunc func(*WebCmd)

// You can custom design your WebCmd for general server request
type WebCmd struct {
	AppName    string   `json:"app"`      // request name
	InputArgs  []string `json:"input"`    // input arguments
	OutputArgs *OutConf `json:"output"`   // output arguments
	Ticket     string   `json:"ticket"`   // count for the request
	Operator   string   `json:"operator"` // service operator like  admin
}
type OutConf struct {
	Offset    int64 `json:"offset"`    //
	Limit     int64 `json:"limit"`     //
	Retention int64 `json:"retention"` //
}

func main() {
	cmdHandler := NewWebHandler(DefaultGoNum, CmdHandleFunc)
	cmdHandler.Serve()

	http.HandleFunc("/test", cmdHandler.HandleRequest)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func NewWebHandler(num int, handler HandlerFunc) *WebHandler {
	webHandler := &WebHandler{}
	webHandler.Init(num, handler)
	return webHandler
}
func (h *WebHandler) Init(num int, handler HandlerFunc) {
	h.handlerNum = num
	h.q = make(chan *WebCmd, DefaultQeueSize)
	h.handler = handler
}

func (h *WebHandler) Serve() {
	// default 10 workers to handle 100 request
	for i := 0; i < h.handlerNum; i++ {
		go h.worker()
	}
}
func (h *WebHandler) worker() {
	for req := range h.q {
		h.handler(req)
	}
}

// CmdHandleFunc handle function to define how to handle request(cmd)
func CmdHandleFunc(cmd *WebCmd) {
	log.Println(rand.Int63n(time.Now().UnixNano()))
	log.Println("cmd handled", cmd)
}

func (h *WebHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {

	var cmd WebCmd
	jsonDec := json.NewDecoder(r.Body)
	if err := jsonDec.Decode(&cmd); err != nil {
		http.Error(w, "Invalid request content", http.StatusBadRequest)
		return
	}
	name := cmd.AppName
	log.Println("Name : ", name)
	fmt.Fprint(w, cmd.AppName)
	w.Write([]byte(cmd.AppName))
	h.q <- &cmd
}
