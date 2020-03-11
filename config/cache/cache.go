package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Emp struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

var Cache = cache.New(5*time.Minute, 5*time.Minute)

func SetCache(key string, emp interface{}) bool {
	Cache.Set(key, emp, cache.NoExpiration)
	return true
}

func GetCache(key string) (*Emp, bool) {
	var found bool
	emp, found := Cache.Get(key)
	if found {
		response := new(Emp)
		response.Name = key
		response.Data = emp

		return response, found
	}
	return nil, found
}
