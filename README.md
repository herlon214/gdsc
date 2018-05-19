# Go Docker Service Clone (aka GDSC)

## Usage

```
$ go get github.com/herlon214/gdsc
$ gdsc
Usage: gdsc --from FROM --name NAME --image IMAGE

Options:
  --from FROM            service that will be cloned if there is no service with the given --name
  --name NAME            service name that will be deployed
  --image IMAGE          new docker image url
  --help, -h             display this help and exit
```

## How it works
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
|website_master|production|master|http://my.website/|
|website_develop|stagging|develop|http://beta.my.website/|

**Developer A** (aka D.A) is working on **feature/a** in the website project.

**Developer B** (aka D.B) is working on **feature/b** in the website project.


D.A just merged his feature to stagging environment (website_develop) because it's done.

D.B wants to test his feature in the stagging environment. But he can't merge to website_develop because it would be merged to master without being done and cause errors (D.B will merge develop to master soon).

> **GDSC** will clone the service **website_develop** with a new name based on the branch's name with the new docker image URL and a new Traefik host.

D.B publishes his **feature/a** to the repository (`git flow feature publish`).

Git repository (gitlab, github) will trigger Jenkins with the update.

Jenkins will build the *Dockerfile* and then call GDSC to create a new service or update the a created service.

> ```
> gdsc website feature/a registry.gitlab.com/mycompany/website:feature_a
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
        sh "gdsc ${env.JOB_NAME} ${env.BRANCH_NAME} REGISTRY_URL:${env.BRANCH_NAME}"
      }
    }
  }
}
```

Lastly, your services will look like:

|Name|Type|Branch|URI|
|----|----|------|---|
|website_master|production|master|http://my.website/|
|website_develop|stagging|develop|http://beta.my.website/|
|website_feature_a|testing|feature/a|http://feature_a.testing.my.website/|


## Published Ports

This project doesn't copy the published ports because Traefik will connect to container port. The published port may cause a `already in use` error.