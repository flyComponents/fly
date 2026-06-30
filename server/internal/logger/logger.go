package logger

import (
	"io"
	"log"
	"os"
)

func New() *log.Logger {

	file, err := os.OpenFile(
		"server.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, file)

	return log.New(mw, "", log.LstdFlags|log.Lshortfile)
}
