package db

import (
	"errors"
	"github.com/davyxu/ulog"
	"github.com/gomodule/redigo/redis"
	"github.com/vmihailenco/msgpack"
)

var (
	ErrModelNotExists = errors.New("model not exists")
)

const (
	Code_DBAccessFailed  = 101 // 与协议的错误号对应,避免引用, 减少关联
	Code_SerializeFailed = 102 // 序列化错误
)

type ModelList struct {
	conn     redis.Conn
	metaList []*ModelMeta
}

func (self *ModelList) Save(modelPtr ModelClass, rawKeyID interface{}) error {

	data, err := msgpack.Marshal(modelPtr)
	if err != nil {
		ulog.Errorf("save marshal failed, %s", err)
		panic(Code_SerializeFailed)
	}

	key := genKey(modelPtr.DBKey(), rawKeyID)

	_, err = self.conn.Do("SET", key, data)

	return err
}

func (self *ModelList) Load(modelPtr ModelClass, rawKeyID interface{}) error {

	key := genKey(modelPtr.DBKey(), rawKeyID)

	data, err := recvData(self.conn.Do("GET", key))
	if err != nil {
		return err
	}

	err = msgpack.Unmarshal(data, modelPtr)
	if err != nil {
		ulog.Errorf("load model failed, %s", err)
		panic(Code_SerializeFailed)
	}

	return nil
}

func (self *ModelList) batchOp(modelPtr ModelClass, rawKeyID interface{}, op int) *ModelMeta {
	key := genKey(modelPtr.DBKey(), rawKeyID)

	var err error
	switch op {
	case ModelSave:
		var data []byte
		data, err = msgpack.Marshal(modelPtr)
		if err != nil {
			ulog.Errorf("batchsave marshal failed, %s", err)
			panic(Code_SerializeFailed)
		}

		err = self.conn.Send("SET", key, data)
	case ModelLoad:
		err = self.conn.Send("GET", key)
	}

	if err != nil {
		ulog.Errorf("batch failed, %s", err)
		panic(Code_DBAccessFailed)
	}

	meta := ModelMeta{
		modelPtr: modelPtr,
		op:       op,
	}

	self.metaList = append(self.metaList, &meta)

	return &meta
}

func (self *ModelList) BatchLoad(modelPtr ModelClass, rawKeyID interface{}) *ModelMeta {
	return self.batchOp(modelPtr, rawKeyID, ModelLoad)
}

func (self *ModelList) BatchSave(modelPtr ModelClass, rawKeyID interface{}) *ModelMeta {
	return self.batchOp(modelPtr, rawKeyID, ModelSave)
}

func recvData(reply interface{}, err error) ([]byte, error) {

	data, err := redis.Bytes(reply, err)
	if err == redis.ErrNil {
		return nil, ErrModelNotExists
	} else if err != nil {
		ulog.Errorf("load model failed, %s", err)
		panic(Code_DBAccessFailed)
	}

	if len(data) == 0 {
		return nil, ErrModelNotExists
	}

	return data, nil
}

func (self *ModelList) Flush() {

	self.conn.Flush()

	for _, meta := range self.metaList {

		data, err := recvData(self.conn.Receive())
		meta.e = err

		if err != nil {
			continue
		}

		if meta.op == ModelLoad {
			err = msgpack.Unmarshal(data, meta.modelPtr)
			if err != nil {
				ulog.Errorf("flush unmarshal failed, %s", err)
				panic(Code_SerializeFailed)
			}
		}
	}
}

func (self *ModelList) IndexMeta(index int) *ModelMeta {
	return self.metaList[index]
}

func (self *ModelList) Index(index int) (ModelClass, error) {
	meta := self.metaList[index]
	return meta.modelPtr, meta.e
}

func (self *ModelList) MustIndex(index int) ModelClass {
	meta := self.metaList[index]

	if meta.e == nil {
		return meta.modelPtr
	}

	return nil
}

func (self *ModelList) Count() int {
	return len(self.metaList)
}

func NewModelList(conn redis.Conn) *ModelList {
	return &ModelList{conn: conn}
}
