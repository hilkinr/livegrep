package server

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/nelhage/livegrep/client"
	"net/http"
)

type innerError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type replyError struct {
	Err innerError `json:"error"`
}

func replyJson(w http.ResponseWriter, status int, obj interface{}) {
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	if err := enc.Encode(obj); err != nil {
		glog.Warningf("writing http response, data=%s err=%s",
			asJSON{obj},
			err.Error())
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	replyJson(w, status, &replyError{Err: innerError{code, message}})
}

func parseQuery(r *http.Request) client.Query {
	params := r.URL.Query()
	return client.Query{
		params.Get("line"),
		params.Get("file"),
		params.Get("repo"),
	}
}

type replySearch struct {
	Info    *client.Stats    `json:"info"`
	Results []*client.Result `json:"results"`
}

func (s *server) ServeAPISearch(w http.ResponseWriter, r *http.Request) {
	backendName := r.URL.Query().Get(":backend")
	backend := s.bk[backendName]
	if backend == nil {
		writeError(w, 400, "bad_backend", fmt.Sprintf("Unknown backend: %s", backendName))
		return
	}

	q := parseQuery(r)

	if q.Line == "" {
		writeError(w, 400, "bad_query", "You must specify a 'line' regex.")
		return
	}

	cl := <-backend.Clients
	defer backend.CheckIn(cl)

	search, err := cl.Query(&q)
	if err != nil {
		writeError(w, 500, "internal_error",
			fmt.Sprintf("Talking to backend: %s", err.Error()))
		return
	}

	reply := &replySearch{}

	for r := range search.Results() {
		reply.Results = append(reply.Results, r)
	}

	reply.Info, err = search.Close()
	if err != nil {
		writeError(w, 500, "internal_error",
			fmt.Sprintf("Talking to backend: %s", err.Error()))
		return
	}

	replyJson(w, 200, reply)
}
