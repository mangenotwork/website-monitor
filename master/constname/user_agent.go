package constname

import (
	"math/rand"
	"time"
)

type UserAgentType int

const (
	PCAgent UserAgentType = iota + 1
	WindowsAgent
	LinuxAgent
	MacAgent
	AndroidAgent
	IosAgent
	PhoneAgent
	WindowsPhoneAgent
	UCAgent
)

var UserAgentMap map[int]string = map[int]string{
	1:  "Mozilla/5.0 (Windows NT 6.2; WOW64; rv:21.0) Gecko/20100101 Firefox/21.0",                                                                                                                                          //Firefox on Windows
	2:  "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.94 Safari/537.36",                                                                                                      //Chrome on Windows
	3:  "Mozilla/5.0 (compatible; WOW64; MSIE 10.0; Windows NT 6.2)",                                                                                                                                                        //Internet Explorer 10
	4:  "Opera/9.80 (Windows NT 6.1; WOW64; U; en) Presto/2.10.229 Version/11.62",                                                                                                                                           //Opera on Windows
	5:  "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/533.20.25 (KHTML, like Gecko) Version/5.0.4 Safari/533.20.27",                                                                                          //Safari on Windows
	6:  "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:21.0) Gecko/20130331 Firefox/21.0",                                                                                                                                      //Firefox on Ubuntu
	7:  "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/535.11 (KHTML, like Gecko) Ubuntu/11.10 Chromium/27.0.1453.93 Chrome/27.0.1453.93 Safari/537.36",                                                                       //Chrome on Ubuntu
	8:  "Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_6; en-US) AppleWebKit/533.20.25 (KHTML, like Gecko) Version/5.0.4 Safari/533.20.27",                                                                                 //Safari on Mac
	9:  "Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.9.168 Version/11.52",                                                                                                                                 //Opera on Mac
	10: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.93 Safari/537.36",                                                                                           //Chrome on Mac
	11: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:21.0) Gecko/20100101 Firefox/21.0",                                                                                                                                 //Firefox on Mac
	12: "Mozilla/5.0 (Linux; Android 4.1.1; Nexus 7 Build/JRO03D) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166  Safari/535.19",                                                                               //Nexus 7 (Tablet)
	13: "Mozilla/5.0 (Linux; U; Android 4.0.4; en-gb; GT-I9300 Build/IMM76D) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",                                                                       //Samsung Galaxy S3 (Handset)
	14: "Mozilla/5.0 (Linux; U; Android 2.2; en-gb; GT-P1000 Build/FROYO) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",                                                                            //Samsung Galaxy Tab (Tablet)
	15: "Mozilla/5.0 (Android; Mobile; rv:14.0) Gecko/14.0 Firefox/14.0",                                                                                                                                                    //Firefox on Android Mobile
	16: "Mozilla/5.0 (Android; Tablet; rv:14.0) Gecko/14.0 Firefox/14.0",                                                                                                                                                    //Firefox on Android Tablet
	17: "Mozilla/5.0 (Linux; Android 4.0.4; Galaxy Nexus Build/IMM76B) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.133 Mobile Safari/535.19",                                                                    //Chrome on Android Mobile
	18: "Mozilla/5.0 (Linux; Android 4.1.2; Nexus 7 Build/JZ054K) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166 Safari/535.19",                                                                                //Chrome on Android Tablet
	19: "Mozilla/5.0 (iPhone; CPU iPhone OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",                                                                            //iPhone
	20: "Mozilla/5.0 (iPad; CPU OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",                                                                                     //iPad
	21: "Mozilla/5.0 (iPhone; CPU iPhone OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",                                                                            //Safari on iPhone
	22: "Mozilla/5.0 (iPad; CPU OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",                                                                                     //Safari on iPad
	23: "Mozilla/5.0 (iPhone; CPU iPhone OS 6_1_4 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) CriOS/27.0.1453.10 Mobile/10B350 Safari/8536.25",                                                                    //Chrome on iPhone
	24: "Mozilla/5.0 (compatible; MSIE 10.0; Windows Phone 8.0; Trident/6.0; IEMobile/10.0; ARM; Touch; NOKIA; Lumia 920)",                                                                                                  //Windows Phone 8
	25: "Mozilla/5.0 (compatible; MSIE 9.0; Windows Phone OS 7.5; Trident/5.0; IEMobile/9.0; SAMSUNG; SGH-i917)",                                                                                                            //Windows Phone 7.5
	26: "User-Agent, UCWEB7.0.2.37/28/999",                                                                                                                                                                                  //UC无
	27: "User-Agent, NOKIA5700/ UCWEB7.0.2.37/28/999",                                                                                                                                                                       //UC标准
	28: "User-Agent, Openwave/ UCWEB7.0.2.37/28/999",                                                                                                                                                                        //UCOpenwave
	29: "User-Agent, Mozilla/4.0 (compatible; MSIE 6.0; ) Opera/UCWEB7.0.2.37/28/999",                                                                                                                                       //UC Opera
	30: "Mozilla/5.0 (Windows; U; Windows NT 6.1; ) AppleWebKit/534.12 (KHTML, like Gecko) Maxthon/3.0 Safari/534.12",                                                                                                       //傲游3.1.7在Win7+ie9,高速模式
	31: "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E)",                    //傲游3.1.7在Win7+ie9,IE内核兼容模式:
	32: "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E; SE 2.X MetaSr 1.0)", //搜狗3.0在Win7+ie9,IE内核兼容模式
	33: "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.3 (KHTML, like Gecko) Chrome/6.0.472.33 Safari/534.3 SE 2.X MetaSr 1.0",                                                                            //搜狗3.0在Win7+ie9,高速模式
	34: "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E)",                    //360浏览器3.0在Win7+ie9
	35: "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/13.0.782.41 Safari/535.1 QQBrowser/6.9.11079.201",                                                                                        //QQ浏览器6.9(11079)在Win7+ie9,极速模式
}

