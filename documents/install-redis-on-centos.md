---
title: "Linux: 给Centos 7.3装Node/Redis踩坑记"
tags:
  - Server
  - Nodejs
categories:
  - Tech
date: 2017-08-30 19:32:12
---

### 前言·吐槽

真的是气得要死，本来装个Redis是非常简单的事情，但是RedHat偏要增加我的工作量，真的是很气TAT

安装Redis倒是容易，但是随之而来的Nodejs以及连接库的安装，让我真的哭笑不得。本来以为5分钟能搞定的事情硬是花了我半个多小时。为防止接下来一段时间里面继续有人踩坑，就稍微写篇软文吐槽一下。

<!-- More -->

### Redis

这东西很好装，到Redis的官网下载相关软件包，编译即可。

值得一提的是，在网上找到了比较好的init.d配置脚本，因为这个脚本出现的太平凡，且没有人索引出处。我想Qoute也无从下手。于是就只能在此一说，还希望网友们在引用别人文章内容的时候指明出处。（TAT

脚本如下：测试服务器已下线

这样，只要把默认的配置文件放在`/etc`目录下，修改为用Deamon启动后，就可以直接通过`/etc/init.d/redis start`来启动redis-server

```bash
  cp ./redis.conf /etc/redis.conf
```

### Nodejs - 坑

虽然我有一段时间支持过“一切皆编译”。但是没过多久我就成了编译安装的反对者之一。虽然编译安装带来了依赖稳定性以及灵活性和可移植性。但是它会大大的消耗系统运算性能。不仅费时，而且会影响邻居（当然，较好的KVM和独服除外）。何况像我这样比较追求系统稳定性的人不会对系统进行大改，因此完全可以信任发行版的软件包管理器。

但是，这次似乎有点打脸。

```bash
  yum install -i node npm --enablerepo=epel
```

这是最简单的安装最新版本nodejs的方法，但是，这次的问题是，抛出了一个错误：缺乏依赖 http-parser

> Requires: http-parser >= 2.7.0

这就很纳闷了，难道是因为阿里云的mirror出了问题？

查了一下，发现http-parser原本确实存在于epel包中，于是我一脸懵逼。

俗话说得好，搜索引擎：从百度到自杀。还是只有Google才能给我答案...

在centos官方bug反馈中，找到了这一个帖子<a href="#sup1"><sup>[1]</sup></a>，中间有一段话：

> http-parser was added to the RedHat Base repository for 7.4, there fore EPEL removed it.

权贵红毛迟早药丸。你不知道，centos用户也会使用epel包吗。。？

解决方法是，手动安装`http-parser`的rpm包（或者编译）

https://kojipkgs.fedoraproject.org/packages/http-parser/2.7.1/3.el7/x86_64

个人建议把devel和标准包都装了，因为node的`node-gyp`编译库需要devel-http-parser的支援。

### 后言

到此，踩坑结束。

RedHat药丸。

以上。

### Works Cited
<sup id="sup1">[1]</sup> https://bugs.centos.org/view.php?id=13669&nbm=1
