---
title: "Gentoo: 給Aliyun的Gentoo升级"
tags:
  - Server
  - Gentoo
categories:
  - Tech
date: 2017-01-18 14:37:23
---

 > 已经下定决心要好好玩玩Linux了。虽然我并没有勇气格电脑装Linux（对于习惯Mac的人来说似乎不太舍得这么做）。我想尽量装个双系统，或者在Cloud Platform上装一个深坑一点的Linux系统。比如ArchLinux 和 Gentoo。
 
 于是我就选择了Aliyun的Gentoo系统，打算先玩玩练练手。（毕竟现在Aliyun坑得只能用来研究了...除了性能，各种不如腾讯云【当然说不定只是一时的】）
 
 但是最近准备給Aliyun按Docker和BBR的时候，突然发现按不上去。检查以后才知道是内核与系统版本太低（2014年的版本！！？）于是，虽然对内核与Linux运作机制遠不了解，但也只能硬着头皮上了。
 
 <!--More-->
 
首先，我还是得感谢一下Aliyun。虽然它坑，但似乎它对自己的运算性能非常有把握.. 一是，第一次看到能在VPS上运行高运算量的Gentoo（毕竟用它的包管理器最常见的就是下载编译）很多服务商都不太敢引入Gentoo。二是，它居然没给系统升级到较新..TAT..这明摆着是想让我们自己升级，无疑又要消耗很多计算资源【噫

### 尝试

为了升级Gentoo.. 失败了很多次.. 很多坑都是出在内核升级上

毕竟我们要将内核从` 3.12 -> 4.4 `跨度很恐怖（当然，之后为了开启BBR，我又将它升级到了 4.9.3）

我按照官方Wiki一步一步来：

先 `emerge --sync` 更新包管理源

然后 `emerge -avq genkernel` 下载内核更新工具（用于initramfs更新）

再 `emerge -avq gentoo-sources` 下载内核源码

    如果想直接越过稳定版本跳至 ~arch 下最新的版本，可以手动fetch内核源码并解压到 /usr/src 目录中

4.4阿里云镜像地址：http://mirrors.aliyun.com/gentoo/distfiles/linux-4.4.tar.xz

4.9阿里云镜像地址：http://mirrors.aliyun.com/gentoo/distfiles/linux-4.9.tar.xz

然后就可以配置内核了

我按照以前的一点点非常单薄的知识，执行以下代码：

```bash
    make -j4 oldconfig \ 
    make -j4 \ 
    make install \ 
    make modules_install \ 
    genkernel initramfs
```

oldconfig中，如果有想加入的新功能就自行钩上，比如如果想开启Docker功能，可以按照Neil的方法，把相关功能开启。[传送们](https://nrechn.de/post/solve-docker-fail-to-start/)

其他不变就一路回车。

最后，阿里云用户请修改grub（阿里云的grub很让人抓狂.. 全部不启用ln软连接，升级起来很麻烦，这时候只要手动把kernel和initramfs的对应项改成ln软连接就行了

```bash
    linux   /boot/vmlinuz root=UUID=xxxxxxxx(你自己的Disk UUID) ro single 
    echo    'Loading initial ramdisk ...'
    initrd  /boot/initramfs
```

PS：一共三处

这样每次升级只需要 `ln -s /boot/xxx-?.?.?-gentoo /boot/xxx` 就好了，不用再次更改grub，非常方便

接下来，reboot

### 失败

然后，很奇怪的事情发生了，ssh怎么都连不上去。

好在Aliyun提供pin接口，可以用net console直接链接。

检查防火墙，检查sshd，统统没有问题。

偶然一次，打算 `emerge` 的时候，发现了错误 -- 没有连上网

整个服务器都是掉线状态。

尝试了 /etc/init.d/net.eth0 之后，找到问题：没有相关驱动，也就是说，还是内核的问题

### 成功

就这么把这个问题放了一天，后来在和Neil TG的时候，Neil推荐我尝试不要用oldconfig来fetch存在的信息，而是直接用localmodconfig，虽然这个经常用来精简系统233，但是它的功能正是我想要的（它会自动检测当前系统加载的modules并且精简或增加.config的选项配置

然后再次编译，reboot，成功。

接下来只要更新python, portage然后直接更新整个系统树就好了。

至于到底是什么驱动使得我们的VPS无法连接到网络，我就不打算一一对着看了（就是偷懒嘛。

于是就把两份.config先后传到git上，用git来检测一下哪里的config出了问题（

以上。
