---
title: "AudioStream: 用 Python 建立自己的音乐流"
tags:
  - AudioStream
  - Music
categories:
  - Tech
date: 2017-08-02 14:51:32
---

> 我发现这几天真的高产如母猪了（ TAT
> 好的，作为昨天研究的后续，今天碰到的问题很少，但是钉子很多。

<!-- More -->

### 准备

昨天全部都是按照官方文档来的，因此也没有贴 Works Cited，今天侧重于“踩坑”，就会有很多 References。

#### 安装 python-shout

`PythonShout`是一个 ICE 自己的 python 库，用于方便我们用 python 代替 ICES 进行歌曲流推送。安装方式很简单。

```bash
pip install python-shout
```

它的前置库在装完 ICECAST 之后都会有，因此就没多大问题了。

参考了以下的代码<a href="#sup1"><sup>[1]</sup></a>，按照昨天的配置修改参数后，`python start.py`启动流推送器。

```python
hostname ="localhost"
port= 8000
password = "??????!!!!!!"

import shout
import sys
import threading
from glob import glob
from random import shuffle,choice

class RunStream (threading.Thread):
   def __init__ (self, channel_mount, music_directory, station_url, genre,name,         description, bitrate="128", samplerate="44100", channels="5",music_format="mp3", ogv=0):
    #connection to icecast
    global hostname,port,password
    self.song_conter= 0
    self.s = shout.Shout()
    self.s.audio_info = {shout.SHOUT_AI_BITRATE:bitrate,   shout.SHOUT_AI_SAMPLERATE:samplerate, shout.SHOUT_AI_CHANNELS:channels}
    self.s.name = name
    self.s.url = station_url
    self.s.mount = channel_mount
    self.s.port = port
    self.ogv = ogv
    self.s.password = password
    self.s.genre = genre
    self.music_directory = music_directory
    self.s.description = description
    self.s.host = hostname
    self.s.format = music_format #using mp3 but it can also be ogg vorbis
    print self.s.open()
    threading.Thread.__init__ (self)

   #checking directories for files to stream
   def scan_directories(self):
      self.files_array = glob(self.music_directory+"/*.[mM][Pp]3")  + glob(self.music_directory+"/*/*.[mM][Pp]3") + glob(self.music_directory+"/*/*/*.[mM][Pp]3")   #checks the specified directory down to the third depth
      print str(len(self.files_array))+" files" #display number of matching files found
      shuffle(self.files_array) # randomize playlist

   def run (self):
      while 1: #infinity
        self.scan_directories() # rescan dir, maybe in time you add some new songs
    self.song_counter = 0   
    for e in self.files_array:
           self.write_future()
           self.sendfile(e)
           self.song_counter = self.song_counter + 1

   def format_songname(self,song): # format song name - on filename (strip "mp3", change _ to " ". Formatting name of song for writing into a text file
      result = song.split("/")[-1].split(".")
      result = ".".join(result[:len(result)-1]).replace("_"," ").replace("-"," - ")
  return result

   def write_future(self): #write playlist
      filename = self.s.mount.replace("/","")+"-current.txt"
      fa = open(filename,"w")
      aid = self.song_counter
      pos = 7 # CHANGE if you want more songs in future playlist
      for s in self.files_array[aid:]:
         fa.write(self.format_songname(s)+"\n")
         pos = pos - 1
         if (pos==0):
            break
      if (pos>0):
         for s in self.files_array[:pos+1]:
            fa.write(self.format_songname(s)+"\n")
      fa.close()   

   def sendfile(self,fa):
      print "opening file %s" % fa
      f = open(fa)
      self.s.set_metadata({'song': self.format_songname(fa)})
      buf = f.read(1024)
      while buf:
         self.s.send(buf)
         buf = f.read(1024)
         self.s.sync()
      f.close()

#running the first stream
RunStream(channel_mount = "/nekostream", music_directory = "/home/neko/nekoneko/music", station_url = "http://yuuno.cc", genre = "Neko",name = "NekoNya", description = "Nya Nya Nya").start()
```

> 真的受不了原作者的奇葩缩进，真的是很奇葩了（3格缩进）。<a href="#sup1"><sup>[1]</sup></a>

但是，居然收到了一个 `Segment Fault` 错误，经过检查，发现是在 `self.s.send(buf)` 这一行出错了。推送器连接成功，但是无法发送 binary 数据。

出门右转 Google，于是看到了这篇文章<a href="#sup2"><sup>[2]</sup></a>。

它说到 ->

> The issue was indeed a compile problem of libshout as cox pointed out. But it only worked in debian 7 and not ubuntu 12. I think the reason why is because I did not install libogg in ubuntu I only installed vorbis which I thought was the same thing. I also installed mp3 codecs just in case.

大致意思就是，这个是官方库的问题，没有对 RedHat Linux 进行支持？

不过也罢，python-shout最后的更新时间是 2012 年。。（在这个音乐已经小得不需要用流传输的时代，连 ICE 都抛弃了audio stream

不知道拐了几个右转弯，总算找到了一个有心人，在 GitHub 上 po 了自己修改后的代码。看了修改之后才知道，之前的问题，只是少了 `PY_SSIZE_T_CLEAN` 这个 define 导致C库没和Python对接上（有点怀疑 ICE 是不是解散了 逃）而且，这个 GitCommit 就没有人访问过，最后的修改时间是 2年前 <a href="#sup3"><sup>[3]</sup></a>。

clone 下来之后安装就OK了

```
python setup.py install
```

### Docker版本

大小约为 200MB 出头

`docker attach {ID}` 查看。里面大概是这个样子的 ～ 

```bash
...
[System] Start Playing /home/ice/music/鴇沢直 - BRYNHILDR IN THE DARKNESS -EJECTRO Extended-.mp3 ...
[System] Sync 0.0305725362633 ...   
```

表示了正在推送的歌曲和推送进度。

以上

### Works Cited

<sup id="sup1">[1]</sup> http://www.it1352.com/306021.html
<sup id="sup2">[2]</sup> https://stackoverflow.com/questions/27654208/shout-python-segmentation-fault-how-can-i-fix-this
<sup id="sup3">[3]</sup> https://github.com/fergalmoran/python-shout