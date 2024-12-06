package dao

import (
	"bufio"
	"fmt"
	"github.com/mangenotwork/common/conf"
	"github.com/mangenotwork/common/log"
	"github.com/mangenotwork/common/utils"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"website-monitor/master/constname"
	"website-monitor/master/entity"
)

type MonitorLogEr interface {
	Write(hostId, mLog string)
	ReadLog(hostId, day string) []*entity.MonitorLog
	ToMonitorLogObj(str string) *entity.MonitorLog
	DataFormat(mLog *entity.MonitorLog) string
	DeleteLog(hostId string) error
	ReadAll(hostId, day string) ([]*entity.MonitorLog, error)
	LogList(hostId string) ([]string, error)
	LogListDay(hostId string) ([]string, error)
	Upload(hostId, day string) (string, error)
}

func NewMonitorLogDao() MonitorLogEr {
	return new(monitorLogDao)
}

type monitorLogDao struct {
}

// 写日志
func (m *monitorLogDao) Write(hostId, mLog string) {
	logPath, err := conf.YamlGetString("logPath")
	if err != nil {
		logPath = "./log/"
	}

	fileName := logPath + hostId + "_" + utils.NowDateLayout(constname.DayLayout) + ".log"
	var file *os.File

	if !utils.Exists(fileName) {
		_ = os.MkdirAll(logPath, 0666)
		file, _ = os.Create(fileName)
	} else {
		file, _ = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	}

	defer func() {
		_ = file.Close()
	}()

	_, err = io.WriteString(file, mLog+"\n")
	if err != nil {
		log.Error("写入日志错误：", err)
		return
	}

}

func (m *monitorLogDao) DataFormat(mLog *entity.MonitorLog) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s|%d|%d|%s|%d|%d|%s|%d|%s|%s|%s|%s|%s|",
		mLog.LogType, mLog.Time, mLog.HostId, mLog.Host, mLog.UriType, mLog.Uri, mLog.UriCode, mLog.UriMs,
		mLog.ContrastUri, mLog.ContrastUriCode, mLog.ContrastUriMs, mLog.Ping, mLog.PingMs, mLog.Msg, mLog.AlertType,
		mLog.MonitorName, mLog.MonitorIP, mLog.MonitorAddr)
}

func (m *monitorLogDao) ReadLog(hostId, day string) []*entity.MonitorLog {
	logPath, err := conf.YamlGetString("logPath")
	if err != nil {
		logPath = "./log/"
	}

	if day == "" {

		list, err := m.LogListDay(hostId)
		if err != nil || len(list) < 1 {
			log.Info(err)
			return make([]*entity.MonitorLog, 0)
		}

		day = list[0]
	}

	fileName := logPath + hostId + "_" + day + ".log"
	log.Info("fileName = ", fileName)

	f, err := os.Open(fileName)
	if err != nil {
		log.ErrorF("open file error:%s", err.Error())
	}

	defer func() {
		_ = f.Close()
	}()

	data := make([]*entity.MonitorLog, 0)
	buff := make([]byte, 0, 4096)
	char := make([]byte, 1)
	stat, _ := f.Stat()
	filesize := stat.Size()
	cursor := 0
	count := 0
	maxCount := 300

	for {
		cursor -= 1
		_, _ = f.Seek(int64(cursor), io.SeekEnd)
		_, err = f.Read(char)
		if err != nil {
			log.Error(err)
			break
		}

		if char[0] == '\n' {
			if len(buff) > 0 {
				revers(buff)
				d := m.ToMonitorLogObj(string(buff))
				if d != nil {
					data = append(data, d)
				}

				count++
				if count == maxCount {
					break
				}
			}

			buff = buff[:0]

		} else {
			buff = append(buff, char[0])
		}

		if int64(cursor) == -filesize {
			break
		}

	}

	return data
}

