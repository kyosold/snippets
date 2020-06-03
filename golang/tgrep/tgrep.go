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
	Verbose      int
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

const StepC float64 = 0.2

func usage(prog string) {
	fmt.Println("----------------------------------------------")
	fmt.Println("Usage:")
	fmt.Printf(" %s -st=15 -et=22 [pattern] [file]\n", prog)
	fmt.Println()
	fmt.Println("options:")
	fmt.Printf("  -st 从指定的小时开始查找\n")
	fmt.Printf("  -et 到指定的小时结束查找, 不写默认到文件结尾\n")
	fmt.Printf("  -i 不区分大小写\n")
	fmt.Printf("  -v 显示详细信息\n")
	fmt.Println()
	fmt.Println("Others:")
	fmt.Printf("  1. 如果文件是gzip，先解压再查询，如:\n")
	fmt.Printf("    gzip -dvc abc.0.gz > abc.0\n")
	fmt.Printf("    tgrep -st=8 -et=10 'pattern' abc.0\n")
	fmt.Println("----------------------------------------------")
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
	verbose := flag.Int("v", 0, "Verbose")

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
	if *verbose >= 0 {
		startIdx++
	}

	pattern = os.Args[startIdx]
	dir = os.Args[startIdx+1:]

	fmt.Println("**************************")
	fmt.Println("Start Hour:", *startTime)
	fmt.Println("End Hour:", *endTime)
	fmt.Println("Ignore Case:", *ignoreCase)
	fmt.Println("Verbose:", *verbose)
	fmt.Println("Pattern:", pattern)
	fmt.Println("Files:", dir)
	fmt.Println("**************************")

	mr := Matcher{}
	mr.StartHour = *startTime
	mr.EndHour = *endTime
	mr.Pattern = pattern
	mr.ShowFileName = false
	mr.Verbose = *verbose
	// if *verbose == 1 {
	// mr.Verbose = true
	// } else {
	// mr.Verbose = false
	// }
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
			go search(mr, &wg)
		} else if cType == "application" {
			if cSubType == "x-gzip" {
				fmt.Printf("[Error]: file:%s Is GZIP file\n", mr.File)
				fmt.Printf("  Run 'gzip -dvc %s > tmp.log'\n", mr.File)
			}
		}
	}
	wg.Wait()
}

