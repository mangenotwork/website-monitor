package business

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/ginHelper"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	gt "github.com/mangenotwork/gathertool"
	udp "github.com/mangenotwork/udp_comm"
	"sync"
	"time"
)

const (
	URIHealth        = "Health"
	URIRandom        = "Random"
	URIPoint         = "Point"
	LogTypeInfo      = "Info"
	LogTypeAlert     = "Alert"
	LogTypeError     = "Error"
	AlertTypeNone    = ""
	AlertTypeErr     = "err"
	AlertTypeCode    = "code"
	AlertTypeTimeout = "timeout"
)

var GlobalName = "监视器"

func Initialize(client *udp.Client) {

	// 监测器获取ip地址信息
	GetMyIP()

	// 拉取website列表和监测点数据
	GetWebsiteAll()

	// 启动监测
	go func() {
		timer := time.NewTimer(time.Second * 1) //初始化定时器
		for {
			select {
			case <-timer.C:
				AllWebsiteData.Range(func(key, value any) bool {
					website := value.(*WebsiteItem)
					website.Conn = client
					Business(website)
					return true
				})
				timer.Reset(time.Second * 1)
			}
		}
	}()

	// 启动定时器 15分钟一次，定时拉取监测网站url
	go func() {
		timer := time.NewTimer(time.Minute * 15) //初始化定时器
		for {
			select {
			case <-timer.C:
				log.Info("执行定期器，定时拉取监测网站url")
				AllWebsiteData.Range(func(key, value any) bool {
					website := value.(*WebsiteItem)
					SetWebsiteUrlData(website.HostID)
					return true
				})
				timer.Reset(time.Minute * 15)
			}
		}
	}()
}

/*
Business

	每一轮监测会先进行确认当前网络，再执行三次请求监测

	确认网络
	1. ping设置的ip,响应时间大于 1000ms 视为网络不好本轮不执行监测
	2. 请求对照组，响应时间大于设置的时间 视为网络不好本轮不执行监测

	三次请求监测
	1. 监测host确认网站是否存活
	2. 监测网站采集到的url随机选择一个进行监测，如果是新添加的网站采集url会延后; 如果存在强行休息1s再执行
	3. 监测设置的网站监测点依次进行监测; 如果存在强行休息1s再执行
*/

func Business(item *WebsiteItem) {
	item.RateItem--
	if item.RateItem <= 0 {
		// 计算频率复位
		item.RateItem = item.MonitorRate
		log.Info("执行 " + item.Host)
		// 报警数据初始化
		// 日志数据
		mLog := &MonitorLog{
			MonitorName: GlobalName, // TODO 获取监测器id和昵称
			LogType:     "Info",
			Time:        utils.NowDate(),
			HostId:      item.HostID,
			Host:        item.Host,
			ContrastUri: item.ContrastUrl,
			Ping:        item.Ping,
			MonitorIP:   MyIP.IP,
			MonitorAddr: MyIP.Address,
		}

		// ping一下，检查当前网络环境
		_, pingRse := item.PingActive(mLog)
		if !pingRse {
			// 网络环境异常不执行监测
			return
		}

		// 请求对照组，对照组有问题不执行监测
		if item.ContrastActive(mLog) {
			return
		}

		// 监测生命URI
		item.MonitorHealthUri(mLog)

		// 随机URI监测  需要获取所有网站Url
		item.MonitorRandomUri(mLog)

		// 判断是否执行监测点监测
		point, pointLen := GetWebSitePoint(item.HostID)
		log.Info("监测点个数 ==> ", pointLen)
		if pointLen > 0 {
			log.Info("执行监测点监测")
			if item.LoopPoint >= pointLen {
				item.LoopPoint = 0
			}
			pointUrl := point[item.LoopPoint]
			item.LoopPoint++
			log.Info("本次监测的监测点是 = ", pointUrl)
			// 执行监测点监测
			item.MonitorPointUri(mLog, pointUrl)
		}

	}
}

var (
	MasterHTTP          = ""
	GetAllWebsiteAPI    = "/data/all/website"
	GetWebsitePointAPI  = func(id string) string { return fmt.Sprintf("/data/website/point/%s", id) }
	GetWebsiteAPI       = func(id string) string { return fmt.Sprintf("/data/website/%s", id) }
	GetWebsiteAllUrlAPI = func(id string) string { return fmt.Sprintf("/data/allurl/%s", id) }
)

