---
title: "Git: 使用大文件版本管理系統 git-lfs"
tags:
  - Git
categories:
  - Tech
date: 2017-12-18 19:07:10
---

### 簡介

首先先介紹一下什麼是 git-lfs. 它的全稱是 Git Large File Storage. 用一句話翻譯它的作用, 就是 git 對大文件或者 binary 文件的解決方案.

這是 git-lfs 的官方鏈接: https://git-lfs.github.com/

 其中有 git 對 git-lfs 更加詳細的描述和解釋, 那我就不再多解釋這個工具啦 ~\(≧▽≦)/~, 主要講一下安裝方法和使用體驗 (與坑).
 
 ### 安裝
 
 如果不怕麻煩的話, 可以直接下載官網的 binary 文件, [官網](https://git-lfs.github.com/) 會有一個大大的 Download. ( 別告訴我你看不到 QAQ
 
 或者你可也以用 brew 等工具裝
 
 ```bash
brew install git-lfs
 ```

再或者下載 [源碼](https://codeload.github.com/github/git-lfs/zip/master) 自行編譯, GitHub 託管地址是: https://github.com/git-lfs/git-lfs

安裝完成之後, 需要把它添加到 git 工具中.

```bash
git lfs install
```

如果只想對某幾個特定的 repo 開啟 lfs, 可以在項目目錄添加 --local 來局部安裝
```bash
git lfs install --local
```

至此, git-lfs 安裝完成

### 使用方法

 git-lfs 的大致原理是在 commit 和 push 之前添加鉤子, 匹配提交的文件是否符合 lfs 的要求, 如果符合, 就從原本的 commit 和 push 中移除, 添加到 lfs 的上傳列表, 并用一個配置文件在原 repo 的位置中替換掉原文件.

所以, 首先我們要建立一個 git-lfs 匹配規則, 命令為 `git lfs track "SOMETHING"`. 比如, 我希望把所有項目中的 `wav` 文件和項目中 `/map/` 目錄下的所有文件添加到 lfs 而非 git repo, 我們只需要如下命令:

```bash
git lfs track "*.wav"
git lfs track "/map/"
```

它的匹配規則和 `.gitignore` 類似, 添加完之後, 我們可以看到項目根目錄生成了 `.gitattributes` 文件. 這就是我們的 記錄文件. 除此以外, 我們還可以通過命令 `git lfs track` 查看目前添加的所有規則.

之後, 我們只需要常規地 `add`, `commit`, `push` 就行了. 被 track 的文件會自動被 git push 到單獨的 storage server.

### 坑

emmmm, 因為 GitHub 很摳門, 連嘗試的機會都不給我們, 明確不額外付費就 1KB 都不給, 吾只能用 bitbucket 來測試 git-lfs. 結果就遇到了坑.

首先, 如果 project 啟用了 lfs, 那麼就不能用傳統的方法 clone, 而需要  `git lfs clone ...`  才行, 不然, 我們會發現只能拉下來用於替換原本的大文件的那些配置文件.

除此之外, 單純的刪除那些被 commit 進入 lfs 的文件, 並不能將它從 remote repo 中刪除. (我不知道這是 bug 還是 emmm `New Feature`). 總之, 我只能單獨進入 bitbucket 的 lfs 管理界面才能刪除那些文件.

### 總結

總之, git 確實提出了一個大文件, 二進制文件嵌入 git project 的解決方案. 不過可以看出, 各大 git 版本管理平台對它的支持還不是特別到位. (如果之前那個 `bug` 真的是 `New Feature` , 那就說明它確實是非常不成熟了).

但是, 從它的效果來說, 也確實可以減少 git 倉庫對檢查大文件的壓力, 並且能夠將它從歷史版本中移除, 通過簡單的 sha256 code 就能讓 git 管理大文件的版本問題.

以上
