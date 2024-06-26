#!/bin/bash

# Startup script to initialize and boot a go-ethereum instance.
#
# This script assumes the following files:
#  - `geth` binary is located in the filesystem root
#  - `genesis.json` file is located in the filesystem root (mandatory)
#  - `chain.rlp` file is located in the filesystem root (optional)
#  - `blocks` folder is located in the filesystem root (optional)
#  - `keys` folder is located in the filesystem root (optional)
#
# This script assumes the following environment variables:
#
#  - HIVE_BOOTNODE                enode URL of the remote bootstrap node
#  - HIVE_NETWORK_ID              network ID number to use for the eth protocol
#  - HIVE_NODETYPE                sync and pruning selector (archive, full, light)
#
# Forks:
#
#  - HIVE_FORK_HOMESTEAD          block number of the homestead hard-fork transition
#  - HIVE_FORK_DAO_BLOCK          block number of the DAO hard-fork transition
#  - HIVE_FORK_DAO_VOTE           whether the node support (or opposes) the DAO fork
#  - HIVE_FORK_TANGERINE          block number of Tangerine Whistle transition
#  - HIVE_FORK_SPURIOUS           block number of Spurious Dragon transition
#  - HIVE_FORK_BYZANTIUM          block number for Byzantium transition
#  - HIVE_FORK_CONSTANTINOPLE     block number for Constantinople transition
#  - HIVE_FORK_PETERSBURG         block number for ConstantinopleFix/PetersBurg transition
#  - HIVE_FORK_ISTANBUL           block number for Istanbul transition
#  - HIVE_FORK_MUIRGLACIER        block number for Muir Glacier transition
#  - HIVE_FORK_BERLIN             block number for Berlin transition
#  - HIVE_FORK_LONDON             block number for London
#
# Clique PoA:
#
#  - HIVE_CLIQUE_PERIOD           enables clique support. value is block time in seconds.
#  - HIVE_CLIQUE_PRIVATEKEY       private key for clique mining
#
# Other:
#
#  - HIVE_MINER                   enable mining. value is coinbase address.
#  - HIVE_MINER_EXTRA             extra-data field to set for newly minted blocks
#  - HIVE_LOGLEVEL                client loglevel (0-5)
#  - HIVE_GRAPHQL_ENABLED         enables graphql on port 8545
#  - HIVE_LES_SERVER              set to '1' to enable LES server

# Set default values
HIVE_LOGLEVEL=${HIVE_LOGLEVEL:-3}
HIVE_NETWORK_ID=${HIVE_NETWORK_ID:-167001}

# Immediately abort the script on any error encountered
set -e

FLAGS="--datadir /geth/data --state.scheme=path"
# Configure http.
FLAGS="$FLAGS --http --http.api=admin,debug,eth,miner,net,personal,txpool,web3"
# Configure ws.
FLAGS="$FLAGS --ws --ws.api=admin,debug,eth,miner,net,personal,txpool,web3"
# Configure auth rpc
FLAGS="$FLAGS --authrpc.port 8551 --allow-insecure-unlock --authrpc.jwtsecret /geth/jwt.hex"
# Configure log level
FLAGS="$FLAGS --verbosity=$HIVE_LOGLEVEL"
# If a specific network ID is requested, use that
FLAGS="$FLAGS --networkid $HIVE_NETWORK_ID"

# It doesn't make sense to dial out, use only a pre-set bootnode.
if [ "$HIVE_BOOTNODE" != "" ]; then
  FLAGS="$FLAGS --bootnodes=$HIVE_BOOTNODE"
else
  FLAGS="$FLAGS --nodiscover"
fi

# Handle any client mode or operation requests
if [ "$HIVE_NODETYPE" == "archive" ]; then
  FLAGS="$FLAGS --syncmode full --gcmode archive"
fi
if [ "$HIVE_NODETYPE" == "full" ]; then
  FLAGS="$FLAGS --syncmode full"
fi
if [ "$HIVE_NODETYPE" == "snap" ]; then
  FLAGS="$FLAGS --syncmode snap"
fi

# Don't immediately abort, some imports are meant to fail
set +e

echo "Running taiko-geth with flags $FLAGS"
geth "${FLAGS}"
