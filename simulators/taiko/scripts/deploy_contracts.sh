#!/bin/bash

# stop docker compose
docker compose -f docker/docker-compose.yml up l1_node l2_pacaya -d --wait
trap "docker compose -f docker/docker-compose.yml down" EXIT SIGINT SIGTERM ERR

sh ./scripts/ontake/deploy_l1_contract.sh
sh ./scripts/pacaya/upgrade_l1_pacaya.sh

# Get env
sh scripts/get_env.sh
# Get txs
sh scripts/get_txs.sh


