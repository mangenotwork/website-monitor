package dao

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
)

// Table
const (
	MailTable               = "mail_table"
	WebSiteTable            = "website_table"
	WebSiteInfoTable        = "website_Info_table"
	WebsiteAlarmRuleTable   = "website_alarm_rule_table"
	WebsiteScanCheckUpTable = "website_scan_checkup_table"
	WebSiteURITable         = "website_url_table"
	WebSiteUrlPointTable    = "website_url_point_table"
	WebSiteAlertTable       = "website_alert_table"
	IPTable                 = "ip_table"
	MasterConfTable         = "master_conf_table"
	HistoryWebsiteTDKITable = "History_WebsiteTDKI_table"
	HistoryIpTable          = "History_Ip_table"
	HistoryNsLookUpTable    = "History_NsLookUp_table"
	HistoryWhoisTable       = "History_Whois_table"
	HistoryICPTable         = "History_ICP_table"
	HistoryPingTable        = "History_Ping_table"
	HistorySSLTable         = "History_SSL_table"
	HistoryWebsiteInfoTable = "History_WebsiteInfo_table"
)

// KeyName
const (
	MailConfKeyName           = "mail_conf"
	MonitorErrInfoKeyName     = "MonitorErrInfo"
	IncrementKeyName          = "Increment"
	MasterConfKeyName         = "master_conf"
	HistoryWebsiteTDKIKeyName = "History_WebsiteTDKI_KEY"
	HistoryIpKeyName          = "History_Ip_KEY"
	HistoryNsLookUpKeyName    = "History_NsLookUp_KEY"
	HistoryWhoisKeyName       = "History_Whois_KEY"
	HistoryICPKeyName         = "History_ICP_KEY"
	HistoryPingKeyName        = "History_Ping_KEY"
	HistorySSLKeyName         = "History_SSL_KEY"
	HistoryWebsiteInfoKeyName = "History_WebsiteInfo_KEY"
)

var (
	DBPath = "./data.db"
	Tables = []string{
		MailTable, WebSiteTable, WebSiteInfoTable, WebsiteAlarmRuleTable, WebsiteScanCheckUpTable,
		WebSiteURITable, WebSiteUrlPointTable, WebSiteAlertTable, IPTable, MasterConfTable,
		HistoryWebsiteTDKITable, HistoryIpTable, HistoryNsLookUpTable, HistoryWhoisTable, HistoryICPTable,
		HistoryPingTable, HistorySSLTable, HistoryWebsiteInfoTable,
	}
	DB     = NewLocalDB(DBPath, Tables)
	ISNULL = fmt.Errorf("ISNULL")
)

type LocalDB struct {
	Path   string
	Tables []string
	Conn   *bolt.DB
}

func NewLocalDB(path string, tables []string) *LocalDB {
	return &LocalDB{
		Path:   path,
		Tables: tables,
	}
}

func (ldb *LocalDB) Init() {
	db, err := bolt.Open(ldb.Path, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		_ = db.Close()
	}()
	for _, table := range ldb.Tables {
		log.Info("检查数据表 : ", table)
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(table))
			if b == nil {
				_, err = tx.CreateBucket([]byte(table))
				if err != nil {
					log.Panic(err)
				}
			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
	}
}

func (ldb *LocalDB) Open() {
	ldb.Conn, _ = bolt.Open(ldb.Path, 0600, nil)
}

func (ldb *LocalDB) Get(table, key string, data interface{}) error {
	ldb.Open()
	defer func() {
		_ = ldb.Conn.Close()
	}()
	return ldb.Conn.View(func(tx *bolt.Tx) error {
	R:
		b := tx.Bucket([]byte(table))
		if b == nil {
			err := ldb.ClearTable(table)
			if err != nil {
				return err
			}
			goto R
		}
		bt := b.Get([]byte(key))
		if len(bt) < 1 {
			return ISNULL
		}
		err := json.Unmarshal(bt, data)
		if err != nil {
			return err
		}
		return nil
	})
}

func (ldb *LocalDB) Set(table, key string, data interface{}) error {
	ldb.Open()
	defer func() {
		_ = ldb.Conn.Close()
	}()
	value, err := utils.AnyToJsonB(data)
	if err != nil {
		return err
	}
	return ldb.Conn.Update(func(tx *bolt.Tx) error {
	R:
		b := tx.Bucket([]byte(table))
		if b == nil {
			err = ldb.ClearTable(table)
			if err != nil {
				return err
			}
			goto R
		}
		err = b.Put([]byte(key), value)
		if err != nil {
			return err
		}
		return nil
	})
}

func (ldb *LocalDB) Delete(table, key string) error {
	ldb.Open()
	defer func() {
		_ = ldb.Conn.Close()
	}()
	return ldb.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b == nil {
			return fmt.Errorf("未获取到表")
		}
		if err := b.Delete([]byte(key)); err != nil {
			return err
		}
		return nil
	})
}

func (ldb *LocalDB) ClearTable(table string) error {
	ldb.Open()
	defer func() {
		_ = ldb.Conn.Close()
	}()
	return ldb.Conn.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(table))
	})
}
