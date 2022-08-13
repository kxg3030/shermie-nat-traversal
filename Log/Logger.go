package Log

import (
	"log"
	"os"
	"sync"
)

var Log *log.Logger

type Logger struct {
	once *sync.Once
}

func NewLogger() *Logger {
	return &Logger{once: &sync.Once{}}
}

func (i *Logger) Initialize() {
	i.once.Do(func() {
		Log = log.New(os.Stdout, "", log.LstdFlags)
	})
}
