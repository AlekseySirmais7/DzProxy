package main

import (
	"fmt"
	"log"
	"net/http"
)

func testparam(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["port"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'port' is missing")
		fmt.Fprintf(w, "simple page\n")
		return
	}
	key := keys[0]
	fmt.Fprintf(w, "simple page\n and user param port with value:"+key+"\n")
}

func main() {
	http.HandleFunc("/testparam", testparam)
	log.Println("testparam server uses 5001 port")
	http.ListenAndServe(":5001", nil)
}
