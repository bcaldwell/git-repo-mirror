# git-repo-mirror

## Set up

Create a config.yml to configure mirror repos. An example config.yml
```yaml
---
  -
    repo: "repo to mirror from"
    mirror_repo: "repo to mirror to"
    url: "/hook url (optional)"
    name: "clone name(optional)"
    force: "force push (optional)"
    dir: "directory to save as (optional)"

  -
    repo: gitlab.com/username/repo
    mirror_repo: github.com/username/repo
    url: /repo_name
    name: some_name
    force: false
    dir: repos_save_location
    # repo will be saved at repos_save_location/some_name

```

###Docker-compose
Example docker compose files
```yaml
version: '2'
services:
  web:
    image: benjamincaldwell/git-repo-mirror:beta
    ports:
      - "8080:8080"
    volumes:
      - ./config.yml:/etc/git-repo-mirror/config.yml
      - ~/.ssh:/root/.ssh
    environment:
      - SECRET=github-or-giblab_secret
```

## Environment Variables

- BASEDIR: base save location (default is `repos`)
<!-- - LOGLEVEL: amount of logging -->
- SECRET: secret used by github and/or gitlab to verify the webhook
- CERT: ssl certificate
- KEY: ssl private keys
- PORT: port to run on
- CRON: cron string to run the cron job on (default `"* * 1 * * *"` which is 1 hour)
- CONFIG_FILE: config file location

## Setting up on a server
### SSH access
For this to work the machine or container that it is running in must have push and read access for the repos. For private repos the suggested method is using ssh keys. 
