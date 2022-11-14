package logger

import (
	"fmt"
	"log"
)

var _ Logger = (*Default)(nil)

const (
	DebugLvl Level = iota
	InfoLvl
	WarnLvl
	ErrorLvl
	FatalLvl
)

type Default struct {
	lvl Level
}

func (d *Default) Debug(args ...any) {
	d.print(DebugLvl, args...)
}

func (d *Default) Debugf(format string, args ...any) {
	d.printf(DebugLvl, format, args...)
}

func (d *Default) Info(args ...any) {
	d.print(InfoLvl, args...)
}

func (d *Default) Infof(format string, args ...any) {
	d.printf(InfoLvl, format, args...)
}

func (d *Default) Warn(args ...any) {
	d.print(WarnLvl, args...)
}

func (d *Default) Warnf(format string, args ...any) {
	d.printf(WarnLvl, format, args...)
}

func (d *Default) Error(args ...any) {
	d.print(ErrorLvl, args...)
}

func (d *Default) Errorf(format string, args ...any) {
	d.printf(ErrorLvl, format, args...)
}

func (d *Default) Fatal(args ...any) {
	d.print(FatalLvl, args...)
}

func (d *Default) Fatalf(format string, args ...any) {
	d.printf(FatalLvl, format, args...)
}

func (d *Default) print(lvl Level, args ...any) {
	if d.lvl > lvl {
		return
	}
	if lvl == FatalLvl {
		log.Fatal(args...)
		return
	}
	log.Println(args...)
}

func (d *Default) printf(lvl Level, format string, args ...any) {
	d.print(lvl, fmt.Sprintf(format, args...))
}
