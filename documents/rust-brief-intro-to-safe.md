---
title: 'Rust: Brief Intro to SAFE'
date: 2019-12-15T01:54:06.136Z
categories:
  - Notes
tags:
  - Rust
  - Notes
---
### 简介

首先，这篇文章仅仅是我在学习中的一些见解，若有错误，望在评论区指正。

学习 Rust 的起因仅仅是因为希望在明年最后一学期的 *CS Capstone* 中运用它强大的 *WASM* 生态以构建一个跨平台（包括浏览器端）的虚拟机。虽说除去 Rust 还有很多选择，比如 Golang，C++。不过，这些选项在漫长的选型过程中逐渐被抛弃：

 - **C++:** 蜕变后的 C++ 确实是一个很好的选项，不过它的轮子很杂乱（我们不可能用4个月时间去研究轮子）。换句话说它的社区没有想象中的这么友好。
 - **Golang:** 这个可以说是我的本命语言，简明的设计和强大的协程成为我的开发首选。但是，如果要提到 *WASM*，它就是心有余而力不足了。使用标准 GO 1.13 编译出来的 *WASM* 体积庞大，但若使用 TinyGO，又生态欠佳，于是又平添许多不确定性（况且合作的组员似乎不会Go）。

因此，种种因素使我们打算挑战一下 Rust 的学习曲线，心中对 Rust 好奇的迷雾也随着知识的积累，渐渐揭开。

### 入“门”

过去，如果要说出 *“我可以两天入门这款语言”* 这种大话，我还是有自信去圆上它的。而当我打开那本《圣经》，我就明白，Rust 确实是一个从气势上就可以吓退很多人的语言。

不得不承认，我敬佩它的设计者，Rust 开创了一个 “安全” 的尝试。概念上而言，它自己定义了一套安全理论，而语义而言，它又以包容的态度，囊括了主流的代码语义理论。接下来，我就以一个菜鸟的角度一一解释一下：

#### 设计理论

安全、快速，是 Rust 的创立之本。如果要用一个形象的描述，那就是 “无 ptr”，“无 GC” 的二无政策。

何谓安全？对于编译语言来说，让用户有机会绕过二进制文件掌控范围的语言，即为不安全。换言之，你查看或修改了不属于你的东西，就是不安全。对于那些运行时语言（python，node）来说，只要不是解释器的漏洞，它们一般是安全的。因为它们一般情况下不能操控指针，它们的作用域始终是在可控范围内的（当然，一些通过 binding 运行二进制文件的库是被排除在这个讨论外的）。对于 *ptr* 这一不可控的东西，Rust 有独到的见解，它使用 “智能指针” Box，Ref 等来掌控被指向内存的控制权和作用域，从某种程度上减少了不可控的范围。对于这部分的内容，由于我没有细看，也不闭着眼睛说话，就点到为止。

何谓快速？如果我们不管垃圾，做出所有针对硬件的编译优化，那就已经很快了。但是一个稳定的语言，不能不管垃圾。所以，Rust 提出了全新的 *“作用域”* 的解读：一旦一般的变量所在的作用域结束了，它的生命周期也会被终止，除非它的所有权被转移到了未被终止的作用域中。这样，就不需要总是检查（或是无法检查造成memory leak）内存，而完成对内存中垃圾的完美回收♻️。

#### 语义理论

现在新的语言都有自己特殊的语义，比如 Go，语义鲜明，你看一眼就知道这就是 Go。而 Rust，不知道是为了方便用户还是其他，它在语义的定义上比较随性。换言之，你在写 Rust 的时候，有 Java 的影子，有 Python 的影子，Swift 的影子，甚至还有 Haskell 的影子。

当然，天下语义一大抄，只要编码者习惯即可，这不是什么大问题。作为一个从 Go 过来的人来说，一开始看是有点不习惯，不过，当经营大项目的时候，终究还是以注解为主，谁会在团队协作的时候写一写精简过度晦涩难懂的语法来刁难人呢？

### 上手

看了 3 天的《圣经》，我觉得是时候在我忘记之前学的内容之前，先写一写小 demo 巩固一下了。于是，给自己出了几道题：

1. 实现一个简易 html DOM 生成器
2. 实现一个 linked list
3. 实现一个可持久化（disk io）的 hashmap

分三个阶段完成，第一个在阅读完 ownership 后，第二个在阅读完 trait 和 box 之后，第三个在阅读完 fearless concurrency 和 multithread/lock 之后。

于是乎就有了第一题：

我的想法其实很基础，每一个 DOM struct 都可以有生成一个 immutable 的 child，而它的 scope 就继承了这个 DOM 的 scope，虽然这个逻辑会根据 DOM Tree 的复杂程度构建许多个 DOM 元素，不过它完美的解决了管理权和拥有权的问题，每个 DOM 只需要管辖自己所属的子元素。

代码如下：

```rust
pub struct DOM {
    pub indent: String,
    pub counter: usize,
}

impl DOM {
    fn new(&self) -> DOM {
        DOM {
            indent: self.indent.clone(),
            counter: self.counter + 1,
        }
    }
    fn print<T>(&self, tag: String, props: Option<Vec<(String, String)>>, next: T) where T: Fn(DOM) {
        let mut prop_str = String::new();
        if let Some(_props) = props {
            for (key, value) in _props {
                prop_str.push_str(&*format!(" {}=\"{}\"", key, value));
            }
        }
        println!("{}<{}{}>", self.indent.repeat(self.counter), tag, prop_str);
        next(self.new());
        println!("{}</{}>", self.indent.repeat(self.counter), tag);
    }
    fn raw(&self, input: String) {
        println!("{}{}", self.indent.repeat(self.counter), input);
    }
}

fn none(_: DOM) {}

fn main() {
    let d = DOM{
        indent: String::from("\t"),
        counter: 0,
    };
    d.print(String::from("html"), None, |d| {
        d.print(String::from("head"), None, none);
        d.print(String::from("body"), None, |d| {
            d.print(String::from("h1"), Some(vec![(String::from("style"), String::from("margin: 0 auto;"))]), |d| {
                d.raw(String::from("Hello, World!"))
            });
            d.print(String::from("script"), Some(vec![(String::from("type"), String::from("javascript"))]), |d| {
                d.raw(format!("window.alert(\"{}\");", "Hello, World"))
            });
        });
    })
}
```

生成结果：

```html
<html>
	<head>
	</head>
	<body>
		<h1 style="margin: 0 auto;">
			Hello, World!
		</h1>
		<script type="javascript">
			window.alert("Hello, World");
		</script>
	</body>
</html>
```
