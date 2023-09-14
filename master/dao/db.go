package dao

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
)

const (
	WebSiteTable        = "website_table"
	MailTable           = "mail_table"
	MailConf            = "mail_conf"
	WebSiteURITable     = "website_uri_table"
	WebSitePointTable   = "website_point_table"
	WebSiteAlertTable   = "website_alert_table"
	MonitorErrInfoKey   = "MonitorErrInfo"
	MonitorErrInfoTable = "monitor_err_info_table"
	IncrementTable      = "increment_table"
	IncrementKey        = "Increment"
	IPTable             = "ip_table"
	MasterConfTable     = "master_conf_table"
	MasterConfKey       = "master_conf"
)

var (
	DBPath = "./data.db"
	Tables = []string{WebSiteTable, MailTable, WebSiteURITable, WebSitePointTable,
		WebSiteAlertTable, MonitorErrInfoTable, IncrementTable, IPTable, MasterConfTable}
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
