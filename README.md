![Build and test workflow](https://github.com/Reconmap/cli/workflows/Build%20and%20test%20workflow/badge.svg)

# Reconmap CLI

Command line interface for the Reconmap pentest automation and reporting tool.

```
$ rmap
Reconmap pentest automation tool.

Usage: rmap [OPTIONS] COMMAND

Commands
 - get clients|projects|tasks|vulnerabilities
 - create clients|projects|tasks|vulnerabilities
 - import
 - run
 - upload

Find out more information at https://reconmap.org/.
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
