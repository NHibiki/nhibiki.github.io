---
title: "JavaScript: Js 中深拷贝的探讨"
tags:
  - JavaScript
  - DeepCopy
categories:
  - Tech
date: 2019-03-01 14:27:22
---

### 前情提要

最近在继续加深自己对 *Javascript* 的理解（各种方面），主要通过复习 MDN 和 查看很多大佬的博客为主。发现一个非常 universal 的问题：他们都提到了 *deepCopy*。

我们知道，在这种脚本语言里面，拷贝大多都是 reference，而并非内存上的拷贝（手法类似于 C 或者 Go 里面的指针拷贝）。很多语言都给出了“官方”解决方案：如 *Python* 有一个 *copy* 库，下面存在 `copy.copy`（浅拷贝）和 `copy.deepcopy`（深拷贝）的方法。但是，在 *Javascript*中，官方并没有提供这类工具库。一方面觉得可能是因为 ECMA 懒，但这种指责确实是“站着说话不腰疼”。自己觉得更有可能的原因是另一方面：*Javascript*的对象都是基于原型的（prototype），一股脑的深拷贝会存在很大的依赖问题。（但其实这也不是借口，比较可以从 constructor 重建，只是浪费资源罢了）

网上探讨这个问题的文章很多，但很多都给出了非常 *dummy* 的 code（假到不用跑就知道有 bug...）。所以这篇文章诞生，但是，毕竟所有技术文章都是有时效性的，所以读者在借鉴这里的 *code* 的时候，请先观察下我的攥稿日期，以免造成不必要的麻烦。

### 分析

首先，在 *Javascript* 中新建一个 Object 不是一个难事，但很杂乱。已知的新建方法数不胜数：

1. 直接 `{}`
2. 使用 `Object.create()` 创建
3. 使用 `class` 构造原型，并 `new 出来`
4. 使用 `function` 构造上下文，并在外部定义原型
5. ...

总之是很多很多的。

这个时候，很多人已经发现网上流行的 *JSON* 大法多么不适了 -- *JSON* 支持的 *datatype* 实在少的可怜：

1. Number
2. Object（其实是 `HashMap<String, Object>`，这里包括 null）
3. Boolean
4. Array

就这样 LOL。它连 `function`，`Map` 与 `Set` 都不支持，更不要说 `__proto__` 或者被单独定义过的（`Object.defineProperty`）*Object* 了。

另外，细心的朋友可能会发现，除此之外 *JSON* 还有一个坑：对于如果 *Object* 内部出现圆环（如下），则就会 *stringify* 出错：

```js
a = {};
a.a = a;
JSON.stringify(a); // ERROR
```

### 着手

基本的思路还是老样子：递归拷贝。这里偷懒了，直接给出源代码。

当然，这份源代码虽然考虑了很多情况，但并不代表它一定能 cover all。还是得具体问题具体分析（比如，如果涉及到 **Vue** 中的 *Observer* 模型，这种拷贝就会出问题）。

```js
function deepCopy(target) {
    /* initialize, using hashmap to record recursive properties */
    const hashMap = new Map();

    /* where deepcopy happens */
    const doCopy = function doCopy(_target) {
        if (_target === null) {
            /* check if it is null */
            return null;
        } else if (typeof _target !== 'object') {
            /* is value variable */
            return _target;
        } else if (hashMap.has(_target)) {
            /* check if it is recurssive properties */
            return hashMap.get(_target);
        } else {
            const returns = new _target.constructor();
            hashMap.set(_target, returns);
            if (_target instanceof Array) {
                /* is numeric-iterable array or list */
                Object.assign(returns, _target.map(doCopy));
            } else if (_target instanceof Set) {
                /* is Set datatype */
                doCopy([..._target]).forEach(v => returns.add(v));
            } else {
                /* is other object datatype */
                if (_target instanceof Map) {
                    const iter = _target.keys();
                    let curr;
                    while (!(curr = iter.next()).done) {
                        returns.set(doCopy(curr.value), doCopy(_target.get(curr.value)));
                    }
                } else {
                    /* otherwise, extends from prototype */
                    const temp = Object.create({});
                    temp.__proto__ = _target.__proto__;
                    /* copy other properties */
                    Object.entries(_target).forEach(([key, value]) => {
                        temp[key] = doCopy(value);
                    });
                    Object.assign(returns, temp);
                }
            }
            return returns;
        }
    };

    return doCopy(target);

}
```

这段代码中最头疼的就是 *recurrsive variables*，一开始想通过 `Map` 实现查重，并用`Object.assign` 一锅端，直接更新栈上所有的 *reference*。但是遇到了我之前所说的原型问题（*prototype*）相反的问题。因为很多自定义（或者官方定义的）都采用了局部变量，并没有统一的方法实现自身的 `clone()` 导致了，外部根本无法复制。（比如 `Set` 和 `Map`）解决方法是从新一条一条添加数据。

但是，如果是用户自定义的私有变量（闭包变量），就完全没有办法解决了。

这里放上一个测试案例：

```js
// define
const a = {
    c: "hello",
    d: function() { console.log(this.c); }
};
a.a = a;

// test
const _a = deepCopy(a);

a.c = "world";
a.d();
// > "world"
_a.d();
// > "hello"
_a.a === _a;
// > true
```

### 小结

至此，对于深拷贝的研究就告一段落了，有其他想法或者改进的朋友可以在下面留言。

对于总结出来的终极解决方法：还是请每个程序猿在实现自己的类的时候，良心写😂，增加一个 `clone()` 方法。

以上
