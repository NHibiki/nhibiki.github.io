---
title: 'Multipass: A light-weight VM helper'
date: 2020-05-23T16:01:54.484Z
categories:
  - Tech
tags:
  - Container
  - Ubuntu
  - Server
---

### 前语

可以说，**虚拟化/容器** 是近几年最火的技术名词之一（二？）了。从 docker 正式开启容器的知识普及，到 CoreOS/Rancher 推动容器 OS 的进程，最后又有 Docker Swarm 和 K8S 为始点的集群改革，整个基础设施界在翻天覆地的变化。

一方面，互联网无止境的扩张和 Distributed Systems 的实践对系统/硬件/网络资源的调度需求越来越大；另一方面，区块链/Raft（etcd）等的技术试验（成熟化）又让可信化基础设施与自动化编排大量地释放了生产力。
这些因素导致了大多数技术也开始向高效率虚拟化与区中心化的指数级发展。

而今天，我要记录的是我使用 [Multipass](https://multipass.run) 的心得。在这个 golang 横飞的时代，multipass 坚持使用了 C++ 作为开发语言[（repo）](https://github.com/canonical/multipass)，可以说确实很会选型（这里不讨论 golang 与 C++ 的优劣，个人拙见：编译语言中善用编译器特性比脱离情景地 benchmark 更有说服力）
而它的编写语言也暗示了 Multipass 的定位：它不会去争夺硝烟中容器霸主的地位，它的目的只有一个：快速又安全地构建本地虚拟化。换言之，它是虚拟化的上层工具，而并非单纯的容器。

### 安装

这点我想我可以略过了，Multipass 的安装非常简单。我使用的测试机是 DigitalOcean 2C4G Ubuntu 18.04，并且使用 snap 对其进行安装和测试。

### 测试环节

大部分的测试我都借鉴了[这篇博客](https://www.freshbrewed.science/ubuntu-multipass-better-than-docker/index.html)，在测试中，原作者成功地运行了 [K3S](https://k3s.io) 一个于 K8S 兼容 API 的轻量级集群容器管理工具，并通过 `kubectl` 成功初始化集群并搭建了 [k8s-dashboard](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/)。而我反过来测试了 docker swarm，可以说是预料得到的可以运作（这句话好像有点拗口（￣▽￣））。

于是，Multipass 是一个操作如 docker 般简单，但有能灵活而高效地编排 vm 的工具。根据上面那篇文章，它在 Mac 下使用 HyperKit 运行虚拟化，查看了官方手册后了解它在 Linux 下使用的是 QEMU 进行虚拟化，而 Windows 是 Hyper-V（这也解释了为什么官方要求 Windows 10 Pro 或者 Enterprise 才能运行 Multipass）。
既然它创建的 Ubuntu 虚拟机都是 vm，那除非一些直接与 CPU 优化或者 GPU 资源挂钩的场景，大多数情况都与完整的系统无异了（它具有完全独立的内核）。

因此，我们其实不难判断，理论上而言，在 Windows 和 Mac 下用 Multipass 起 Ubuntu 实例的效率应当比 Docker 高（毕竟 Docker 需要一层虚拟化+一层容器化）。
而这也正好对应了官方的一大宣传：一句命令就能初始化的完美开发环境。

除此之外，为了对标（贴合）`docker-compose.yml` 等初始化配置文件带来的便捷，Multipass 也有完整的 `cloud-init` 支持。在初始化 vm 时，我们可以手动传入 `cloud-init` 脚本对虚拟机进行初始化，从而提高构建效率（那我是不是可以理解，之后会和 docker 一样允许 socks 或者 grpc 的全自动编排呢？）。

### 总结

这次的测试仅仅是尝了一个鲜，让我知道了还有这样一个方便的虚拟化工具，并且它的势头不减容器化。目前对于我而言，它的使用场景更适合放在开发上 -- 它可以完美替代 WSL，也可以省去 Mac 下 docker 那些烦心事（我可以直接在 Ubuntu VM 中安装 docker）。

因此，对我而言 Multipass 目前只有在开发者社区中的有着独到的优势。它今后的路会怎样，能不能真正在众多虚拟产品中站稳脚跟，就很难预料了。