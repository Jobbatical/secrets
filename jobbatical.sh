#!/usr/bin/env bash

V=$(cat ./VERSION)

source build-all

cd target/$V

echo "Execute file decryption using binary"

SECRETS_KEY_LOCATION=global SECRETS_KEY_RING=immi-project-secrets SECRETS_REPO_HOST=github.com SECRETS_ORG=jobbatical ./secrets-darwin-amd64 open --open-all