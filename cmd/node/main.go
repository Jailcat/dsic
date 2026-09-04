package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("DSIC node starting...")

	port := "8080"
	if p := os.Getenv("DSIC_PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "this node is alive on DSIC")
	})

	fmt.Printf("node listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
