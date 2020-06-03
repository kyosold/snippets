## 说明:
```
Usage:
 tgrep -st=15 -et=22 [pattern] [file]
options:
  -v 显示详细信息
  -st 指定的小时时间开始查找
  -et 指定的小时时间结束查找, 不写默认到文件结尾
  -i 不区分大小写
Others:
  1. 如果文件是gzip，先解压再查询，如:
    gzip -dvc abc.0.gz > abc.0
    tgrep -st=8 -et=10 'pattern' abc.0
```

## 例子:
- 在文件 `abc.log` 中查询 `8点-13点` 的日志，关键词为 `113322`:
```
./tgrep -v=1 -st=8 -et=13 '113322' abc.log
```
