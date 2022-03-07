package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wkj9893/masky/internal/log"
	"github.com/wkj9893/masky/internal/masky"
)

var (
	client *masky.Client
	config masky.ClientConfig
)

func configs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(c))

	case http.MethodPatch:
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
		log.Info(fmt.Sprintf("change config to: %+v", config))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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

func Start(c *masky.Client) error {
	client = c
	config = client.GetConfig()
	http.Handle("/", http.FileServer(http.Dir("../../web/build")))
	http.HandleFunc("/api/configs", configs)
	http.HandleFunc("/api/cache", cache)
	return http.ListenAndServe("127.0.0.1:1081", nil)
}
