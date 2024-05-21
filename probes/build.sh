#!/bin/bash

docker build . -t wreulicke/probes-tests
docker tag wreulicke/probes-tests localhost:5001/wreulicke/probes-tests
docker push localhost:5001/wreulicke/probes-tests