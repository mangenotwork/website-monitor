package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
)

// Table
const (
	MailTable                = "mail_table"
	WebSiteTable             = "website_table"
	WebSiteInfoTable         = "website_Info_table"
	WebsiteAlarmRuleTable    = "website_alarm_rule_table"
	WebsiteScanCheckUpTable  = "website_scan_checkup_table"
	WebSiteURITable          = "website_url_table"
	WebSiteUrlPointTable     = "website_url_point_table"
	WebSiteAlertTable        = "website_alert_table"
	IPTable                  = "ip_table"
	MasterConfTable          = "master_conf_table"
	HistoryWebsiteTDKITable  = "History_WebsiteTDKI_table"
	HistoryIpTable           = "History_Ip_table"
	HistoryNsLookUpTable     = "History_NsLookUp_table"
	HistoryWhoisTable        = "History_Whois_table"
	HistoryICPTable          = "History_ICP_table"
	HistoryPingTable         = "History_Ping_table"
	HistorySSLTable          = "History_SSL_table"
	HistoryWebsiteInfoTable  = "History_WebsiteInfo_table"
	AlertTable               = "alert_table"
	AlertWebsiteTable        = "alert_website_table"
	RequestTable             = "Request_Table"             // Request
	RequestDirTable          = "RequestDir_Table"          // RequestDir
	RequestCookieTable       = "RequestCookie_Table"       // RequestCookie
	RequestGlobalHeaderTable = "RequestGlobalHeader_Table" // RequestGlobalHeader
	RequestNowListTable      = "RequestNowList_Table"      // RequestNowList 导航栏上的请求
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
		HistoryPingTable, HistorySSLTable, HistoryWebsiteInfoTable, AlertTable, AlertWebsiteTable,
		RequestTable, RequestDirTable, RequestCookieTable, RequestGlobalHeaderTable, RequestNowListTable,
	}
	DB           = NewLocalDB(DBPath, Tables)
	ISNULL       = fmt.Errorf("ISNULL")
	errDBNotInit = errors.New("数据库未初始化")
)

// LocalDB 单进程内只保留一个 *bolt.DB 连接。
// Bolt 的读事务（View）可并发，写事务（Update）串行；
// 每次操作重新 bolt.Open 会拿独占文件锁，导致所有读写排队，必须避免。
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

// Init 打开数据库并创建缺失的桶，进程退出前不再关闭。
func (ldb *LocalDB) Init() {
	if ldb.Conn != nil {
		return
	}

	db, err := bolt.Open(ldb.Path, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Panic(err)
	}
	ldb.Conn = db

	err = db.Update(func(tx *bolt.Tx) error {
		for _, table := range ldb.Tables {
			if tx.Bucket([]byte(table)) != nil {
				continue
			}
			log.Info("创建数据表 : ", table)
			if _, err := tx.CreateBucket([]byte(table)); err != nil {
				return fmt.Errorf("创建数据表 %s 失败: %w", table, err)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	log.Info("数据库已打开: ", ldb.Path)
}

// Close 仅在进程退出时调用
func (ldb *LocalDB) Close() {
	if ldb.Conn == nil {
		return
	}
	if err := ldb.Conn.Close(); err != nil {
		log.Error("关闭数据库失败: ", err)
	}
	ldb.Conn = nil
}

func (ldb *LocalDB) conn() (*bolt.DB, error) {
	if ldb.Conn == nil {
		return nil, errDBNotInit
	}
	return ldb.Conn, nil
}

// View 只读事务，可并发调用
func (ldb *LocalDB) View(fn func(tx *bolt.Tx) error) error {
	db, err := ldb.conn()
	if err != nil {
		return err
	}
	return db.View(fn)
}

// Update 读写事务，进程内串行
func (ldb *LocalDB) Update(fn func(tx *bolt.Tx) error) error {
	db, err := ldb.conn()
	if err != nil {
		return err
	}
	return db.Update(fn)
}

func (ldb *LocalDB) Get(table, key string, data interface{}) error {
	return ldb.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(table))
		if b == nil {
			return fmt.Errorf("%s 表不存在", table)
		}

		bt := b.Get([]byte(key))
		if len(bt) < 1 {
			return ISNULL
		}

		return json.Unmarshal(bt, data)
	})
}

func (ldb *LocalDB) Set(table, key string, data interface{}) error {
	value, err := utils.AnyToJsonB(data)
	if err != nil {
		return err
	}

	return ldb.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(table))
		if b == nil {
			var createErr error
			b, createErr = tx.CreateBucket([]byte(table))
			if createErr != nil {
				return fmt.Errorf("创建数据表 %s 失败: %w", table, createErr)
			}
		}

		return b.Put([]byte(key), value)
	})
}

func (ldb *LocalDB) Delete(table, key string) error {
	return ldb.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(table))
		if b == nil {
			return fmt.Errorf("%s 表不存在", table)
		}

		return b.Delete([]byte(key))
	})
}

func (ldb *LocalDB) ClearTable(table string) error {
	return ldb.Update(func(tx *bolt.Tx) error {
		if tx.Bucket([]byte(table)) == nil {
			return nil
		}
		return tx.DeleteBucket([]byte(table))
	})
}
