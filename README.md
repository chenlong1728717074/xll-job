# xll-job

#### 介绍
基于go语言的轻量级分布式任务调度平台

1.xxl-job的调用严重依赖于数据库,所有业务都是进行扫库

2.PowerJob解决了xxl-job的问题,但是它的用户体系很不友好

二、进度规划

1.xll-job调度采用grpc+pb,因为grpc支持的语言很多,接入起来也很简单

2.采用gin会使项目更加轻量级

3.前期采用Robfig进行cron库,后期打算自己实现实现或者采用cronexpr 

4.前期预计只支持cron进行定时任务调用,后期支持cron和固定速度两种任务模式

5.最后可能会调研rust的可行性,因为rust性能更好,并且rust支持线程+协程的调度方式,可以使系统更加稳定

三、Robfig缺点:

1.Robfig所有的任务管理都是单机管理,对分布式调度十分的不友好

2.Robfig只支持六位数的cron表达式,目前项目下七位数的表达式都会被转换成六位(去除表达式最后一位年份位)

四、cron表达式区别:

1.六位: 秒 分钟 小时 天（月） 月份 星期几

2.七位: 秒 分钟 小时 天（月） 月份 星期几 年份


五.下面是一些go的cron库

https://github.com/reugn/go-quartz

https://github.com/gorhill/cronexpr

#### 软件架构
软件架构说明
gin+grpc

#### 安装教程

1.  git clone https://gitee.com/a-little-dragon/xll-job
2.  cd xll-job
3.  go build
4.  xll-job.exe

#### 使用说明

1.  xxxx
2.  xxxx
3.  xxxx

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  Gitee 官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解 Gitee 上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
5.  Gitee 官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
