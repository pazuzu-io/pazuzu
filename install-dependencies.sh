#!/usr/bin/env sh
PAZUZU_DIR=$GOPATH/src/github.com/zalando-incubator/pazuzu
SWAGGER_DIR=$PAZUZU_DIR/swagger
SWAGGER_FILE=$SWAGGER_DIR/swagger.yaml
PAZUZU_REGISTRY_DIR=$PAZUZU_DIR/../pazuzu-registry/
PAZUZU_REGISTRY_SWAGGER_LOCAL=$PAZUZU_REGISTRY_DIR/src/main/resources/api/swagger.yaml
PAZUZU_REGISTRY_SWAGGER_UPSTREAM="https://raw.githubusercontent.com/zalando-incubator/pazuzu-registry/master/src/main/resources/api/swagger.yaml"

get_swagger_api_definition() {
  if [ ! -d $SWAGGER_DIR ]
  then
    mkdir $SWAGGER_DIR
  fi

  if [ -z $SWAGGER_LOCAL ]
  then
    # by default try to get swagger definition from local copy of pazuzu-registry repo
    SWAGGER_LOCAL=1
  fi
  if [ 1 -eq $SWAGGER_LOCAL && -f $PAZUZU_REGISTRY_SWAGGER_LOCAL ]
  then
    cp $PAZUZU_REGISTRY_SWAGGER_LOCAL $SWAGGER_FILE
  else
    curl $PAZUZU_REGISTRY_SWAGGER_UPSTREAM -o $SWAGGER_FILE
  fi
}

go get -u github.com/kardianos/govendor
go get -u github.com/go-swagger/go-swagger/cmd/swagger

cd $PAZUZU_DIR
govendor sync

get_swagger_api_definition

cd $SWAGGER_DIR
swagger generate client -f swagger.yaml
