package main

import (
	"fmt"
	"log"
	"mail"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-ini/ini"
)

type CONF struct {
	host string
	from string
	to   string
	eml  string

	tls      bool
	isAuth   bool
	user     string
	password string

	attachList []string
}

func usage(prog string) {
	fmt.Println("Usage:")
	fmt.Printf("  %s [config file]\n", prog)
	fmt.Printf("  %s [cofnig file] -d\n", prog)
}

func main() {
	argsWithProg := os.Args
	if len(argsWithProg) != 2 && len(argsWithProg) != 3 {
		usage(filepath.Base(argsWithProg[0]))
		return
	}
	verbose := true
	if len(argsWithProg) == 3 {
		if argsWithProg[2] == "-d" {
			verbose = false
		}
	}

	cfgFile := argsWithProg[1]
	cfg, err := ini.Load(cfgFile)
	if err != nil {
		fmt.Println("Read config file fail:", err)
		return
	}

	smtpSection, err := cfg.GetSection("SMTP")
	if err != nil {
		fmt.Println("Parse ini section 'SMTP' file fail:", err)
		return
	}
	host := smtpSection.Key("Host").String()
	from := smtpSection.Key("From").String()
	to := smtpSection.Key("To").String()
	eml := smtpSection.Key("Eml").String()
	tls := smtpSection.Key("Tls").String()
	isAuth := smtpSection.Key("IsAuth").String()
	user := smtpSection.Key("User").String()
	password := smtpSection.Key("Password").String()
	attach := smtpSection.Key("addAttach").String()

	headSection, err := cfg.GetSection("ADDHEADER")
	if err != nil {
		fmt.Println("Parse ini section 'ADDHEADER' fail fail:", err)
		return
	}
	//fmt.Println(headSection.KeyStrings())

	fmt.Println("**************************************")
	fmt.Println("[SMTP]")
	fmt.Println("  Eml:", eml)
	fmt.Println("  Host:", host)
	fmt.Println("  IsAuth:", isAuth)
	fmt.Println("  User:", user)
	fmt.Println("  Password:", password)
	fmt.Println("  From:", from)
	fmt.Println("  To:", to)
	fmt.Println("  Tls:", tls)
	fmt.Println("  Attach:", attach)
	fmt.Println()

	fmt.Println("[ADDHEADER]")
	for _, key := range headSection.KeyStrings() {
		fmt.Printf("  %s: %s\n", key, headSection.Key(key).String())
	}
	fmt.Println("**************************************")

	toList := strings.Split(to, ",")
	btls := false
	if tls == "1" {
		btls = true
	}
	var attachList []string
	if len(attach) > 0 {
		attachList = strings.Split(attach, ",")
	}

	reg := regexp.MustCompile(`^{\$(\S*?)\$}$`)

	header := make(map[string]string)
	for _, key := range headSection.KeyStrings() {
		result := reg.FindAllStringSubmatch(key, -1)
		if len(result) == 0 {
			header[key] = headSection.Key(key).String()
		} else {
			hk := result[0][1]
			if hk == "Message-ID" {
				header[hk] = mail.CreateMessageID()
			}
		}
	}

	smtp, err := mail.NewSmtp(host, btls, from, toList)
	if err != nil {
		log.Fatal(err)
	}
	if isAuth == "1" {
		smtp.SetAuth(user, password)
	}
	smtp.SetVerbose(verbose)
	qid, err := smtp.SendMailWithEmlOrAttach(eml, header, attachList)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(qid)
		fmt.Println("Finished.")
	}
}
