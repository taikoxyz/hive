.PHONY: hive_alpine_dependency hive_taiko_dependency taiko

all: hive_alpine_dependency hive_taiko_dependency taiko

hive_alpine_dependency:
	docker build --no-cache -t hive_alpine_dependency:latest --target hive_alpine_dependency .

hive_taiko_dependency:
	docker build --no-cache -t hive_taiko_dependency:latest --target hive_taiko_dependency .

taiko:
	docker build --no-cache -t hive/simulators/taiko:latest .
