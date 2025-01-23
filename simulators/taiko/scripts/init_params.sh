#!/bin/bash

# taiko default version
TAIKO_VERSION=${TAIKO_VERSION:-ontake}

case "$TAIKO_VERSION" in
    ontake)
        echo "TAIKO_VERSION is $TAIKO_VERSION"
        ./scripts/$TAIKO_VERSION/deploy_l1_contract.sh
        ./scripts/$TAIKO_VERSION/abigen.sh
        ./scripts/$TAIKO_VERSION/get_env.sh
        ;;
    pacaya)
        echo "TAIKO_VERSION is $TAIKO_VERSION"
        ./scripts/$TAIKO_VERSION/deploy_l1_contract.sh
        ./scripts/$TAIKO_VERSION/abigen.sh
        ./scripts/$TAIKO_VERSION/get_env.sh
        ;;
    *)
        echo "Invalid TAIKO_VERSION: $TAIKO_VERSION"
        exit 1
        ;;
esac