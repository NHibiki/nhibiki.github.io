---
title: "AudioStream: 使用 ICECAST 和 ICES 配合Nodejs做一个 Audio Stream"
tags:
  - AudioStream
  - Music
categories:
  - Tech
date: 2017-08-01 20:21:19
---

> 所谓的 Audio Stream 就是一个音频流媒体服务，可以让连接到流媒体的用户听到一样的东西。

### 简介

虽然 socket.io 已经提供了很好的流媒体方案，但是鉴于 socket.io流 本身更适用于视频或者画面流，对音频流不支持，于是就有了下面这篇文章。在 Centos 7 中，搭建一个 ICE 系列产品的生产环境。

用到的工具：

> yum wget tar gcc gcc-c++ make

<!--More-->

### 安装方法

先说在前面，如果不想安装，可以使用我已经封装好的 docker，但是因为这个 docker 并不是用于发布的，因此没有精简，所有的工具都在里面 23333，甚至于包括了 yum 的 repo缓存。完整的大小是 600 MB，docker 压缩之后是 207 MB。

#### 安装

先安装必要库文件，这些都是在默认 repo 中的，如果有强迫症，也可以自行编译。

```bash
yum install -y libxslt-devel libogg-devel libvorbis-devel libcurl-devel libshout-devel
```

然后再安装我们需要用到的一些软件

首先是 ICECAST，是提供流媒体服务的服务器。

```bash
wget http://downloads.xiph.org/releases/icecast/icecast-2.4.3.tar.gz
tar zxvf icecast-2.4.3.tar.gz
rm icecast-2.4.3.tar.gz -rf
cd icecast-2.4.3
```

常规安装方法，顺便指定一下 prefix

```bash
./configure --prefix=/usr/local/icecast
make && make install
cd ..
```

接下来是 LAME 库，一个解码 mpeg3 文件的东西，安装方法同前

```bash
wget https://downloads.sourceforge.net/project/lame/lame/3.99/lame-3.99.5.tar.gz
tar zxvf lame-3.99.5.tar.gz
rm lame-3.99.5.tar.gz -rf
cd lame-3.99.5
./configure --prefix=/usr/local/lame
make && make install
cd ..
```

最后一个是 ICES，也是 ICE 系列产品，用于向 STREAM 推送文件

```bash
wget http://downloads.us.xiph.org/releases/ices/ices-0.4.tar.gz
tar zxvf ices-0.4.tar.gz
rm ices-0.4.tar.gz -rf
cd ices-0.4
```

安装 ICE 时要指定一下 lame 的目录，不然会搜不到

```bash
./configure --prefix=/usr/local/ices --with-lame=/usr/local/lame
make && make install
cd ..
```

连接文件

```bash
ln /usr/local/ices/bin/ices /usr/bin/ices
ln /usr/local/icecast/bin/icecast /usr/bin/icecast
```

这样，所有的依赖以及软件都装完了。

#### 配置

因为 ICECAST 的权限限制，我们尽量不要用在 chroot 用户组的用户运行 STREAM，因此，新建用户

```bash
groupadd nekobc
useradd nekobc -m -g nekobc -G users,wheel,audio -s /bin/bash
passwd nekobc
mkdir /home/nekobc
chown nekobc:nekobc /home/nekobc
```

下面要配置一下 ICECAST 和 ICE 的文件，创建音乐文件目录

```bash
cp /usr/local/icecast/etc/icecast.xml ./icecast.xml
cp /usr/local/ices/etc/ices.conf.dist ./ices.xml
mkdir music
```

其中 `icecast.xml` 文件中要修改的有如下部分：

```xml
<source-password></source-password> # 给 ICES 的密码
<relay-password></relay-password>   # 可以不管，建议改掉初始密码
<admin-user></admin-user>           # 网络管理用户名
<admin-password></admin-password>   # 网络管理密码
<logdir>/home/nekobc/log</logdir>

<listen-socket>
  <port>8010</port>
  <bind-address>0.0.0.0</bind-address>
  <mountpoint>/stream</mountpoint>
</listen-socket>

# 启用以下配置，把注释拿掉
<changeowner>
  <user>nekobc</user>
  <group>nekobc</group>
</changeowner>
```

然后是 `ices.xml` 文件的设置：

```xml
<File>list.txt</File>      # 指定列表文件
<Randomize>0</Randomize>   # 随机播放
<Background>1</Background> # 后台运行
<Hostname>localhost</Hostname>
<Port>8010</Port>
<Password></Password>
<Mountpoint>/stream</Mountpoint>
<Public>1</Public>
                           # 要和 ICSCAST 一致
```

其他的描述信息可以更具自己的需求改

最后，整理音乐列表，可以用这个命令

```bash
find / -name "*.mp3" > /home/nekobc/list.txt
```

启动服务

```bash
icecast -b -c icecast.xml
ices -c ices.xml
```

浏览器里面访问 `http://{Your_Url}:8010/stream` 就可以了。

以上
