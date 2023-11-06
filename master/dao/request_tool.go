package dao

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"website-monitor/master/entity"
)

type RequestToolEr interface {
	SetGlobalHeader(list []*entity.RequesterGlobalHeader) error
	GetGlobalHeader() ([]*entity.RequesterGlobalHeader, error)
	DelGlobalHeader(key string) error
}

func NewRequestTool() RequestToolEr {
	return new(requestToolDao)
}

type requestToolDao struct{}

func (r *requestToolDao) SetGlobalHeader(list []*entity.RequesterGlobalHeader) error {
	var err error
	for _, v := range list {
		err = DB.Set(RequestGlobalHeaderTable, v.Key, v)
	}
	return err
}

func (r *requestToolDao) GetGlobalHeader() ([]*entity.RequesterGlobalHeader, error) {
	conn := GetDBConn()
	defer func() {
		_ = conn.Close()
	}()
	list := make([]*entity.RequesterGlobalHeader, 0)
	err := conn.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(RequestGlobalHeaderTable))
		if b == nil {
			return ISNULL
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			data := &entity.RequesterGlobalHeader{}
			//log.Info("k = ", string(k))
			err := json.Unmarshal(v, &data)
			if err != nil {
				return err
			}
			list = append(list, data)
		}
		return nil
	})
	return list, err
}

func (r *requestToolDao) DelGlobalHeader(key string) error {
	return DB.Delete(RequestGlobalHeaderTable, key)
}
