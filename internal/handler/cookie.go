package handler

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/haatos/markdown-blog/internal"
	"github.com/labstack/echo/v4"
)

var secureCookie *securecookie.SecureCookie

var (
	charset               = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890-_|!/"
	seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
)

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generateRandomKey(length int) string {
	return stringWithCharset(length, charset)
}

func NewSecureCookie() {
	var hashKey []byte
	var blockKey []byte
	if internal.Settings.Domain != "localhost" {
		// generate random keys
		// store keys in .env
		hk, hkOk := os.LookupEnv("HASH_KEY")
		bk, bkOk := os.LookupEnv("BLOCK_KEY")

		if hkOk {
			hashKey = []byte(hk)
		} else {
			hashKey = []byte(generateRandomKey(32))
			writeToDotenv("HASH_KEY", string(hashKey))
		}
		if bkOk {
			blockKey = []byte(bk)
		} else {
			blockKey = []byte(generateRandomKey(24))
			writeToDotenv("BLOCK_KEY", string(blockKey))
		}
	} else {
		// keys for development/testing, cookies will not become invalid
		// when server is restarted
		hashKey = []byte("ASDFASDFASDFASDFASDFASDFASDFASDF")
		blockKey = []byte("ASDFASDFASDFASDFASDFASDF")
	}
	secureCookie = securecookie.New(hashKey, blockKey)
}

func writeToDotenv(name, value string) {
	f, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := f.Write([]byte(name + `=` + value + `\n`)); err != nil {
		log.Fatal(err)
	}
}

func setSessionCookie(c echo.Context, data map[string]any) error {
	encoded, err := secureCookie.Encode(internal.SessionCookie, data)
	if err != nil {
		return err
	}
	setCookie(c, internal.SessionCookie, encoded, time.Now().UTC().Add(internal.Settings.SessionExpires))
	return nil
}

func setCookie(c echo.Context, name, value string, expires time.Time) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   internal.Settings.Domain != "localhost",
		Domain:   internal.Settings.Domain,
		Expires:  expires,
	}
	c.SetCookie(cookie)
}

func PostCookieBar(c echo.Context) error {
	value := c.Request().URL.Query().Get("value")

	cookie := &http.Cookie{
		Name:     "cookiebar",
		Value:    value,
		Path:     "/",
		Expires:  time.Now().UTC().Add(10 * 365 * 24 * time.Hour),
		Domain:   internal.Settings.Domain,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
	return nil
}

/*
middleware to add cookie status to request context with ctxKey("cookiebar")
in order check whether user has accepted using all cookies, or only strictly necessary cookies
in case only strictly necessary cookies are being used, this is unnecessary
*/
func CookieBarMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("cookiebar")
		if err == nil {
			c.Set("cookiebar", cookie.Value)
		}
		return next(c)
	}
}
