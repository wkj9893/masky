package client

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/wkj9893/masky/internal/log"
)

var api struct {
	sync.Mutex

	config *Config
	index  int
}

func setConfig(c *Config) {
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
	if i < 0 || i >= len(api.config.Proxies) {
		log.Error("please choose index within the proxies")
		return
	}
	api.index = i
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c, err := json.Marshal(api.config)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(c))
	case http.MethodPut:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))
			return
		}
		var config Config
		err = json.Unmarshal(body, &config)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
		setConfig(&config)
		log.Info(fmt.Sprintf("change config to:%+v", config))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}
func handleProxies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c, err := json.Marshal(api.config.Proxies)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(c))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

var server http.Server

func StartApi(config *Config) {
	addr := fmt.Sprintf(":%v", config.Port+1)

	m := http.NewServeMux()
	m.Handle("/", http.FileServer(http.Dir("../../web/build")))
	m.HandleFunc("/api/config", handleConfig)
	m.HandleFunc("/api/proxy", handleProxies)

	log.Info(fmt.Sprintf("start api server at:%v", addr))
	if err := http.ListenAndServe(addr, m); err != nil {
		log.Error(err)
	}
}
