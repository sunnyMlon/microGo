## 微服务网关项目

### 约法三章

1. 学习笔记以markdown的方式记录
    - [markdown语法](https://www.jianshu.com/p/191d1e21f7ed)
    - [markdown在线](https://tool.lu/markdown)
    
2. 代码部分日更到github上，包含学习笔记及代码
    - [个人项目入口](https://github.com/sunnyMlon/microGo)
    
3. 复杂部分，需要以流程图 | UML图 | 思维导图方式，多维展示
    - ![go学习路线思维导图](https://i.loli.net/2021/07/17/YOS4xJUKylajzVd.jpg)
    - [在线传图地址](https://sm.ms/)
    
 
## 1_2021717_microGo学习笔记

### 一、功能点预览
1. doshboard
    - 大盘指标
    - qps 访问量 
    - 曲线图 饼状图
2. 服务列表
    - http服务 tcp服务 grpc服务
    - ip白名单 ip黑名单 限流  轮询方式()
3. 租户列表
    -  租户增删改查
    -  行为统计
    
### 二、学习意义
1. 亿级流量的网站架构中，网关是一个核心，包括Es、redis的中台服务需要跨业务线的流量统计等
2. 后端开发者，网关是业务搬砖和后台架构的分水岭
3. golang借助本身的高并发优势，能够更好的用于企业的保性能架构

### 三、正课学习
+ 2.1 为何学习网络基础
    - OSI七层网络协议 应用层 表示层  会话层 传输层 网络层 数据链路层 物理层
+ 2.2 经典协议与数据包
    - 以太网手部(18字节) + Ips首部(20~60) + TCP首部（20）+ 应用数据 + 以太网尾部
    - 最大1500字节 TCP首部和应用数据在抓包的时候比较常用(TCP段)
    - 应用数据基于底层HTTP协议的支持  报文 包含 请求 + 响应
    - websocket 握手协议 底层信用http协议建立一个长连接 然后通过二进制流来传递
    - websocket data协议 Fin Rsv opcode Mask payloadlen Masking-key PayloadData
+ 2.3 TCP的三次握手、和四次挥手
    - 保证连接的双工和可控 ， 通过数据的重传来保证
    - syn seq - ack 
    - tcpdump
+ 2.4 为啥timeWait要等待2MSL?
    - (Maximum Segment Lifetime, 30秒~1分钟)
    - 保证TCP协议的全双工连接并能够可靠的关闭
    - 保证这次连接的重复数据段从网络中消失
    
      为何会出现大量的closeWait?
    - 首先closewait一般出现在被关闭方
    - 并发请求太多导致
    - 被动关闭方没有及时释放占用的内存
    - [实际例子](demo/base/close_wait_test/server/main.php)
        1. 问题1 ** git.apache.org/thrift.git@v0.13.0: Get "https://proxy.golang.org/git.apache.org/thrift.git/@v/v0.13.0.mod":** 
           ```golang
              //因为go的服务管理代理被墙，需要切换为国内的代理服务
              sudo  go env -w GO111MODULE=on
              go env -w GOPROXY=https://goproxy.cn
           ```
+ 2.5 Tcp为啥需要流量控制？
    - 由于通讯双方网速不同，通讯任一方发送过快都会导致对方消息处理不过来，所以需要把数据放到    缓存区中
    - 如果缓存区满了，发送方还在发送，接受方只能把数据包丢弃，因此需要控制发送的速率
    - 缓存区剩余大小称之为接受窗口，用变量win表示。如果win=0,则发送方停止发送
    
+ 2.6 Tcp为啥需要拥塞控制？
    - 流量控制和拥塞控制是两个概念，拥塞控制是调解网络的负载。
    - 接收方网络资源繁忙，因未及时响应ACK导致发送方重传大量的数据，会导致网络更加拥堵
    - 拥塞控制是动态的调整缓冲区的大小，而非只是依赖缓冲区的大小来确定窗口
    - ![慢开始和拥塞避免](https://i.loli.net/2021/07/17/Sg9hcTzKbu7A8WX.jpg)
    - ![快重传和快恢复](https://i.loli.net/2021/07/17/oGbvFUHX9MQ5dlE.jpg)

2.7 为何会出现沾包和拆包？
    - 应用程序写入的数据大于套接字缓冲区的大小，这将就会发生拆包
    - 应用程序写入的数据小于套接字缓冲区的大小，网卡将应用多次写入的数据发送到网络上，这将会发生沾包
    - 进行MSS(最大的报文长度)大小的TCP分段，当TCP报文长度-tcp头部长度 > Mss将发生拆包
    - ![图解沾包和拆包](https://i.loli.net/2021/07/17/tk3WHZFXV7DBQmA.jpg)
    
    如何获取完整应用的数据报文？
    - 使用带消息头的协议，头部写入包长度，然后再读取包内容
    - 设置定长的消息，每次读取定长的内容，长度不够时空位补固定字符
    - 设置消息边界，服务端从网络流中按消息边界分离出消息内容，一般使用'\n'
    - 使用更为复杂的协议，比如 json protobuf(天然的支持边界，完整的结构体)

2.8 自定义消息格式实现装包与拆包

    ```
        graph LR
        A(如何获取完整的数据报文)-->B1(定义数据格式msg_header+content_len+content)
        A-->B2(编码encode)
        A-->B3(解码decode)
        A-->B4(tcp_client) 
        B4-->C1(连接服务器) 
        B4-->C2(数据编码)
        A-->B5(tcp_server) 
        B5-->C3(监听接口) 
        B5-->C4(接受请求) 
        B5-->C5(创建独立协程) 
        B5-->C6(数据解码)
    ```
    
 
