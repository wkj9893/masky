package client

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

func index(w http.ResponseWriter, r *http.Request) {
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

func startApi(config Config) error {
	http.HandleFunc("/", index)
	// TODO
	// http.HandleFunc("/api/setting", setting)
	if config.AllowLan {
		return http.ListenAndServe(fmt.Sprintf(":%v", config.Port+1), nil)
	}
	return http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", config.Port+1), nil)
}
