# AI Infra 工程师（容器方向）知识点拆解

> 来源：[`jd.txt`](./jd.txt)（阿里巴巴 2027 届应届生岗位，更新于 2026-08-06）
>
> 目标：把 JD 中的职责和要求拆成可学习、可实践、可验证的知识单元。本文不是简单的关键词清单，而是一张围绕 AI 工作负载构建云原生基础设施的能力地图。

## 1. 岗位本质

这个岗位要解决的核心问题是：**让大规模 AI 训练、推理和 Agent 工作负载在 Kubernetes 上运行得更快、更稳、更省资源。**

可以将岗位能力归纳为四层：

1. **系统基础层**：Linux、网络、文件系统、进程、容器、Go、分布式系统。
2. **资源编排层**：Kubernetes、Controller、Scheduler、GPU 调度、弹性伸缩、资源隔离。
3. **AI 数据路径层**：模型加载、Checkpoint、并行文件 I/O、RDMA/RoCE、分布式通信。
4. **生产工程层**：高可用、可观测性、性能分析、镜像服务、智能运维和工程规范。

判断是否真正掌握一个知识点，可以使用下面四个标准：

- 能解释它解决什么问题以及为什么需要它。
- 能沿调用链或数据路径解释底层机制。
- 能用指标和工具定位瓶颈，而不是凭感觉调参。
- 能实现一个最小系统，并通过测试或基准数据证明效果。

## 2. 计算机基础与编程能力

### 2.1 数据结构与算法

#### 必备知识

- 时间复杂度、空间复杂度、均摊分析和常数开销。
- 数组、链表、栈、队列、哈希表、堆、并查集、Trie、树和图。
- 排序、二分、双指针、滑动窗口、前缀和、贪心、回溯和动态规划。
- BFS、DFS、拓扑排序、最短路、最小生成树和强连通分量。
- 位运算、概率算法、一致性哈希、布隆过滤器和 Count-Min Sketch。
- 并发场景中的无锁队列、环形缓冲区、优先队列和延迟队列。

#### 与岗位的连接

- Scheduler 的待调度队列、优先级队列和拓扑搜索。
- 集群资源索引、节点过滤、缓存以及调度复杂度控制。
- 拓扑感知放置、Gang Scheduling 和装箱问题。
- 海量监控数据的聚合、近似统计与异常检测。

#### 验证方式

- 能分析一个调度算法的时间复杂度和空间复杂度。
- 能解释为何资源调度通常是 NP-hard 问题，以及工程上如何用启发式算法近似求解。
- 能实现带优先级、退避和公平性的工作队列。

### 2.2 Go 语言（核心）

#### 语言与运行时

- 类型系统、接口、组合、方法集、泛型和错误处理。
- Slice、Map、String 的内存布局和扩容行为。
- Goroutine、Channel、select、Context 和 sync 包。
- Go 调度器的 GMP 模型、抢占、syscall 和 work stealing。
- 栈增长、逃逸分析、内存分配器和垃圾回收。
- Race Detector、pprof、trace、benchmark 和 fuzz test。

#### 工程实践

- 模块管理、依赖治理、代码生成和交叉编译。
- 面向接口设计，但避免不必要的抽象。
- 超时、取消、重试、限流、熔断和优雅退出。
- 并发安全、背压、资源泄漏和 Goroutine 泄漏排查。
- 单元测试、集成测试、表驱动测试和基准测试。

#### Kubernetes 开发相关

- client-go：ClientSet、Dynamic Client、Informer、Lister 和 WorkQueue。
- controller-runtime：Manager、Reconciler、Cache、Webhook 和 Leader Election。
- API Machinery：Group/Version/Kind、Scheme、Codec、RESTMapper 和 Discovery。
- Kubernetes API 的乐观并发控制、resourceVersion、watch 和 patch。

#### 验证方式

- 实现一个有超时、重试、限流、指标和优雅退出的并发服务。
- 使用 pprof 找出 CPU、内存或 Goroutine 瓶颈，并提交前后对比数据。
- 实现一个 Kubernetes Controller，而不只是套用脚手架。

### 2.3 其他语言

- **Python**：自动化、数据处理、基准测试，以及 PyTorch 训练/推理工作负载构造。
- **C/C++**：系统调用、内存管理、eBPF、容器 Runtime、网络与存储热路径。
- **Java**：大型工程平台和服务端生态；理解 JVM 资源模型有助于排查混部问题。
- **Shell**：进程、网络、磁盘和容器问题的快速诊断；避免脆弱的生产脚本。

## 3. Linux 系统基础

### 3.1 进程、线程与调度

- 进程地址空间、线程模型、上下文切换和系统调用。
- fork/exec、信号、僵尸进程、PID namespace 和 subreaper。
- Linux 调度类、CFS、实时调度、CPU affinity 和 NUMA affinity。
- Load Average、运行队列、上下文切换和 CPU steal 的含义。
- cgroup v1/v2 的层级模型及 CPU、Memory、IO、PIDs 控制器。
- CPU request/limit 与 CFS bandwidth throttling 的关系。

### 3.2 内存管理

- 虚拟内存、页表、缺页异常、TLB、Huge Page 和 Transparent Huge Page。
- Page Cache、匿名页、mmap、共享内存和 Copy-on-Write。
- NUMA、内存亲和性、跨 NUMA 访问成本。
- cgroup memory.high、memory.max、OOM Killer 和 PSI。
- pinned memory、GPU 显存、Unified Memory 和主机到设备的数据传输。
- 内存碎片、泄漏、抖动以及模型加载时的峰值内存。

### 3.3 文件系统与块 I/O

- VFS、inode、dentry、文件描述符和 mount namespace。
- ext4/XFS、OverlayFS、copy-up、写时复制和 inode 压力。
- Page Cache、Direct I/O、Buffered I/O、mmap 和 fsync 语义。
- Block layer、I/O scheduler、queue depth、IOPS、吞吐和延迟。
- io_uring、异步 I/O、零拷贝和并行读取。
- 本地盘、网络盘、对象存储和分布式文件系统的差异。

### 3.4 Linux 网络栈

- Ethernet、ARP、IP、ICMP、TCP、UDP、DNS 和 HTTP/2 基础。
- Socket、连接队列、拥塞控制、重传、MTU 和分片。
- Netfilter、iptables/nftables、conntrack、路由和策略路由。
- Network Namespace、veth、bridge、VXLAN、IPVS 和 eBPF datapath。
- RSS/RPS/RFS、IRQ affinity、GRO/GSO/TSO 和零拷贝。
- 延迟、吞吐、PPS、丢包、抖动和 tail latency 的测量。

### 3.5 系统诊断工具

- 进程与资源：`top`、`pidstat`、`vmstat`、`mpstat`、`sar`、`free`。
- 文件与 I/O：`iostat`、`iotop`、`lsof`、`strace`、`blktrace`。
- 网络：`ss`、`ip`、`ethtool`、`tcpdump`、`nstat`、`tc`。
- 性能：`perf`、FlameGraph、`bpftrace`、BCC 和 libbpf。
- 硬件：`numactl`、`lscpu`、`nvidia-smi`、DCGM 和 PCIe 拓扑工具。

## 4. 容器原理与运行时

### 4.1 容器隔离机制

- Namespace：PID、Mount、Network、IPC、UTS、User、Cgroup 和 Time。
- cgroup v2 的资源记账、限制、委托和压力指标。
- Capability、seccomp、AppArmor/SELinux 和 Rootless Container。
- 容器并不是虚拟机：共享内核带来的性能收益与隔离边界。

### 4.2 OCI 与镜像

