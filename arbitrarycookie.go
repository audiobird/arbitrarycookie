package arbitrarycookie

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"log"
	"net/http"
	"strings"
)

var (
	ErrCorrupt = errors.New("cookiecart: Cart is corrupt")

	cookieName = "ArbitraryCookie"
)

func Init(name string) {
	cookieName = name
}

func Save(data any, w http.ResponseWriter) {
	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(data)
	if err != nil {
		log.Println(err)
		return
	}

	cookie := http.Cookie{
		Name:  cookieName,
		Value: base64.RawURLEncoding.EncodeToString(buf.Bytes()),
	}

	http.SetCookie(w, &cookie)
}

func Read(data any, r *http.Request) error {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil
		}
		return err
	}
	t, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Println("Failed to Decode Cookie")
		return nil
	}

	reader := strings.NewReader(string(t))

	if err := gob.NewDecoder(reader).Decode(data); err != nil {
		log.Println("Failed to Decode Gob")
		return nil
	}

	return nil
}
