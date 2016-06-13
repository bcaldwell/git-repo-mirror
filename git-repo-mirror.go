package main

import (
  "fmt"
  "net/http"
  "os/exec"
  "log"
  "os"
  "github.com/ghodss/yaml"
  "io/ioutil"
  "strings"
)

// var saveLocation string = "repos/" should make this work

type webhook struct {
  Url string
  Repo string
  Mirror_repo string
  Name string
}

func (hook *webhook) init () {
  if hook.Name == "" {
    parts := strings.Split(hook.Repo, "/")
    hook.Name = parts[len(parts) - 1]
  }
  if hook.Url == "" {
    parts := strings.Split(hook.Name, ".git")
    hook.Url = "/" + parts[0]
  }

  hook.createRoute()
  hook.setUpRepo()
}


func (hook *webhook) createRoute () {
  http.HandleFunc(hook.Url, hook.ServeHTTP)
}

func (hook *webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  go hook.mirrorRepo()
  fmt.Fprintf(w, hook.Repo)
}

func (hook *webhook) setUpRepo () {
  repoExist, _ := exists (hook.Name)
  if !repoExist {
    runCmd ("git",  []string{"clone", "--mirror", hook.Repo})
  }
  runCmd ("git", []string{"remote", "set-url", "--push", "origin", hook.Mirror_repo}, hook.Name)
  hook.mirrorRepo()
}

func (hook *webhook) mirrorRepo () {
  runCmd ("git", []string{"fetch", "-p","origin"}, hook.Name)
  runCmd ("git", []string{"push", "--mirror"}, hook.Name)
}


func parseYamlConfig () []webhook{
  hooks := make([]webhook, 0)

  data, err := ioutil.ReadFile("/etc/git-repo-mirror/config.yml")
  // data, err := ioutil.ReadFile("config.yml")
  handleError(err);

  err = yaml.Unmarshal(data, &hooks)
  handleError(err)

  fmt.Printf("%+v", hooks)

  return hooks
}


func main() {

    hooks := parseYamlConfig()

    for _, hook := range hooks {
      hook.init()
    }

    http.HandleFunc("/", handler)

    fmt.Println("Starting server on part 8080")
    http.ListenAndServe(":8080", nil)

}

func handler(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "success")
}

func runCmd (cmd string, args []string, dir ...string) (string){
  fmt.Println(args)
  command := exec.Command(cmd, args...)
  if len(dir) > 0 && dir[0] != "" {
    command.Dir = dir[0]
  }
  out, err := command.Output()
  handleError(err)

  fmt.Printf("%s\n", out)
  return string(out)
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
