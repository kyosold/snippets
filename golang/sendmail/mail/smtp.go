package mail

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	DIS int = 0
	TLS int = 1
	SSL int = 2
)

type CONF struct {
	host string
	from string
	to   string
	eml  string

	cryto    int
	isAuth   bool
	user     string
	password string

	attachList []string
}

type step int

const (
	GREET    step = 0
	EHELO    step = 1
	AUTH     step = 2
	USER     step = 3
	PASS     step = 4
	FROM     step = 5
	TO       step = 6
	DATA     step = 7
	EOH      step = 8
	EOM      step = 9
	QUIT     step = 10
	CLOSE    step = 11
	STARTTLS step = 12
	NEW      step = 99
)

type SMTP struct {
	CurStep  step
	UseTLS   int
	Host     string
	User     string
	Password string
	EnvFrom  string
	EnvTO    []string
	Verbose  bool
}

type MailMime struct {
	Head map[string]string
	Body string
}

// host: server_address:port
func NewSmtp(host string, useTLS int, from string, to []string) (SMTP, error) {
	smtp := SMTP{}
	smtp.Verbose = false
	smtp.Host = host
	smtp.UseTLS = useTLS
	smtp.EnvFrom = from
	smtp.EnvTO = to

	return smtp, nil
}

/*
func (smtp SMTP) SendMailWithMailMime(mail MailMime) (string, error) {
	res := ""

	var smtpStatus step = GREET
	toi := 0

	conn, err := net.Dial("tcp", smtp.Host)
	if err != nil {
		fmt.Println("Connect Fail:", err)
		return res, err
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		var line = make([]byte, 8192)
		_, err := reader.Read(line)
		if err != nil {
			fmt.Println("Read Fail:", err)
			conn.Close()
			break
		}

		buf := string(line)

		if smtp.Verbose {
			fmt.Printf("< %s", buf)
		}

		code := buf[:3]

		// 判断返回码
		if res, _ := resCode(code, smtpStatus); res != true {
			fmt.Println(smtpStatus, "Server Fail:", buf)
			conn.Close()
			break
		}

		switch smtpStatus {
		case GREET:
			smtpStatus = EHELO
		case EHELO:
			if strings.Contains(buf, "250-AUTH") {
				smtpStatus = AUTH
			} else {
				smtpStatus = FROM
			}
		case AUTH:
			smtpStatus = USER
		case USER:
			smtpStatus = PASS
		case PASS:
			smtpStatus = FROM
		case FROM:
			smtpStatus = TO
		case TO:
			if toi >= len(smtp.EnvTO) {
				smtpStatus = DATA
			}
		case DATA:
			smtpStatus = EOM
		case EOM:
			smtpStatus = QUIT
		case QUIT:
			smtpStatus = CLOSE
		}

		if smtpStatus == EHELO {
			hname, err := os.Hostname()
			if err != nil {
				fmt.Println(err)
				return res, nil
			}

			text := "EHLO " + hname + "\r\n"
			if smtp.Verbose {
				fmt.Printf("> %s", text)
			}

			err = sendSmtp(conn, text)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}

		} else if smtpStatus == AUTH {
			text := "AUTH LOGIN\r\n"
			if smtp.Verbose {
				fmt.Printf("> %s", text)
			}

			err := sendSmtp(conn, text)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}

		} else if smtpStatus == USER {
			text := smtp.User + "\r\n"
			if smtp.Verbose {
				fmt.Printf("> %s", text)
			}

			err := sendSmtp(conn, text)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}

		} else if smtpStatus == PASS {
			text := smtp.Password + "\r\n"
			if smtp.Verbose {
				fmt.Printf("> %s", text)
			}

			err := sendSmtp(conn, text)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}

		} else if smtpStatus == FROM {
			text := "MAIL FROM:<" + smtp.EnvFrom + ">\r\n"
			if smtp.Verbose {
				fmt.Printf("> %s", text)
			}
			err := sendSmtp(conn, text)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}

		} else if smtpStatus == TO {
			text := "RCPT TO:<" + smtp.EnvTO[toi] + ">\r\n"
			if smtp.Verbose {
				fmt.Printf("> %s", text)
			}
			err := sendSmtp(conn, text)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}
			toi++

		} else if smtpStatus == DATA {
			text := "DATA\r\n"
			if smtp.Verbose {
				fmt.Printf("> %s", text)
			}
			err := sendSmtp(conn, text)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}

		} else if smtpStatus == EOM {
			for k, v := range mail.Head {
				hline := fmt.Sprintf("%s: %s\r\n", k, v)
				if smtp.Verbose {
					fmt.Printf("> %s", hline)
				}
				err := sendSmtp(conn, hline)
				if err != nil {
					fmt.Println("Write Fail:", err)
					return res, err
				}
			}
			text := "\r\n"
			sendSmtp(conn, text)

			err := sendSmtp(conn, mail.Body)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}
			conn.Write([]byte("\r\n.\r\n"))

			if smtp.Verbose {
				fmt.Printf("> \r\n")
				fmt.Printf("> %s", mail.Body)
				fmt.Printf("\r\n.\r\n")
			}

		} else if smtpStatus == QUIT {
			res = buf

			text := "QUIT\r\n"
			if smtp.Verbose {
				fmt.Printf("> %s", text)
			}
			sendSmtp(conn, text)

		} else if smtpStatus == CLOSE {
			conn.Close()
			break
		}
	}

	return res, nil
}
*/

