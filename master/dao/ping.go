package dao

import (
	"fmt"
	"time"

	"github.com/mangenotwork/common/utils"
	gt "github.com/mangenotwork/gathertool"
)

type PingRse struct {
	Item int    `json:"item"`
	Date string `json:"date"`
	Ms   string `json:"ms"`
	Err  string `json:"err"`
}

func Ping(ip string) []*PingRse {
	rse := make([]*PingRse, 0)

	for i := 0; i < 4; i++ {
		ping, err := gt.Ping(ip)
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}

		rse = append(rse, &PingRse{
			Item: i + 1,
			Date: utils.NowDate(),
			Ms:   fmt.Sprintf("%.2f ms", float64(ping)/1000000),
			Err:  errStr,
		})

		time.Sleep(500 * time.Millisecond)

	}

	return rse
}