- OCI Image Spec、Runtime Spec 和 Distribution Spec。
- 镜像 Manifest、Config、Layer、Content Addressable Storage 和 digest。
- OverlayFS 联合挂载、镜像解压、快照和容器可写层。
- 多阶段构建、镜像瘦身、可复现构建、SBOM 和镜像签名。
- Registry API、鉴权、分层上传、垃圾回收和镜像复制。

### 4.3 Runtime 调用链

- Kubernetes → kubelet → CRI → containerd/CRI-O → shim → runc 的完整链路。
- containerd 的 plugin、content store、snapshotter、metadata store 和 shim。
- runc 创建 namespace、cgroup、rootfs 并启动进程的过程。
- CRI 的 RuntimeService、ImageService 和 Pod Sandbox 模型。
- CNI、CSI、Device Plugin 分别在哪个阶段介入。

### 4.4 轻量化与高密运行时

- 容器启动时间分解：调度、拉镜像、解压、挂载、创建 sandbox、启动进程和应用初始化。
- Lazy Pulling、远程 Snapshotter、按需加载和镜像预热。
- sandbox 复用、进程级隔离、MicroVM 和用户态内核的取舍。
- Kata Containers、gVisor、Firecracker 与 runc 的隔离/性能边界。
- 高密部署下的 PID、文件描述符、内存和控制面开销。

#### 验证方式

- 手工使用 namespace、cgroup 和 pivot_root 构建一个最小容器。
- 跟踪一次 Pod 创建，从 API 请求定位到容器主进程。
- 对镜像拉取和容器启动链路分段打点，找出冷启动主导因素。

## 5. Kubernetes 架构与核心机制

### 5.1 控制面

- API Server 的认证、鉴权、准入、版本转换和持久化链路。
- etcd 的 Raft、一致性、revision、watch、压缩、碎片整理和备份恢复。
- Scheduler 的队列、调度周期、绑定周期和 Framework 插件。
- Controller Manager 的控制循环、期望状态、Leader Election 和限速队列。
- 声明式 API、Level-triggered Reconciliation 和最终一致性。

### 5.2 节点面

- kubelet 的 Pod Lifecycle、PLEG、Status Manager、Volume Manager 和 Eviction Manager。
- Static Pod、Pod Sandbox、Probe、QoS 和驱逐。
- kube-proxy 的 iptables/IPVS/nftables 模式及替代方案。
- Node Allocatable、系统预留、拓扑管理和设备管理。

### 5.3 API 与扩展机制

- CRD 设计、版本演进、默认值、校验和 status/spec 分离。
- Informer 的 List-Watch、本地缓存、事件处理和 resync。
- Admission Webhook、聚合 API 和 API Priority and Fairness。
- Operator 的幂等性、终结器、OwnerReference 和垃圾回收。

### 5.4 可靠性与故障排查

- 控制面高可用、etcd quorum 和故障域。
- API Server、etcd、Scheduler、Controller 和 kubelet 的 SLI/SLO。
- Pending、CrashLoopBackOff、ImagePullBackOff、NotReady 和 Unknown 状态排查。
- 大集群中的 API QPS、watch fan-out、对象数量和控制器风暴。

#### 验证方式

- 从零实现一个带 CRD、Controller、Webhook、状态机和故障恢复测试的 Operator。
- 能说明删除一个带 Finalizer 的自定义资源时，每一步发生了什么。
- 能通过审计日志、事件、组件日志和指标定位 Pod 长时间 Pending 的根因。

## 6. AI 容器调度

### 6.1 Scheduler Framework

- QueueSort、PreFilter、Filter、PostFilter、PreScore、Score、Reserve、Permit、PreBind、Bind 和 PostBind。
- Scheduling Cycle 与 Binding Cycle 的并发边界。
- 调度缓存、Assume、并行打分、失败重试和退避。
- Extender、自定义 Scheduler、Scheduler Plugin 和调度 CRD 的取舍。

### 6.2 GPU 资源模型

- GPU 型号、显存容量、计算能力、功耗和健康状态。
- Kubernetes Extended Resource 与 Device Plugin API。
- NVIDIA Device Plugin、GPU Operator 和 DCGM Exporter。
- GPU 独占、time-slicing、MPS、MIG 和 vGPU 的适用边界。
- PCIe、NVLink、NVSwitch、NUMA、NIC 和 GPU 的物理拓扑。
- GPU Direct RDMA、GPUDirect Storage 及其对放置策略的影响。

### 6.3 拓扑感知调度

- kubelet Topology Manager、CPU Manager、Memory Manager 和 Device Manager。
- NUMA 对齐、PCIe Root Complex、NIC-GPU 亲和和跨节点网络拓扑。
- NodeResourceTopology、Topology Manager Policy 和拓扑信息采集。
- 训练任务的同机、同交换机、同机架与跨故障域放置策略。
- 拓扑打分中性能、可用性与碎片率的权衡。

### 6.4 批任务与协同调度

- Gang Scheduling、All-or-Nothing、Co-scheduling 和调度屏障。
- Queue、Quota、Priority、Preemption、Fair Share 和 Dominant Resource Fairness。
- Volcano、Kueue、Kubeflow Training Operator 的职责边界。
- Job、JobSet、MPIJob、PyTorchJob 和 RayJob 的生命周期。
- Backfill、Reservation、Capacity Scheduling 和多租户队列。

### 6.5 资源碎片治理

- CPU、内存、GPU、显存和拓扑多维装箱问题。
- Bin Packing 与 Spread 策略对利用率、可靠性和后续调度的影响。
- GPU 型号/显存异构导致的不可替代资源碎片。
- Resource Reservation、Descheduler、任务迁移和整理策略。
- 调度质量指标：成功率、排队时间、碎片率、GPU 利用率和作业完成时间。

#### 验证方式

- 编写一个 Scheduler Framework 插件，按 NUMA、GPU-NIC 距离和剩余碎片打分。
- 用可重复的集群负载比较默认调度器与自定义策略的排队时间和利用率。
- 能解释抢占如何产生级联影响，以及如何避免低优先级任务长期饥饿。

## 7. AI 训练与推理工作负载

### 7.1 GPU 与模型计算基础

- CUDA 执行模型、Kernel、Stream、Event 和异步执行。
- HBM、L2 Cache、Shared Memory、PCIe 和 NVLink 的带宽层级。
- 计算受限与访存受限；Arithmetic Intensity 和 Roofline Model。
- FP32、TF32、FP16、BF16、INT8、FP8 和量化的性能/精度取舍。
- 显存构成：参数、梯度、优化器状态、激活值、KV Cache 和临时 Buffer。

### 7.2 分布式训练

- 数据并行、张量并行、流水线并行、专家并行和序列并行。
- AllReduce、AllGather、ReduceScatter、AllToAll 和 Broadcast。
- Ring/Tree 集合通信算法，以及延迟与带宽成本模型。
- NCCL、Gloo、MPI、Rendezvous、Rank、World Size 和 Process Group。
- DDP、FSDP、ZeRO 1/2/3、梯度累积和 Activation Checkpointing。
- 通信计算重叠、Straggler、训练容错和弹性训练。

### 7.3 推理服务

- Prefill 与 Decode 阶段的计算和内存特征。
- 首 Token 延迟 TTFT、单 Token 延迟 TPOT、吞吐、并发和 P99。
- Continuous Batching、Dynamic Batching、PagedAttention 和 KV Cache 管理。
- Tensor Parallel、Pipeline Parallel、Speculative Decoding 和模型量化。
- vLLM、Triton Inference Server、TensorRT-LLM 等系统的职责与差异。
- 在线推理、离线推理和训练任务混部的隔离策略。

