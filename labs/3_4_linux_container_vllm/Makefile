SHELL := /usr/bin/env bash

LAB_CONTAINER ?= linux-container-lab
LAB_IMAGE ?= m.daocloud.io/docker.io/library/busybox:1.36.1
VLLM_CONTAINER ?= vllm-lab

.PHONY: start inspect shell stop vllm-start vllm-inspect vllm-logs vllm-stop

start:
	docker run --detach --rm \
		--name "$(LAB_CONTAINER)" \
		--hostname container-lab \
		--cpus 0.5 \
		--memory 128m \
		--pids-limit 64 \
		--shm-size 64m \
		"$(LAB_IMAGE)" \
		sh -c 'while true; do sleep 3600; done'

inspect:
	./scripts/inspect-container.sh "$(LAB_CONTAINER)"

shell:
	docker exec -it "$(LAB_CONTAINER)" sh

stop:
	docker stop "$(LAB_CONTAINER)"

vllm-start:
	VLLM_CONTAINER="$(VLLM_CONTAINER)" ./scripts/run-vllm.sh

vllm-inspect:
	./scripts/inspect-container.sh "$(VLLM_CONTAINER)"

vllm-logs:
	docker logs --follow "$(VLLM_CONTAINER)"

vllm-stop:
	docker stop "$(VLLM_CONTAINER)"
