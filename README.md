# Go Docker Service Clone (aka GDSC)

Clone docker services in a easy way.

## Usage

```
$ go get github.com/herlon214/gdsc
$ gdsc -h
Utilities to deal with services in docker.

Usage:
  gdsc [command]

Available Commands:
  export      Export a service to a json file
  help        Help about any command
  import      Import a JSON file and create a service using it
  upsert      Copy a service and create a new one overriding some props

Flags:
      --api-url string   Docker api url (default "http://127.0.0.1:2375")
  -h, --help             help for gdsc

Use "gdsc [command] --help" for more information about a command.
```

## Create / Update behaviors
There are two behaviors when executed: create / update.

**Create**

This behavior will be triggered if there's no service created with name `--name`. If triggered, it will copy the whole `--copy-from` service, secrets, hosts, mounts, environment vars, labels, creating a new one with name `--name` and the given `--image`.

It will change the `traefik.frontend.rule` label in the new service setting a new host to the service. The new host will be `http://new_service_name.domain_from_arg`. The new service name is cleaned replacing all non-words (`/\W/`) with "`_`" to be a bit web friendly. 

In other words: if you specify `--name my_service_feature/1 --domain mycompany.org`, the service name will be set to `my_service_feature_1` and the host will be `http://my_service_feature_1.mycompany.org`.

**Update**

This behavior will be triggered if there is a service with the name given (`--name`), it will only update the docker image of this service with the given one (`--image`).

## How to build

If you have docker installed you can type `make build` and the binary output will appear in `dist/` folder.

## Published Ports issue

This project doesn't copy the **published ports** because Traefik will connect to container port. The published port may cause a `already in use` error.

## How it works in a CI stack
Think about the scenario below and you can understand what's the purpose of this project.

Your stack:
- [git flow](https://github.com/nvie/gitflow)
- [Jenkins Multibranch Pipeline](https://jenkins.io/)
- [Docker](https://www.docker.com/)
- [Traefik](https://traefik.io/)
- [Docker Swarm](https://github.com/docker/swarm)

Your services (based on git flow branches):

|Name|Type|Branch|URI|
|----|----|------|---|
|website_master|production|master|http://mycompany.org/|
|website_develop|stagging|develop|http://beta.mycompany.org/|

**Developer A** (aka D.A) is working on **feature/a** in the website project.

**Developer B** (aka D.B) is working on **feature/b** in the website project.


D.A just merged his feature to stagging environment (website_develop) because it's done.

D.B wants to test his feature in the stagging environment. But he can't merge to website_develop because it would be merged to master without being done and cause errors (D.B will merge develop to master soon).

> **GDSC** will clone the service **website_develop** with a new name based on the branch's name with the new docker image URL and a new Traefik host.

D.B publishes his **feature/a** to the repository (`git flow feature publish`).

Git repository (gitlab, github) will trigger Jenkins with the update.

Jenkins will build the *Dockerfile* and then call GDSC to create a new service or update a created service.

> ```
> gdsc --copy-from voting_example_develop --name voting_example_feature/a --image registry.gitlab.com/mycompany/website:feature_a --domain mycompany.org
> ```

Here is a Jenkinsfile example:

```
pipeline {
  stages {
    stage ("Clone") {
      steps {
        checkout scm
      }
    }
    stage ("Build") {
      steps {
        sh "docker build --no-cache -t REGISTRY_URL:${env.BRANCH_NAME} ."
        sh "docker push REGISTRY_URL:${env.BRANCH_NAME}"
      }
    }
    stage ("Deploy") {
      steps {
        sh "gdsc --copy-from ${env.JOB_NAME} --name ${env.JOB_NAME}${env.BRANCH_NAME} --image REGISTRY_URL:${env.BRANCH_NAME} --domain mycompany.org"
      }
    }
  }
}
```

Lastly, your services will look like:

|Name|Type|Branch|URI|
|----|----|------|---|
|website_master|production|master|http://mycompany.org/|
|website_develop|stagging|develop|http://beta.mycompany.org/|
|website_feature_a|testing|feature/a|http://feature_a.mycompany.org/|
