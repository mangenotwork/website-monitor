package dao

import (
	"bufio"
	"bytes"
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
	f, err := os.Open(fileName)
	if err != nil {
		log.ErrorF("open file error:%s", err.Error())
		return make([]*entity.MonitorLog, 0)
	}
	defer func() {
		_ = f.Close()
	}()

	stat, err := f.Stat()
	if err != nil || stat.Size() == 0 {
		return make([]*entity.MonitorLog, 0)
	}

	// 从文件尾部按块倒读，凑够 maxCount 行即可，避免逐字节 Seek（1MB 文件约 15 万次系统调用）
	const maxCount = 300
	const chunkSize = 32 * 1024

	pos := stat.Size()
	buf := make([]byte, 0, chunkSize)
	for pos > 0 && bytes.Count(buf, []byte{'\n'}) <= maxCount {
		n := int64(chunkSize)
		if pos < n {
			n = pos
		}
		pos -= n

		chunk := make([]byte, n)
		if _, err := f.ReadAt(chunk, pos); err != nil && err != io.EOF {
			log.Error(err)
			break
		}
		buf = append(chunk, buf...)
	}

	lines := bytes.Split(buf, []byte{'\n'})
	data := make([]*entity.MonitorLog, 0, maxCount)
	for i := len(lines) - 1; i >= 0 && len(data) < maxCount; i-- {
		line := bytes.TrimSpace(lines[i])
		if len(line) == 0 {
			continue
		}
		if d := m.ToMonitorLogObj(string(line)); d != nil {
			data = append(data, d)
		}
	}
	return data
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

func (m *monitorLogDao) DeleteLog(hostId string) error {
	var firstErr error
	for _, path := range m.logFiles(hostId) {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			log.Error(err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

// logFiles 只匹配指定站点的日志文件，避免 Walk 整个日志目录
func (m *monitorLogDao) logFiles(hostId string) []string {
	logPath, err := conf.YamlGetString("logPath")
	if err != nil {
		logPath = "./log/"
	}
	matches, _ := filepath.Glob(filepath.Join(logPath, hostId+"_*.log"))
	return matches
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
		if len(line) > 0 {
			if d := m.ToMonitorLogObj(strings.TrimRight(string(line), "\r\n")); d != nil {
				data = append(data, d)
			}
		}
		if e != nil {
			if e != io.EOF {
				log.Error(e)
				err = e
			}
			break
		}
	}

	return data, err
}

func (m *monitorLogDao) LogList(hostId string) ([]string, error) {
	list := make([]string, 0)
	for _, path := range m.logFiles(hostId) {
		list = append(list, filepath.Base(path))
	}
	sort.Strings(list)
	return list, nil
}

func (m *monitorLogDao) LogListDay(hostId string) ([]string, error) {
	list := make([]string, 0)
	for _, path := range m.logFiles(hostId) {
		parts := strings.SplitN(filepath.Base(path), "_", 2)
		if len(parts) < 2 {
			continue
		}
		list = append(list, strings.TrimSuffix(parts[1], ".log"))
	}
	sort.Slice(list, func(i, j int) bool { return list[i] > list[j] })
	return list, nil
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
