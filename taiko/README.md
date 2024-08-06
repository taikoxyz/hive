# taiko

1. Support bindings which are used in simulator/taiko.
2. Support scripts to deploy taiko contract and get txs and env.
3. Support other common taiko go modules.

## Update dependencies

```shell
make update
```

## Update bindings

```shell
cd taiko
sh scripts/abigen.sh
```

## Get txs and env for taiko contract.

* build hive/clients/taiko/geth:latest docker image

```shell
cd taiko
docker build --no-cache -t docker/geth:latest ./docker/geth
```

* start geth

```shell
docker run --name taiko_contracts_geth --rm -it -p 8545:8545 -p 8546:8546 -p 3500:3500 docker/geth:latest
```

* deploy taiko contract

```shell
cd taiko
sh scripts/deploy_l1_contract.sh
```

* get txs

```shell
cd taiko
sh scripts/get_txs.sh
```

* get env

```shell
cd taiko
sh scripts/get_env.sh
```

* Exit the running geth container and then the two files(.env l1contract_txs.txt) are what we want to get.

```shell
docker stop taiko_contracts_geth
```