### 7.4 Agent 工作负载

- Agent 的无状态计算与有状态会话、工具执行和沙箱隔离。
- 短时突发、长尾执行、依赖安装和按需环境创建等负载特点。
- 冷启动组成：镜像、模型、依赖、代码、沙箱、网络和缓存。
- Agent 池化、Warm Pool、预创建 Pod、预拉镜像和模型预热。
- 会话状态、检查点、任务幂等、超时、取消和重试。

## 8. AI 高性能存储

### 8.1 Kubernetes 存储

- PV、PVC、StorageClass、AccessMode、VolumeMode 和拓扑约束。
- CSI Controller/Node Plugin、Sidecar、挂载与卸载流程。
- Dynamic Provisioning、Volume Binding、Attach、Mount 和扩容。
- Local PV、HostPath、网络块存储、文件存储和对象存储。
- Snapshot、Clone、数据保护和故障恢复。

### 8.2 AI 数据访问模式

- 训练样本的大量小文件、顺序读、随机读和元数据压力。
- 多 Worker 同时读取导致的带宽放大和热点问题。
- 模型权重加载、分片格式、反序列化和页缓存行为。
- Checkpoint 的大文件写入、周期性突发、全局同步和恢复时间。
- 推理模型、Tokenizer、Adapter 与 KV Cache 的存储生命周期。

### 8.3 Checkpoint 优化

- 同步与异步 Checkpoint、增量 Checkpoint 和分片 Checkpoint。
- 计算节点到本地盘、共享存储和对象存储的分级落盘。
- 写入流水线、并行化、压缩、去重和后台上传。
- 一致性、原子提交、校验、保留策略和失败恢复。
- RPO、RTO、保存耗时、训练阻塞时间和恢复耗时指标。

### 8.4 缓存与存储卸载

- 内存、本地 NVMe、分布式缓存、共享存储、对象存储的多级缓存。
- Cache Aside、Read Through、预取、淘汰、一致性和缓存击穿。
- FUSE、用户态文件系统、Sidecar 和独立数据服务的开销。
- 数据加载与 GPU 计算重叠、Pinned Memory 和异步预取。
- GPUDirect Storage、RDMA Storage 及 CPU bypass 路径。

#### 验证方式

- 构造模型加载和 Checkpoint 基准，记录吞吐、P99 延迟、CPU 利用率和训练停顿。
- 实现一个分片并行、异步上传且可原子恢复的最小 Checkpoint 流程。
- 使用 eBPF/perf/存储指标判断瓶颈在应用、VFS、块设备还是网络。

## 9. AI 高性能网络

### 9.1 Kubernetes 网络模型

- Pod IP、Service、EndpointSlice、Ingress/Gateway API 和 NetworkPolicy。
- CNI 调用链、IPAM、Overlay/Underlay 和路由模式。
- kube-proxy、IPVS、eBPF Service 转发以及东西向流量路径。
- SR-IOV、Multus、多网卡与设备直通。

### 9.2 RDMA 与 RoCE

- RDMA Verbs：Queue Pair、Completion Queue、Memory Region 和 Work Request。
- Send/Receive、Read/Write、零拷贝和 Kernel Bypass。
- InfiniBand 与 RoCEv2 的差异。
- RoCE 无损网络：PFC、ECN、DCQCN、拥塞与 Head-of-Line Blocking。
- RDMA CM、GID、MTU、路由、网卡固件和驱动。
- RDMA Device Plugin、SR-IOV CNI 和容器内设备暴露。

### 9.3 分布式训练网络

- NCCL 拓扑发现、通信通道、Ring/Tree 和网卡选择。
- GPU-NIC 亲和、NUMA 绑定和 GPUDirect RDMA。
- Rail-optimized 网络、多平面网络和链路聚合。
- Incast、拥塞、丢包、慢节点和集合通信尾延迟。
- 网络 SLI：带宽利用率、RTT、重传、ECN 标记、PFC pause、丢包和 NCCL bus bandwidth。

### 9.4 内核态路径优化

- socket 路径中的拷贝、软中断、协议栈和调度开销。
- XDP、TC eBPF、AF_XDP、DPDK 和 io_uring 的适用场景。
- IRQ/队列/CPU/NUMA 亲和，RSS 和 busy polling。
- 零拷贝、批处理、连接复用和拥塞控制调优。

#### 验证方式

- 画出一个 Pod 中 NCCL 流量从 GPU 到远端 GPU 的完整数据路径。
- 使用 `iperf3`、`qperf`、`ib_write_bw`、NCCL Tests 建立分层基准。
- 能区分应用延迟、容器网络开销、TCP/RDMA 问题和物理链路问题。

## 10. 弹性调度与冷启动优化

### 10.1 弹性伸缩基础

- HPA、VPA、Cluster Autoscaler 和 KEDA 的控制对象与时序。
- 指标采集延迟、稳定窗口、容忍度、冷却时间和抖动。
- Reactive、Predictive 和 Scheduled Scaling。
- 应用扩容、Pod 扩容、节点扩容和 GPU 资源供给的协同。
- Scale from Zero 的容量发现、调度与冷启动问题。

### 10.2 多级缓存与预加载

- Registry 缓存、节点镜像缓存、远程 Snapshotter 和 P2P 分发。
- 模型缓存、权重分片、内存映射和本地 NVMe 缓存。
- Runtime/Sandbox Warm Pool、容器池和进程池。
- 流量预测、热点识别、预热收益与闲置资源成本。
- 缓存命中率、淘汰准确率和预热放大风险。

### 10.3 端到端启动链路

- 请求到达 → 容量决策 → Pod 创建 → 调度 → 节点扩容 → 镜像/模型准备 → 容器启动 → Readiness → 接流量。
- 为每一阶段建立 trace/span 和延迟预算。
- 区分控制面瓶颈、资源供给瓶颈、数据准备瓶颈和应用初始化瓶颈。
- 用 P50/P95/P99 和冷/温/热启动分类评价优化效果。

#### 验证方式

- 建立冷启动火焰图或瀑布图，并量化每项优化贡献。
- 实现一个基于队列长度和预测负载的弹性控制器。
- 验证突发流量、预测错误、节点故障和缓存未命中时的系统行为。

## 11. 资源管控与多租户

### 11.1 Kubernetes 资源语义

- Request、Limit、QoS Class、LimitRange 和 ResourceQuota。
- CPU 可压缩资源与内存不可压缩资源的差异。
- Pod Overhead、Init Container、Sidecar 和临时存储记账。
- PriorityClass、Preemption、PodDisruptionBudget 和 Eviction。

### 11.2 多租户隔离

- Namespace、RBAC、ServiceAccount 和 API 隔离。
- NetworkPolicy、存储权限、Secret 管理和镜像供应链。
- cgroup、CPU pinning、NUMA、Cache、内存带宽和 I/O 干扰。
- GPU 显存/算力隔离、MIG/MPS/time-slicing 的隔离强度。
- Noisy Neighbor、资源超卖、混部和干扰检测。

### 11.3 作业管理规范

- 作业模板、资源声明、队列、优先级、重试和超时。
- 失败语义、幂等、断点续训、退出码与重启策略。
- 配额借用、回收、公平性、成本归属和审计。
- 训练数据、模型、日志和 Checkpoint 生命周期。

## 12. 分布式系统与平台架构

### 12.1 分布式系统基础

- CAP、PACELC、线性一致性、顺序一致性和最终一致性。
- 共识、Leader Election、Lease、Quorum 和脑裂。
- 超时、重试、指数退避、抖动、幂等和去重。
- Fail-stop、网络分区、时钟偏移、部分失败和级联故障。
- 负载均衡、限流、熔断、背压和过载保护。