func (smtp SMTP) SendMailWithEmlOrAttach(eml string, header map[string]string, attachList []string) (string, error) {
	res := ""

	emlFP, err := os.Open(eml)
	if err != nil {
		fmt.Println(err)
		return res, err
	}
	defer emlFP.Close()
	readerEml := bufio.NewReader(emlFP)
	bufEml := make([]byte, 4096)

	hname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		return res, nil
	}

	isTLS := 0

	var smtpStatus step = GREET
	toi := 0

	var conn net.Conn
	if smtp.UseTLS == SSL {
		conn, err = tls.Dial("tcp", smtp.Host, nil)
		if err != nil {
			fmt.Println(err)
			return res, err
		}
	} else {
		conn, err = net.Dial("tcp", smtp.Host)
		if err != nil {
			fmt.Println(err)
			return res, err
		}
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		var line = make([]byte, 8192)
		_, err := reader.Read(line)
		if err != nil {
			fmt.Println("Read Fail:", err)
			conn.Close()
			break
		}

		buf := string(line)

		if smtp.Verbose {
			fmt.Printf("< %s", buf)
		}

		code := buf[:3]

		// 判断返回码
		if res, _ := resCode(code, smtpStatus); res != true {
			fmt.Println(smtpStatus, "Server Fail:", buf)
			conn.Close()
			break
		}

		switch smtpStatus {
		case GREET:
			smtpStatus = EHELO
		case STARTTLS:
			tlsConf := &tls.Config{
				ServerName:         hname,
				InsecureSkipVerify: true,
			}
			conn = tls.Client(conn, tlsConf)
			reader = bufio.NewReader(conn)
			isTLS = 2
			smtpStatus = EHELO
		case EHELO:
			if smtp.UseTLS == TLS && isTLS == 0 && strings.Contains(buf, "250-STARTTLS") {
				// use tls
				isTLS = 1
				smtpStatus = STARTTLS
			} else if strings.Contains(buf, "250-AUTH") {
				smtpStatus = AUTH
			} else {
				smtpStatus = FROM
			}
		case AUTH:
			smtpStatus = USER
		case USER:
			smtpStatus = PASS
		case PASS:
			smtpStatus = FROM
		case FROM:
			smtpStatus = TO
		case TO:
			if toi >= len(smtp.EnvTO) {
				smtpStatus = DATA
			}
		case DATA:
			smtpStatus = EOM
		case EOM:
			smtpStatus = QUIT
		case QUIT:
			smtpStatus = CLOSE
		}

		if smtpStatus == EHELO {

			text := "EHLO " + hname + "\r\n"
			err = sendSmtp(conn, text, smtp.Verbose)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}

		} else if smtpStatus == STARTTLS {
			text := "STARTTLS\r\n"
			err := sendSmtp(conn, text, smtp.Verbose)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}

		} else if smtpStatus == AUTH {
			text := "AUTH LOGIN\r\n"
			err := sendSmtp(conn, text, smtp.Verbose)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}

		} else if smtpStatus == USER {
			text := smtp.User + "\r\n"
			err := sendSmtp(conn, text, smtp.Verbose)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}

		} else if smtpStatus == PASS {
			text := smtp.Password + "\r\n"
			err := sendSmtp(conn, text, smtp.Verbose)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}

		} else if smtpStatus == FROM {
			text := "MAIL FROM:<" + smtp.EnvFrom + ">\r\n"
			err := sendSmtp(conn, text, smtp.Verbose)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}

		} else if smtpStatus == TO {
			text := "RCPT TO:<" + smtp.EnvTO[toi] + ">\r\n"
			err := sendSmtp(conn, text, smtp.Verbose)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}
			toi++

		} else if smtpStatus == DATA {
			text := "DATA\r\n"
			err := sendSmtp(conn, text, smtp.Verbose)
			if err != nil {
				fmt.Println("Write Fail:", err)
				return res, err
			}

		} else if smtpStatus == EOM {
			if len(header) > 0 { // 写自定义头
				for key, val := range header {
					hline := key + ": " + val + "\r\n"
					err = sendSmtp(conn, hline, smtp.Verbose)
					if err != nil {
						fmt.Println("Write Fail:", err)
						return res, err
					}
				}
			}

			// 写附件
			if len(attachList) > 0 {
				mhead := "Date: " + CreateDateString() + "\r\n"
				mhead = mhead + "From: " + smtp.EnvFrom + "\r\n"
				mhead = mhead + "To: " + strings.Join(smtp.EnvTO, ",") + "\r\n"
				mhead = mhead + "Subject: " + ConvStringToMimeString(filepath.Base(attachList[0])) + "\r\n"
				boundaryID := "------=_Part_" + CreateMessageID()
				mhead = mhead + "Content-Type: multipart/mixed;\r\n"
				mhead = mhead + "\tboundary=\"" + boundaryID + "\"\r\n"
				mhead = mhead + "MIME-Version: 1.0\r\n"

				err = sendSmtp(conn, mhead, smtp.Verbose)
				if err != nil {
					fmt.Println("Write Fail:", err)
					return res, err
				}

				for _, attach := range attachList {
					name := ConvStringToMimeString(filepath.Base(attach))
					boundaryHead := "\r\n--" + boundaryID + "\r\n"
					boundaryHead = boundaryHead + "Content-Type: application/octet-stream;\r\n\tname=\"" + name + "\"; charset=UTF-8\r\n"
					boundaryHead = boundaryHead + "Content-Transfer-Encoding: base64\r\n"
					boundaryHead = boundaryHead + "Content-Disposition: attachment;\r\n\tfilename=\"" + name + "\"\r\n\r\n"

					err = sendSmtp(conn, boundaryHead, smtp.Verbose)
					if err != nil {
						fmt.Println("Write Fail:", err)
						return res, err
					}

					attText, err := getAttachConextMimeWithFile(attach)
					if err != nil {
						fmt.Println("Error:", err)
						return res, err
					}

					err = sendSmtp(conn, attText, smtp.Verbose)
					if err != nil {
						fmt.Println("Write Fail:", err)
						return res, err
					}

				}
				boundaryEnd := "\r\n\r\n--" + boundaryID + "--\r\n"
				err = sendSmtp(conn, boundaryEnd, smtp.Verbose)
				if err != nil {
					fmt.Println("Write Fail:", err)
					return res, err
				}
			} else {
				// 写邮件正文
				for {
					n, err := readerEml.Read(bufEml)
					if err != nil && err != io.EOF {
						fmt.Println("Read eml file fail:", err)
						return res, err
					}
					if n == 0 {
						break
					}

					err = sendSmtp(conn, string(bufEml), smtp.Verbose)
					if err != nil {
						fmt.Println("Write Fail:", err)
						return res, err
					}
				}
			}

			conn.Write([]byte("\r\n.\r\n"))

			if smtp.Verbose {
				fmt.Printf("\r\n.\r\n")
			}

		} else if smtpStatus == QUIT {
			res = buf

			text := "QUIT\r\n"
			sendSmtp(conn, text, smtp.Verbose)

		} else if smtpStatus == CLOSE {
			conn.Close()
			break
		}
	}

	return res, nil
}