func search(m Matcher, wg *sync.WaitGroup) {
	defer wg.Done()

	bHourMatched := false

	gid := GetUUID()
	if m.Verbose > 0 {
		fmt.Printf("%s: %s [%s]\n", gid, m.File, m.ContentType)
		fmt.Println("--------------------------")
	}

	// Get start pos
	spos, slineTime, err := getPos(m.File, m.StartHour, m.Verbose, gid, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("Start Pos:", spos, "LineTime:", slineTime)

	// Get end pos
	var epos int64 = 0
	var elineTime LineTime
	if m.EndHour > 0 {
		if m.Verbose > 0 {
			fmt.Println()
		}
		epos, elineTime, err = getPos(m.File, m.EndHour, m.Verbose, gid, 1)
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
	// fmt.Println("End Pos:", spos, "LineTime:", slineTime)

	fp, err := os.Open(m.File)
	if err != nil {
		fmt.Println(gid, m.File, ":", err)
		return
	}
	defer fp.Close()

	fmt.Println("--------------------------")
	fmt.Printf("%s: Pattern Start:[%d][%d:%d] --- End:[%d][%d:%d] Need Read[%s]\n",
		gid, spos, slineTime.Hour, slineTime.Minute, epos, elineTime.Hour, elineTime.Minute, Kscal(epos-spos))
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

		if bHourMatched == false {
			lineTime := getLineTime(sline)
			if m.StartHour != lineTime.Hour {
				continue
			} else {
				bHourMatched = true
			}
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

// flag: 0: 取 matchHour 前面pos； 1：取 matchHour 后面的pos
func getPos(f string, matchHour int, v int, gid string, flag int) (int64, LineTime, error) {
	// resPos := int64(0)
	isA := false
	isB := false
	var lineTime LineTime
	var lineTimeLast LineTime

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

	spos := int64(0)
	epos := fileInfo.Size()
	cpos := epos / 2

	// if v > 0 {
	// fmt.Printf("  %s: start:%d cur:%d filesize:%d\n", gid, spos, cpos, epos)
	// }
	fp.Seek(cpos, os.SEEK_SET)
	reader := bufio.NewReader(fp)
	reader.ReadString('\n') // 为了下一行是完整行
	for {
		sline, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(gid, f, ":", err)
			break
		}

		lineTime = getLineTime(sline)
		if matchHour <= lineTime.Hour { // 在上面
			isA = true
			if v > 0 {
				fmt.Printf("  %s: A MatchHour[%d] LineHour:[%d] start:%d cur:%d end:%d\n", gid, matchHour, lineTime.Hour, spos, cpos, epos)
			}

			if isB {
				epos = cpos
				// lineTime = lineTimeLast
				break
			}

			// lineTimeLast = lineTime
			epos = cpos
			cpos = epos / 2
			fp.Seek(cpos, os.SEEK_SET)
			reader = bufio.NewReader(fp)
			reader.ReadString('\n')

		} else if matchHour > lineTime.Hour { // 在下面
			isB = true
			if v > 0 {
				fmt.Printf("  %s: B MatchHour[%d] LineHour:[%d] start:%d cur:%d end:%d\n", gid, matchHour, lineTime.Hour, spos, cpos, epos)
			}

			if isA {
				spos = cpos
				// lineTime = lineTimeLast
				break
			}
			if spos == cpos {
				break
			}

			// lineTimeLast = lineTime
			spos = cpos
			cpos = (epos-cpos)/2 + cpos
			fp.Seek(cpos, os.SEEK_SET)
			reader = bufio.NewReader(fp)
			reader.ReadString('\n')

		}
	}
	lineTimeLast = lineTime

	step := int64(float64(epos-spos) * StepC)
	if v > 0 {
		fmt.Printf("  %s: Start Setp MatchHour[%d] LineHour:[%d] step:%d spos:%d --- epos:%d\n", gid, matchHour, lineTime.Hour, step, spos, epos)
	}

	if step > 100 {
		cpos = spos + step
		fp.Seek(cpos, os.SEEK_SET)
		reader = bufio.NewReader(fp)
		reader.ReadString('\n')
		for {
			sline, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println(gid, f, ":", err)
				break
			}

			lineTime = getLineTime(sline)
			if matchHour > lineTime.Hour { // 继续向下找
				if v > 0 {
					fmt.Printf("  %s: [STEP] MatchHour[%d] LineHour:[%d] cur:%d\n", gid, matchHour, lineTime.Hour, cpos)
				}

				if cpos > epos {
					fmt.Printf("  %s: [SETP] E cpos(%d) > epos(%d), exit.\n", cpos, epos)
					break
				}

				lineTimeLast = lineTime
				cpos = cpos + step
				fp.Seek(cpos, os.SEEK_SET)
				reader = bufio.NewReader(fp)
				reader.ReadString('\n')
			} else {

				if flag == 0 {
					if v > 0 {
						fmt.Printf("  %s: [STEP] E MatchHour[%d] LineHour:[%d] cpos:%d\n", gid, matchHour, lineTime.Hour, cpos)
					}
					cpos = cpos - step
					lineTime = lineTimeLast
					break
				} else {
					if matchHour < lineTime.Hour {
						if v > 0 {
							fmt.Printf("  %s: [STEP] E MatchHour[%d] LineHour:[%d] cpos:%d\n", gid, matchHour, lineTime.Hour, cpos)
						}
						break
					}
					if v > 0 {
						fmt.Printf("  %s: [STEP] MatchHour[%d] LineHour:[%d] cpos:%d\n", gid, matchHour, lineTime.Hour, cpos)
					}
					cpos = cpos + step
					fp.Seek(cpos, os.SEEK_SET)
					reader = bufio.NewReader(fp)
					reader.ReadString('\n')
				}
			}
		}
	} else {
		if flag == 0 {
			cpos = spos
		} else {
			cpos = epos
		}
	}

	return cpos, lineTime, nil
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
	ret.Month = ConverMonthToIntFromString(month)
	ret.Day, _ = strconv.Atoi(day)

	timeList := strings.Split(time, ":")
	ret.Hour, _ = strconv.Atoi(timeList[0])
	ret.Minute, _ = strconv.Atoi(timeList[1])
	ret.Second, _ = strconv.Atoi(timeList[2])

	return ret
}
