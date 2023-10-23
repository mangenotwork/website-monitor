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
- master 存储所有网站监测日志
- master 分析监测日志并产生报警
- master 统计报警并发送邮件通知
- master 报警信息
- master 任务表
- master 测试
- master 压力测试


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

