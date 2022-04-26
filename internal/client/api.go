package client

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"path"
	"sync"

	"github.com/google/uuid"
	"github.com/wkj9893/masky/internal/log"
)

var api struct {
	sync.RWMutex

	config *Config
	index  int
}

func getConfig() *Config {
	api.RLock()
	defer api.RUnlock()
	return api.config
}

func setConfig(c *Config) {
	api.Lock()
	defer api.Unlock()
	api.config = c
}

func getIndex() (string, uuid.UUID) {
	api.RLock()
	defer api.RUnlock()
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
		c, err := json.Marshal(getConfig())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(c))
	case http.MethodPatch:
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

func StartApi() {
	port := getConfig().Port + 1
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if p == "/config" || p == "/proxies" {
			http.ServeFile(w, r, "../../web/build/index.html")
			return
		}
		http.ServeFile(w, r, path.Join("../../web/build", r.URL.Path))
	})
	m.HandleFunc("/api/config", handleConfig)

	log.Info(fmt.Sprintf("start http server at http://localhost:%v", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), m); err != nil {
		log.Error(err)
	}
}
