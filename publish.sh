#!/usr/bin/env bash
set -euox pipefail

PRIVATE_ECR_URL="509869530682.dkr.ecr.us-west-2.amazonaws.com"
aws ecr get-login-password --region us-west-2 | ko login --username AWS --password-stdin $PRIVATE_ECR_URL

export KO_DOCKER_REPO="$PRIVATE_ECR_URL/internal"
ko build -t latest -B ./cmd/demo-w
