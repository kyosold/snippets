# 使用方法:

**1. 创建 `ctlog` 对象:**
```go
ctlog, err := NewCTLog("./this.log", false)
if err != nil {
    fmt.Println(err)
    return
}
```
- 参数:
   - 日志文件
   - 是否使用系统`log`, 如果你希望日志中带有函数名,行号等，设置为`false`


**2. 设置日志输出等级:**
```go
ctlog.SetLevel(CTINFO)
```
参数说明:

参数 | 值
--|:-:|--
CTEMERG | 0
CTALERT | 1
CTCRIT | 2
CTERR | 3
CTWARNING | 4
CTNOTICE | 5
CTINFO | 6
CTDEBUG | 7

**3. 调用方式:**
```
ctlog.Debug("%s: this is Debug", "sj")
ctlog.Info("%s: this is Info", "sj")
ctlog.Notice("%s: this is Notice", "sj")
ctlog.Warning("%s: this is Warning", "sj")
ctlog.Error("%s: this is Error", "sj")
ctlog.Crit("%s: this is Crit", "sj")
ctlog.Alert("%s: this is Alert", "sj")
ctlog.Emerg("%s: this is Emerg", "sj")
```

## Example:
```go
package main

import "fmt"

func main() {
    ctlog, err := NewCTLog("./this.log", false)
    if err != nil {
        fmt.Println(err)
        return
    }
    ctlog.SetLevel(CTINFO)

    ctlog.Info("call show")

    show(ctlog)
}

func show(ctlog *Ctlog) {
    ctlog.Debug("%s: this is Debug", "sj")
    ctlog.Info("%s: this is Info", "sj")
    ctlog.Notice("%s: this is Notice", "sj")
    ctlog.Warning("%s: this is Warning", "sj")
    ctlog.Error("%s: this is Error", "sj")
    ctlog.Crit("%s: this is Crit", "sj")
    ctlog.Alert("%s: this is Alert", "sj")
    ctlog.Emerg("%s: this is Emerg", "sj")
}
```
