---
title: 'Golang: Iterator Parameter Cheating'
date: 2020-06-05T09:45:54.484Z
categories:
  - Tech
tags:
  - Golang
  - Server
---

### 前语

这篇文章没什么技术含量，仅仅是对自己的一个警示。对大多数阅读这篇文章的同学来说，本文讲述的问题你们可能早已知晓，因此我也就直奔主题。
毕竟这篇文章所提及的问题不仅仅在 `Golang` 而可能发生在所有支持 async 的场景中。（haskell 等纯函数式编程语言除外）

### 问题

我们先来看我之前的一段代码：（里面是有问题的）

```golang
for rule := range RuleStore.GroupRulesIteratorWithValidations(m.Chat.ID) {

    // Some other code goes here ...

    // pack the function to be run later
    handlerFunc := func() {
        for _, a := range rule.Action {
            if a.Instance.Handle(m, argv) == 0 {
                if a.OnErr != nil {
                    a.OnErr.Instance.Handle(m, argv)
                }
                break
            }
        }
    }

    // run asyncly if not on lambda
    if executor.Synchronous() {
        handlerFunc()
    } else {
        go handlerFunc()
    }

}
```

这是来自我的 telegram bot 的一段代码，大致意思是根据 群聊id 匹配当前群聊开启的机器人策略，并依次执行。其中 `m` 是用户发送的消息实力，`argv` 是 Context 中的一个 memcache，用于保存临时变量。
`rule` 则是一个机器人策略，它包含 `Trigger` 和 `Action`，也就是说多个用户行为可以触发一个策略，并执行复数的 `Action` 来回应该行为。（例如：一个 helloworld 的策略包含一个触发器和一个执行器，
触发器为 **"被回复消息"**，执行器为 **"发送消息 Helloworld"**。这就是个简单的策略。）

为了应对高并发场景，我设计了一个简单的分布式，各个 worker 从消息队列里取出最新的 message 然后分别执行。若当前 worker 支持使用协程，则开启一个新的协程跑策略（因为有 worker 是部署在 aws lambda 上，超时会被直接挂起，因此这类 worker 不能使用协程），这就是代码的最后一段表示的内容。

乍一看其实代码也没啥逻辑上的问题，但当我部署了代码，机器人就炸了，可复现的行为是：一则消息居然触发了它本不该触发的策略...

### 排查

其实对于一个 js 老手来说，这个问题实在是不该犯。对于一个写闭包就像拉拉链这么简单的我而言，我又又又栽在了变量作用域上。（是的，你找到问题了吗？）

在最外层 `for` 循环中，每一次循环迭代器会给 `rule` 赋值，因此接下来的代码可以用 `rule` 表示所有既存的策略。

这时候，问题就来了：如果顺序执行，那就算跑到世界末日也不会出错。因为 `for` 内调用的 `rule` 就是这一层迭代器吐出的 `rule`。然而，并发的世界，就不是这样了：

因为 `handlerFunc` 被延迟执行，所以它当中的 `rule` 只有运行时才会动态获取，也就是，如果在执行 `handlerFunc` 时，当前 `for` 已经结束，那它的 `rule` 就已经改变了。

因此，有两种解决方案：（我只想到了两种）

- 缓存这个 `rule`，`snapshotRule := rule`。在 `handlerFunc` 中使用 `snapshotRule`。因为 `snapshotRule` 的作用域是内部循环，因此对于不同循环而言它是不同的。
- 当作参数传入，变成 `handlerFunc(rule *rules.Rule)`

### 解密

这里我写了一个小小程序来解释我上面内容的正确性，大家感兴趣不妨跑一下：

```golang
package main

type Unit struct {
	Number int
}

func Num(max int) <-chan *Unit {
	channel := make(chan *Unit)
	go func() {
		for i := 0; i < max; i++ {
			num := &Unit{i}
			channel <- num
		}
		close(channel)
	}()
	return channel
}

func main() {
	for i := range Num(10) {
		println(&i, i, i.Number)
	}
}
```

最后的结果可能是（内存基址不一样，结果会有一定的变化）：

```bash
0xc000036760 0xc00001a088 0
0xc000036760 0xc00001a090 1
0xc000036760 0xc00001a098 2
0xc000036760 0xc00001a0a0 3
0xc000036760 0xc00001a0a8 4
0xc000036760 0xc00001a0b0 5
0xc000036760 0xc00001a0b8 6
0xc000036760 0xc00001a0c0 7
0xc000036760 0xc00001a0c8 8
0xc000036760 0xc00001a0d0 9
```

因此以上理论也得以证实