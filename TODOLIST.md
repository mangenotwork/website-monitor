#### v0.0.1

监测功能，基础页面，监测器，基础工具等

- [ok] master 登录 
- [ok] master 邮件 
- [ok] master 设计"网站" 实体 
- [ok] 证书信息获取，dns扫描  
- [ok] IP属地获取，ipc查询 
- [ok] tool : 获取网站的T, D, K, 图标  
- [ok] tool : ip信息查询 
- [ok] tool : 查询dns  
- [ok] tool : Whois查询  
- [ok] tool : 查询备案  
- [ok] tool : 在线ping 
- [ok] tool : 获取证书  
- [ok] tool : 网站信息获取  
- [ok] master 添加网站 
- [ok] master 监测网站列表 
- [ok] master 获取网站基础信息 
- [ok] master 删除网站监测 
- [ok] master 扫描网站Url 
- [ok] master 扫描网站 css, js url  
- [ok] master 扫描网站 其他静态资源 url 
- [ok] master 定期扫描网站  
- [ok] monitor 每次启动拉一次监测任务表 
- [ok] monitor 监测网站-Host基础 
- [ok] master 监测日志存储和读取 
- [ok] master 监测点设置  
- [ok] master 通知拉取网站监测点
- [ok] master 删除网站监测
- [ok] master 增加网站监测 改为指定增加网站监测
- [ok] master 删除网站监测 改为指定网站删除
- [ok] master 修改网站监测配置 指定修改
- [ok] master 网站扫描的Url
- [ok] master 查看日志列表
- [ok] master 监测图表
- [ok] monitor 定时拉取网站的url
- [ok] monitor 执行网站监测点监测
- [ok] monitor 执行网站随机url监测
- [ok] master 获取monitor在线情况
- [ok] master 监测器列表
- [ok] monitor 获取ip环境地址
- [ok] monitor 获取系统信息
- [ok] master 首页未部署监测器提示
- [ok] master 监测器信息展示
- [ok] master 存储所有网站监测日志
- [ok] master 分析监测日志并产生报警
- [ok] master 发送报警邮件通知
- [ok] master 报警信息列表接口 
- [ok] master 报警信息列表数量接口 指定网站获取接口
- [ok] master 报警信息列表，新增导航
- [ok] master 首页报警信息
- [ok] master 监测列表报警信息
- [ok] master 报警消息删除与清空
- [ok] master 首页网站列表查看报警消息
- [ok] v0.0.1 测试
- [ok] v0.0.1 改bug
- [ok] v0.0.1 git发布  2023-10-31

#### v0.0.2

主要开发请求器功能，参考postman  apiPost

- [ok] master 请求器实体设计
- [ok] user agent 列表
- [ok] master 全局开关邮件通知
- [ok] master 导航栏 快速请求
- [ok] master 请求器 页面设计
- [ok] master 请求器 页面设计 几种Body界面设计
- [ok] master 请求器 交互设计
- [ok] master 全局Header相关接口

#### v0.0.3

主要是优化和改bug

- 增加代码可读性 [ok]
- 代码评审和优化
- UI界面优化
- 交互逻辑优化
- [优化]面板数据加载很慢
- 修改bug

#### v0.0.4

主要是提升监控的准确性和监控维度的调优







[bug]
4. ipc信息读取不到 (有反爬，可以找替代或者解决反爬)
```
替换地址:
https://www.aichaicp.com/latest
https://m.chaicp.com/icp.html
https://icplishi.com/
http://freeicp.com/
https://icp.5118.com/
```
15. 修改了超时时间 monitor 视乎没有起效果，报警信息还是之前的超时做对比
17. unexpected EOF 请求错误忽略
```
2023-11-02 09:52:36 https://www.97654.com请求超时 : Get "https://www.97654.com": context deadline exceeded (Client.Timeout exceeded while awaiting headers) ;
```
18. [ok] 随机监测的url 需要排除无效标签，如下
```
https://www.8300.cn/zst/ssqjbw4;
```



