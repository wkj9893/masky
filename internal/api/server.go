package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
)

var (
	client *masky.Client
	config masky.ClientConfig
)

func setting(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPatch {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
		err = json.Unmarshal(body, &config)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
		client.SetConfig(config)
		log.Info(fmt.Sprintf("change config to:%+v", config))

	} else {
		c, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(c))
	}
}

func cache(w http.ResponseWriter, r *http.Request) {
	cache, err := client.MarshalCache()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(cache))
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := path.Join("./web/build", r.URL.Path)
	if _, err := os.Stat(name); err != nil {
		if data, err := os.ReadFile("./web/build/index.html"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, string(data))
		}
		return
	}
	http.ServeFile(w, r, name)
}

func Start(c *masky.Client) error {
	client = c
	config = client.GetConfig()
	http.HandleFunc("/", handler)
	http.HandleFunc("/api/setting", setting)
	http.HandleFunc("/api/cache", cache)
	if config.AllowLan {
		return http.ListenAndServe(fmt.Sprintf(":%v", config.Port+1), nil)
	}
	return http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", config.Port+1), nil)
}
