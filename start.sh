#!/usr/bin/env bash

export GOPATH='/data/DevelopmentRoot/GoLang'
echo -ne "\ec" >> /tmp/compile-dashboard-temp.log
go run *.go
