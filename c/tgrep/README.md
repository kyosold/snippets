## 使用方式:
```bash
./tgrep -v -s[开始时间] -e[结束时间] 'Pattern' [文件]
```

- 查询：在文件 abc.log `16:03:14` 点的关键字为 'aaabbb' 的结果
```
./tgrep -v -s16 -e17 'aaabbb' ./abc.log
```

## 使用说明:
```
----------------------------------------------
Usage:
  ./tgrep -v -s15 -e17 [pattern] [file]

Options:
  -s: 从指定的小时开始查找
  -e: 到指定的小时结束查找, 不写默认到文件结尾
  -i: 不区分大小写
  -v: 显示详细信息
  -S: 显示文件名

Others:
  1. 如果文件是gzip，先解压再查询，如:
    a. gzip -dvc abc.0.gz > abc.0
    b. tgrep -st=8 -et=10 'pattern' abc.0
----------------------------------------------
```