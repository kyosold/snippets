package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	xid "github.com/rs/xid"
)

func ReadLinesAllContentWithFile(filename string) ([]string, error) {
	return ReadLinesOffsetN(filename, 0, -1)
}

func ReadLinesOffsetN(filename string, offset uint, n int) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []string{""}, err
	}
	defer f.Close()

	var ret []string
	reader := bufio.NewReader(f)
	for i := 0; (i < (n + int(offset))) || (n < 0); i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return ret, err
			}
		}

		if i < int(offset) {
			continue
		}

		ret = append(ret, strings.Trim(line, "\n"))
	}

	return ret, nil
}

func getUUID() string {
	guid := xid.New()
	return guid.String()
}

func kscal(size int64) string {
	var GB int64 = 1024 * 1024 * 1024
	var MB int64 = 1024 * 1024
	var KB int64 = 1024

	var ret string
	if size > GB {
		ret = fmt.Sprintf("%.2f GB", float64(size/GB))
	} else if size > MB {
		ret = fmt.Sprintf("%.2f MB", float64(size/MB))
	} else if size > KB {
		ret = fmt.Sprintf("%.2f KB", float64(size/KB))
	} else {
		ret = fmt.Sprintf("%.2f B", size)
	}
	return ret
}