func GetWebsiteAll() {
	log.Info("启动获取 website ")
	ctx, err := gt.Get(MasterHTTP + GetAllWebsiteAPI)
	if err != nil {
		log.Error(err)
		return
	}
	list := make([]*Website, 0)
	err = AnalysisData(ctx.Json, &list)
	if err != nil {
		log.Error(err)
		return
	}
	// 清空 websiteItem
	EmptyAllWebsiteData()
	for _, v := range list {
		log.Info("v = ", v)
		item := &WebsiteItem{v, 0, nil, 0}
		item.Add()
		// 获取监测点
		GetWebsitePoint(v.HostID)
		// 获取所有url
		SetWebsiteUrlData(v.HostID)
	}
}

func GetWebsite(hostID string) {
	ctx, err := gt.Get(MasterHTTP + GetWebsiteAPI(hostID))
	if err != nil {
		log.Error(err)
		return
	}
	data := &Website{}
	log.Info(ctx.Json)
	err = AnalysisData(ctx.Json, &data)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(data)
	item := &WebsiteItem{data, 0, nil, 0}
	item.Add()
	// 获取监测点
	GetWebsitePoint(data.HostID)
}

func DelWebsite(hostID string) {
	AllWebsiteData.Delete(hostID)
}

func GetWebsitePoint(hostID string) {
	ctx, err := gt.Get(MasterHTTP + GetWebsitePointAPI(hostID))
	if err != nil {
		log.Error(err)
		return
	}
	data := &WebSitePoint{}
	err = AnalysisData(ctx.Json, &data)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(data)
	WebsitePointDataMap.Store(hostID, data)
}

type IPInfo struct {
	IP      string
	Address string
}

const GetMyIPInfoUrl = "https://www.ip.cn/api/index?ip=&type=0"

var MyIP *IPInfo
var MyIPOnce sync.Once

func GetMyIP() *IPInfo {
	MyIPOnce.Do(func() {
		MyIP = getIPAddr()
	})
	return MyIP
}

func getIPAddr() *IPInfo {
	ctx, _ := gt.Get(GetMyIPInfoUrl)
	ip, _ := gt.JsonFind2Str(ctx.Json, "/ip")
	address, _ := gt.JsonFind2Str(ctx.Json, "/address")
	log.Info("获取到ip属地= ", address)
	return &IPInfo{
		IP:      ip,
		Address: address,
	}
}

func GetMasterHTTP() string {
	url, err := conf.YamlGetString("masterHTTP")
	if err != nil {
		log.Error(err)
	}
	return url
}

func AnalysisBody(jsonStr string) (any, error) {
	body := &ginHelper.ResponseJson{}
	err := json.Unmarshal([]byte(jsonStr), &body)
	if err != nil {
		return nil, err
	}
	return body.Date, nil
}

func AnalysisData(jsonStr string, data any) error {
	body, err := AnalysisBody(jsonStr)
	if err != nil {
		return err
	}
	bodyStr, err := json.Marshal(body)
	if err != nil {
		log.Error(err)
		return err
	}
	err = json.Unmarshal(bodyStr, &data)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// MonitorLog 监测日志
type MonitorLog struct {
	LogType         string // Info  Alert Error
	Time            string
	HostId          string
	Host            string
	UriType         string // 监测的URI类型 Health:根URI,健康URI  Random:随机URI  Point:监测点URI
	Uri             string // URI
	UriCode         int    // URI响应码
	UriMs           int64  // URI响应时间
	ContrastUri     string // 对照组URI
	ContrastUriCode int    // 对照组URI响应码
	ContrastUriMs   int64  // 对照组URI响应时间
	Ping            string
	PingMs          int64
	Msg             string
	AlertType       string // code  timeout
	MonitorName     string // 监测器名称
	MonitorIP       string // 监测器 ip地域信息
	MonitorAddr     string
}

func Struct2Buf(any2 any) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, any2)
	if err != nil {
		log.Error(err)
	}
	return buf.Bytes()
}

func request(url string) (int, int64, error) {
	ctx, err := gt.Get(url, gt.Header{
		"Accept-Encoding": "gzip, deflate, br",
		"Referer":         url,
	},
		gt.RetryTimes(1), // 重试设为0
	)
	if err != nil {
		log.Error(err)
		return 0, 0, err
	}
	return ctx.StateCode, ctx.Ms.Milliseconds(), nil
}
