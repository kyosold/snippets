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

