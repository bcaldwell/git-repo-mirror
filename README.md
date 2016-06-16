# git-repo-mirror

##Set up

Create a config.yml to configure mirror repos. An example config.yml
```yaml
---
  -
    Repo: "repo to mirror from"
    Mirror_repo: "repo to mirror to"
    Url: "/hook url (optional)"
    Name: "clone name(optional)"

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
      - SECRET=github/giblab_secret
```

## Setting up on a server
### SSH access
For this to work the machine or container that it is running in must have push and read access for the repos. For read access great deploy keys (works for both github and gitlab). For read access either create a new deploy user or use deploy keys with push access (works in github)