### 12.2 高可用与可扩展性

- 无单点、故障域、冗余、优雅降级和容量保护。
- 水平扩展、数据分片、热点治理和状态迁移。
- SLI、SLO、错误预算、容量规划和故障演练。
- RTO、RPO、备份恢复和跨区域容灾。
- 控制面与数据面分离，快路径与慢路径分离。

### 12.3 平台工程

- 声明式 API、平台抽象、Golden Path 和自助服务。
- API/CRD 的兼容性、版本治理和废弃策略。
- 多集群管理、联邦调度、集群资源画像和故障转移。
- 发布、灰度、回滚、Feature Gate 和配置治理。
- 成本、利用率、可靠性、性能和研发效率之间的权衡。

## 13. 可观测性与性能工程

### 13.1 可观测性

- Metrics、Logs、Traces、Events 和 Continuous Profiling。
- RED（Rate、Errors、Duration）与 USE（Utilization、Saturation、Errors）。
- Prometheus 数据模型、PromQL、抓取、聚合和高基数问题。
- OpenTelemetry 的 Context Propagation、Trace、Metric 和 Collector。
- GPU 指标、容器指标、节点指标、Kubernetes 状态和业务指标关联。
- 告警的可行动性、去重、抑制、分级和 SLO 告警。

### 13.2 性能分析方法

- 先建立基线和假设，再进行单变量实验。
- 平均值、百分位、直方图、吞吐、并发和排队时间。
- Coordinated Omission、预热、噪声、重复实验和置信区间。
- CPU profiling、off-CPU profiling、内存分析、锁竞争和 I/O 分析。
- eBPF 的 kprobe、uprobe、tracepoint、fentry/fexit 和 CO-RE。
- 从硬件计数器、内核、Runtime、容器、Kubernetes 到应用的关联分析。

### 13.3 容量与效率指标

- GPU SM 利用率、显存利用率、Tensor Core 利用率和功耗。
- Model FLOPS Utilization、训练吞吐、Step Time 和 Job Completion Time。
- 推理 TTFT、TPOT、Tokens/s、并发、P99 和单位 Token 成本。
- CPU/内存/网络/存储饱和度和集群资源碎片率。
- 调度等待、镜像拉取、模型加载和扩缩容耗时。

#### 验证方式

- 对一个跨层性能问题给出证据链，而不是只展示单一监控面板。
- 写一份包含基线、环境、负载、假设、实验、数据和结论的性能报告。
- 能说明优化吞吐为何可能恶化尾延迟，以及如何确定合理目标。

## 14. 镜像基础服务

### 14.1 Registry 架构

- OCI Distribution API、Blob、Manifest、Tag、Digest 和 Upload Session。
- 元数据与 Blob 数据分离、对象存储后端和缓存层。
- 多副本、负载均衡、强弱一致性、垃圾回收和容量规划。
- 身份认证、授权、租户隔离、审计和配额。

### 14.2 大规模镜像分发

- Registry 限流、节点并发拉取和出口带宽瓶颈。
- P2P 分发、分层缓存、预热和热点镜像识别。
- Lazy Pulling、eStargz、Nydus、OverlayBD 等按需加载思路。
- 大模型镜像与模型权重解耦，避免镜像层频繁失效。
- 拉取成功率、首字节延迟、解压速度、缓存命中率和带宽节省。

### 14.3 供应链安全

- 最小基础镜像、漏洞扫描、SBOM、签名与验证。
- Provenance、可复现构建和准入策略。
- Secret 泄漏、恶意镜像、Tag 漂移和依赖投毒防护。

## 15. 智能运维 Agent

### 15.1 运维自动化闭环

- Signal → Detection → Correlation → Diagnosis → Decision → Action → Verification。
- 告警聚合、事件关联、变更关联和拓扑关联。
- Runbook 自动化、动作幂等、超时、回滚和审批。
- 人在回路、置信度阈值、权限边界和审计记录。

### 15.2 根因定位

- 症状、直接原因、根本原因与促成因素的区分。
- 依赖拓扑、时序相关、因果推断和变更影响分析。
- 指标、日志、Trace、事件和配置的多源证据融合。
- 故障知识库、历史案例检索和诊断假设排序。

### 15.3 LLM/Agent 工程

- Tool Calling、状态管理、任务规划、重试和终止条件。
- RAG、知识时效、证据引用和幻觉控制。
- Prompt Injection、工具权限、敏感信息和命令执行安全。
- 离线评测、故障回放、成功率、误操作率、MTTD 和 MTTR。
- 自动修复必须具备限权、灰度、回滚、验证和熔断机制。

## 16. 开源与工程协作

### 16.1 推荐关注的项目

- Kubernetes：核心架构、Scheduler、kubelet、API Machinery、client-go。
- containerd、runc、CRI-O：容器 Runtime 与生命周期。
- Volcano、Kueue：批调度、队列和配额。
- Kubeflow Training Operator、JobSet、Ray：AI 作业编排。
- NVIDIA GPU Operator、Device Plugin、DCGM Exporter：GPU 管理。
- Cilium、Multus、SR-IOV CNI：容器网络与高性能网络。
- CSI、JuiceFS、Alluxio 等：存储接口与缓存思路。
- OpenKruise：工作负载、原地升级和高级发布能力。

### 16.2 开源贡献能力

- 阅读贡献指南、开发环境、测试体系和代码所有权。
- 从 issue 复现、最小化问题、定位根因到提交测试与修复。
- 设计提案、API Review、向后兼容和版本演进。
- 清晰描述问题、测试证据、性能数据和方案权衡。
- 持续的小型高质量贡献比一次大型但不可维护的改动更有说服力。

## 17. 知识优先级

### P0：岗位基础，必须能深入回答

- Go 并发、Runtime、性能分析和工程测试。
- Linux 进程、内存、文件系统、网络栈、cgroup 和 namespace。
- 容器镜像、OCI、containerd、runc、CRI、CNI 和 CSI 调用链。
- Kubernetes 控制面、kubelet、Controller、Informer 和 Scheduler Framework。
- GPU 资源模型、Device Plugin、拓扑感知和 Gang Scheduling。
- AI 训练/推理的显存、通信、存储和伸缩特征。
- Prometheus、OpenTelemetry、perf、eBPF 和分层故障定位。

### P1：形成岗位竞争力

- Volcano/Kueue/Kubeflow 中至少一个项目的源码与实践。
- RDMA/RoCE、NCCL、SR-IOV 和 GPUDirect RDMA。
- Checkpoint、模型加载、多级缓存和高性能数据读取。
- 镜像加速、Lazy Pulling、P2P 分发和冷启动优化。
- GPU 共享、资源碎片、队列公平性与多租户隔离。
- 高可用、SLO、容量规划和大规模集群控制面优化。

### P2：差异化方向

- containerd snapshotter、轻量沙箱或 Runtime 热路径开发。
- GPUDirect Storage、存储卸载和异步 Checkpoint。
- 预测式弹性、跨集群调度和资源供给协同。
- 基于 eBPF 的跨层性能观测平台。
- 有安全边界、评测和自愈闭环的智能运维 Agent。

## 18. 推荐项目闭环

### 项目一：GPU 拓扑感知批调度器

- 定义带 GPU、显存、NUMA、NIC 拓扑需求的作业 API。
- 实现 Gang Scheduling、队列配额和拓扑打分插件。
- 构造异构 GPU 与多队列负载。
- 对比默认策略的排队时间、作业完成时间、利用率和碎片率。
- 加入故障、抢占、饥饿和并发调度测试。

### 项目二：AI 工作负载冷启动加速

