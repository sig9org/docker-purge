# docker-purge

Stop and delete all containers.

## Install

You can download binaries from the [release page](https://github.com/sig9org/docker-purge/releases) on this repository. Rename the relevant binary for your OS to `docker-purge` and copy it to `$HOME/.docker/cli-plugins`.

Or copy it into one of these folders for installing it system-wide:

- `/usr/local/lib/docker/cli-plugins` OR `/usr/local/libexec/docker/cli-plugins`
- `/usr/lib/docker/cli-plugins` OR `/usr/libexec/docker/cli-plugins`

(might require to make the downloaded file executable with `chmod +x`)

## Quick Start

To stop all containers, do the following:

```sh
# docker purge
CONTAINER ID  IMAGE       STATUS         NAMES               PURGED
0bf662d0106c  sig9/nginx  Up 42 seconds  dreamy_babbage      Done
90d18e5f5dad  alpine      Up 52 seconds  flamboyant_hellman  Done
```

To stop all containers and delete all images, do the following:

```sh
# docker purge --images
CONTAINER ID  IMAGE       STATUS        NAMES            PURGED
9d00c8e0666d  sig9/nginx  Up 7 seconds  zealous_hawking  Done
f70bd26035ae  alpine      Up 9 seconds  sleepy_bell      Done

REPOSITORY  TAG     IMAGE ID      SIZE  PURGED
sig9/nginx  latest  eaa3878ec1ef  33MB  Done
alpine      latest  c059bfaa849c  5MB   Done
```
