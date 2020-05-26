package db

const (
	ModelLoad = 1
	ModelSave = 2
)

type ModelMeta struct {
	modelPtr ModelClass
	op       int
	ctx      interface{}
	e        error
}

func (self *ModelMeta) WithContext(ctx interface{}) {
	self.ctx = ctx
}

func (self *ModelMeta) Context() interface{} {
	return self.ctx
}

func (self *ModelMeta) Model() ModelClass {
	return self.modelPtr
}

func (self *ModelMeta) Error() error {
	return self.e
}
