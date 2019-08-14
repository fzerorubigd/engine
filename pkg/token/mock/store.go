package mock

import (
	"errors"
	"fmt"
	"time"

	"elbix.dev/engine/pkg/random"
	"elbix.dev/engine/pkg/token"
)

type provider map[string]map[string]interface{}

func (p provider) Delete(token string) {
	delete(p, token)
}

func (p provider) Store(data map[string]interface{}, exp time.Duration) (string, error) {
	key := <-random.ID
	cp := make(map[string]interface{}, len(data)+1)
	for i := range data {
		cp[i] = data[i]
	}
	cp["exp"] = time.Now().Add(exp).Unix()
	p[key] = cp
	return key, nil
}

func (p provider) Fetch(token string) (map[string]interface{}, error) {
	data, ok := p[token]
	if !ok || data == nil {
		return nil, errors.New("key is invalid, no data")
	}

	exp := data["exp"].(int64)
	if time.Now().Unix() >= exp {
		return nil, errors.New("key is invalid, expired" + fmt.Sprint(data))
	}

	ret := make(map[string]interface{}, len(data))
	for i := range data {
		if i != "exp" {
			ret[i] = data[i]
		}
	}

	return ret, nil
}

// NewMockStorage is the mock provider
func NewMockStorage() token.Provider {
	return make(provider)
}
