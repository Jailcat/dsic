package registry

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
}

type Accounts struct {
	mu       sync.RWMutex
	accounts map[string]Account
	file     string
}

func NewAccounts(file string) *Accounts {
	a := &Accounts{
		accounts: make(map[string]Account),
		file:     file,
	}
	a.load()
	return a
}

func hash(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

func (a *Accounts) Register(username, password, domain string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, acc := range a.accounts {
		if acc.Domain == domain {
			return fmt.Errorf("domain already claimed")
		}
		if acc.Username == username {
			return fmt.Errorf("username already taken")
		}
	}

	a.accounts[username] = Account{
		Username: username,
		Password: hash(password),
		Domain:   domain,
	}
	a.save()
	return nil
}

func (a *Accounts) Login(username, password string) (Account, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	acc, ok := a.accounts[username]
	if !ok {
		return Account{}, fmt.Errorf("account not found")
	}
	if acc.Password != hash(password) {
		return Account{}, fmt.Errorf("wrong password")
	}
	return acc, nil
}

func (a *Accounts) save() {
	data, _ := json.Marshal(a.accounts)
	os.WriteFile(a.file, data, 0600)
}

func (a *Accounts) load() {
	data, err := os.ReadFile(a.file)
	if err != nil {
		return
	}
	json.Unmarshal(data, &a.accounts)
}
