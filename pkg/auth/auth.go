package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	gossh "golang.org/x/crypto/ssh"
)

type Auth interface {
	PubKey(gossh.PublicKey) (string, bool)
	Password(user, password string) (string, bool)
}

type PubKeyAuth struct {
	PubKey string `json:"pubKey"`
	Alias  string `json:"alias"`
}

type PasswordAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Alias    string `json:"alias"`
}

type DefaultAuth struct {
	PubKeys     []*PubKeyAuth   `json:"pubKeys"`
	Passwords   []*PasswordAuth `json:"passwords"`
	pubKeyMap   map[string]string
	passwordMap map[string]string
}

func Update(a *DefaultAuth) *DefaultAuth {
	a.pubKeyMap = make(map[string]string)
	for _, item := range a.PubKeys {
		out, _, _, _, err := gossh.ParseAuthorizedKey([]byte(item.PubKey))
		if nil != err {
			continue
		}
		hash := sha256.Sum256(out.Marshal())
		k := hex.EncodeToString(hash[:])
		a.pubKeyMap[k] = item.Alias
	}

	a.passwordMap = make(map[string]string)
	for _, item := range a.Passwords {
		a.passwordMap[fmt.Sprintf("%s:%s", item.Username, item.Password)] = item.Alias
	}
	return a
}

func (d *DefaultAuth) PubKey(key gossh.PublicKey) (string, bool) {
	hash := sha256.Sum256(key.Marshal())
	k := hex.EncodeToString(hash[:])
	a, found := d.pubKeyMap[k]
	return a, found
}

func (d *DefaultAuth) Password(user, password string) (string, bool) {
	key := fmt.Sprintf("%s:%s", user, password)
	a, found := d.passwordMap[key]
	return a, found
}
