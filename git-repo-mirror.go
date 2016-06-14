package main

import (
  "fmt"
  "net/http"
  "log"
  "os"
  "github.com/ghodss/yaml"
  "io/ioutil"
)

// var saveLocation string = "repos/" should make this work

func main() {

    hooks := parseYamlConfig()

    for _, hook := range hooks {
      hook.init()
    }

    http.HandleFunc("/", handler)

    fmt.Println("Starting server on part 8080")
    http.ListenAndServe(":8080", nil)

}

func parseYamlConfig () []webhook{
  hooks := make([]webhook, 0)

  configFile := os.Getenv("CONFIG_FILE")
  if (configFile == ""){
    configFile = "config.yml"
  }

  // data, err := ioutil.ReadFile("/etc/git-repo-mirror/config.yml")
  data, err := ioutil.ReadFile(configFile)
  handleError(err);

  err = yaml.Unmarshal(data, &hooks)
  handleError(err)

  fmt.Printf("%+v\n\n", hooks)

  return hooks
}

func handler(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "success")
}

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func handleError (err error) {
  if err != nil {
    log.Fatal(err)
  }
}
