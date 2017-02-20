#!/bin/bash -x
set -e

cd ${TRAVIS_BUILD_DIR:?"Required ENV variable: TRAVIS_BUILD_DIR"}
VERSION=${TRAVIS_TAG?"Required ENV variable: TRAVIS_TAG"}
RELEASE=pazuzu_${VERSION#v}

mkdir -p $RELEASE/{darwin,linux,windows}_amd64/

for i in darwin linux ; do
    GOOS=$i GOARCH=amd64 go build -v ./cli/pazuzu
    mv pazuzu ${RELEASE}/${i}_amd64
done

GOOS=windows GOARCH=amd64 go build -v ./cli/pazuzu
mv pazuzu.exe ${RELEASE}/windows_amd64

zip -r ${RELEASE}.zip ${RELEASE}

fi
