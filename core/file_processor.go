package core

type Processor interface {
}
type FileToProcessor interface {
	GetProcessor(filename string) Processor
}

type TransLog struct {
}
