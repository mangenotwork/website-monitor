package dao

import (
	"fmt"
)

type History interface {
	Set(value any) error
	Get() ([]any, error)
	Clear() error
}

const (
	FidWebsiteTDKI = 1 // 获取网站的T, D, K, 图标
	FidIp          = 2 // ip信息查询
	FidNsLookUp    = 3 // 查询dns
	FidWhois       = 4 // Whois查询
	FidICP         = 5 // 查询备案
	FidPing        = 6 // 在线ping
	FidSSL         = 7 // 获取证书
	FidWebsiteInfo = 8 // 网站信息获取
)

func NewHistory(toolId int) (History, error) {
	switch toolId {
	case FidWebsiteTDKI:
		return new(historyWebsiteTDKI), nil
	case FidIp:
		return new(historyIP), nil
	case FidNsLookUp:
		return new(historyNsLookUp), nil
	case FidWhois:
		return new(historyWhois), nil
	case FidICP:
		return new(historyICP), nil
	case FidPing:
		return new(historyPing), nil
	case FidSSL:
		return new(historySSL), nil
	case FidWebsiteInfo:
		return new(historyWebsiteInfo), nil
	}
	
	return nil, fmt.Errorf("没有这个工具的历史记录")
}

func sizeValue(list []any, i any) []any {
	if len(list) > 100 {
		list = list[1:len(list)]
	}
	list = append(list, i)
	return list
}

type historyWebsiteTDKI struct{}

func (h *historyWebsiteTDKI) Set(value any) error {
	list, _ := h.Get()
	list = sizeValue(list, value)
	return DB.Set(HistoryWebsiteTDKITable, HistoryWebsiteTDKIKeyName, list)
}

func (h *historyWebsiteTDKI) Get() ([]any, error) {
	value := make([]any, 0)
	err := DB.Get(HistoryWebsiteTDKITable, HistoryWebsiteTDKIKeyName, &value)
	return value, err
}

func (h *historyWebsiteTDKI) Clear() error {
	list := make([]any, 0)
	return DB.Set(HistoryWebsiteTDKITable, HistoryWebsiteTDKIKeyName, list)
}

type historyIP struct{}

func (h *historyIP) Set(value any) error {
	list, _ := h.Get()
	list = sizeValue(list, value)
	return DB.Set(HistoryIpTable, HistoryIpKeyName, list)
}

func (h *historyIP) Get() ([]any, error) {
	value := make([]any, 0)
	err := DB.Get(HistoryIpTable, HistoryIpKeyName, &value)
	return value, err
}

func (h *historyIP) Clear() error {
	list := make([]any, 0)
	return DB.Set(HistoryIpTable, HistoryIpKeyName, list)
}

type historyNsLookUp struct{}

func (h *historyNsLookUp) Set(value any) error {
	list, _ := h.Get()
	list = sizeValue(list, value)
	return DB.Set(HistoryNsLookUpTable, HistoryNsLookUpKeyName, list)
}

func (h *historyNsLookUp) Get() ([]any, error) {
	value := make([]any, 0)
	err := DB.Get(HistoryNsLookUpTable, HistoryNsLookUpKeyName, &value)
	return value, err
}

func (h *historyNsLookUp) Clear() error {
	list := make([]any, 0)
	return DB.Set(HistoryNsLookUpTable, HistoryNsLookUpKeyName, list)
}

type historyWhois struct{}

func (h *historyWhois) Set(value any) error {
	list, _ := h.Get()
	list = sizeValue(list, value)
	return DB.Set(HistoryWhoisTable, HistoryWhoisKeyName, list)
}

func (h *historyWhois) Get() ([]any, error) {
	value := make([]any, 0)
	err := DB.Get(HistoryWhoisTable, HistoryWhoisKeyName, &value)
	return value, err
}

func (h *historyWhois) Clear() error {
	list := make([]any, 0)
	return DB.Set(HistoryWhoisTable, HistoryWhoisKeyName, list)
}

type historyICP struct{}

func (h *historyICP) Set(value any) error {
	list, _ := h.Get()
	list = sizeValue(list, value)
	return DB.Set(HistoryICPTable, HistoryICPKeyName, list)
}

func (h *historyICP) Get() ([]any, error) {
	value := make([]any, 0)
	err := DB.Get(HistoryICPTable, HistoryICPKeyName, &value)
	return value, err
}

func (h *historyICP) Clear() error {
	list := make([]any, 0)
	return DB.Set(HistoryICPTable, HistoryICPKeyName, list)
}

type historyPing struct{}

func (h *historyPing) Set(value any) error {
	list, _ := h.Get()
	list = sizeValue(list, value)
	return DB.Set(HistoryPingTable, HistoryPingKeyName, list)
}

func (h *historyPing) Get() ([]any, error) {
	value := make([]any, 0)
	err := DB.Get(HistoryPingTable, HistoryPingKeyName, &value)
	return value, err
}

func (h *historyPing) Clear() error {
	list := make([]any, 0)
	return DB.Set(HistoryPingTable, HistoryPingKeyName, list)
}

type historySSL struct{}

func (h *historySSL) Set(value any) error {
	list, _ := h.Get()
	list = sizeValue(list, value)
	return DB.Set(HistorySSLTable, HistorySSLKeyName, list)
}

func (h *historySSL) Get() ([]any, error) {
	value := make([]any, 0)
	err := DB.Get(HistorySSLTable, HistorySSLKeyName, &value)
	return value, err
}

func (h *historySSL) Clear() error {
	list := make([]any, 0)
	return DB.Set(HistorySSLTable, HistorySSLKeyName, list)
}

type historyWebsiteInfo struct{}

func (h *historyWebsiteInfo) Set(value any) error {
	list, _ := h.Get()
	list = sizeValue(list, value)
	return DB.Set(HistoryWebsiteInfoTable, HistoryWebsiteInfoKeyName, list)
}

func (h *historyWebsiteInfo) Get() ([]any, error) {
	value := make([]any, 0)
	err := DB.Get(HistoryWebsiteInfoTable, HistoryWebsiteInfoKeyName, &value)
	return value, err
}

func (h *historyWebsiteInfo) Clear() error {
	return DB.Set(HistoryWebsiteInfoTable, HistoryWebsiteInfoKeyName, make([]any, 0))
}
