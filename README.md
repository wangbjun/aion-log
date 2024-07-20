# 永恒之塔网游战斗日志分析系统。

基于永恒之塔客户端日志Chat.log来分析玩家信息，比如种族、职业、伤害等数据。

测试地址：https://d1i452fxut2vqb.cloudfront.net

这个项目后端是golang,主要是负责解析数据入库、提供接口，前端是antd框架，主要是几个页面数据的展示。

纯属个人娱乐，如果有人想用的话，我简单说一下， 前提是你是懂Web开发的，小白就算了，我这个不是面向小白的。

# 数据库
数据库使用了sqlLite,无需安装，开箱即用，默认位于项目目录下aion.db

# 命令
这是一个命令行应用，主要包含几个命令
```
Aion Chatlog Analyse System

Usage:
  aion [command] [flags]
  aion [command]

Available Commands:
  class       Classify Player Info
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  httpServer  Start A Http Server
  parse       Parse Chatlog file
  rank        Rank Player Info

Flags:
  -c, --conf string   config file (default "app.ini")
  -h, --help          help for aion

Use "aion [command] --help" for more information about a command.
```
主要是parse，用于解析日志文件

其次就是httpServer,用于启动http接口服务

假设你的日志文件位于/data/chat.log，那么你的命令就是
```
go run main.go parse -f /data/chat.log
```
# web页面
前端文件位于文件夹**frontend**，需要先安装依赖
```
npm install && npm run dev
```
如果你要部署，那就```npm run build```，懂前端开发的人自然知道我在说什么。。。不懂我也不好再解释！