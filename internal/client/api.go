package client

import (
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

var api struct {
	sync.Mutex

	config *Config
	index  int
}

func SetConfig(c *Config) {
	api.Lock()
	defer api.Unlock()
	api.config = c
}

func getIndex() (string, uuid.UUID) {
	api.Lock()
	defer api.Unlock()
	p := api.config.Proxies
	i := api.index
	if i == 0 {
		i = rand.Intn(len(p)-1) + 1
	}
	return p[i].Server[0], p[i].ID
}

func setIndex(i int) {
	api.Lock()
	defer api.Unlock()
	//	TODO check i with length of config proxies
	api.index = i
}
