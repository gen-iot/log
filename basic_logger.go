package log

import "fmt"

type basicLogger struct {
}

func (this *basicLogger) Print(v ...interface{}) {
	fmt.Print(v...)
}
func (this *basicLogger) Println(v ...interface{}) {
	fmt.Println(v...)
}
func (this *basicLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
