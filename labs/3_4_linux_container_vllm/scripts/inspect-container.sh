#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: $0 CONTAINER" >&2
  exit 2
fi

readonly CONTAINER_NAME="$1"

if ! docker inspect "$CONTAINER_NAME" >/dev/null 2>&1; then
  echo "container not found: $CONTAINER_NAME" >&2
  exit 1
fi

readonly CONTAINER_PID="$(docker inspect --format '{{.State.Pid}}' "$CONTAINER_NAME")"
readonly IMAGE_NAME="$(docker inspect --format '{{.Config.Image}}' "$CONTAINER_NAME")"

print_heading() {
  printf '\n== %s ==\n' "$1"
}

print_namespaces() {
  printf '%-9s %-22s %-22s\n' "namespace" "host shell" "container init"
  for namespace_name in cgroup ipc mnt net pid user uts; do
    printf '%-9s %-22s %-22s\n' \
      "$namespace_name" \
      "$(readlink "/proc/self/ns/$namespace_name")" \
      "$(docker exec "$CONTAINER_NAME" readlink "/proc/self/ns/$namespace_name")"
  done
}

print_cgroup_file() {
  local cgroup_root="$1"
  local filename="$2"
  local filepath="$cgroup_root/$filename"

  if [[ -r "$filepath" ]]; then
    printf '%-18s %s\n' "$filename" "$(<"$filepath")"
  fi
}

print_cgroup_limits() {
  local cgroup_path
  cgroup_path="$(awk -F: '$1 == "0" { print $3 }' "/proc/$CONTAINER_PID/cgroup")"
  local cgroup_root="/sys/fs/cgroup$cgroup_path"

  printf 'path: %s\n' "$cgroup_path"
  print_cgroup_file "$cgroup_root" cpu.max
  print_cgroup_file "$cgroup_root" memory.max
  print_cgroup_file "$cgroup_root" memory.current
  print_cgroup_file "$cgroup_root" pids.max
  print_cgroup_file "$cgroup_root" pids.current
}

print_heading "identity"
printf 'name:  %s\nimage: %s\nhost PID: %s\n' \
  "$CONTAINER_NAME" "$IMAGE_NAME" "$CONTAINER_PID"

print_heading "namespaces"
print_namespaces

print_heading "cgroup v2"
cat "/proc/$CONTAINER_PID/cgroup"
print_cgroup_limits

print_heading "configured resource limits"
docker inspect --format \
  'nano CPUs={{.HostConfig.NanoCpus}} memory={{.HostConfig.Memory}} pids={{.HostConfig.PidsLimit}} shm={{.HostConfig.ShmSize}}' \
  "$CONTAINER_NAME"

print_heading "filesystem and mounts"
docker inspect --format 'snapshotter/storage driver={{.Driver}}' "$CONTAINER_NAME"
docker inspect --format \
  '{{range .Mounts}}{{printf "%s -> %s (%s, %s)\n" .Source .Destination .Type .Mode}}{{end}}' \
  "$CONTAINER_NAME"
docker exec "$CONTAINER_NAME" df -h /dev/shm 2>/dev/null || true

print_heading "network"
docker inspect --format \
  '{{range $name, $network := .NetworkSettings.Networks}}{{printf "%s ip=%s gateway=%s\n" $name $network.IPAddress $network.Gateway}}{{end}}' \
  "$CONTAINER_NAME"
docker inspect --format '{{json .NetworkSettings.Ports}}' "$CONTAINER_NAME"

print_heading "isolation and security"
docker inspect --format \
  'pid={{.HostConfig.PidMode}} ipc={{.HostConfig.IpcMode}} network={{.HostConfig.NetworkMode}} privileged={{.HostConfig.Privileged}}' \
  "$CONTAINER_NAME"
docker inspect --format \
  'cap-add={{json .HostConfig.CapAdd}} cap-drop={{json .HostConfig.CapDrop}} security-opt={{json .HostConfig.SecurityOpt}}' \
  "$CONTAINER_NAME"
docker exec "$CONTAINER_NAME" grep -E \
  '^(Cap(Inh|Prm|Eff|Bnd|Amb)|NoNewPrivs|Seccomp):' /proc/1/status

print_heading "image layers"
docker image inspect --format '{{range .RootFS.Layers}}{{println .}}{{end}}' "$IMAGE_NAME"
