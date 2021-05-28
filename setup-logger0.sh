#!/bin/bash

set -e
set -x

cd ../logger0
export $(cat .env)
export LOGGER0_CONTROL_TOKEN=$(go run ./cmd/logger0 gen-token)
export LOGGER0_CONTROL_URL=http://logger0.a0core.net:9091

cd ../auth0-cli

set +x
set +e
