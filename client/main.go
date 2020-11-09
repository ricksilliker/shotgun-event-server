package client

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/work", dccWorker)
	http.ListenAndServe(":8090", nil)
}

func dccWorker(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}