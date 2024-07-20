package logger

import (
	"database/sql"
	"fmt"
	"log"
)

type Logger struct {
	Caller string
}

func (l Logger) Fatalln(err error) {
	msg := fmt.Sprintf("[-] Error in %s\n", l.Caller)
	log.Fatalln(msg, err)
}

func (l Logger) Println(err error) {
	msg := fmt.Sprintf("[-] Error in %s\n", l.Caller)
	log.Println(msg, err)
}

func (l Logger) AppendError(err error) error {
	if err != nil && err != sql.ErrNoRows {
		err = fmt.Errorf("[-] Error in %s\n%w", l.Caller, err)
	}
	return err
}
