#!/bin/bash

go get github.com/jfeliu007/goplantuml/parser
go get github.com/jfeliu007/goplantuml/cmd/goplantuml
sudo apt install graphviz default-jre-headless -y

dir="$(go env GOPATH)/src/github.com/jfeliu007/goplantuml"

go run $dir/cmd/goplantuml/main.go \
-ignore ../../unused -ignore ../../experimental \
-recursive \
-show-aggregations -show-aliases -show-compositions -show-connection-labels -show-implementations \
../../ | java -jar ./plantuml.jar -pipe > out.png
