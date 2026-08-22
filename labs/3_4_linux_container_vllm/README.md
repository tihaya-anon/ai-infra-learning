# Linux Container Boundaries to vLLM Lab

这个 Lab 先用一个很小的 BusyBox 容器建立 Linux 容器化的基础模型，再把同一组观察方法应用到 vLLM。重点不是记 Docker 命令，而是能回答：vLLM 进程被哪些 Linux 机制隔离和限制，这些机制怎样影响模型加载、共享内存、请求处理和故障表现。

## 学习目标

完成后应能：

- 解释容器为什么是受 namespace、cgroup、mount 和安全策略约束的宿主机进程，而不是轻量虚拟机。
- 从容器找到宿主机 PID，并比较 PID、mount、network、IPC 和 UTS namespace。
- 读取 cgroup v2 的 CPU、内存和 PID 限额，区分配置值、当前使用量和 OOM 行为。
- 解释 OverlayFS、bind mount、端口映射和 `/dev/shm` 的数据路径与生命周期。
- 把以上机制映射到 vLLM 的 API Server、Engine/Worker、模型权重、Tokenizer、KV Cache 和 GPU 设备。
- 用证据判断 vLLM 问题来自应用配置、Linux 资源边界还是 GPU/驱动层。

## 前置条件

基础部分不需要 GPU：

- Linux，使用 cgroup v2。
- Docker Engine、GNU Make、Bash、`readlink` 和 `awk`。
- 当前用户可以访问 Docker daemon。

vLLM 扩展部分额外需要 NVIDIA GPU、兼容的宿主机驱动、NVIDIA Container Toolkit，以及足够容纳模型权重、CUDA Context 和 KV Cache 的显存。首次运行还要访问镜像仓库和模型仓库。

```bash
cd labs/3_4_linux_container_vllm
```

## 心智模型

容器的核心对象仍然是 Linux 进程：

```text
Docker CLI
  -> Docker daemon / OCI runtime
    -> clone/unshare + mount + cgroup + capabilities + seccomp
      -> 宿主机上的容器 init 进程
        -> vLLM API Server / Engine / Worker 进程
          -> CUDA driver -> GPU
```

| Linux 机制 | 它控制什么 | 对 vLLM 的直接影响 |
| --- | --- | --- |
| PID namespace | 进程编号和可见性 | 容器内 PID 1、API Server 与 Worker 的进程关系 |
| Mount namespace | 文件系统视图 | 镜像层、模型缓存、配置、日志和临时文件 |
| Network namespace | 网卡、路由和端口 | `8000` 监听、端口映射、Service/Gateway 数据路径 |
| IPC namespace 与 tmpfs | IPC 对象和 `/dev/shm` | PyTorch 多进程共享数据，尤其是 Tensor Parallel |
| cgroup v2 | CPU、内存、PID 和 I/O 的计量与限制 | Tokenization、调度线程、Pinned Memory、OOM 和 Throttling |
| Device 与安全策略 | 设备访问、syscall 和进程权限 | `/dev/nvidia*`、驱动调用、调试能力和攻击面 |
| OverlayFS 与 volume | 只读镜像层和可写数据 | 大镜像、模型权重缓存、冷启动和数据保留 |

## 1. 建立基础容器

启动一个有明确资源上限的容器：

```bash
make start
docker ps --filter name=linux-container-lab
make inspect
```

它有 `0.5 CPU`、`128 MiB` 内存、`64` 个 PID 和 `64 MiB /dev/shm`。这些不是 BusyBox 自己实现的限制，而是 Docker 将配置转换成 OCI Runtime 配置，最终落实到 namespace、mount 和 cgroup。

保留 `make inspect` 的输出。后面运行 vLLM 时会使用同一个检查器，直接对比两个容器的 Linux 边界。

## 2. 进程与 namespace

比较容器内外看到的 PID：

```bash
docker exec linux-container-lab ps
docker inspect --format '{{.State.Pid}}' linux-container-lab
```

第一个命令通常把容器 init 显示为 PID 1，第二个命令给出同一进程在宿主机 PID namespace 中的 PID。进程没有复制两份；不同 namespace 只是给同一个内核 task 提供不同视图。

