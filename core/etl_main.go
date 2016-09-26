package core

import (
	"fmt"
)

type ETLMain struct {
}

func (e *ETLMain) fileDeleted(file string) {
	fmt.Println("fileDeleted ", file)
}

func (e *ETLMain) fileCreated(file string) {
	fmt.Println("fileCreated ", file)
}
