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

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Masky")
}

func logs(w http.ResponseWriter, r *http.Request) {

}

func configs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c, err := json.Marshal(config)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(c))

	case http.MethodPost:
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

func Start(c *masky.Client) {
	client = c
	config = client.GetConfig()
	http.HandleFunc("/", hello)
	http.HandleFunc("/logs", logs)
	http.HandleFunc("/configs", configs)
	http.HandleFunc("/cache", cache)
	if err := http.ListenAndServe("127.0.0.1:2022", nil); err != nil {
		log.Error(err)
	}
}