`make inspect` 输出中：

- `pid`、`mnt`、`net`、`ipc` 和 `uts` 通常与当前宿主机 shell 不同。
- `user` 是否不同取决于 Docker 是否启用 user namespace/remap 或 rootless 模式。
- namespace 负责隔离“看见什么”，不负责限制“最多能用多少”；资源限制属于 cgroup。

进入该容器的 network namespace：

```bash
CONTAINER_PID=$(docker inspect --format '{{.State.Pid}}' linux-container-lab)
sudo nsenter --target "$CONTAINER_PID" --net ip address
sudo nsenter --target "$CONTAINER_PID" --net ip route
```

思考：为什么 `docker exec` 和 `nsenter` 都能观察容器，但前者通过 Runtime 执行新进程，后者直接让进程加入目标 namespace？

### 与 vLLM 的关系

vLLM 容器内并非只有一个抽象的“模型进程”。API Server、Engine 和分布式 Worker 可能形成多进程结构。排障时先确定容器 PID、宿主机 PID、线程和 GPU Context 的对应关系，才能把 `perf`、eBPF、CPU 调度和 GPU 指标关联到正确执行单元。

## 3. cgroup：资源限制不是隔离视图

查看 cgroup v2 的生效值：

```bash
make inspect
docker stats --no-stream linux-container-lab
```

关注：

- `cpu.max`：配额和周期；`50000 100000` 表示每 `100 ms` 最多运行 `50 ms`，即 `0.5 CPU`。
- `memory.max`：硬上限；`memory.current` 是当前记账量。
- `pids.max`：可创建的 task 数上限，线程也计入 PID controller。
- `max` 表示本层没有有限上限，不代表祖先 cgroup 或宿主机没有约束。

观察 CPU Throttling：

```bash
docker exec -d linux-container-lab sh -c 'while true; do :; done'
sleep 3
docker stats --no-stream linux-container-lab

CONTAINER_PID=$(docker inspect --format '{{.State.Pid}}' linux-container-lab)
CGROUP_PATH=$(awk -F: '$1 == "0" {print $3}' "/proc/$CONTAINER_PID/cgroup")
cat "/sys/fs/cgroup$CGROUP_PATH/cpu.stat"
```

容器应接近 `50%` CPU，而不是占满一个完整核心；`cpu.stat` 中的 `nr_throttled` 和 `throttled_usec` 应随负载增长。

### 与 vLLM 的关系

GPU 推理不等于 CPU 不重要。HTTP 解析、Chat Template、Tokenization、请求调度、输出解码和多进程协调都消耗 CPU。CPU Throttling 可能表现为 GPU 出现空洞、TTFT/TPOT 上升，但显存和 SM 利用率并未达到预期。容器内看到很多 CPU 也不代表 cgroup 配额允许持续使用它们。

内存必须分层理解：

- cgroup `memory.max` 主要约束 Host RAM，不是 GPU HBM。
- 模型权重可能同时经过 page cache、Host RAM 和 GPU HBM。
- Pinned Memory、Tokenizer、Python Heap 和共享内存都会增加 Host RAM 压力。
- GPU OOM 与 cgroup OOM 是不同故障，应分别查容器状态、内核事件和 GPU 日志。

## 4. 镜像层、可写层与模型缓存

查看镜像层和容器存储驱动：

```bash
docker image inspect m.daocloud.io/docker.io/library/busybox:1.36.1 \
  --format '{{json .RootFS.Layers}}' | jq
docker inspect linux-container-lab --format 'snapshotter/storage driver={{.Driver}}'
```

验证容器可写层的生命周期：

```bash
docker exec linux-container-lab sh -c 'echo transient > /tmp/lifecycle.txt'
docker exec linux-container-lab cat /tmp/lifecycle.txt
make stop
make start
docker exec linux-container-lab test ! -e /tmp/lifecycle.txt
```

容器被删除后，可写层随之消失。镜像层仍在本地缓存，但运行时写入的数据不能被当作持久存储。

### 与 vLLM 的关系

vLLM 镜像包含用户态 CUDA/Python 依赖，模型缓存则通过 bind mount 保留在容器生命周期之外：

