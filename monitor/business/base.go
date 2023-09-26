package business

import (
	"encoding/json"
	"fmt"
	"github.com/mangenotwork/beacon-tower/udp"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/ginHelper"
	"github.com/mangenotwork/common/log"
	gt "github.com/mangenotwork/gathertool"
	"time"
)

func Initialize(client *udp.Client) {
	go func() {
		GetWebsite()
	}()
}

func Business(client *udp.Client) {
	log.Info("业务...")
	go func() {
		for {
			time.Sleep(10 * time.Second)
			rse, err := client.Get("conn/test", []byte("test"))
			if err != nil {
				udp.Error(err)
				return
			}
			udp.Info("get 请求返回 = ", string(rse))
		}
	}()
}

var (
	MasterHTTP          = ""
	GetAllWebsiteAPI    = "/data/all/website"
	GetWebsiteAllUrlAPI = func(id string) string { return fmt.Sprintf("/data/allurl/%s", id) }
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
	for _, v := range list {
		log.Info("v = ", v)
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
