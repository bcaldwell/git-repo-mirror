package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ghodss/yaml"
)

func main() {

	hooks := parseYamlConfig()

	for _, hook := range hooks {
		hook.init()
	}

	http.HandleFunc("/", handler)

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

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "Not Found")
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
