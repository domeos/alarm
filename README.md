domeos/alarm
============

## Notice

domeos/alarm模块是以open-falcon原生alarm模块为基础，为适应DomeOS监控报警需求而设计修改的，包名已修改为`github.com/domeos/alarm`

原生open-falcon系统中，judge通过hbs读取portal中设置的报警策略，判断报警事件，把报警event写入redis。alarm从redis读取event，做相应处理，
可能是发报警短信、邮件，也可能是callback某个http地址。生成的短信、邮件再次写入redis，sender专门负责来发送。

在DomeOS中，不再单独部署portal模块，报警策略信息由DomeOS报警模块提供并写入portal库，供hbs读取。alarm从redis中获取报警event后将保存至DomeOS数据库，
以提供DomeOS前端展示报警信息。主要修改如下：

- 所有关于version的操作：提供falcon和domeos两种结果，对应修改了control脚本中pack部分获取版本的命令
- alarm的http页面去除config选项
- SafeEvents的put和delete中加入对DomeOS数据库alarm_event_info_draft表的对应操作(添加、更新、删除)
- DomeOS中选择忽略报警时调用alarm的/event/solve接口，同时更新alarm内存与DomeOS数据库
- 报警的consume过程，获取Action由调用portal的/api/action/<action_id>接口改为调用DomeOS的/api/alarm/action/wrap/<actionId>接口
- 报警的consume过程，获取用户组对应的用户信息由调用uic的/team/users接口改为调用DomeOS的/api/alarm/group/users/wrap接口
- 无需单独安装link组件，改为将link信息存储在DomeOS数据库的alarm_link_info表中，短信报警时显示为一个链接
- 存储link接口/api/alarm/link/store，获取link接口/api/alarm/link/{linkId}

## Installation

```bash
# set $GOPATH and $GOROOT
mkdir -p $GOPATH/src/github.com/domeos
cd $GOPATH/src/github.com/domeos
git clone https://github.com/domeos/alarm.git
cd alarm
go get ./...
./control build
./control start
```

## Configuration

- database: DomeOS数据库地址，需提供用户名、密码、地址与对应端口
- maxIdle: MySQL连接池最大空闲连接数
- http: 监听的http端口
- queue: 要发送的短信、邮件写入的队列，需要与sender配置一致
- redis: highQueues和lowQueues区别是是否做报警合并，默认配置是P0/P1不合并，收到之后直接发出；>=P2做报警合并
- api: 配置DomeOS的Server地址

## Run In Docker Container

首先构建domeos/alarm镜像：

```bash
sudo docker build -t="domeos/alarm:latest" ./docker/
```

启动docker容器：
```bash
sudo docker run -d --restart=always \
    -p <_alarm_http_port>:9912 \
    -e DATABASE="\"<_domeos_db_user>:<_domeos_db_passwd>@tcp(<_domeos_db_addr>)/domeos?loc=Local&parseTime=true\"" \
    -e REDIS_ADDR="\"<_redis>\"" \
    -e API_DOMEOS="\"<_domeos_server>\"" \
    --name alarm \
    pub.domeos.org/domeos/alarm:1.0
```

参数说明：

- _alarm_http_port: alarm服务http端口，主要用于状态检测、调试等。
- _domeos_db_user: DomeOS中MySQL数据库的用户名。
- _domeos_db_passwd: DomeOS中MySQL数据库的密码。
- _domeos_db_addr: DomeOS中MySQL数据库的地址，格式为IP:Port。
- _redis: 用于报警的redis服务地址，格式为IP:Port。
- _domeos_server: DomeOS的server地址。

样例：

```bash
sudo docker run -d --restart=always \
    -p 9912:9912 \
    -e DATABASE="\"root:root@tcp(10.16.42.199:3306)/domeos?loc=Local&parseTime=true\"" \
    -e REDIS_ADDR="\"10.16.42.199:6379\"" \
    -e API_DOMEOS="\"http://domeos.example.com\"" \
    --name alarm \
    pub.domeos.org/domeos/alarm:1.0
```

验证：

通过curl -s localhost:<_alarm_http_port>/health命令查看运行状态，若运行正常将返回ok。

DomeOS仓库中domeos/alarm对应版本：pub.domeos.org/domeos/alarm:1.0
