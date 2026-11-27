package business

import (
	"time"

	"website-monitor/master/dao"

	"github.com/mangenotwork/common/log"
)

func Collect() {

	// todo 采集周期需要可配置
	timer := time.NewTimer(time.Hour * 1) //初始化定时器

	for {

		select {

		case <-timer.C:
			log.Info("采集周期...")

			// 读取站点扫描数据信息
			website, _, _ := dao.NewWebsite().SelectList()
			for _, v := range website {

				scan, _ := dao.NewWebsite().GetScanCheckUp(v.HostID)
				if scan == nil || scan.ScanRate <= 0 {
					continue
				}
				tNow := time.Now().Unix()

				// 当前时间 - 创建时间 -> 计算出小时差 |  小时差 余 扫描频率 == 0 就执行扫描
				if ((tNow-v.Created)/(60*60))%scan.ScanRate == 0 {

					// Go 1.19 循环变量复用，必须显式复制，否则所有 goroutine 都扫最后一个站点
					v, scan := v, scan
					go func() {
						dao.Scan(v.Host, v.HostID, scan.ScanDepth)
					}()

				}

			}

			timer.Reset(time.Hour * 1)
		}

	}
}
