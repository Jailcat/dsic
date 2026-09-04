package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Jailcat/dsic/pkg/registry"
)

func main() {
	fmt.Println("DSIC node starting...")

	reg := registry.New()

	port := "8080"
	if p := os.Getenv("DSIC_PORT"); p != "" {
		port = p
	}

	domain := "unnamed.node"
	if d := os.Getenv("DSIC_DOMAIN"); d != "" {
		domain = d
	}

	reg.Register(domain, "127.0.0.1", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "this node is alive on DSIC as %s", domain)
	})

	http.HandleFunc("/registry", func(w http.ResponseWriter, r *http.Request) {
		data, err := reg.ToJSON()
		if err != nil {
			http.Error(w, "failed to serialize registry", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(json.RawMessage(data))
	})

	fmt.Printf("node listening on port %s as %s\n", port, domain)
	http.ListenAndServe(":"+port, nil)
}
