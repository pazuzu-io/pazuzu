# pazuzu
[![Travis BuildStatus](https://travis-ci.org/zalando-incubator/pazuzu.svg?branch=master)](https://travis-ci.org/zalando-incubator/pazuzu)
[![Stories in Ready](https://badge.waffle.io/zalando/pazuzu.png?label=ready&title=Ready)](https://waffle.io/zalando/pazuzu)


# What is Pazuzu?
Pazuzu is a tool that builds Docker images from feature snippets, while
resolving all dependencies between them. One of the common use cases is
Continuous Integration environment, where jobs require specific tooling present
for building and testing. Pazuzu can significantly ease that process, by
letting user choose from a wide selection of predefined Dockerfile snippets
that represent those dependencies (e.g. Golang, Python, Android SDK, customized
NPM installs).


## Building Pazuzu
1. Make sure you setup Go according to: https://golang.org/doc/install#install

2. Clone Pazuzu repository:
  ```bash
  git clone git@github.com:zalando-incubator/pazuzu.git  $GOPATH/src/github.com/zalando-incubator/pazuzu
  ```

3. Install govendor and go-swagger:
  ```bash
  go get -u github.com/kardianos/govendor
  ```
  ```bash
  go get -u github.com/go-swagger/go-swagger/cmd/swagger
  ```

4. Sync vendor packages:
  ```bash
  cd $GOPATH/src/github.com/zalando-incubator/pazuzu/
  govendor sync
  ```

5. Generate code from swagger.yaml:
  ```bash
  cd $GOPATH/src/github.com/zalando-incubator/pazuzu/swagger
  ```
  ```bash
  swagger generate client -f swagger.yaml
  ```

6. Make sure that the tests are passing:
  ```bash
  cd $GOPATH/src/github.com/zalando-incubator/pazuzu/
  ```
  ```bash
  go get -t -v
  ```
  ```bash
  go test -v ./...
  ```

7. Build command-line utilities:
  ```bash
  cd $GOPATH/src/github.com/zalando-incubator/pazuzu/cli/pazuzu
  ```
  ```bash
  go build
  ```

8. Install pazuzu command globally [Optional]:
  ```bash
  go install
  ```

## Usage

Basically, pazuzu CLI tool has 4 subcommands:
- `config` - configure pazuzu tool
- `project` - configure and build project image
- `search` - search for available features inside the repository

### Search features

`pazuzu search` is used to check which features are actually available in configured repository (substring matching is used to match features names):

  ```bash
  pazuzu search node
  ```
  ```bash
  pazuzu search ja
  ```

### Configure project features

`pazuzu project` command is used to configure the project definition.

  ```bash
  pazuzu project add node,java
  ```

If `Pazuzufile` already exists in the directory, `pazuzu project` takes it as a base. If not, it generates
the new one based on the given features and default base image specified in the configuration.

`add` subcommand adds features to an existing set of features.
`remove` subcommand removes features from an existing set of features.
`list` subcommand lists an existing set of features.

  ```bash
  pazuzu project add node,java  # add node and java features to a project
  ```
  ```bash
  pazuzu project remove java    # remove java from a project
  ```

`-d` (or `--directory`) option sets the working directory to a specified path.

  ```bash
  pazuzu project add node -d /tmp
  ```

  In the given example, Node.js feature will be added to the list of features specified in `/tmp/Pazuzufile`
  (if it exists) and the output files will be saved back to `/tmp/`

### Clean

`pazuzu project clean` step removes `Pazuzufile`, `Dockerfile` and `test.bats`.

  ```bash
  pazuzu project clean
  ```

### Build Docker image

`pazuzu project build` is responsible for a final step - building and validating the Docker image.

```
pazuzu project build -n hellodocker -d /tmp
```

`-n` (or `--name`) option sets the name for the created Docker image.

`-d` (or `--directory`) option sets the working directory where `Dockerfile` is located.

### Configuration

`pazuzu config` provides a set of tools to configure pazuzu CLI. Configurations are stored in ` ~/pazuzu-cli.yaml` .

```bash
pazuzu config list  # lists all configurations
```
```bash
pazuzu config set registry.port 8080 # sets value of registy.port parameter
```
```bash
pazuzu confi get registy.port  # gets value of registy.port parameter
```

## Initial setup

Currenly pazuzu supports only registry as a storage.

See: [Pazuzu Registry](https://github.com/zalando-incubator/pazuzu-registry)

### Base image

Base image can be also set using `pazuzu config` command.

```bash
pazuzu config set base ubuntu:16.04
```

## Helpers

- Switch on verbose mode using `-v/--verbose`:
  ```bash
	pazuzu -v project add node,npm
	```
- Getting help message:
	```bash
	pazuzu help

	NAME:
	   pazuzu - Build Docker features from pazuzu-registry
     ...
  ```

## Development environment installation (macOS)

- Check if Go is installed:
  ```bash
  $ go version
  ```

- Set up `\$GOPATH` variable in your profile (e.g. to `~/go`)

- Clone Pazuzu repository:
  ```bash
  git clone git@github.com:zalando-incubator/pazuzu.git  $GOPATH/src/github.com/zalando-incubator/pazuzu
  ```
- Build command-line tools
  ```bash
  cd $GOPATH/src/github.com/zalando-incubator/pazuzu/cli/pazuzu
  ```
  ```bash
  go build
  ```
- Run tests
  ```bash
  cd $GOPATH/src/github.com/zalando-incubator/pazuzu/cli/pazuzu
  ```
  ```bash
  go test ./...
  ```


## License

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
