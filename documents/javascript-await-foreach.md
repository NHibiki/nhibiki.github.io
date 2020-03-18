---
title: "JavaScript: 当 forEach 遇上 await 的坑"
tags:
  - JavaScript
categories:
  - Tech
date: 2018-06-22 22:49:22
---

### 前言

当 JavaScript 进入 ES6+ 时代后，它的语法糖和性能的不凡表现（相比于其他的脚本语言），都使它逐渐站稳了脚跟，成为一个不可替代的流行语言。
然而，不得不承认，ES6+ 毕竟大多还是草案。
许多的函数设计都没有想象中的那么科学，也常常会因为没有仔细考量而出很多 new Feature（bug）。
在我写新的静态博客引擎时，就遇到了很多问题。以下这个问题我最为记忆犹新，它不是 bug。个人觉得偏向设计上的一个缺陷。

### 撞坑

因为在我的 **meta生成器** 中，我使用到了一个非常优秀的语义编译插件 **unified**，它可以让我们非常方便地转换语言规则。
将一种语言设计变化成另一种。如 **markdown** 变为 **rehype** （或者说 html）或是，**txt**，**LaTeX** 等等。
其中，为了不阻塞线程，**unified** 采用 pipe + Promise 的方式，将编译好的文字异步返回。

由于在处理每一篇文章的时候，我都必须遍历每一个 .md 文件。并为他们一一生成文章 meta 后才能进行更后面的操作。
（如，统计，排序，等等）因此，为了避免我的主线程进入异步模式，我决定直接对 **unified** 操作进行 `await` 处理。
代码大概是这样的：

```javascript
/* 遍历所有 .md 文件 */
mkDir.forEach(md => {

    /* 此处省略秀的一手好操作 */

    const rawHTMLContent = await pipeHTML(rawMarkdownContent);

});
```

当然，最后还得把真个 `main` 函数改为 async 以支援 await 的调用。

并用如下的方式启动 `main` 函数：

```javascript
main().then(() => {process.exit(0);});
```

但是，信心满满地执行之后，发现系统报错：

`await should be used in async function`

这我就纳闷了，明明 `main` 已经是 async 函数了呀。想了一会儿，反应过来 forEach 本质是分次执行。
调用 `md => {}` 匿名函数。这个函数并不是 async 于是。将上面第一行改为：

```javascript
/* 遍历所有 .md 文件 */
mkDir.forEach(async md => {

    /* 此处省略秀的两手好操作 */

});
```

这次不报错了，但是，我惊讶地发现 `const rawHTMLContent ...` 这一行包括之后的居然都没被执行！

顺着上次错误的思路分析。forEach 中所有函数都是 async 也就是这个 forEach 本身就变成了一个异步函数 2333。真是然我哭笑不得喵 TAT

为了解决这个方法。最好的方法，直接放弃 forEach 改用它的替代物（虽然丑了一些）：

```javascript
for (md of mkDir) {

    /* 此处省略秀的两手好操作 */

}
```

至此，坑填上。

### 后记

后来还是不死心🤣，就是要用优雅的 forEach 于是，就自己封装一个，专门接收烫山芋：

```javascript
Array.prototype.forEachAwait = async function (fn) {
    for (let i in this) await fn(this[i], i);
}
```

执行的代码或许是这样的：

```javascript
/* 首先你需要一个需要被等待的函数 */
async function getValue() {
    return new Promise(resolve => setTimeout(()=>{
        resolve(Math.random())
    }), 10000);
}

/* 它可以这样被执行 */
await [1,2,3,4,5].forEachAwait(async v => {
    let a = await getValue();
    console.log(v, a);
});
console.log("done");

/* 输出可以为             */
/* 1 0.2001572090840873 */
/* 2 0.9472712652910735 */
/* 3 0.1854800495690847 */
/* 4 0.7421414244802222 */
/* 5 0.0497587313905893 */
/* done                 */

/* 当然你可以这么对比 */
[1,2,3,4,5].forEach(async v => {
    let a = await getValue(); 
    console.log(v, a);
});
console.log("done");

/* 输出可以为             */
/* done                 */
/* 1 0.2036789330780133 */
/* 2 0.0226663336755419 */
/* 3 0.0726829453933766 */
/* 4 0.6238320663197987 */
/* 5 0.9275859417197259 */
```

以上
