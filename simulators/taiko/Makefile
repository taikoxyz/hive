.PHONY: taiko log_geth log_driver log_proposer log_prover

all: taiko

deploy_contracts:
	sh scripts/deploy_l1_contract.sh

taiko:
	docker build --no-cache -t hive/simulators/taiko:latest .

hive_prysmctl:
	docker build --no-cache -t hive/simulators/hive_prysmctl:latest .

log_geth:
	docker ps -a | grep "taiko-geth" | awk '{print $1}' | xargs docker logs -f

log_taiko:
	docker ps -a | grep "taiko:latest" | awk '{print $1}' | xargs docker logs -f

log_driver:
	docker ps -a | grep "driver" | awk '{print $1}' | xargs docker logs -f

log_proposer:
	docker ps -a | grep "proposer" | awk '{print $1}' | xargs docker logs -f

log_prover:
	docker ps -a | grep "prover" | awk '{print $1}' | xargs docker logs -f