## 2_20210718_moicroGo学习笔记

+ 2.8 代码解析
    - 定义数据的编码规范 unpack 
    - 数据的写入 tcp
    - 数据的服务器 tcpServer
 
 + 2.9 Golang如何创建UDP的服务器和客户端？
     ```
         gragh LR
         A(UDP服务器和客户端的代码演示)-->B1(服务端)
         A-->B2(客户端)
         B1-->C1(step1 监听服务器)
         B1-->C2(step2 循环读取消息内容)
         B1-->C3(step3 回复数据)
         B2-->C4(Step1 连接服务器)
         B2-->C5(Step2 发送数据)
         B2-->C6(Step3 接收数据)
     ```
     - 代码 udpServer
     - 代码 udpClient

