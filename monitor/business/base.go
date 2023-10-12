package business

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/mangenotwork/beacon-tower/udp"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/ginHelper"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	gt "github.com/mangenotwork/gathertool"
	"sync"
	"time"
)

const (
	URIHealth    = "Health"
	URIRandom    = "Random"
	URIPoint     = "Point"
	LogTypeInfo  = "Info"
	LogTypeAlert = "Alert"
	LogTypeError = "Error"
)

var GlobalName = "监视器"

func Initialize(client *udp.Client) {

	// 监测器获取ip地址信息
	GetMyIP()

	// 拉取website列表和监测点数据
	GetWebsite()

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

}

/*
Business

	每一轮监测会先进行确认当前网络，再执行三次请求监测

	确认网络
	1. ping设置的ip,响应时间大于 1000ms 视为网络不好本轮不执行监测
	2. 请求对照组，响应时间大于设置的时间 视为网络不好本轮不执行监测

	三次请求监测
	1. 监测host确认网站是否存活
	2. 监测网站采集到的url随机选择一个进行监测，如果是新添加的网站采集url会延后
	3. 监测设置的网站监测点依次进行监测
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

		// TODO 随机URI监测  需要获取所有网站Url
		// isAlert2 := item.MonitorRandomUri(masterConf, mLog, alert)

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
			// TODO... 执行监测点监测
		}

	}
}

var (
	MasterHTTP         = ""
	GetAllWebsiteAPI   = "/data/all/website"
	GetWebsitePointAPI = func(id string) string { return fmt.Sprintf("/data/website/point/%s", id) }
)

func GetWebsite() {
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
	}
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
	})
	if err != nil {
		log.Error(err)
		return 0, 0, err
	}
	return ctx.StateCode, ctx.Ms.Milliseconds(), nil
}