#### 简介

我和小伙伴们经常要搬运YouTube的视频，大家用的工具各种各样，但都是下载到个人电脑，不方便管理。于是乎想出这个服务器端的下载器。小伙伴们在YouTube上快乐的挑选视频，遇到合适的只需要点一下鼠标就全部下载到中央服务器，过程如下

看到视频时候点击视频右下角的按钮添加到一个playlist

<img src="http://owo5nif4b.bkt.clouddn.com/QQ20180115-171401a.png" width="400">

下载器自动下载

<img src="http://owo5nif4b.bkt.clouddn.com/QQ20180115-164419@2x.png" width="400">

#### 支持环境

Mac 和 Linux

#### 使用方法

1. [下载](https://github.com/zhanghjster/youtube-downloader/releases) youtube-downloader

2. 安装you-get https://you-get.org

3. 设置youtube api的访问权限, [这里](https://developers.google.com/youtube/v3/getting-started)

   * 在[谷歌开发者终端](https://console.developers.google.com/)，先创建一个project然后在这个project里创建一个“oauth api”的credential, 

   <img src="http://owo5nif4b.bkt.clouddn.com/QQ20180115-171638@2x.png">

   * 在credential列表里找到刚创建的credential，点击打开它，然后点"download json"下载验证文件, 保存为"client_secret.json" ，下载程序运行时要用

     <img src="http://owo5nif4b.bkt.clouddn.com/QQ20180115-162651@2x.png">

   * 然后在“Library”里激活 ”youtube data api"

   <img src="http://owo5nif4b.bkt.clouddn.com/QQ20180115-161628@2x.png">

   * 登录youtube账号后，在任何一个视频右下角下面点击“+”的那个按钮创建一个“playlist”

     <img src="http://owo5nif4b.bkt.clouddn.com/QQ20180115-163355@2x.png" width=400>

   * 在youtube账号的的playlist列表里点击刚创建的playlist， 取出url里"&list=xxxxx"后面的“xxxxx"参数，这是playlist的id，程序运行时要用

4. 运行程序

   ~~~shell
   $ youtube-downloader --playlist PL84X8jD5ofVWwYZLCNIPw6mDsrNAkbfbT \
   	--sock-proxy 127.0.0.1:1080 \
   	--secret ~/Documents/client_secret.json \
   	--video-dir /tmp/ \
   	--index-dir ./youtube-download			
   ~~~

   参数说明：

   ~~~
   --playlyst 		为youtube账号里创建的playlist的id
   --sock-proxy 	翻墙使用
   --video-dir   	视频的保存地址, 默认为运行目录下的 ./video/
   --index-dir 	一些索引文件的保存目录，默认为运行目录下的 ./.index/
   --concurrent	并发个数，默认为1
   ~~~

   第一次运行时会弹出需要授权访问playlist的信息, 按照提示将url拷贝到浏览器执行授权

   <img src="http://owo5nif4b.bkt.clouddn.com/QQ20180115-165638@2x.png" width="600">

   在浏览里执行授权后的最后一页，将授权码copy到上面刚执行的命令下面，然后回车

   <img src="http://owo5nif4b.bkt.clouddn.com/QQ20180115-165914@2x.png" width="600">

   下载器开始运行 

   <img src="http://owo5nif4b.bkt.clouddn.com/QQ20180115-170128@2x.png" width="600">

   ​

#### TODO

* 支持多个账号多个playlist同时下载
* 支持youtube访问的页面授权