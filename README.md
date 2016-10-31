# pazuzu
[![Travis BuildStatus](https://travis-ci.org/zalando-incubator/pazuzu.svg?branch=master)](https://travis-ci.org/zalando-incubator/pazuzu)
[![Stories in Ready](https://badge.waffle.io/zalando/pazuzu.png?label=ready&title=Ready)](https://waffle.io/zalando/pazuzu)

A note to the next team
---------------
Pazuzu is a big enough project to keep you busy for more than 2 weeks 
of onboarding, so please pick your tasks wisely and leave the rest to 
the next team. Good place to start picking tasks is already existing 
TODO list: https://github.com/zalando-incubator/pazuzu/wiki/TODO

You can find overview of the whole project: 
https://github.com/zalando-incubator/pazuzu/wiki/Overview-of-Pazuzu-project

These flow diagrams may help you understand better how pazuzu is going 
to be used from command line: https://github.com/zalando-incubator/pazuzu/wiki/CLI-commands-flow

Good luck! :)

What is Pazuzu?
---------------
Pazuzu is a tool that builds Docker images from feature snippets, while
resolving all dependencies between them. One of the common use cases is
Continuous Integration environment, where jobs require specific tooling present
for building and testing. Pazuzu can significantly ease that process, by
letting user choose from a wide selection of predefined Dockerfile snippets
that represent those dependencies (e.g. Golang, Python, Android SDK, customized
NPM installs).

## Build
Make sure you setup Go acording to: https://golang.org/doc/install#install
```
go get -v
go build
```

## Usage

Basicaly, pazuzu cli tool has 3 steps:
- **search** - search for available features
- **compose** - compose **Pazuzufile** with desired features
- **build** - create Dockerfile out of Pazuzufile

## Search features
List step is used to check what features are actually available in configured snippet provider
```
$ pazuzu search [regexp]
```
Mask is a valid regexp, i.e.:
```
$ pazuzu search "node-v4.6*|java8|mvn"
```
Mask `node-v4.6*|java8|mvn` will use following features:
- latest `node` version that match regexp **node-v4.6***
- latest `java8` release available in feature snippets provider(default java8 )

## Compose features
Compose step actually creates features file out of specified 
```
$ pazuzu compose <space-separated-feature-names> 
```
### Available options
- **--with-base NAME** - specify base image
- **-o/--out FILENAME** - specify name of the output features file. Default is Pazuzu file
### Examples:
`$ pazuzu compose node java8` - generate Pazuzufile from the list of features

`$ pazuzu compose node java8 --out node-und-java8.yml` - Save everything to **node-und-java8.yml** file

## Build Dockerfile
This step aims to actually create **Dockerfile** out of the snippets configured for features.

`$ pazuzu build . `  - Builds a **Dockerfile** from **Pazuzufile** located in the current directory.

` $ pazuzu build <node-with-babel.yml>` - Builds a **Dockerfile** from specified feature file.

**NOTE:** build command gives a sample command of how to run `docker
build` 

## Run docker build
Execute `docker build` command to actually create image

---
## Installation and Configuration
All set configuration will be stored ` ~/.pazuzu/config`

-  Setup snippets provider:

    ```
    $ pazuzu config set git.url git@github.com:zalando-incubator/pazuzu.git
    $ pazuzu config set github.url https://github.com/zalando-incubator/pazuzu.git
    ```
- Setup base image

  ```
  $ pazuzu config set base-image ubuntu:16.04
  ```
  
## Helpers

- Switch on verbose mode using `-v/--verbose`:
    ```
	$ pazuzu -v compose node npm
	```
- Getting help message:
	```
	$ pazuzu help
	NAME:
	   pazuzu - Build Docker features from pazuzu-registry
	
	USAGE:
	   pazuzu [global options] command [command options] [arguments...]
	
	VERSION:
	   0.1
	
	COMMANDS:
	     search   search for features in registry
	     compose  Compose Pazuzufile out of the selected features
	     build    build Dockerfile out of Pazuzufile
	     config   Configure pazuzu
	     help, h  Shows a list of commands or help for one command
	
	GLOBAL OPTIONS:
	   --verbose, -v  Verbose output
	   --help, -h     show help
	   --version      Print version
	
	```

---
License
---

The MIT License (MIT)
Copyright © 2016 Zalando SE, https://tech.zalando.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the “Software”), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