- 为 Pod 启动链路增加端到端 Trace。
- 实现镜像预拉取、远程 Snapshotter 或 P2P 分发中的一种。
- 实现节点模型缓存和 Warm Pool。
- 比较冷、温、热三类启动的 P50/P95/P99。
- 量化加速收益、缓存成本和最差情况下的退化。

### 项目三：异步分片 Checkpoint 服务

- 将训练进程写入与持久化上传解耦。
- 支持分片、并行上传、校验、原子提交和恢复。
- 设计本地 NVMe 与对象存储两级路径。
- 注入节点宕机、上传失败、重复请求和部分分片损坏。
- 评价训练阻塞时间、吞吐、RPO 和 RTO。

### 项目四：AI 集群跨层性能诊断

- 汇聚 Kubernetes、节点、容器、GPU、NCCL 和应用指标。
- 用 eBPF 采集调度、网络和 I/O 热路径证据。
- 针对 GPU 利用率下降建立自动化诊断规则。
- 覆盖 CPU throttling、I/O 阻塞、网络拥塞和慢 Worker 场景。
- 输出可复现的故障报告与根因证据链。

## 19. 面试准备检查表

- 能在白板上画出 Kubernetes 创建 Pod 到容器进程启动的完整链路。
- 能解释一次分布式训练的数据路径和集合通信路径。
- 能比较 CPU、内存、GPU、网络、存储瓶颈的典型症状和工具。
- 能设计 GPU 拓扑感知调度插件，并说明复杂度、缓存和并发问题。
- 能设计秒级弹性方案，并量化预热成本与容量风险。
- 能设计大规模 Checkpoint 方案，并解释一致性和失败恢复。
- 能说明 RDMA/RoCE 在容器中的设备、网络和权限配置链路。
- 能针对一个真实性能问题给出基线、假设、实验和数据。
- 能讲清一个复杂项目中的取舍、失败尝试和可验证结果。
- 能展示至少一个与 Kubernetes、容器或 AI Infra 相关的可运行项目或开源贡献。

## 20. 避免“只会关键词”

以下回答通常不足以证明掌握：

- “Kubernetes 负责容器编排”，但说不清控制循环和 Pod 创建链路。
- “RDMA 是零拷贝”，但说不清内存注册、Queue Pair 和容器设备暴露。
- “用缓存加速模型加载”，但没有一致性、淘汰、容量和命中率设计。
- “用 eBPF 排查性能”，但不知道挂载点、事件模型和观测开销。
- “提高 GPU 利用率”，但没有区分计算、显存、通信和数据供给瓶颈。
- “实现自动扩缩容”，但没有端到端时延、抖动和资源供给分析。

最终应形成的能力不是背诵组件名，而是完成这条闭环：

**理解 AI 负载 → 建立资源与性能模型 → 设计云原生机制 → 实现并测试 → 用数据验证 → 处理生产故障。**

## 21. 学习资料索引

> 选择原则：优先官方文档、项目源码文档、原始论文和大学公开课。带版本号的 Kubernetes、NVIDIA、PyTorch、KEDA 等文档可能随版本变化，阅读时应同时确认页面上的版本与特性状态。

### 21.1 数据结构、算法与 Go

