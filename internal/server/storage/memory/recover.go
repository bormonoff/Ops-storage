package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ops-storage/internal/server/logger"
)

func InitRecover(relPath string, interval int, restore bool) {
	path, err := createPath(relPath)
	if err != nil {
		logger.MainLog.Errorf("Create recover file error, continue without saving: %s", err.Error())
	}
	store.recFilePath = path
	if restore {
		store.tryLoad()
	}
	store.runDump(interval)
}

func (s *memStorage) storeToFile() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sRepr, err := json.Marshal(s)
	if err != nil {
		logger.MainLog.Errorf("Store file error: %s", err.Error())
	}

	file, err := os.OpenFile(s.recFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.MainLog.Errorf("Open file error: %s", err.Error())
	}
	defer file.Close()
	file.Write(sRepr)
}

func (s *memStorage) tryLoad() {
	tFile, err := os.OpenFile(s.recFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.MainLog.Errorf("Skip recover, load recover file error: %s", err.Error())
	}

	defer tFile.Close()

	fileinfo, err := tFile.Stat()
	if err != nil {
		logger.MainLog.Errorf("Get file stat error: %s", err.Error())
	}

	var buff = make([]byte, fileinfo.Size())
	_, err = tFile.Read(buff)
	if err != nil {
		logger.MainLog.Errorf("Read file error: %s", err.Error())
	}

	err = json.Unmarshal(buff, s)
	if err != nil {
		logger.MainLog.Errorf("Skip recover, parse recover file error: %s", err.Error())
	}
}

func (s *memStorage) runDump(interval int) {
	err := createFile(s.recFilePath)
	if err != nil {
		logger.MainLog.Errorf("Skip recover, create file error: %s", err.Error())
		return
	}

	go func() {
		for {
			time.Sleep(time.Duration(interval) * time.Second)
			logger.MainLog.Info("Saving data to the file...")
			s.storeToFile()
			logger.MainLog.Info("Saving complete.")
		}
	}()
}

func createPath(relPath string) (string, error) {
	var pth strings.Builder

	ex, _ := os.Executable()
	pth.WriteString(filepath.Dir(ex))

	dir, fName := filepath.Split(relPath)

	for _, fld := range strings.Split(dir, "/") {
		if fld == "" {
			continue
		}
		os.Mkdir(fmt.Sprintf("%s/%s", pth.String(), fld), 0744)
		pth.WriteString("/" + fld)
	}
	pth.WriteString("/" + fName)

	return pth.String(), nil
}

func createFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}