func (s *SMTP) SetAuth(user string, pwd string) {
	s.User = base64.StdEncoding.EncodeToString([]byte(user))
	s.Password = base64.StdEncoding.EncodeToString([]byte(pwd))
}
func (s *SMTP) SetVerbose(v bool) {
	s.Verbose = v
}

func sendSmtp(conn net.Conn, text string, verbose bool) error {
	_, err := conn.Write([]byte(text))
	if err != nil {
		fmt.Println(err)
		return err
	}
	if verbose {
		fmt.Printf("> %s", text)
	}
	return nil
}

func resCode(code string, sp step) (bool, error) {
	switch sp {
	case GREET:
	case STARTTLS:
		if code != "220" {
			return false, nil
		}
	case EHELO:
		if code != "250" {
			return false, nil
		}
	case AUTH:
		if code != "334" {
			return false, nil
		}
	case USER:
		if code != "334" {
			return false, nil
		}
	case PASS:
		if code != "235" {
			return false, nil
		}
	case FROM:
		if code != "250" {
			return false, nil
		}
	case TO:
		if code != "250" {
			return false, nil
		}
	case DATA:
		if code != "354" {
			return false, nil
		}
	case EOM:
		if code != "250" {
			return false, nil
		}
	case QUIT:
		if code != "221" {
			return false, nil
		}

	}

	return true, nil
}

// ----------------------

func NewMailMime(from string, tolist []string) (MailMime, error) {
	mm := MailMime{}
	mm.Head = make(map[string]string)

	dateStr := CreateDateString()
	mm.Head["Date"] = dateStr

	mm.Head["From"] = from
	mm.Head["To"] = strings.Join(tolist, ",")

	msgID := CreateMessageID()
	mm.Head["Message-ID"] = msgID

	mm.Head["Content-Type"] = "text/plain; charset=\"UTF-8\""
	mm.Head["Content-Transfer-Encoding"] = "base64"

	return mm, nil
}

func (mm *MailMime) SetSubject(subject string) {
	subB64 := base64.StdEncoding.EncodeToString([]byte(subject))
	s := "=?utf-8?B?" + subB64 + "?="
	mm.Head["Subject"] = s
}

func (mm *MailMime) SetHeader(key string, val string) {
	mm.Head[key] = val
}

// body 不需要做编码，使用原始字符串就可以
func (mm *MailMime) SetBody(body string) {

	// base64 enc
	bodyB64 := base64.StdEncoding.EncodeToString([]byte(body))

	mm.Body = bodyB64
}