- [MIT 6.006 Introduction to Algorithms](https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-spring-2020/)：包含课程视频、讲义、习题和考试，适合系统补齐数据结构与算法。
- [Go 官方文档入口](https://go.dev/doc/)：语言规范、Effective Go、诊断和教程的总入口。
- [A Tour of Go](https://go.dev/tour/)：适合快速建立语法、接口、泛型和并发基础。
- [Effective Go](https://go.dev/doc/effective_go)：Go 的命名、接口、错误处理和并发惯用法。
- [The Go Memory Model](https://go.dev/ref/mem)：理解 happens-before、Channel、锁和原子操作的必读资料。
- [Go Diagnostics](https://go.dev/doc/diagnostics)：pprof、trace、运行时统计和性能诊断入口。
- [Data Race Detector](https://go.dev/doc/articles/race_detector)：并发竞态检测方法、限制和典型案例。
- [Go Blog: Profiling Go Programs](https://go.dev/blog/pprof)：通过 pprof 分析 CPU 和内存瓶颈的完整示例。

### 21.2 Linux 系统基础

- [Linux Kernel Administrator's Guide](https://www.kernel.org/doc/html/latest/admin-guide/index.html)：内核管理文档总入口，覆盖内存、块设备、文件系统、CPU 和 namespace。
- [Control Group v2](https://www.kernel.org/doc/html/latest/admin-guide/cgroup-v2.html)：cgroup v2 的权威设计和各控制器接口。
- [Linux namespaces 手册](https://man7.org/linux/man-pages/man7/namespaces.7.html)：各类 namespace 的关系、生命周期和 API。
- [Linux memory management documentation](https://www.kernel.org/doc/html/latest/mm/index.html)：页表、分配器、回收、NUMA 和内存管理子系统。
- [Linux networking documentation](https://www.kernel.org/doc/html/latest/networking/index.html)：内核网络栈、驱动、XDP 和协议相关文档。
- [Linux filesystems documentation](https://www.kernel.org/doc/html/latest/filesystems/index.html)：VFS、OverlayFS、ext4 等文件系统机制。
- [perf Wiki](https://perf.wiki.kernel.org/index.php/Main_Page)：perf 事件、采样和性能分析资料。
- [eBPF Userspace API](https://docs.kernel.org/userspace-api/ebpf/index.html)：eBPF 程序、Map 和用户态 API 的内核文档入口。
- [BPF Documentation](https://docs.kernel.org/bpf/index.html)：Verifier、程序类型、BTF、CO-RE 等内核侧资料。
- [Brendan Gregg Linux Performance](https://www.brendangregg.com/linuxperf.html)：Linux 性能工具、USE 方法和分析图谱。

### 21.3 容器、OCI 与 Runtime

- [OCI Runtime Specification](https://github.com/opencontainers/runtime-spec)：容器配置、状态和生命周期规范。
- [OCI Image Specification](https://github.com/opencontainers/image-spec)：Manifest、Config、Layer、Descriptor 和 Image Index。
- [OCI Distribution Specification](https://github.com/opencontainers/distribution-spec)：镜像和 OCI Artifact 的 Push、Pull 与内容发现协议。
- [runc](https://github.com/opencontainers/runc)：OCI Runtime Spec 的参考实现，适合结合源码理解容器创建。
- [containerd Architecture](https://github.com/containerd/containerd/blob/main/docs/historical/design/architecture.md)：containerd 的服务、元数据、内容和执行子系统。
- [containerd Runtime v2](https://github.com/containerd/containerd/blob/main/docs/runtime-v2.md)：containerd、shim 和 runc 的交互模型。
- [containerd CRI Architecture](https://github.com/containerd/containerd/blob/main/docs/cri/architecture.md)：kubelet、CRI、CNI、Sandbox 和业务容器的创建链路。
- [Kubernetes Container Runtime Interface](https://kubernetes.io/docs/concepts/containers/cri/)：Kubernetes 对容器运行时的接口定义。
- [Kata Containers Architecture](https://github.com/kata-containers/kata-containers/blob/main/docs/design/architecture/README.md)：轻量虚拟机容器的隔离与架构设计。
- [gVisor Architecture Guide](https://gvisor.dev/docs/architecture_guide/intro/)：用户态内核、系统调用拦截和沙箱边界。
- [Firecracker Design](https://github.com/firecracker-microvm/firecracker/blob/main/docs/design.md)：MicroVM 的设备模型、启动和安全设计。

### 21.4 Kubernetes 架构与 Controller

- [Kubernetes Concepts](https://kubernetes.io/docs/concepts/)：官方概念文档总入口。
- [Cluster Architecture](https://kubernetes.io/docs/concepts/architecture/)：API Server、etcd、Scheduler、Controller Manager、kubelet 和 kube-proxy。
- [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)：控制循环与期望状态模型。
- [Kubernetes API Concepts](https://kubernetes.io/docs/reference/using-api/api-concepts/)：resourceVersion、watch、并发控制、分页和 Patch。
- [Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)：CRD、自定义 Controller 与 API Aggregation 的选择。
- [Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)：用 CRD 与控制器封装领域运维知识。
- [client-go](https://github.com/kubernetes/client-go)：ClientSet、Informer、Lister 和 WorkQueue 的官方 Go 客户端源码。
- [sample-controller](https://github.com/kubernetes/sample-controller)：使用 client-go 编写 Controller 的最小官方示例。
- [controller-runtime FAQ](https://github.com/kubernetes-sigs/controller-runtime/blob/main/FAQ.md)：Reconcile、Cache、Client 和事件处理的常见设计问题。
- [Programming Kubernetes CRDs](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)：Kubernetes API 约定，设计 CRD 时应重点阅读。
- [etcd API guarantees](https://etcd.io/docs/v3.7/learning/api_guarantees/)：线性一致性、事务与 Watch 保证。
- [Operating etcd clusters for Kubernetes](https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/)：etcd 备份、恢复、扩容和运维入口。

### 21.5 Kubernetes 调度与批任务

- [Kubernetes Scheduler](https://kubernetes.io/docs/concepts/scheduling-eviction/kube-scheduler/)：默认调度器的基本职责和选择流程。
- [Scheduling Framework](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/)：Scheduling Cycle、Binding Cycle 和插件扩展点。
- [Scheduler Configuration](https://kubernetes.io/docs/reference/scheduling/config/)：调度 Profile、插件启停和参数配置。
- [Resource Bin Packing](https://kubernetes.io/docs/concepts/scheduling-eviction/resource-bin-packing/)：按资源利用率进行装箱打分。
- [Pod Topology Spread Constraints](https://kubernetes.io/docs/concepts/scheduling-eviction/topology-spread-constraints/)：可用区、节点和自定义拓扑域上的分布策略。
- [Topology-Aware Workload Scheduling](https://kubernetes.io/docs/concepts/scheduling-eviction/topology-aware-scheduling/)：PodGroup 与拓扑域协同放置；阅读时注意功能成熟度。
- [Volcano Scheduler Overview](https://volcano.sh/docs/scheduler/overview/)：Gang、DRF、Binpack、Priority 等批调度插件。
- [Volcano Queue](https://volcano.sh/docs/concepts/queue/)：队列、保证资源、权重和回收语义。
- [Kueue Concepts](https://kueue.sigs.k8s.io/docs/concepts/)：ClusterQueue、LocalQueue、Cohort、Flavor 和 Admission Check。
- [Kueue Fair Sharing](https://kueue.sigs.k8s.io/docs/concepts/fair_sharing/)：队列公平共享和资源借用。
- [JobSet](https://jobset.sigs.k8s.io/docs/overview/)：分布式 AI/HPC 作业的多 Job 编排 API。
- [Kubeflow Trainer](https://www.kubeflow.org/docs/components/trainer/)：PyTorch、MPI 等分布式训练作业在 Kubernetes 上的管理。

### 21.6 GPU、设备管理与拓扑

- [CUDA Programming Guide](https://docs.nvidia.com/cuda/cuda-programming-guide/contents.html)：Kernel、执行层级、内存、Stream 和多 GPU 编程。
- [CUDA Programming Model](https://docs.nvidia.com/cuda/cuda-programming-guide/01-introduction/programming-model.html)：CPU/GPU 异构系统、SM、内存和互连的入门章节。
- [NVIDIA Multi-Instance GPU User Guide](https://docs.nvidia.com/datacenter/tesla/mig-user-guide/)：MIG 的隔离、实例划分和管理模型。
- [Kubernetes Device Plugins](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/device-plugins/)：节点设备发现、注册和分配接口。
- [NVIDIA Kubernetes Device Plugin](https://github.com/NVIDIA/k8s-device-plugin)：GPU、MIG、Time-Slicing 和 MPS 的 Kubernetes 集成。
- [NVIDIA GPU Operator](https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/latest/)：驱动、Device Plugin、Toolkit、MIG Manager 和 DCGM 的自动化管理。
- [GPU Operator with MIG](https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/gpu-operator-mig.html)：在 Kubernetes 中配置和切换 MIG 策略。
- [Kubernetes Topology Manager](https://kubernetes.io/docs/tasks/administer-cluster/topology-manager/)：CPU、内存和设备的 NUMA 对齐策略。
- [Kubernetes CPU Manager](https://kubernetes.io/docs/tasks/administer-cluster/cpu-management-policies/)：静态 CPU 独占和 CPU affinity。
- [NVIDIA DCGM Documentation](https://docs.nvidia.com/datacenter/dcgm/latest/)：GPU 健康、诊断、性能指标和 Profiling。

### 21.7 分布式训练

- [PyTorch Distributed Overview](https://pytorch.org/tutorials/beginner/dist_overview.html)：DDP、FSDP、Tensor Parallel 和 DeviceMesh 的选择入口。
- [PyTorch Distributed Tutorials](https://docs.pytorch.org/tutorials/distributed.html)：官方分布式训练教程集合。
- [DistributedDataParallel](https://docs.pytorch.org/docs/stable/generated/torch.nn.parallel.DistributedDataParallel.html)：DDP 的同步、进程模型与配置语义。
- [FullyShardedDataParallel](https://docs.pytorch.org/docs/stable/fsdp.html)：参数、梯度和优化器状态分片机制。
- [PyTorch Tensor Parallel](https://docs.pytorch.org/docs/stable/distributed.tensor.parallel.html)：基于 DeviceMesh 的张量并行 API。
- [PyTorch Elastic](https://docs.pytorch.org/docs/stable/distributed.elastic.html)：Rendezvous、成员变化、容错和弹性启动。
- [NCCL User Guide](https://docs.nvidia.com/deeplearning/nccl/user-guide/index.html)：Communicator、集合通信、拓扑发现和故障排查。
- [NCCL Collective Operations](https://docs.nvidia.com/deeplearning/nccl/user-guide/docs/usage/collectives.html)：AllReduce、AllGather、ReduceScatter、AllToAll 等操作的准确语义。
- [NCCL Tests](https://github.com/NVIDIA/nccl-tests)：测试集合通信正确性、算法带宽和总线带宽。
- [Megatron-LM](https://github.com/NVIDIA/Megatron-LM)：张量、流水线、数据、专家和序列并行的工程实现。
- [DeepSpeed ZeRO](https://www.deepspeed.ai/tutorials/zero/)：ZeRO 各阶段的状态分片和内存优化。

### 21.8 大模型推理与 Agent 工作负载

- [vLLM Documentation](https://docs.vllm.ai/en/stable/)：PagedAttention、高吞吐 Serving 和分布式推理入口。
- [vLLM Architecture Overview](https://docs.vllm.ai/en/latest/design/arch_overview/)：Engine、Worker、Model Runner 和在线服务架构。
- [vLLM Optimization and Tuning](https://docs.vllm.ai/en/latest/configuration/optimization/)：KV Cache、批处理、并发和显存相关调优。
- [NVIDIA Triton Inference Server](https://docs.nvidia.com/deeplearning/triton-inference-server/user-guide/docs/)：Dynamic Batching、模型仓库、调度和指标。
- [TensorRT-LLM Architecture](https://nvidia.github.io/TensorRT-LLM/architecture/overview.html)：推理 Runtime、并行策略和优化组件。
- [Hugging Face Text Generation Inference](https://huggingface.co/docs/text-generation-inference/index)：连续批处理、量化和分布式模型 Serving 实践。
- [OpenAI Agents SDK: Agents](https://openai.github.io/openai-agents-python/agents/)：Agent、Runner、Tool、Handoff 和 Guardrail 的基础模型。
- [OpenAI Agents SDK: Tools](https://openai.github.io/openai-agents-python/tools/)：工具调用、审批、超时和错误处理。
- [OpenAI Agents SDK: Testing](https://openai.github.io/openai-agents-python/testing/)：Agent 编排、工具与 Guardrail 的确定性测试方法。

### 21.9 Kubernetes 与 AI 存储

- [Kubernetes Storage](https://kubernetes.io/docs/concepts/storage/)：Volume、PV、PVC、StorageClass、Snapshot 等文档入口。
- [Persistent Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)：PV/PVC 生命周期、访问模式和回收策略。
- [Storage Classes](https://kubernetes.io/docs/concepts/storage/storage-classes/)：动态供应、拓扑绑定和存储参数。
- [Container Storage Interface Specification](https://github.com/container-storage-interface/spec)：CSI Identity、Controller 和 Node Service 的协议规范。
- [CSI Developer Documentation](https://kubernetes-csi.github.io/docs/)：Sidecar、部署拓扑、Snapshot 和驱动开发。
- [PyTorch Distributed Checkpoint](https://docs.pytorch.org/docs/stable/distributed.checkpoint.html)：分布式模型和优化器状态的并行保存/加载 API。
- [PyTorch Asynchronous Checkpointing](https://docs.pytorch.org/tutorials/recipes/distributed_async_checkpoint_recipe.html)：降低 Checkpoint 对训练主循环的阻塞。
- [PyTorch DataLoader](https://docs.pytorch.org/docs/stable/data.html)：多进程读取、Pinned Memory、Prefetch 和 Worker 模型。
- [NVIDIA GPUDirect Storage](https://docs.nvidia.com/gpudirect-storage/index.html)：GPU 与存储之间直接 DMA 的入口文档。
- [GPUDirect Storage Design Guide](https://docs.nvidia.com/gpudirect-storage/design-guide/index.html)：GDS 软件栈、拓扑和数据路径设计。

### 21.10 Kubernetes 网络、RDMA 与 RoCE

- [Kubernetes Services, Load Balancing, and Networking](https://kubernetes.io/docs/concepts/services-networking/)：Service、EndpointSlice、Ingress、Gateway 和 NetworkPolicy 总入口。
- [CNI Specification](https://github.com/containernetworking/cni/blob/main/SPEC.md)：Runtime 调用网络插件的标准接口。
- [Cilium Architecture](https://docs.cilium.io/en/stable/overview/component-overview/)：基于 eBPF 的容器网络、Service 和策略实现。
- [Cilium kube-proxy replacement](https://docs.cilium.io/en/stable/network/kubernetes/kubeproxy-free/)：eBPF Service 转发路径和部署模式。
- [Multus CNI](https://k8snetworkplumbingwg.github.io/multus-cni/docs/quickstart.html)：为 Pod 附加多张网络接口。
- [SR-IOV Network Device Plugin](https://github.com/k8snetworkplumbingwg/sriov-network-device-plugin)：在 Kubernetes 中发现和分配 VF 资源。
- [RDMA CNI](https://github.com/k8snetworkplumbingwg/rdma-cni)：把 RDMA 设备放入 Pod Network Namespace。
- [RDMA Aware Networks Programming User Manual](https://docs.nvidia.com/rdma-aware-networks-programming-user-manual-1-7.pdf)：Verbs、QP、CQ、MR、Send/Receive 和 RDMA CM。
- [NVIDIA RoCE Documentation](https://docs.nvidia.com/networking/display/nvidiamlnxofeddocumentationv24101140105lts/rdma%2Bover%2Bconverged%2Bethernet%2B%28roce%29)：RoCEv1/v2、GID、PFC 和流量分类。
- [NCCL GPU Troubleshooting](https://docs.nvidia.com/deeplearning/nccl/user-guide/docs/troubleshooting/gpu_troubleshooting.html)：PCIe、NVLink、ACS、容器拓扑发现等问题。
- [GPUDirect RDMA Documentation](https://docs.nvidia.com/cuda/gpudirect-rdma/)：GPU 显存与第三方 PCIe 设备之间的 DMA 路径。

### 21.11 弹性伸缩与冷启动

- [Kubernetes Autoscaling Workloads](https://kubernetes.io/docs/concepts/workloads/autoscaling/)：HPA、VPA 和工作负载伸缩总览。
- [Horizontal Pod Autoscaling](https://kubernetes.io/docs/concepts/workloads/autoscaling/horizontal-pod-autoscale/)：算法、指标、容忍区间、稳定窗口和行为配置。
- [Node Autoscaling](https://kubernetes.io/docs/concepts/cluster-administration/node-autoscaling/)：节点供给、调度约束和自动扩容交互。
- [KEDA Concepts](https://keda.sh/docs/latest/concepts/)：事件驱动伸缩、ScaledObject 和 Scale-to-Zero。
- [Knative Serving Autoscaling](https://knative.dev/docs/serving/autoscaling/)：并发驱动伸缩、Scale-to-Zero 和激活链路。
- [containerd Stargz Snapshotter](https://github.com/containerd/stargz-snapshotter)：eStargz、远程 Snapshot 和 Lazy Pulling。
- [Stargz Snapshotter Overview](https://github.com/containerd/stargz-snapshotter/blob/main/docs/overview.md)：按需取块、Prefetch 和本地缓存机制。
- [Nydus RAFS](https://nydus.dev/docs/concepts/rafs/)：面向容器镜像的按需加载、块级缓存和校验。
- [Dragonfly Documentation](https://d7y.io/docs/)：大规模镜像、模型和文件的 P2P 分发。
- [OpenTelemetry Traces](https://opentelemetry.io/docs/concepts/signals/traces/)：为端到端冷启动瀑布图建立 Trace/Span 模型。

### 21.12 资源管控与多租户

- [Resource Management for Pods and Containers](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/)：Request、Limit、cgroup 和扩展资源语义。
- [Pod Quality of Service Classes](https://kubernetes.io/docs/concepts/workloads/pods/pod-qos/)：Guaranteed、Burstable、BestEffort 和驱逐优先级。
- [Resource Quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/)：Namespace 级对象与计算资源配额。
- [Limit Ranges](https://kubernetes.io/docs/concepts/policy/limit-range/)：单个 Pod、Container、PVC 的默认值和边界。
- [Pod Priority and Preemption](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-priority-preemption/)：优先级、抢占和非抢占式 PriorityClass。
- [Pod Disruptions](https://kubernetes.io/docs/concepts/workloads/pods/disruptions/)：PDB、自愿中断和高可用约束。
- [Multi-tenancy](https://kubernetes.io/docs/concepts/security/multi-tenancy/)：软多租户、硬多租户及隔离边界。
- [NVIDIA GPU Time-Slicing](https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/latest/gpu-sharing.html)：共享 GPU 的配置、隔离差异和限制。
- [Kubernetes Security Checklist](https://kubernetes.io/docs/concepts/security/security-checklist/)：API、工作负载、网络和认证授权基线。

### 21.13 分布式系统、高可用与平台工程

- [Designing Data-Intensive Applications course materials](https://github.com/ept/ddia-references)：《Designing Data-Intensive Applications》各章的论文与资料索引。
- [Raft Paper](https://raft.github.io/raft.pdf)：共识、Leader Election、日志复制和成员变更的原始论文。
- [etcd Learning](https://etcd.io/docs/v3.7/learning/)：Raft、API 保证、数据模型和工程设计。
- [etcd API Guarantees](https://etcd.io/docs/v3.7/learning/api_guarantees/)：线性一致性、Watch 顺序和租约语义。
- [Google SRE Book](https://sre.google/sre-book/table-of-contents/)：SLO、容量、过载、故障处理和发布工程。
- [Google SRE: Service Level Objectives](https://sre.google/sre-book/service-level-objectives/)：SLI、SLO、SLA 与错误预算。
- [Google SRE: Handling Overload](https://sre.google/sre-book/handling-overload/)：限流、负载保护和降级策略。
- [Google SRE: Addressing Cascading Failures](https://sre.google/sre-book/addressing-cascading-failures/)：级联故障、重试风暴、慢启动和容量问题。
- [CNCF Platforms White Paper](https://tag-app-delivery.cncf.io/whitepapers/platforms/)：平台工程、能力抽象和开发者体验。
- [Kubernetes Multi-cluster SIG](https://multicluster.sigs.k8s.io/)：多集群 API、工作组和相关项目入口。

### 21.14 可观测性与性能工程

- [Prometheus Data Model](https://prometheus.io/docs/concepts/data_model/)：时序数据、Metric Name、Label 与 Cardinality。
- [PromQL Basics](https://prometheus.io/docs/prometheus/latest/querying/basics/)：查询、聚合、Range Vector 和函数。
- [Prometheus Instrumentation Practices](https://prometheus.io/docs/practices/instrumentation/)：在线服务、离线任务和批处理的指标设计。
- [Prometheus Histograms and Summaries](https://prometheus.io/docs/practices/histograms/)：延迟分布、分位数和聚合取舍。
- [OpenTelemetry Concepts](https://opentelemetry.io/docs/concepts/)：Signals、Context Propagation、Sampling 和 Semantic Conventions。
- [OpenTelemetry Signals](https://opentelemetry.io/docs/concepts/signals/)：Metrics、Logs、Traces、Baggage 和 Profiles。
- [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/)：Receiver、Processor、Exporter 和部署拓扑。
- [Linux perf Examples](https://www.brendangregg.com/perf.html)：CPU、调用栈、Off-CPU 和硬件事件分析示例。
- [FlameGraph](https://github.com/brendangregg/FlameGraph)：火焰图生成工具与解释资料。
- [BPF Performance Tools Book Examples](https://github.com/brendangregg/bpf-perf-tools-book)：BCC/bpftrace 性能工具示例。
- [Google SRE: Monitoring Distributed Systems](https://sre.google/sre-book/monitoring-distributed-systems/)：四个黄金信号、白盒/黑盒监控和告警设计。
- [Google SRE: Service Level Objectives](https://sre.google/sre-book/service-level-objectives/)：从用户体验反推性能与可靠性指标。
- [NVIDIA DCGM Exporter](https://github.com/NVIDIA/dcgm-exporter)：将 GPU/DCGM 指标暴露给 Prometheus。

### 21.15 镜像服务与供应链安全

- [OCI Distribution Specification](https://github.com/opencontainers/distribution-spec/blob/main/spec.md)：Registry API 的端点、错误码和一致性要求。
- [CNCF Distribution](https://github.com/distribution/distribution)：OCI Registry 的参考实现。
- [Harbor Documentation](https://goharbor.io/docs/)：企业级 Registry 的复制、扫描、签名、配额和高可用实践。
- [ORAS Documentation](https://oras.land/docs/)：通过 OCI Registry 分发镜像之外的 Artifact。
- [Dragonfly Documentation](https://d7y.io/docs/)：P2P 镜像分发、预热和大规模下载加速。
- [Nydus Documentation](https://nydus.dev/docs/)：块级镜像格式、远程加载和运行时集成。
- [Sigstore Cosign: Signing Containers](https://docs.sigstore.dev/cosign/signing/signing_with_containers/)：容器镜像的签名、验证与证明。
- [SLSA Specification](https://slsa.dev/spec/)：构建来源、供应链完整性和成熟度模型。
- [SPDX Specification](https://spdx.dev/specifications/)：SBOM 的开放标准和数据模型。
- [Kubernetes Image Verification](https://kubernetes.io/docs/tasks/administer-cluster/verify-signed-artifacts/)：Kubernetes 发布制品签名验证示例。

### 21.16 智能运维 Agent

- [Google SRE: Effective Troubleshooting](https://sre.google/sre-book/effective-troubleshooting/)：Triage、观察、假设、实验与根因定位方法。
- [Google SRE: Managing Incidents](https://sre.google/sre-book/managing-incidents/)：事件角色、协调、沟通和状态管理。
- [Google SRE: Postmortem Culture](https://sre.google/sre-book/postmortem-culture/)：无责复盘、根因与行动项。
- [OpenAI Agents SDK: Tools](https://openai.github.io/openai-agents-python/tools/)：把诊断、查询和修复动作建模为受控 Tool。
- [OpenAI Agents SDK: Guardrails](https://openai.github.io/openai-agents-python/guardrails/)：输入、输出和工具执行前后的约束。
- [OpenAI Agents SDK: Tracing](https://openai.github.io/openai-agents-python/tracing/)：记录模型调用、工具调用、Handoff 和自定义 Span。
- [OpenAI Agents SDK: Human-in-the-loop examples](https://openai.github.io/openai-agents-python/examples/)：工具审批、状态序列化和人工介入模式。
- [MITRE ATLAS](https://atlas.mitre.org/)：AI 系统攻击技术与安全威胁建模，可用于设计 Agent 权限边界。

### 21.17 开源贡献

- [Kubernetes Contributor Guide](https://www.kubernetes.dev/docs/guide/)：社区结构、开发流程、测试和 PR 规范。
- [Kubernetes Community](https://kubernetes.io/community/)：SIG、Slack、会议和贡献入口。
- [Kubernetes Enhancements](https://github.com/kubernetes/enhancements)：KEP 流程、功能生命周期和设计示例。
- [Kubernetes Code Generation](https://github.com/kubernetes/code-generator)：API、Client、Informer 和 Lister 代码生成工具。
- [Volcano Contributing Guide](https://github.com/volcano-sh/volcano/blob/master/CONTRIBUTING.md)：Volcano 的开发与贡献流程。
- [Kueue Contributing Guide](https://github.com/kubernetes-sigs/kueue/blob/main/CONTRIBUTING.md)：Kueue 的开发、测试和贡献入口。
- [containerd Contributing Guide](https://github.com/containerd/containerd/blob/main/CONTRIBUTING.md)：containerd 社区、DCO 和代码贡献规范。

### 21.18 推荐阅读顺序

1. **第一阶段：系统与语言**：21.1 → 21.2 → 21.3，建立 Go、Linux、OCI 和 Runtime 的底层模型。
2. **第二阶段：Kubernetes 核心**：21.4 → 21.5 → 21.12，掌握控制循环、调度和资源治理。
3. **第三阶段：AI 工作负载**：21.6 → 21.7 → 21.8，理解 GPU、训练通信和推理服务。
4. **第四阶段：高性能数据路径**：21.9 → 21.10 → 21.11，学习存储、RDMA/RoCE 和冷启动。
5. **第五阶段：生产工程**：21.13 → 21.14 → 21.15 → 21.16，形成可靠性、可观测和自动化闭环。
6. **持续进行**：21.17，边学习边复现 issue、阅读源码并尝试小型开源贡献。
