#!/usr/bin/env sh
go get -u github.com/kardianos/govendor
go get -u github.com/go-swagger/go-swagger/cmd/swagger
cd $GOPATH/src/github.com/zalando-incubator/pazuzu/
govendor sync
cd $GOPATH/src/github.com/zalando-incubator/pazuzu/swagger
swagger generate client -f swagger.yaml
