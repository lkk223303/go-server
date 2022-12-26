package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	DefaultQeueSize = 100
	// DefaultGoNum    = 10
)

type WebHandler struct {
	q          chan *WebCmd
	handler    HandlerFunc
	handlerNum int
}
type HandlerFunc func(*WebCmd)

type OutConf struct {
	Offset    int64 `json:"offset"`    //
	Limit     int64 `json:"limit"`     //
	Retention int64 `json:"retention"` //
}

// You can custom design your WebCmd for general server request
type WebCmd struct {
	AppName    string   `json:"app"`      // request name
	InputArgs  []string `json:"input"`    // input arguments
	OutputArgs *OutConf `json:"output"`   // output arguments
	Ticket     string   `json:"ticket"`   // count for the request
	Operator   string   `json:"operator"` // service operator like  admin
	Fanout     FanOut
}
type FanOut struct {
	Write chan string
	Read  chan string
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

func (h *WebHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {

	var cmd WebCmd
	jsonDec := json.NewDecoder(r.Body)
	if err := jsonDec.Decode(&cmd); err != nil {
		http.Error(w, "Invalid request content", http.StatusBadRequest)
		return
	}
	name := cmd.AppName

	if cmd.InputArgs[0] == "0" {
		rc := make(chan string, 1)
		rc <- "Fanout Reading"
		cmd.Fanout.Read = rc
	} else if cmd.InputArgs[0] == "1" {
		wc := make(chan string, 1)
		wc <- "Fanout writing"
		cmd.Fanout.Write = wc
	}

	log.Println("Name : ", name)
	// fmt.Fprint(w, cmd.AppName)
	w.Write([]byte(cmd.AppName + "\n"))

	h.q <- &cmd
}

func (h *WebHandler) PrintHello(w http.ResponseWriter, r *http.Request) {

	var cmd WebCmd
	var rsp string
	jsonDec := json.NewDecoder(r.Body)
	if err := jsonDec.Decode(&cmd); err != nil {
		http.Error(w, "Invalid request content", http.StatusBadRequest)
		return
	}
	name := cmd.AppName

	log.Println("Hello ", name)
	rsp = fmt.Sprintf("Hello %s %s\n", r.Method, name)

	w.Write([]byte(rsp))
	h.q <- &cmd
}
