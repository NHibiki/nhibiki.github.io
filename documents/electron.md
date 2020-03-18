---
title: "Node.js: Electron和坑"
tags:
  - Electron
  - JavaScript
categories:
  - Tech
date: 2016-12-10 17:10:22
---

### 关于Electron

我想我大概不用仔细介绍什么是Electron了。简单的说，Electron可以用Html，CSS，Javascript制作Cross Platform的Native APP。它使用Chromium和Nodejs作为解释器。

食用Electron之前，请保证自己的电脑

 - 连接到了国际广域网
 
 - 拥有至少1G的存储空间
 
 - 不是拖拉机级别的处理器和内存
 
 - 至少1秒钟之内不会蓝屏
 
 - 装有Node.js并且没有把Terminal/bash/cmd之类的删掉
 
 - 其他...

<!--More-->

### 首先是安装

没装Nodejs的请自行安装....

最简单的方式是通过 `npm install -g electron` 来解决。

以下工具一起实用最佳：（以ヒビキ使用的Mac OS X为准）

 - `npm install -g electron-packager` -- 一个封装工具（可以把Chromium引擎和Nodejs与应用包装起来）
 
 - `wget https://dl.winehq.org/wine-builds/macosx/i686/winehq-staging-1.9.23.pkg` -- For Darwin(OS X) Only，为了打包 win32 必须用到wine. 对于GNU/Liunx，请自行下载或者编译wine.
 
 ### 创建最简单的应用
 
 #### 编辑命令
 
 通过 `npm init` 新建新的 manifest.json 文件。打开并在script后添加如下命令，修改后效果如下 - 
 
 ```javascript
 "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1",
    "start": "electron .", //如果是不是-g方式安装，请改成 "start": "./node_modules/.bin/electron .",
    "packager": "electron-packager ./ MyApp --out ./OutApp --overwrite --all --icon=./ico.icns"
  }，
 ```
 
 这就些命令可以通过 `npm run-script` 执行
 
 对于 `"packager"` 的特殊说明：
 
 `./` Workdir
 
 `MyApp` 为显示的应用名
 
 `--out ./OutApp` 为输出目录  
 
 `--overwrite` 表示覆盖旧版本
 
 `--all` 为生成的 platform，all为全平台，当然你可以把它替换为 `--platform=darwin --arch=x64` 来限定平台（这个是x64的Max OS X）
 
 `--icon` 包图标
 
> 特殊说明：
> Max OS X 必须使用 icns 的专用图标格式， 其他的请使用 ico 格式。

#### 入口脚本

```javascript
"use strict";
const electron = require('electron');
const {app} = electron;
const {BrowserWindow} = electron;

let win; // 这样可以保证当次脚本推出之后回收窗口

function createWindow() {
  win = new BrowserWindow({width: 600,
                           height: 400,
                           resizable: false,
                           maximizable: false,
                           alwaysOnTop: true,
                           fullscreenable: false,
                           title: "喵"
                         });
  win.loadURL(`file://${__dirname}/view/index.html`);
  //win.webContents.openDevTools();
  app.setName("喵");
  win.on('closed', () => {
    win = null;
  });
}

app.on('ready', createWindow);
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  if (win === null) {
    createWindow();
  }
});
```

然后，把你之前写的网页应用复制到 `./view` 下，保证入口是 `index.html` 就好啦！

### Electron的常用函数和坑

#### win 创建选项

- 'width'／'height' 宽度／高度
- 'title' 窗口标题
- 'alwaysOnTop' 是否永久置顶
- 'maximizable' 是否可最大化
- 'resizable' 是否可拉伸
- 'fullscreen'／'fullscreenable' 全屏显示／是否可以全屏（后者用于Mac OS）
- 'icon' 图标，<div style="color:red">（注意！！选项，只在Windows和Linux下生效，而且，它是窗口图标，不是应用图标QwQ）</div>
- 'show' 是否显示（默认是True）
- 'closable' 是否可以关闭2333（对Linux无效）
- 'movable' 是否可以移动（对Linux无效）
- 'x'／'y' 初始化的窗口位置（默认是屏幕正当中）

#### win.on Listener

 - 'resize' 窗口大小改变
 - 'move' 窗口移动
 - 'show' / 'hide' / 'focus' 窗口被显示／隐藏／鼠标或者Tab切换选中
 - 'closed' 窗口关闭
 - 'responsive' / 'unresponsive' 页面是否响应
 - 还有很多参考[官方API](http://electron.atom.io/docs/api/browser-window/)

#### 坑·javascript环境切换

如果你发现你的非Nodejs部分的脚本无法运行那么请参考以下方法：

因为javascript是通过Nodejs执行，所以，一定要在 `<script></script>` Tag 前后加入如下内容：

```html
    <script>if (typeof module === 'object') {window.module = module; module = undefined;}</script>
    <script src="........"></script>
    <script>if (window.module) module = window.module;</script>
```

因为在Nodejs执行页面Javascript时，极少数的JS会判断运行环境，所以我们要在运行页面脚本之前建立一个“虚拟”的页面环境，并在结束后还原Nodejs环境（神坑）

#### 参考资料

 - http://electron.atom.io/docs/api/browser-window/
 - http://blog.csdn.net/sinat_25127047/article/details/51418682
 - http://stackoverflow.com/questions/31529772/how-to-set-app-icon-for-electron-atom-shell-app

以上 ～