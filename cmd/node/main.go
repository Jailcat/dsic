package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Jailcat/dsic/pkg/registry"
)

func main() {
	fmt.Println("DSIC node starting...")

	reg := registry.New()
	accounts := registry.NewAccounts("accounts.json")

	port := "8080"
	if p := os.Getenv("DSIC_PORT"); p != "" {
		port = p
	}

	domain := "unnamed.node"
	if d := os.Getenv("DSIC_DOMAIN"); d != "" {
		domain = d
	}

	var peers []string
	if p := os.Getenv("DSIC_PEERS"); p != "" {
		peers = strings.Split(p, ",")
	}

	reg.Register(domain, "127.0.0.1", port)

	gossip := registry.NewGossip(reg, peers, port)
	gossip.Start()

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

	http.HandleFunc("/claim", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", 405)
			return
		}
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Domain   string `json:"domain"`
			Address  string `json:"address"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad request", 400)
			return
		}
		if err := accounts.Register(body.Username, body.Password, body.Domain); err != nil {
			http.Error(w, err.Error(), 409)
			return
		}
		reg.Register(body.Domain, body.Address, port)
		w.WriteHeader(201)
		fmt.Fprintf(w, "domain %s claimed", body.Domain)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", 405)
			return
		}
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad request", 400)
			return
		}
		acc, err := accounts.Login(body.Username, body.Password)
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(acc)
	})

	fmt.Printf("node listening on port %s as %s\n", port, domain)
	http.ListenAndServe(":"+port, nil)
}
