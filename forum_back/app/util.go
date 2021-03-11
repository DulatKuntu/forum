package app

import (
	"errors"
	"time"
)

func (app *Application) containsToken(token string) (int, error) {
	for key, val := range app.Cookies {
		if val.Value == token {
			if val.Expires.Sub(time.Now()) <= 0 {
				return key, errors.New("token expired")
			}
			return key, nil
		}
	}
	return 0, errors.New("not authorized")
}
