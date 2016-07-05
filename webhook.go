package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"gopkg.in/robfig/cron.v2"
)

type webhook struct {
	Url         string
	Repo        string
	Mirror_repo string
	Name        string
	Force       bool
	Dir         string
}

func (hook webhook) init() {
	if hook.Name == "" {
		parts := strings.Split(hook.Repo, "/")
		hook.Name = parts[len(parts)-1]
	}
	if hook.Url == "" {
		parts := strings.Split(hook.Name, ".git")
		hook.Url = "/" + parts[0]
	}
	if hook.Repo == "" || hook.Mirror_repo == "" {
		err := errors.New("webhook configuration must contain both repo and mirror_repo options")
		log.Fatal(err)
	}

	go hook.createCron()
	hook.createRoute()
	hook.setUpRepo()
}

func (hook webhook) createCron() {
	interval := envDefault("CRON", "* * 1 * * *")

	if strings.ToLower(interval) == "false" {
		return
	}
	c := cron.New()
	c.AddFunc(interval, hook.mirrorRepo)
	c.Start()
}

func (hook webhook) createRoute() {
	http.HandleFunc(hook.Url, hook.ServeHTTP)
}

func (hook *webhook) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if verifyRequest(req) {
		go hook.mirrorRepo()
		fmt.Fprintf(res, hook.Repo)
	} else {
		http.Error(res, "400 Bad Request - Missing X-GitHub-Event Header", http.StatusBadRequest)
		return
	}
}

func (hook *webhook) setUpRepo() {
	repoExist, _ := exists(hook.savePath())
	if !repoExist {
		fmt.Println("Cloning", hook.Repo)
		runCmd("git", []string{"clone", "--mirror", hook.Repo}, hook.saveDir())
	} else {
		fmt.Println("Setting pull remote to ", hook.Mirror_repo)
		runCmd("git", []string{"remote", "set-url", "origin", hook.Repo}, hook.savePath())
	}
	fmt.Println("Setting push remote to ", hook.Mirror_repo)
	runCmd("git", []string{"remote", "set-url", "--push", "origin", hook.Mirror_repo}, hook.savePath())
	hook.mirrorRepo()
}

func (hook *webhook) mirrorRepo() {

	fmt.Println("Pulling", hook.Repo)
	runCmd("git", []string{"fetch", "-p", "origin"}, hook.savePath())
	fmt.Println("Pushing", hook.Mirror_repo)
	gitPushArgs := []string{"push", "--mirror"}
	if hook.Force {
		gitPushArgs = append(gitPushArgs, "-f")
	}
	runCmd("git", gitPushArgs, hook.savePath())
}

func (hook *webhook) saveDir() (dir string) {
	dir = hook.Dir
	if dir == "" {
		dir = envDefault("BASEDIR", "repos")
	}
	os.MkdirAll(dir, 0755)
	return
}

func (hook *webhook) savePath() string {
	return path.Join(hook.saveDir(), hook.Name)
}

func verifyRequest(req *http.Request) bool {
	body, err := ioutil.ReadAll(req.Body)
	handleError(err)

	secret := os.Getenv("SECRET")
	if secret != "" {
		const signaturePrefix = "sha1="
		const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

		sig := req.Header.Get("X-Hub-Signature")
		gitlabToken := req.Header.Get("X-Gitlab-Token")

		if sig == "" && gitlabToken != "" {
			if gitlabToken == secret {
				return true
			} else {
				return false
			}
		} else if len(sig) != signatureLength || !strings.HasPrefix(sig, signaturePrefix) {
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

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("Git output: %s\n", string(outs))
	}
}

func runCmd(cmd string, args []string, dir ...string) {
	command := exec.Command(cmd, args...)
	if len(dir) > 0 && dir[0] != "" {
		command.Dir = dir[0]
	}

	output, err := command.CombinedOutput()
	handleError(err)
	if err != nil {
		printOutput(output)
	}
}