func revers(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func (m *monitorLogDao) ToMonitorLogObj(str string) *entity.MonitorLog {
	strList := strings.Split(str, "|")

	if len(strList) < 18 {
		return nil
	}

	return &entity.MonitorLog{
		LogType:         strList[0],
		Time:            strList[1],
		HostId:          strList[2],
		Host:            strList[3],
		UriType:         strList[4],
		Uri:             strList[5],
		UriCode:         utils.AnyToInt(strList[6]),
		UriMs:           utils.AnyToInt64(strList[7]),
		ContrastUri:     strList[8],
		ContrastUriCode: utils.AnyToInt(strList[9]),
		ContrastUriMs:   utils.AnyToInt64(strList[10]),
		Ping:            strList[11],
		PingMs:          utils.AnyToInt64(strList[12]),
		Msg:             strList[13],
		AlertType:       strList[14],
		MonitorName:     strList[15],
		MonitorIP:       strList[16],
		MonitorAddr:     strList[17],
		ClientIP:        strList[18],
	}
}

// DeleteLog 删除日志
func (m *monitorLogDao) DeleteLog(hostId string) error {
	logPath, err := conf.YamlGetString("logPath")
	if err != nil {
		logPath = "./log/"
	}

	return filepath.Walk(logPath, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return err
		}

		fileName := info.Name()
		fid := strings.Split(fileName, "_")
		if len(fid) > 0 && fid[0] == hostId {
			log.Info("fileName = ", fileName, path)
			err = os.Remove(path)
			if err != nil {
				log.Error(err)
			}
		}

		return err
	})
}

func (m *monitorLogDao) ReadAll(hostId, day string) ([]*entity.MonitorLog, error) {
	logPath, err := conf.YamlGetString("logPath")
	if err != nil {
		logPath = "./log/"
	}

	filePath := logPath + hostId + "_" + day + ".log"
	data := make([]*entity.MonitorLog, 0)
	log.Info("filePath = ", filePath)

	f, err := os.Open(filePath)
	if err != nil {
		return data, err
	}

	defer func() {
		_ = f.Close()
	}()

	r := bufio.NewReader(f)

	for {
		line, e := r.ReadBytes('\n')
		if e == nil {

			d := m.ToMonitorLogObj(string(line))
			if d != nil {
				data = append(data, d)
			}
		}

		if e != nil && e != io.EOF {
			log.Error(e)
			err = e
		}

		if e == io.EOF {
			break
		}

	}

	return data, err
}

func (m *monitorLogDao) LogList(hostId string) ([]string, error) {
	logPath, err := conf.YamlGetString("logPath")
	if err != nil {
		logPath = "./log/"
	}

	list := make([]string, 0)
	err = filepath.Walk(logPath, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return err
		}

		fileName := info.Name()
		fid := strings.Split(fileName, "_")
		if len(fid) > 0 && fid[0] == hostId {
			log.Info("fileName = ", fileName, path)
			list = append(list, fileName)
		}

		return err
	})

	return list, err
}

func (m *monitorLogDao) LogListDay(hostId string) ([]string, error) {
	logPath, err := conf.YamlGetString("logPath")
	if err != nil {
		logPath = "./log/"
	}

	list := make([]string, 0)

	err = filepath.Walk(logPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return err
		}

		fileName := info.Name()
		fid := strings.Split(fileName, "_")
		if len(fid) > 1 && fid[0] == hostId {
			list = append(list, strings.Replace(fid[1], ".log", "", -1))
		}

		return err
	})

	sort.Slice(list, func(i, j int) bool {
		return i > j
	})

	return list, err
}

func (m *monitorLogDao) Upload(hostId, day string) (string, error) {
	logPath, err := conf.YamlGetString("logPath")
	if err != nil {
		logPath = "./log/"
	}

	filePath := logPath + hostId + "_" + day + ".log"
	if utils.Exists(filePath) {
		return filePath, nil
	}

	return "", fmt.Errorf("日志不存在")
}
