// log.go
package parse

import (
	"log"
	"runtime"
)

func trap_error() func() {
	return func() {
		if err := recover(); err != nil {
			log.Printf("%v", err)
			for i := 0; i < 10; i++ {
				funcName, file, line, ok := runtime.Caller(i)
				if ok {
					log.Printf(" frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
				}
			}
		}
	}
}
