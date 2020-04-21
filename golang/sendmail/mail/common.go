package mail

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	xid "github.com/rs/xid"
)

func CreateDateString() string {
	// 星期
	now := time.Now()

	weekday := fmt.Sprintf("%s", now.Weekday())
	day := now.Day()
	month := fmt.Sprintf("%s", now.Month())
	year := now.Year()
	time := fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())

	hdate := fmt.Sprintf("%s, %02d %s %d %s", weekday[:3], day, month[:3], year, time)

	return hdate
}

func CreateMessageID() string {
	guid := xid.New()

	return guid.String()
}

func ConvStringToMimeString(str string) string {
	b64Str := base64.StdEncoding.EncodeToString([]byte(str))
	res := "=?UTF8?B?" + b64Str + "?="
	return res
}

func getAttachConextMimeWithFile(attach string) (string, error) {
	fileInfo, err := os.Stat(attach)
	if err != nil {
		return "", err
	}
	fileSize := fileInfo.Size()

	fp, err := os.Open(attach)
	if err != nil {
		return "", err
	}
	defer fp.Close()

	bufBytes := make([]byte, fileSize)
	_, err = fp.Read(bufBytes)
	if err != nil {
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(bufBytes)

	maxBytes := 76
	res := make([]byte, 0)
	x := 0
	for i := 0; i < len(b64); i++ {
		if x >= maxBytes {
			res = append(res, byte('\r'))
			res = append(res, byte('\n'))
			x = 0
		}
		res = append(res, b64[i])
		x++
	}

	return string(res), nil
}
