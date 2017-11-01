# go-sync utility

[![GitHub release](https://img.shields.io/github/release/webdevops/go-sync.svg)](https://github.com/webdevops/go-sync/releases)
[![license](https://img.shields.io/github/license/webdevops/go-sync.svg)](https://github.com/webdevops/go-sync/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/webdevops/go-sync.svg?branch=master)](https://travis-ci.org/webdevops/go-sync)
[![Github All Releases](https://img.shields.io/github/downloads/webdevops/go-sync/total.svg)]()
[![Github Releases](https://img.shields.io/github/downloads/webdevops/go-sync/latest/total.svg)]()

Easy project file and database synchronization for developers

Successor for [CliTools Sync](https://github.com/webdevops/clitools) written on Golang

Features
========

General:
- Yaml based configuration files (`gosync.yml` or `.gosync.yml`)
- Automatic cleanup of schemas before restore

Sync:
- Filesync (rsync) from remote servers using SSH
- Create file stubs instead of fetching files from remote (with real images, see ``options.generate-stubs = true``)
- Dump MySQL schemas from remote servers using SSH, Docker and SSH+Docker
- Restore MySQL schema to local MySQL servers or Docker/Docker-Compose containers
- Filtering databases tabes with regexp
- Rsync filters
- Custom exec scripts (startup/finish) on local or remote machine (using SSH)

Deployment:
- Filesync (rsync) from local to remote servers using SSH
- Dump MySQL schemas from local MySQL servers or Docker/Docker-Compose containers
- Filtering databases tabes with regexp
- Rsync filters
- Custom exec scripts (startup/finish) on local or remote machine (using SSH)

Install
=======

The binary file can be found in the [project releases](https://github.com/webdevops/go-sync/releases).

```
DOWNLOAD_VERSION=0.5.2
DOWNLOAD_OS=linux
DOWNLOAD_ARCH=x64

wget -O/usr/local/bin/gosync "https://github.com/webdevops/go-sync/releases/download/${DOWNLOAD_VERSION}/gosync-${DOWNLOAD_OS}-${DOWNLOAD_ARCH}"
chmod +x /usr/local/bin/gosync
```

Help
====

```
Usage:
  gosync [OPTIONS] <command>

Application Options:
  -v, --verbose  verbose mode

Help Options:
  -h, --help     Show this help message

Available commands:
  deploy       Deploy to server
  list         List server configurations
  self-update  Self update
  sync         Sync from server
  version      Show version

```

Example
=======

```
> gosync sync production

:: Initialisation
   -> found configuration file /Users/xxxxxx/Projects/examples/gosync.yml
   -> using production server
   -> using connection Exec[Type:ssh SSH:ssh-user@example.com]
:: Starting exec mode "startup"
   -> executing >> Exec[Type:local Command:date +%s]
   -> executing >> Exec[Type:local Command:date +%s]
   -> executing >> Exec[Type:local Workdir:/ Command:date]
   -> executing >> Exec[Type:local Workdir:/ Command:date]
   -> executing >> Exec[Type:remote Workdir:/ Command:date]
:: Starting sync of Filesystem[Path:/home/xxxxxx/application1/ -> Local:./application1/]
:: Starting sync of Filesystem[Path:/home/xxxxxx/application2/ -> Local:./application2/]
:: Starting sync of Database[Schema:application1 User:mysql-user Passwd:***** -> Schema:test-local]
   -> dropping local database "test-application1"
   -> creating local database "test-application1"
   -> syncing database structure
   -> get list of mysql tables for table filter
   -> syncing database data
:: Starting sync of Database[Schema:application2 User:mysql-user Passwd:***** -> Schema:test]
   -> dropping local database "test-application2"
   -> creating local database "test-application2"
   -> syncing database structure
   -> get list of mysql tables for table filter
   -> syncing database data
:: Starting exec mode "finish"
   -> executing >> Exec[Type:remote Workdir:/ Command:date]
-> finished
```

## Docker support

Docker support
==============

Using the configuration ``connection.docker=configuration`` this command can be
execued with docker containers. If the container id is passed the
container is used without lookup using eg. `docker-compose`.

**docker-compose:**

*CONTAINER* is the name of the docker-compose container.

| DSN style configuration                                             | Description                                                                                     |
|:--------------------------------------------------------------------|:------------------------------------------------------------------------------------------------|
| ``compose:CONTAINER``                                               | Use container with docker-compose in current directory                                          |
| ``compose:CONTAINER;path=/path/to/project``                         | Use container with docker-compose in `/path/to/project` directory                               |
| ``compose:CONTAINER;path=/path/to/project;file=custom-compose.yml`` | Use container with docker-compose in `/path/to/project` directory and `custom-compose.yml` file |
| ``compose:CONTAINER;project-name=foobar``                           | Use container with docker-compose in current directory with project name `foobar`               |
| ``compose:CONTAINER;host=example.com``                              | Use container with docker-compose in current directory with docker host `example.com`           |
| ``compose:CONTAINER;env[FOOBAR]=BARFOO``                            | Use container with docker-compose in current directory with env var `FOOBAR` set to `BARFOO`    |

| Query style configuration                                             | Description                                                                                     |
|:----------------------------------------------------------------------|:------------------------------------------------------------------------------------------------|
| ``compose://CONTAINER``                                               | Use container with docker-compose in current directory                                          |
| ``compose://CONTAINER?path=/path/to/project``                         | Use container with docker-compose in `/path/to/project` directory                               |
| ``compose://CONTAINER?path=/path/to/project&file=custom-compose.yml`` | Use container with docker-compose in `/path/to/project` directory and `custom-compose.yml` file |
| ``compose://CONTAINER?project-name=foobar``                           | Use container with docker-compose in current directory with project name `foobar`               |
| ``compose://CONTAINER?host=example.com``                              | Use container with docker-compose in current directory with docker host `example.com`           |
| ``compose://CONTAINER?env[FOOBAR]=BARFOO``                            | Use container with docker-compose in current directory with env var `FOOBAR` set to `BARFOO`    |
