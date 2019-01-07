package main

import (
	"github.com/davyxu/cellmesh/discovery/memsd/model"
	"os"
	"time"
)

func loadPersistFile(fileName string) {

	fileHandle, err := os.OpenFile(fileName, os.O_RDONLY, 0666)

	// 可能文件不存在，忽略
	if err != nil {
		return
	}

	log.Infoln("Load values...")

	err = model.LoadValue(fileHandle)
	if err != nil {
		log.Errorf("load values failed: %s %s", fileName, err.Error())
		return
	}

	log.Infof("Load %d values", model.ValueCount())
}

func startPersistCheck(fileName string) {

	ticker := time.NewTicker(time.Minute)

	for {

		<-ticker.C

		// 与收发在一个队列中，保证无锁
		model.Queue.Post(func() {

			if model.ValueDirty {

				log.Infoln("Save values...")

				fileHandle, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					log.Errorf("save persist file failed: %s %s", fileName, err.Error())
					return
				}

				err = model.SaveValue(fileHandle)

				if err != nil {
					log.Errorf("save values failed: %s %s", fileName, err.Error())
					return
				}

				log.Infof("Save %d values", model.ValueCount())

				model.ValueDirty = false

			}

		})

	}

}
