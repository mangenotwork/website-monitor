package dao

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"sort"
	"time"
	"website-monitor/master/entity"
)

type RequestToolEr interface {
	Add(data *entity.RequestTool) error
	GetAtID(id string) (*entity.RequestTool, error)
	History() ([]*entity.RequestTool, error)
	HistoryDelete(id string) error
	SetGlobalHeader(list []*entity.RequesterGlobalHeader) error
	GetGlobalHeader() ([]*entity.RequesterGlobalHeader, error)
	DelGlobalHeader(key string) error
	SetRequestNowList(data *entity.RequestNowList) error
	GetRequestNowList() ([]*entity.RequestNowList, error)
}

func NewRequestTool() RequestToolEr {
	return new(requestToolDao)
}

type requestToolDao struct{}

func (r *requestToolDao) Add(data *entity.RequestTool) error {
	err := DB.Set(RequestTable, data.ID, data)
	if err != nil {
		return err
	}
	return r.SetRequestNowList(&entity.RequestNowList{
		Id:     data.ID,
		Method: data.Method,
		Url:    data.Url,
		Name:   data.Name,
		Time:   time.Now().Unix(),
	})
}

func (r *requestToolDao) GetAtID(id string) (*entity.RequestTool, error) {
	data := &entity.RequestTool{}
	err := DB.Get(RequestTable, id, &data)
	return data, err
}

func (r *requestToolDao) History() ([]*entity.RequestTool, error) {
	conn := GetDBConn()
	defer func() {
		_ = conn.Close()
	}()
	list := make([]*entity.RequestTool, 0)
	err := conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RequestTable))
		if b == nil {
			return ISNULL
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			data := &entity.RequestTool{}
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

func (r *requestToolDao) HistoryDelete(id string) error {
	return DB.Delete(RequestTable, id)
}

func (r *requestToolDao) SetRequestNowList(data *entity.RequestNowList) error {
	return DB.Set(RequestNowListTable, data.Id, data)
}

func (r *requestToolDao) GetRequestNowList() ([]*entity.RequestNowList, error) {
	conn := GetDBConn()
	defer func() {
		_ = conn.Close()
	}()
	list := make([]*entity.RequestNowList, 0)
	err := conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RequestNowListTable))
		if b == nil {
			return ISNULL
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			data := &entity.RequestNowList{}
			err := json.Unmarshal(v, &data)
			if err != nil {
				return err
			}
			list = append(list, data)
		}
		return nil
	})
	sort.Slice(list, func(i, j int) bool {
		if list[i].Time > list[j].Time {
			return true
		}
		return false
	})
	list[0].IsNow = true
	return list, err
}

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
