![Build and test workflow](https://github.com/Reconmap/cli/workflows/Build%20and%20test%20workflow/badge.svg)

# Reconmap CLI

Command line interface for the Reconmap pentest automation and reporting tool.

```
$ ./rmap config --api-url https://api.reconmap.org
$ ./rmap login -u admin -p ******
$ ./rmap command run -cid 2 -var Host=soki.com.ar
Reconmap v1.0 - https://reconmap.org

> Downloading docker image 'instrumentisto/nmap'
> Starting container.
> Container started.
> Container 'instrumentisto/nmap' exited.
```

## Build requirements

- Golang 1.15+
- Make

## Build instructions

```shell
$ make
```

## Runtime requirements

- Docker engine with [API version 1.40](https://docs.docker.com/engine/api/v1.40/)

# Troubleshooting

### Error response from daemon: client version 1.XX is too new. Maximum supported API version is 1.40

```shell
export DOCKER_API_VERSION=1.40
```
