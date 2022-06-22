package main

import (
	"fmt"
	"net/http"
)

func upload(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "upload\n")
}
func actions(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "actions\n")

}

func main() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/actions", actions)

	http.ListenAndServe(":8090", nil)
}