### 需求池
- 使用说明
- master 请求相关接口
- master 历史请求相关接口
- master 创建请求目录相关接口
- master 请求保存相关接口
- master 指定 monitor 请求，未指定默认 master请求
- master 请求器 请求记录
- master 请求器 保存请求 增删改查
- master 扫描检查证书过期时间
- master 接口测试 页面
- master 接口测试 脚本编写
- master 接口测试 管理 与 增删改查
- master 接口测试 执行，未指定 monitor 默认 master执行
- monitor 执行接口测试
- master 可视化操作新增或编辑接口测试
- 压力测试开发
- 渗透测试开发
- 扫描等工具的开发
- 接口测试开发


---

其他工具
- tcp|udp端口扫描 
- 半开扫描（TCP SYN）
- MAC地址扫描
- ICMP扫描
- WWW服务扫描


渗透测试相关资料

```
1. 黑盒（Black Box）渗透。黑盒（Black Box）渗透测试又被称之为“zero-knowledge testing”，渗透者完全处于对目标网络系统一无所知的状态。
通常这类测试，只能通过DNS、Web、E-mail等网络对外公开提供的各种服务器进行扫描探测，从而获得公开的信息，以决定渗透的方案与步骤。


2. 漏洞扫描。通过上述的信息收集，在获得目标网络各服务器开放的服务之后，就可以对这些服务进行重点扫描，扫出其所存在的漏洞。常用的扫描工具主要有：
针对操作系统漏洞扫描的工具，包括X-Scan、ISS、Nessus、SSS、Retina等；针对Web网页服务的扫描工具，包括SQL扫描器、文件PHP包含扫描器、上传漏洞扫
描工具，以及各种专业全面的扫描系统，如AppScan、Acunetix Web Vulnerability Scanner（如图1-53所示）等；针对数据库的扫描工具，包括Shadow 
Database Scanner（如图1-54所示）、NGSSQuirreL，以及SQL空口令扫描器等。另外，许多入侵者或渗透测试人员也有自己的专用扫描器，其使用更加个性化。


3. 相关软件工具（如扫描工具X-Scan等）来收集网络系统中各个主机系统的相关信息


4. 上传漏洞扫描
例如：正常的上传路径是“http://www.xxx.com/net op/upload/01.jpg”，但黑客可使用“\0”来构造filepath为“http://www.xxx.com.cn/netop.asp \0/01.jpg”，
这样当服务器接收filepath数据时，就会简单地看到“netop.asp”后面的“\0”后，认为filepath数据到此就结束了，此时上传文件就被保存为www.xxx.com.cn/netop.asp。
利用这个上传漏洞就可以上传ASP木马，再连接上传的网页，进而控制整个网站系统了。


5. SYN攻击利器HGod；  SYN-Killer是由国人编写的一款基于SYN攻击程序的代表作品； 


6.  网络协议与活跃主机发现技术。▯ 基于ARP协议的活跃主机发现技术。▯ 基于ICMP协议的活跃主机发现技术。▯ 基于TCP协议的活跃主机发现技术。
▯ 基于UDP协议的活跃主机发现技术。▯ 基于SCTP协议的活跃主机发现技术。


7. http 报文攻击方式
这里有个有意思的事情，我在查询http报文资料的时候，发现有一种攻击方法，称之为slow header和 slow post，这里来解释下是什么意思。
slow header：表示客户端连接到服务器后，通过慢速度发送数据，但是一直不发送\r\n\r\n，服务器一直在接收，所以始终占着服务器连接，当该种连接过多时，会导致服务器连接数满，从而不能接收新的请求。
slow post: 这里指的是，通过post发送数据，但是将Content-Length设置的很大，还是每次只发送很小的数据，和上述一样，当该种连接过多时候，会导致服务器连接数满，从而不能接收新的请求。

参考: https://www.cnblogs.com/NoneID/p/17513530.html


8. 信息搜索工具: Whois   nslookup命令工具  OneForAll子域名搜集工具
扫描探测工具:  Nmap  Nessus   AWVS   Dirsearch    Nikto
漏洞扫描工具：  sqlmap注入工具   
Webshell管理工具：   冰蝎    中国蚁剑    哥斯拉   
网络抓包分析工具：   Wireshark   Fiddler    tcpdump   


9. 地址扫描探测。主要利用ARP、ICMP请求目标IP或网段，通过回应消息获取目标网段中存活机器的IP地址和MAC地址，进而掌握拓扑结构。


10. 设备指纹探测。根据扫描返回的数据包匹配TCP/IP协议栈指纹来识别不同的操作系统和设备


11. 




```

