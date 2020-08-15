---
title: "JavaScript: 使用 ES6 特性实现真正的 Reflect"
tags:
  - JavaScript
categories:
  - Tech
date: 2018-05-21 17:19:39
---

> 启发问题: 源Topic已经移除

### 题记

无聊之中看到这个题目, 问题中有人使用了 `eval` 来解决『获取函数中的变量名』的问题，但如果我们想获取当前 scope 下的变量名呢? 这就变得非常复杂。于是这里就提出一个解决方案，使得对于所有类型都可以在任何地方实现自身的 reflect。

### 解决过程

先贴代码：

```javascript
function extending(c, v) {
	return v.identifier === "" ? v : new class extends c {
		constructor(...props) {super(props); this.identifier = "";}
    }(v);
}

function define(key, value) {
	this[key] = extending(value.constructor, value);
	this[key].identifier = Math.random();
	window._identifiers = {...window._identifiers};
	window._identifiers[this[key].identifier] = key;
	return value;
}

function refetch (value) {
	return window._identifiers[value.identifier];
}
```

实现思路是这样的。

JavaScript 存在 `mutable` 和 `immutable` 的两类变量，为了给变量设立全局可查找的 `identifier`，最高效的做法就是继承原本的数据类型并且增加一个 identifier 属性。这样，我们就可以通过变量独特的 `identifier` 来追回变量的名字，实现反射。

使用方法如下：

```javascript
// var a = 1; 用以下代码替换，也可以在当前 scope 中声明变量 a
> define("a", 1);
< 1
> console.log(a); //我们可以看到 在控制台中打出的 1，数据类型为 Number
< Number{1}
> refetch(a); //返回字符串 "a" 表示当前变量的名字为 "a"
< "a"
```

### 后记

虽然以上的方法没有完美地运用在所有情况中。比如依赖于 `window` 这个全局 `object`，还有会和 `uglify` 插件严重冲突。（在上面的例子中即可体现）。但这确实也给出了另一种，在大多数情况中都能实现的可行方案。

以上