```text
registry -> vLLM image layers -> container rootfs
model registry -> host model cache -> bind mount -> model loader -> RAM/HBM
```

这形成两个不同的冷启动阶段：拉取镜像与下载/加载模型。优化前必须分别计时。模型缓存还需要容量、并发下载、模型 revision、一致性和淘汰策略，不能只做一个永久不清理的 HostPath。

## 5. 网络与端口映射

基础容器默认加入 Docker bridge，但没有发布端口：

```bash
make inspect
docker network inspect bridge
```

端口发布 `-p 8000:8000` 不会让进程自动监听。应用仍须监听容器内的 `0.0.0.0:8000`；如果只监听 `127.0.0.1`，宿主机的转发路径无法连接它。

### 与 vLLM 的关系

vLLM 的 OpenAI-compatible Server 是 HTTP 数据面入口。在 Kubernetes 中，这条路径通常是：

```text
Client -> Gateway/Ingress -> Service -> Pod network namespace -> vLLM :8000
```

认证、租户限流、重试和公网 TLS 不应仅依赖 vLLM 进程。排查延迟时也要区分客户端排队、Gateway、Service 转发、容器网络和 Engine 内部排队。

## 6. IPC 与 `/dev/shm`

查看基础容器独立的 tmpfs：

```bash
docker exec linux-container-lab df -h /dev/shm
docker inspect --format 'ipc={{.HostConfig.IpcMode}} shm={{.HostConfig.ShmSize}}' linux-container-lab
```

本 Lab 使用私有 IPC namespace 和 `64 MiB` 的 `/dev/shm`。共享内存不是普通镜像层文件，它通常挂载为 tmpfs，并受 Host RAM 和 cgroup 内存共同影响。

### 与 vLLM 的关系

vLLM 使用 PyTorch；多进程执行，尤其 Tensor Parallel，可能使用共享内存交换数据。官方容器文档建议使用 `--ipc=host` 或显式增大 `--shm-size`。本 Lab 的 vLLM 脚本选择私有 IPC namespace 加 `--shm-size=1g`，保留隔离边界，同时让容量可见、可控。

| 配置 | 优点 | 风险/成本 |
| --- | --- | --- |
| `--ipc=host` | 使用宿主机 `/dev/shm`，配置简单 | 共享宿主机 IPC namespace，隔离更弱 |
| `--shm-size=1g` | 保留私有 IPC namespace，容量明确 | 需要按 Worker/并行规模估算并验证 |

不能把增大 `/dev/shm` 当成所有 OOM 的修复。KV Cache 主要在 GPU HBM，cgroup Host Memory、shared memory 和 GPU HBM 是三套相关但不同的容量边界。

## 7. 安全边界

`make inspect` 会显示 capabilities、seccomp、privileged 和 namespace mode。默认 Docker 容器通常仍保留一组 capabilities，并应用 seccomp profile；容器内 UID 0 也不自动等于宿主机完整 root 权限。

不要为了排障长期使用 `--privileged`、`--pid=host`、`--network=host` 和 `--ipc=host`。这些选项可能快速绕过症状，同时移除了需要理解的隔离边界。应先确认具体缺少的设备、syscall、capability 或共享内存容量。

### 与 vLLM 的关系

NVIDIA Container Toolkit 把 GPU 设备和匹配的驱动能力暴露给容器。镜像中的 CUDA Runtime 不能替代宿主机内核驱动。出现 `CUDA unavailable` 时至少区分：

1. 宿主机驱动和 GPU 是否正常。
2. Docker 是否注册 GPU runtime/CDI 设备。
3. 容器是否获得 `/dev/nvidia*` 与所需库。
4. vLLM/PyTorch/CUDA 与 GPU 架构是否兼容。

## 8. 启动并观察 vLLM（需要 NVIDIA GPU）

默认使用官方镜像和小模型：

```bash
make vllm-start
make vllm-logs
```

另一个终端检查 API：

```bash
curl http://127.0.0.1:8000/v1/models | jq
curl http://127.0.0.1:8000/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{
    "model": "Qwen/Qwen3-0.6B",
    "messages": [{"role": "user", "content": "Explain PID namespaces in one sentence."}],
    "temperature": 0,
    "max_tokens": 64
  }' | jq
```

