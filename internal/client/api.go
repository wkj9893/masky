package client

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/google/uuid"
	"github.com/wkj9893/masky/internal/log"
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
	fmt.Println(p[i].Server[0], p[i].ID)
	return p[i].Server[0], p[i].ID
}

func setIndex(i int) {
	api.Lock()
	defer api.Unlock()
	if i < 0 || i >= len(api.config.Proxies) {
		log.Error("please choose index within the proxies")
		return
	}
	api.index = i
}
