package helpers

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func Error(e error) {
	fmt.Printf("%s\n", e.Error())
	syscall.Kill(os.Getpid(), syscall.SIGTERM)
}

func TimestampToString (ts *timestamp.Timestamp) string {
	timestamp, err := ptypes.Timestamp(ts) 
	if err != nil {
		timestamp = time.Now()
	}
	return timestamp.In(time.Local).Format("03:04:05 PM")
}