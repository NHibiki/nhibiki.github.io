---
title: "WildCard: 来自 Let's Encrypt 的无料野卡申请指南"
tags:
  - SSL
categories:
  - Tech
date: 2018-03-16 22:56:16
---

> 16 Mar 2018 更新, 附上 Alpine 中申请 WildCard方法

### 题记

盼了半年的 Let's Encrypt WildCard 终于于昨天(15 Mar 2018)通过了测试, 正式对外开放. 这对于 LEC HTTPS 使用者来说是个炒鸡大的好消息.

WildCard 的优点在于, 它直接签发了通配符域名证书, 这样就不用单独为自己的诸多子域名签证了. 直接一张证书就可以加密自己所有的网站.

### 步骤

首先, 因为这次更新除了允许申请 WildCard 外, 还改写了很多 Policy. 根据 Let's Encrypt 官方博客, 为了不影响之前协议用户的使用, LEC 专门开放了另一个签发入口用于第二版本的证书签发. 地址如下

`https://acme-v02.api.letsencrypt.org/directory`

于是, 我们在使用 certbot 进行签证时, 需要指定签证地址, 以免误用 v1 入口导致签证失败.

#### 第一步

更新 certbot.

<del>
这里我采用了 certbot-auto 脚本, 需要注意的是, 此脚本只支持部分系统, 以完成自动安装. 如果目前你的系统是 alpine 等非特别主流的系统, 还是建议等官方源更新 certbot 或者 将证书迁移 来使用.
</del>

> Alpine 用户可以通过跨越 v3.7 Stable 使用 edge 激进编译包管理列表来安装最新的 certbot

```bash
# 非 Alpine 用户无视这条命令
echo "http://mirrors.aliyun.com/alpine/edge/main" >> /etc/apk/repositories
echo "http://mirrors.aliyun.com/alpine/edge/community" >> /etc/apk/repositories

# 或者不使用 Aliyun 镜像的孩纸们:
echo "http://dl-cdn.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories
echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
```

```bash
wget https://dl.eff.org/certbot-auto
chmod a+x certbot-auto
```

然后执行以下命令自动安装 certbot 环境, 并签发证书

```
./certbot-auto --server https://acme-v02.api.letsencrypt.org/directory -d "*.your.domain" certonly --manual
```

<span style="color:red;">这里要注意的一点是, 一定要有 `--manual` 这一选项, 不然会签发失败. (因为目前各大签证 bot 还没有实现 DNS API 的接入, 所以接下来的验证只能手动完成)</span>

这里将 `*.your.domain` 变成你自己的根域名

与 v1 不同的是, 它不会要求你开放临时服务器或者链接网站根目录来验证域名. v2 wildcard 采用了 DNS TXT Record 认证方式. 在最后一步, 把认证代码指向 _acme-challenge.your.domain 即可完成认证.

至此, 证书申请完毕. 可以到 `/etc/letsencrypt/live/your.domain/fullchain.pem` 获取你的 WildCard 证书.

另外需要强调的一点是, 签发周期依旧是 3月/次, 因此, 目前看来 auto-renew 的工作特别烦琐 ( 除非自己接入 DNS 的 API 自动修改 TXT. 整个脚本并不容易写出来啊 23333

所以, 尝鲜归尝鲜, 目前大多数服务, 还是以 v1 为准吧 QwQ

以上