再运行与基础容器相同的检查：

```bash
make vllm-inspect
docker top vllm-lab -eo pid,ppid,nlwp,pcpu,pmem,comm,args
curl http://127.0.0.1:8000/metrics | head -80
```

启动参数可通过环境变量覆盖：

```bash
VLLM_IMAGE=m.daocloud.io/docker.io/vllm/vllm-openai:<固定版本> \
VLLM_MODEL=Qwen/Qwen3-0.6B \
VLLM_CPUS=2 \
VLLM_MEMORY=12g \
VLLM_SHM_SIZE=2g \
VLLM_PORT=8000 \
make vllm-start
```

`latest` 适合第一次探索，不适合实验报告。正式记录前应换成固定版本，保存 `docker image inspect` 的 image ID/RepoDigest，并记录模型 revision。

模型需要 Hugging Face Token 时，只通过环境变量传入，不要写入脚本或 Git：

```bash
export HF_TOKEN='<token>'
make vllm-start
```

## 9. 对照实验

完成基础观察后做三个单变量实验。每次记录容器配置、vLLM 启动日志、`/metrics`、`docker stats` 和请求延迟。

### 实验 A：CPU 配额

分别用 `VLLM_CPUS=1` 和 `VLLM_CPUS=4` 启动相同模型与负载。比较 `cpu.stat`、GPU 利用率空洞、TTFT、TPOT 和吞吐。结论必须说明 CPU 参与了哪段请求路径，不能只写“CPU 越多越快”。

### 实验 B：shared memory

分别设置较小和充足的 `VLLM_SHM_SIZE`，先在单 GPU 下建立基线；有多 GPU 时再启用 Tensor Parallel。比较启动日志、Worker 行为和 `/dev/shm` 使用量。不要通过共享宿主机 IPC namespace 来制造结果。

### 实验 C：模型缓存与冷启动

清晰区分三种状态：

- 冷：本地没有 vLLM 镜像，也没有模型缓存。
- 温：已有镜像，但没有模型缓存。
- 热：已有镜像和模型缓存，但 Engine 仍需初始化并加载到 HBM。

分别记录镜像准备、模型下载、模型加载、Engine 初始化、健康检查成功和首个请求完成时间。不要在仍有进程运行时删除共享缓存。

## 10. 验收

提交以下产物：

- `make inspect` 与 `make vllm-inspect` 的原始输出，并标注至少五个不同的 Linux 边界。
- 一张从 Docker/OCI 到 vLLM Worker 和 GPU 的进程/资源路径图。
- 一张 cgroup Host Memory、`/dev/shm` 和 GPU HBM 的区别表。
- 三个对照实验中的至少一个，包含环境、单一变量、原始数据、结论和反例。
- 一次故障诊断：证明问题属于 CPU Throttling、Host OOM、shared memory、GPU OOM、模型缓存或网络中的哪一层。

验收标准不是“接口返回 200”，而是能从 Linux 证据解释 vLLM 的运行和故障。

## 清理

```bash
make stop       # 基础容器仍在运行时
make vllm-stop  # vLLM 容器仍在运行时
```

模型缓存在 `.cache/huggingface`，已被仓库 `.gitignore` 忽略。确认没有其他实验复用后再手动清理。

## 参考资料

- [Linux namespaces](https://man7.org/linux/man-pages/man7/namespaces.7.html)
- [Control Group v2](https://docs.kernel.org/admin-guide/cgroup-v2.html)
- [OCI Runtime Specification](https://github.com/opencontainers/runtime-spec)
- [Docker storage drivers](https://docs.docker.com/engine/storage/drivers/)
- [Docker resource constraints](https://docs.docker.com/engine/containers/resource_constraints/)
- [vLLM Docker deployment](https://docs.vllm.ai/en/latest/deployment/docker/)
- [vLLM Architecture Overview](https://docs.vllm.ai/en/latest/design/arch_overview/)
- [vLLM Production Metrics](https://docs.vllm.ai/en/latest/usage/metrics/)
