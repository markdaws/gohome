#!/bin/bash

git pull
go install -v -ldflags "-X main.VERSION=`cat ./VERSION`" ./cmd/...
