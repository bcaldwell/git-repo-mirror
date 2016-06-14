package main

import (
  "fmt"
  "strings"
  "crypto/hmac"
  "crypto/sha1"
  "encoding/hex"
  "net/http"
  "os"
  "os/exec"
  "io/ioutil"
  "gopkg.in/robfig/cron.v2"
)

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

  hook.createCron()
  hook.createRoute()
  hook.setUpRepo()
}

func (hook *webhook) createCron () {
  interval := os.Getenv("CRON")
  if strings.ToLower(interval) == "false" {
    return
  } else if interval == "" {
    interval = "@hourly"
  }
  c := cron.New()
  c.AddFunc(interval, hook.mirrorRepo)
  c.Start()
}

func (hook *webhook) createRoute () {
  http.HandleFunc(hook.Url, hook.ServeHTTP)
}

func (hook *webhook) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  if (verifyRequest(req)){
    go hook.mirrorRepo()
    fmt.Fprintf(res, hook.Repo)
  } else {
    http.Error(res, "400 Bad Request - Missing X-GitHub-Event Header", http.StatusBadRequest)
    return
  }
}


func (hook *webhook) setUpRepo () {
  repoExist, _ := exists (hook.Name)
  if !repoExist {
    fmt.Println("Cloning", hook.Repo)
    runCmd ("git",  []string{"clone", "--mirror", hook.Repo})
  }
  fmt.Println("Setting push remote to ", hook.Mirror_repo)
  runCmd ("git", []string{"remote", "set-url", "--push", "origin", hook.Mirror_repo}, hook.Name)
  hook.mirrorRepo()
}

func (hook *webhook) mirrorRepo () {
  fmt.Println("Pulling", hook.Repo)
  runCmd ("git", []string{"fetch", "-p","origin"}, hook.Name)
  fmt.Println("Pushing", hook.Mirror_repo)
  runCmd ("git", []string{"push", "--mirror"}, hook.Name)
}

func verifyRequest (req *http.Request) bool {
  body, err := ioutil.ReadAll(req.Body)
  handleError(err)

  secret := os.Getenv("SECRET")
  if secret != ""{
    const signaturePrefix = "sh1="
    const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

    sig := req.Header.Get("X-Hub-Signature")

    if sig == "" || len(sig) != signatureLength || !strings.HasPrefix(sig, signaturePrefix) {
      return false
    }

    mac := hmac.New(sha1.New, []byte(secret))
    mac.Write(body)
    expectedMac := mac.Sum(nil)
    expectedSig := "sha1=" + hex.EncodeToString(expectedMac)

    return hmac.Equal([]byte(expectedSig), []byte(sig))
  }
  return true
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
