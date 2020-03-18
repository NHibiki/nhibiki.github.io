---
title: 'Rust: The first attempt to ownership'
date: 2019-12-18T04:31:28.959Z
categories:
  - Notes
tags:
  - Rust
  - Notes
---
> 上一篇 [Rust: Brief Intro to SAFE](https://yuuno.cc/articles/rust-brief-intro-to-safe)

### 前言

这篇文章是继第一篇初探 Rust 之后的续集。在之前那篇博客中，我简介了我对 Rust 的看法，并且定下了我的学习里程碑计划。今天，我会在这篇文章讲述对 Rust 的 LinkedList 系列的实现，以及简介它的部分 Syntax。

首先，为什么要学习用 Rust 实现 LinkedList？

我个人的理解是，Rust 之所以与其他编程语言不同，关键在于它的 *“所有权“* 和 *“生命周期“*。换言之，在编写 Rust 的过程中，对于自己创建的每一个变量，包括临时变量，都要对它的作用范围了如指掌。因此，我们可以找一个相对复杂的案例（譬如 LinkedList）来加深自己对这两个关键词的理解。但或许你会说，LinkedList 又不是复杂环境！不如直接上手 Map（key 与 value 甚至都可以是不同生命周期的案例）。当然！你当然可以！不过 LinkedList 对于初学者而言，是一个 “我们都清楚它的原理，但不清楚它的生命周期“ 的最简单案例。中间并没有复杂的逻辑，我们只需要一心扑在我们想学的知识点即可。这样的做法，更直接，也更容易有收获。（毕竟这个 LinkedList 花了我整整半天的时间去 debug）

### 初版

> 实现 LinkedList Implement 的Stack

这是一个反例，我后悔自己没有仔细规划就开始动手码代码。以至于即便是一个 Stack 最后都变成了 DOP (*debug oriented programming*)。

首先我构建了两个简单的 struct：

```rust
 pub struct Node<T> {
    pub next: Option<Box<Self>>,
    pub value: T,
}
pub struct Stack<T> {
    pub head: Node<T>,
    pub length: usize,
}
```

这个结构再简单不过，Node能**拥有**自己的值，也拥有一个*可空*的指针指向下一个节点。

然后，我们只需要把它的下一个节点新建并放入之间的那个节点中，就能把这些节点连起来了：

```rust
impl<T> Stack<T> {
    fn push(&mut self, v: T) -> &mut Self {
        self.length += 1;
        let mut node = Node::new(v);
        node.next = self.head.next;
        self.head.next = Some(Box::new(node));
        self
    }
}
```

很简洁，很直观，甚至支持链式调用，但是也报了很多错。

编译器这时候不乐意了：虽然你说的都对，但是我没办法帮你复制`*self`。

```rust
node.next = self.head.next;
|
move occurs because `self.head.next` has type `std::option::Option<std::boxed::Box<Node<T>>>`, which does not implement the `Copy` trait
```

原因是，`=` 在 Rust 中有两层语义，一个是赋值，一个是所有权转移或者借用，由于 `Option<Box<Node<T>>>` 并不是基类，从它赋值的时候，Rust 需要知道如何*拷贝*。（譬如我们知道，在 python 和 javascript 中，这个等号代表了新建了一个 reference），而 Rust 并不知道你的 struct 中哪些需要被转移所有权，哪些需要借用，又或者哪些需要直接拷贝（u8, i32这些）。

编译器原本希望我能把它 Copy 出去，但是并不是所有类都能实现 Copy trait（譬如 String），所以，这里最直接的解决方法告诉编译器，我不需要复制任何变量，你把它所有权收走就行 `node.next = self.head.next.take();`，于是，在此之后，这个 `next` 的所有者就从 `self.head` 转移到了 `node` 上。

最后，将它封装起来的成品：

```rust
pub struct Node<T> {
    pub next: Option<Box<Self>>,
    pub value: Option<T>,
}

impl<T> Node<T> {
    fn new(v: T) -> Self {
        Node {
            next: None,
            value: Some(v),
        }
    }
    fn new_empty() -> Self {
        Node {
            next: None,
            value: None,
        }
    }
}

pub struct Stack<T> {
    pub head: Node<T>,
    pub length: usize,
}
impl<T> Stack<T> {
    fn new() -> Self {
        Stack {
            head: Node::new_empty(),
            length: 0,
        }
    }
    fn push(&mut self, v: T) -> &mut Self {
        self.length += 1;
        let mut node = Node::new(v);
        node.next = self.head.next.take();
        self.head.next = Some(Box::new(node));
        self
    }
    fn pop(&mut self) -> Option<T> {
        self.length -= 1;
        let node = self.head.next.take();
        if let Some(pk) = node {
            self.head.next = pk.next;
            pk.value
        } else {
            None
        }
    }
    fn index(&mut self, idx: usize) -> Option<&T> {
        let mut k = &self.head;
        let mut i = idx;
        while let Some(pk) = &k.next {
            k = &pk;
            if i == 0 {
                if let Some(v) = &k.value {
                    return Some(&v);
                } else {
                    return None;
                }
            };
            i -= 1;
        }
        None
    }
    fn for_each<U>(&mut self, p: U) where U: Fn(&T) {
        let mut k = &self.head;
        while let Some(pk) = &k.next {
            k = &pk;
            if let Some(v) = &k.value {
                p(&v);
            } else {
                return;
            }
        }
    }
    fn reverse(&mut self) -> Self {
        let mut new = Self::new();
        while self.length > 0 {
            new.push(self.pop().unwrap());
        }
        new
    }
}
```

以及用于测试的案例：

```rust
fn main() {
    let mut a: Stack<&str> = Stack::new();
    a.push("Hi,")
     .push("this")
     .push("is")
     .push("fancy");
    let b = a.pop();
    a.push("a")
     .push("trial!");
    a.for_each(|v| {
        print!("{} ", v);
    });
    println!("\nGot b={} from a[{}]", b.unwrap(), &a.length);
    println!("Got a[3]={}", a.index(3).unwrap_or(&"undefined"));
    println!("Got a[5]={}", a.index(5).unwrap_or(&"undefined"));
    let mut a = a.reverse();
    println!("\nReversed:");
    a.for_each(|v| {
        print!("{} ", v);
    });
}
```

输出结果：
```bash
[Running] cd "/User/nyu/dev/linked-test/src/" && rustc linked.rs && "/User/nyu/dev/linked-test/src/"linked
trial! a is this Hi, 
Got b=fancy from a[5]
Got a[3]=this
Got a[5]=undefined

Reversed:
Hi, this is a trial! 
[Done] exited with code=0 in 0.472 seconds
```

### 问题

如果你在读代码的时候感觉这个代码有问题，那恭喜你！虽然它运行起来没问题，但其实它有很大的又不重要的问题 - -：

- **对值的绝对控制**：这个案例中，Stack 默认索取了被 push 的内容的所有权和生命周期。当然，我们确实可以新建一个引用来解决问题（如上面的案例 `Stack<&str>`），但总是感觉没那么美观。（并且，如果我们希望用 Box 来承接引用，我们就需要手动规定被承接内容的生命周期（`pub next: Option<Box<&'a Self>>`）。）
- **很难实现的 Iterator**：Rust 官方支持 `Iterator` trait 和 `IntoIterator` trait，前者可以使 struct 具有被 `for` 遍历的功能，并且支持 `.iter()` 和 `.iter_mut()`，后者允许从一个类中拓展出一个支持遍历的子类[（参考资料）](https://doc.rust-lang.org/std/iter/trait.IntoIterator.html)。不过，在这个案例中，由于 `next` 都是实实在在存在实体而不是一个 reference，导致构架 Iterator 难度急剧上升（当然可能是我没想到好的方法）。在我自己的实验中，可以拓展 `stack.for_each`，并使用试验性的 Generator 来达成目的，不过毕竟 Generator 还是实验性的，实战中并无法使用。

### 总结

真正的编码过程其实并没有上面展示的这么容易，中间绕了很多弯路，不过大部分的错误都是类似的：

```rust
|
move occurs because *** has type ***, which does not implement the `Copy` trait
```

以及，生命不够长的问题（这么说好像也没什么不对 - -）。

不过，综合来说，还是对自己设定的变量的掌控力不足。

这款语言或许不可能成为世界上最流行的语言，但它对于计算机安全以及可持续化的冲击力是不可估量的。
