package log

type EmptyLogger struct{}

func (l *EmptyLogger) Print(v ...interface{})                 {}
func (l *EmptyLogger) Println(v ...interface{})               {}
func (l *EmptyLogger) Printf(format string, v ...interface{}) {}
