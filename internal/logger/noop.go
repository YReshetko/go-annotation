package logger

var _ Logger = (*Noop)(nil)

type Noop struct{}

func (n Noop) Debug(...any)          {}
func (n Noop) Debugf(string, ...any) {}

func (n Noop) Info(...any)          {}
func (n Noop) Infof(string, ...any) {}

func (n Noop) Warn(...any)          {}
func (n Noop) Warnf(string, ...any) {}

func (n Noop) Error(...any)          {}
func (n Noop) Errorf(string, ...any) {}

func (n Noop) Fatal(...any)          {}
func (n Noop) Fatalf(string, ...any) {}
