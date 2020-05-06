package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// http://mygit.server/group/project?go-get=1
// <meta content='corp.server/group/project git git+ssh://git@corp.server/group/project.git' name='go-import'>
var port = 8080

func main() {
	rtr := mux.NewRouter()
	rtr.PathPrefix("/").HandlerFunc(metaHandler)
	http.Handle("/", rtr)
	portString := fmt.Sprintf(":%d", port)
	log.Printf("Listening %s", portString)
	http.ListenAndServe(portString, nil)
}

func metaHandler(w http.ResponseWriter, req *http.Request) {
	uri := strings.TrimLeft(req.RequestURI, "/")
	if strings.Index(uri, "@v") > 0 {
		log.Printf("410 meta request: %s\n", req.RequestURI)
		w.WriteHeader(http.StatusGone)
		fmt.Fprintf(w, "not found: %s\n", uri)
		return
	}
	if isGet := req.URL.Query().Get("go-get"); isGet != "1" {
		log.Printf("Bogus Request: %s\n", req.RequestURI)
		http.NotFound(w, req)
		return
	}
	meta := fmt.Sprintf("<html><head><meta content='git+ssh://git@%s.git' name='go-import'></head></html>\n", uri)
	log.Print(meta)
	fmt.Fprintf(w, meta)
}

