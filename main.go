package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"encoding/json"

	"github.com/ghodss/yaml"
)

func main() {

	hooks := parseYamlConfig()

	for i := range hooks {
		go hooks[i].init()
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "pong")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(200)
		response, _ := json.MarshalIndent(hooks, "", "\t")
		fmt.Fprintf(w, "Configuration:\n\n%s", string(response))
	})

	port := envDefault("PORT", "8080")

	cert, key := os.Getenv("CERT"), os.Getenv("KEY")

	var err error

	if cert != "" && key != "" {
		fmt.Println("Starting TLS server on port", port)
		err = http.ListenAndServeTLS(":"+port, cert, key, nil)
	} else {
		fmt.Println("Starting server on port", port)
		err = http.ListenAndServe(":"+port, nil)
	}
	handleError(err)
	runCmd("sleep", []string{"1000"})

}

func parseYamlConfig() []webhook {
	hooks := make([]webhook, 0)

	configFile := envDefault("CONFIG_FILE", "config.yml")

	data, err := ioutil.ReadFile(configFile)
	handleError(err)

	err = yaml.Unmarshal(data, &hooks)
	handleError(err)

	fmt.Printf("%+v\n\n", hooks)

	return hooks
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func handleError(err error, fatal ...bool) {
	if err != nil {
		if len(fatal) > 0 && fatal[0] {
			log.Fatal(err)
		} else {
			log.Print(err)
		}
	}
}

func envDefault(env, fallback string) (val string) {
	val = os.Getenv(env)
	if val == "" {
		val = fallback
	}
	return val
}
