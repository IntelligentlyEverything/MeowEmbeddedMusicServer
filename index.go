package main

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "MeowMusicEmbeddedServer")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Printf("[Web Access] Handling request for %s\n", r.URL.Path)
	if r.URL.Path != "/" {
		fileHandler(w, r)
		return
	}
	fmt.Fprintf(w, "<h1>音乐服务器</h1>")
}
