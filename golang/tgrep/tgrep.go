package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Matcher struct {
	StartHour    int
	EndHour      int
	IgnoreCase   bool
	Verbose      bool
	Pattern      string
	File         string
	ContentType  string
	ShowFileName bool
}

type LineTime struct {
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
}

func usage(prog string) {
	fmt.Println("Usage:")
	fmt.Printf(" %s -st=15 -et=22 [pattern] [file]\n", prog)
	fmt.Println("options:")
	fmt.Printf("  -st 指定的小时时间开始查找\n")
	fmt.Printf("  -et 指定的小时时间结束查找, 不写默认到文件结尾\n")
	fmt.Printf("  -i 不区分大小写\n")
	fmt.Printf("  -v 显示详细信息\n")
	fmt.Println("Others:")
	fmt.Printf("  1. 如果文件是gzip，先解压再查询，如:\n")
	fmt.Printf("    gzip -dvc abc.0.gz > abc.0\n")
	fmt.Printf("    tgrep -st=8 -et=10 'pattern' abc.0\n")
	fmt.Println()
}

func main() {
	var pattern string
	var dir []string

	if len(os.Args) < 2 {
		usage(filepath.Base(os.Args[0]))
		return
	}

	startTime := flag.Int("st", -1, "Start Hour")
	endTime := flag.Int("et", -1, "End Hour")
	ignoreCase := flag.Uint("i", 0, "Ignore Case")
	verbose := flag.Uint("v", 0, "Verbose")

	flag.Parse()

	startIdx := 1
	if *startTime > 0 {
		startIdx++
	} else {
		usage(filepath.Base(os.Args[0]))
		return
	}
	if *endTime > 0 {
		startIdx++
	}
	if *ignoreCase > 0 {
		startIdx++
	}
	if *verbose > 0 {
		startIdx++
	}

	pattern = os.Args[startIdx]
	dir = os.Args[startIdx+1:]

	fmt.Println("****************************")
	fmt.Println("Start Hour:", *startTime)
	fmt.Println("End Hour:", *endTime)
	fmt.Println("Ignore Case:", *ignoreCase)
	fmt.Println("Verbose:", *verbose)
	fmt.Println("Pattern:", pattern)
	fmt.Println("Files:", dir)
	fmt.Println("****************************")

	mr := Matcher{}
	mr.StartHour = *startTime
	mr.EndHour = *endTime
	mr.Pattern = pattern
	mr.ShowFileName = false
	if *verbose == 1 {
		mr.Verbose = true
	} else {
		mr.Verbose = false
	}
	if len(dir) > 1 {
		mr.ShowFileName = true
	}
	if *ignoreCase == 1 {
		mr.IgnoreCase = true
		mr.Pattern = strings.ToLower(pattern)
	} else {
		mr.IgnoreCase = false
	}

	var wg sync.WaitGroup

	for _, file := range dir {
		mr.File = file

		ct, err := GetFileContentType(mr.File)
		if err != nil {
			fmt.Println(file, " [Error]:", err)
			continue
		}
		mr.ContentType = ct

		ctlist := strings.Split(ct, "/")
		cType := strings.ToLower(ctlist[0])
		cSubType := strings.ToLower(ctlist[1])

		if cType == "text" {
			wg.Add(1)
			go searchV2(mr, &wg)
		} else if cType == "application" {
			if cSubType == "x-gzip" {
				fmt.Printf("[Error]: file:%s Is GZIP file\n", mr.File)
				fmt.Printf("  Run 'gzip -dvc %s > tmp.log'\n", mr.File)
			}
		}

	}

	wg.Wait()
}