---
网站监控参考
```
https://github.com/argentmoon/host-monitor
https://github.com/dhjz/dwatch
https://github.com/ptonlix/netdog
https://github.com/xxscloud5722/goflame
https://github.com/JiLoveZn/WebMonitor
https://github.com/zmh-program/fymonitor
https://github.com/xuexiangyou/monitor-web
https://github.com/zjw939057120/Website-Monitoring
```

---

集成  nmap , 如下参考资料
```
https://github.com/Ullaakut/nmap
https://blog.csdn.net/m0_56262476/article/details/128728464
https://github.com/lair-framework/go-nmap
https://github.com/lcvvvv/gonmap
https://github.com/ivopetiz/network-scanner
https://github.com/xiaoyaochen/yscan
https://github.com/vus520/go-scan
https://github.com/vdjagilev/nmap-formatter
https://github.com/CTF-MissFeng/NmapTools
https://github.com/marco-lancini/goscan
https://github.com/projectdiscovery/naabu
https://github.com/Ullaakut/Gorsair
https://github.com/Adminisme/ServerScan
https://github.com/hktalent/scan4all
https://github.com/luijait/GONET-Scanner
https://github.com/qq431169079/PortScanner-3
```

---

DOS库 参考
```
https://github.com/grafov/hulk
https://github.com/marant/goloris
https://github.com/IgorHalfeld/lagoinha
https://github.com/UBISOFT-1/AnonymousPAK-DDoS
https://github.com/XORbit01/ddosarmy
https://github.com/Xart3mis/AKILT
https://github.com/jantechner/dos-attacker
https://github.com/a7600999/goformdos
```


---

子域名探索方法
在线接口
暴力枚举
DNS解析
爬虫 Scraping（抓取）


---


API 攻击
API 攻击是对应用程序编程接口 (API) 的恶意使用或破坏。API 安全包括防止攻击者利用和滥用 API 的实践和技术。黑客以 API 为目标，因为它们是现代 Web 应用程序和微服务架构的核心。

API 攻击的例子包括：

注入攻击：当 API 未正确验证其输入并允许攻击者提交恶意代码作为 API 请求的一部分时，就会发生这种类型的攻击。SQL 注入 (SQLi) 和跨站点脚本 (XSS) 是最突出的例子，但还有其他例子。传统上针对网站和数据库的大多数类型的注入攻击也可用于攻击 API。

DoS/DDoS 攻击：在拒绝服务 (DoS) 或分布式拒绝服务 (DDoS) 攻击中，攻击者试图使 API 对目标用户不可用。速率限制可以帮助缓解小规模的 DoS 攻击，但大规模的 DDoS 攻击可以利用数百万台计算机，并且只能通过云规模的反 DDoS 技术来解决。

数据暴露： API 经常处理和传输敏感数据，包括信用卡信息、密码、会话令牌或个人身份信息 (PII)。如果 API 处理数据不正确，如果它很容易被诱骗向未经授权的用户提供数据，以及如果攻击者设法破坏 API 服务器，则数据可能会受到损害。