#!/bin/bash

docker build . -t wreulicke/slow-deployment-test
docker tag wreulicke/slow-deployment-test localhost:5001/wreulicke/slow-deployment-test
docker push localhost:5001/wreulicke/slow-deployment-test