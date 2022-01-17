# mapboxspider

## 简介

通过 golang 脚本，快速下载`mapbox账号`配置的 mapbox 的图源，包括：影像、地名、矢量、高程等。

## 方法

1. 申请`mapbox`账号，获取一个token
2. 参考`config-example.json`，复制一份`config.json`文件，替换相关内容
3. 脚本启动下载
4. 下载完之后，可以用`mbutil`这个工具，将数据转换为`mbtiles`数据库

## 注意
1. 阅读代码，逻辑很简单
2. 一些参数可调，比如层级、CPU等等

## 致谢

1. [colly](http://go-colly.org/)——用于构建 Web 爬虫的 Golang 框架

