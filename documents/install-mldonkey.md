---
title: "Linux: 安装万能下载工具 MlDonkey"
date: 2016-07-15 09:54:37
tags:
  - Server
  - Download
categories:
  - Tech
---

### 简介

说起Linux中的文件下载，在我们的脑海中立刻浮现的可能是 wget / aria2 。这些可谓是 Storige Server 的必备工具了。但是，他们还是有一个缺陷。wget只能下载直连，这就不多说了。aria2的下载范围更广一些，但它也只能解析torrent，没办法下载电驴等其他p2p协议下的文件。

图形系统下的工具我们就不提了，我们需要的是一个命令行+Web管理的下载工具。于是乎，MlDonkey就浮现在我们眼前了。

MlDonkey不是一个新项目，我倒Sourceforge看了一下，最早Push的时间是在 2006-01-05。最近一次更新是在 2014-03-22。但从我个人的体验上看，这个项目还是不错的。

<!-- More -->	

### 安装

走了不少弯路，就像一开始，我看到网上有很多人说，直接用静态编译版（static）就好了，结果就出现了莫名其妙的错误。【谁说的！给我打一顿 TAT

筋疲力竭的我决定还是自己编译2014年的源代码QAQ

首先，你要确保自己的服务器中有 gcc / g++ / make / wget / zlib 存在，如果没有，请安装。（给出的是Ubuntu的命令，其他Linux请自行修改Repository指令）

```bash
apt-get update && apt-get install -y gcc g++ make wget zlib1g-dev
```

然后，下载“最新”的MlDonkey，解压。编译，安装

> 重要提示：zlib一般Linux中都不自带，一般情况下都要手动安装

```bash
wget https://sourceforge.net/projects/mldonkey/files/mldonkey/3.1.5/mldonkey-3.1.5.tar.bz2
//下载
tar -jxvf mldonkey-3.1.5.tar.bz2
//解压
cd mldonkey-3.1.5
./configure
//配置，中间它会要求下载一个附加库，Y回车即可
make
//编译
make install
//安装
```

至此，MlDonkey安装结束。

你可以直接

```bash
mlnet
```

如果安装成功，那么最后一行指令应该是 [dMain] Core started

### 配置

Ctrl + C 结束进程

进入配置目录

```bash
cd ~/.mldonkey
```

编辑downloads.ini (vim/vi/或其他都可以)

找到 allowed_ips = ["127.0.0.1";]  改为 allowed_ips = ["0.0.0.0/0";] 

在42%处找到 shared_directories 【在vim中，你也可以使用查找命令 “/shared_directories”】直接跳转到该位置

把两个incoming的dirname改为自己想要的目录，这里我就设置成了“/home/files”

修改后的片段如下：

```ini
shared_directories = [
{
  dirname = shared
  strategy = all_files
  priority = 0
};
{
  dirname = "/home/files"
  strategy = incoming_files
  priority = 0
};
{
  dirname = "/home/files"
  strategy = incoming_directories
  priority = 0
};
]
```

保存，mkdir建立设置的目录。然后，你就可以用screen命令启动mlnet了

进入 http://你的服务器ip:4080

在命令栏中输入 useradd admin 你的密码，这样，以后用web访问，就需要登陆了

至此，教程结束

以上