// 每种类型设备的useragent的列表
var (
	listPCAgent           = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35}
	listWindowsAgent      = []int{1, 2, 3, 4, 5, 30, 31, 32, 33, 34, 35}
	listLinuxAgent        = []int{6, 7}
	listMacAgent          = []int{8, 9, 10, 11}
	listAndroidAgent      = []int{12, 13, 14, 15, 16, 17, 18}
	listIosAgent          = []int{19, 20, 21, 22, 23}
	listPhoneAgent        = []int{12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	listWindowsPhoneAgent = []int{24, 25}
	listUCAgent           = []int{26, 27, 28, 29}
)

var AgentType = map[UserAgentType][]int{
	PCAgent:           listPCAgent,
	WindowsAgent:      listWindowsAgent,
	LinuxAgent:        listLinuxAgent,
	MacAgent:          listMacAgent,
	AndroidAgent:      listAndroidAgent,
	IosAgent:          listIosAgent,
	PhoneAgent:        listPhoneAgent,
	WindowsPhoneAgent: listWindowsPhoneAgent,
	UCAgent:           listUCAgent,
}

// GetAgent 随机获取那种类型的 user-agent
func GetAgent(agentType UserAgentType) string {
	rand.Seed(time.Now().UnixNano())
	switch agentType {
	case PCAgent:
		if v, ok := UserAgentMap[listPCAgent[rand.Intn(len(listPCAgent))]]; ok {
			return v
		}
	case WindowsAgent:
		if v, ok := UserAgentMap[listWindowsAgent[rand.Intn(len(listWindowsAgent))]]; ok {
			return v
		}
	case LinuxAgent:
		if v, ok := UserAgentMap[listLinuxAgent[rand.Intn(len(listLinuxAgent))]]; ok {
			return v
		}
	case MacAgent:
		if v, ok := UserAgentMap[listMacAgent[rand.Intn(len(listMacAgent))]]; ok {
			return v
		}
	case AndroidAgent:
		if v, ok := UserAgentMap[listAndroidAgent[rand.Intn(len(listAndroidAgent))]]; ok {
			return v
		}
	case IosAgent:
		if v, ok := UserAgentMap[listIosAgent[rand.Intn(len(listIosAgent))]]; ok {
			return v
		}
	case PhoneAgent:
		if v, ok := UserAgentMap[listPhoneAgent[rand.Intn(len(listPhoneAgent))]]; ok {
			return v
		}
	case WindowsPhoneAgent:
		if v, ok := UserAgentMap[listWindowsPhoneAgent[rand.Intn(len(listWindowsPhoneAgent))]]; ok {
			return v
		}
	case UCAgent:
		if v, ok := UserAgentMap[listUCAgent[rand.Intn(len(listUCAgent))]]; ok {
			return v
		}
	default:
		if v, ok := UserAgentMap[rand.Intn(len(UserAgentMap))]; ok {
			return v
		}
	}
	return ""
}