func searchV2(m Matcher, wg *sync.WaitGroup) {
	defer wg.Done()

	gid := GetUUID()
	if m.Verbose {
		fmt.Printf("%s: %s [%s]\n", gid, m.File, m.ContentType)
	}

	spos, slineTime, err := getPos(m.File, m.StartHour, 0, m.Verbose, gid)
	if err != nil {
		fmt.Println(err)
		return
	}

	var epos int64 = 0
	var elineTime LineTime
	if m.EndHour > 0 {
		if m.Verbose {
			fmt.Println()
		}
		// eHourInt, _ := m.EndHour
		epos, elineTime, err = getPos(m.File, m.EndHour, 1, m.Verbose, gid)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fileInfo, err := os.Stat(m.File)
		if err != nil {
			fmt.Println(err)
			return
		}
		epos = fileInfo.Size()
	}

	fp, err := os.Open(m.File)
	if err != nil {
		fmt.Println(gid, m.File, ":", err)
		return
	}
	defer fp.Close()

	fmt.Println("--------------------------")
	fmt.Printf("  %s: Pattern Start:[%s][%d:%d] --- End:[%s][%d:%d] Need Read[%s]\n",
		gid, Kscal(spos), slineTime.Hour, slineTime.Minute, Kscal(epos), elineTime.Hour, elineTime.Minute, Kscal(epos-spos))
	fmt.Println("--------------------------")
	fmt.Println()

	fp.Seek(spos, os.SEEK_SET)
	reader := bufio.NewReader(fp)
	if spos > 0 {
		reader.ReadString('\n')
	}
	for {
		sline, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println(gid, "---EOF---")
				break
			}
			fmt.Println(gid, m.File, ":", err)
			break
		}

		if m.IgnoreCase {
			if strings.Contains(strings.ToLower(sline), m.Pattern) {
				if m.ShowFileName {
					fmt.Printf("%s:%s", m.File, sline)
				} else {
					fmt.Printf("%s", sline)
				}
			}
		} else {
			if strings.Contains(sline, m.Pattern) {
				if m.ShowFileName {
					fmt.Printf("%s:%s", m.File, sline)
				} else {
					fmt.Printf("%s", sline)
				}
			}
		}

		if epos > 0 {
			cc, _ := fp.Seek(0, os.SEEK_CUR)
			if cc > epos {
				break
			}
		}
	}
}

func getLineTime(line string) LineTime {
	month := ""
	day := ""
	time := ""

	ret := LineTime{}
	x := 0
	for i := 0; i < len(line); i++ {

		if x == 0 {
			month = month + string(line[i])
		} else if x == 1 {
			day = day + string(line[i])
		} else if x == 2 {
			time = time + string(line[i])
		} else {
			break
		}

		if line[i] == ' ' {
			for {
				if line[i+1] == ' ' {
					i++
					continue
				}
				break
			}
			x++
			continue
		}
	}

	ret.Month = ConverMaothToIntFromString(month)
	ret.Day, _ = strconv.Atoi(day)

	timeList := strings.Split(time, ":")
	ret.Hour, _ = strconv.Atoi(timeList[0])
	ret.Minute, _ = strconv.Atoi(timeList[1])
	ret.Second, _ = strconv.Atoi(timeList[2])

	return ret
	// return date, day, time
}

// t: 0 start, 1 end
func getPos(f string, matchHour int, t int, v bool, gid string) (int64, LineTime, error) {
	resPos := int64(0)
	isA := false
	isB := false
	var lineTime LineTime

	fileInfo, err := os.Stat(f)
	if err != nil {
		return 0, lineTime, err
	}
	if fileInfo.Size() <= 0 {
		return 0, lineTime, nil
	}

	fp, err := os.Open(f)
	if err != nil {
		return 0, lineTime, err
	}
	defer fp.Close()

	soffset := int64(0)
	eoffset := fileInfo.Size()
	coffset := eoffset / 2

	if v {
		fmt.Printf("  %s: start:%d cur:%d filesize:%d\n", gid, soffset, coffset, eoffset)
	}

	fp.Seek(coffset, os.SEEK_SET)
	reader := bufio.NewReader(fp)
	reader.ReadString('\n') // 为了下一行为完整行
	for {
		sline, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				//fmt.Println(gid, "---EOF---")
				break
			}
			fmt.Println(gid, f, ":", err)
			break
		}

		lineTime = getLineTime(sline)
		if lineTime.Hour > matchHour {
			isA = true
			eoffset = coffset
			coffset = coffset / 2
			fp.Seek(coffset, os.SEEK_SET)
			reader = bufio.NewReader(fp)
			reader.ReadString('\n')

			if v {
				fmt.Printf("  %s: A MatchHour[%d] LineHour:[%d] start:%d cur:%d end:%d\n", gid, matchHour, lineTime.Hour, soffset, coffset, eoffset)
			}
			if isB && isA {
				// resPos = soffset
				break
			}
		} else if lineTime.Hour < matchHour {
			isB = true
			soffset = coffset
			coffset = (eoffset-soffset)/2 + coffset
			fp.Seek(coffset, os.SEEK_SET)
			reader = bufio.NewReader(fp)
			reader.ReadString('\n')

			if v {
				fmt.Printf("  %s: B MatchHour[%d] LineHour:[%d] start:%d cur:%d end:%d\n", gid, matchHour, lineTime.Hour, soffset, coffset, eoffset)
			}
			if soffset == coffset {
				break
			}
			if isA && isB {
				// resPos = eoffset
				break
			}
		} else {
			break
		}
	}

	if t == 0 {
		resPos = soffset
	} else {
		resPos = eoffset
	}

	return resPos, lineTime, nil
}
