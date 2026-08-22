#!/usr/bin/env bash
set -euo pipefail

readonly VLLM_CONTAINER="${VLLM_CONTAINER:-vllm-lab}"
readonly VLLM_IMAGE="${VLLM_IMAGE:-m.daocloud.io/docker.io/vllm/vllm-openai:latest}"
readonly VLLM_MODEL="${VLLM_MODEL:-Qwen/Qwen3-0.6B}"
readonly VLLM_PORT="${VLLM_PORT:-8000}"
readonly VLLM_CPUS="${VLLM_CPUS:-4}"
readonly VLLM_MEMORY="${VLLM_MEMORY:-16g}"
readonly VLLM_SHM_SIZE="${VLLM_SHM_SIZE:-1g}"
readonly HF_CACHE_DIR="${HF_CACHE_DIR:-$(pwd)/.cache/huggingface}"

if docker inspect "$VLLM_CONTAINER" >/dev/null 2>&1; then
  echo "container already exists: $VLLM_CONTAINER" >&2
  exit 1
fi

mkdir -p "$HF_CACHE_DIR"

docker_args=(
  run --detach --rm
  --name "$VLLM_CONTAINER"
  --gpus all
  --cpus "$VLLM_CPUS"
  --memory "$VLLM_MEMORY"
  --pids-limit 512
  --shm-size "$VLLM_SHM_SIZE"
  --publish "$VLLM_PORT:8000"
  --volume "$HF_CACHE_DIR:/root/.cache/huggingface"
)

if [[ -n "${HF_TOKEN:-}" ]]; then
  docker_args+=(--env HF_TOKEN)
fi

docker_args+=(
  "$VLLM_IMAGE"
  --model "$VLLM_MODEL"
  --host 0.0.0.0
)

docker "${docker_args[@]}"

printf 'vLLM is starting on http://127.0.0.1:%s\n' "$VLLM_PORT"
printf 'follow startup: docker logs --follow %s\n' "$VLLM_CONTAINER"
printf 'inspect Linux boundaries: make vllm-inspect\n'
