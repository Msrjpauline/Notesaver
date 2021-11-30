package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/qrtz/nativemessaging"
)

type message struct {
	Text string `json:"text"`
}

type response struct {
	Text    string `json:"text"`
	Success bool   `json:"success"`
}

var strNoteFile string
var nowTime time.Time

func deleteLine(no int) {
	var fileRecords []string
	file, err := os.Open("UserSelection.log")	
	cnt := 1
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)	
	for scanner.Scan() {
		if no != cnt{
			fileRecords = append(fileRecords,scanner.Text())
		}   
		cnt++;
	}
	nfile, err := os.Create("UserSelection.log")         
	defer nfile.Close()
	for _, v := range fileRecords {
		nfile.WriteString(v)
		nfile.WriteString("\n")		
	}     
}

func formThePage() {
	var cnt int = 1
	f, err := os.Open("UserSelection.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f1, err := os.Create(strNoteFile)
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer f1.Close()
	f1.WriteString("<!DOCTYPE html><html><style>table, th, td { border:1px solid black;}</style><body>")
	f1.WriteString("<h1>My Saved Text</h1><table style=\"width:95%\">")
	f1.WriteString("<tr><th>MessageName</th><th>Messge</th></tr>")

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		strWords := strings.Split(scanner.Text(), "DEL")
		if len(strWords) != 3 {
			continue
		}
		f1.WriteString("<tr><td><h2>" + strWords[0] + "</h2>" + "<p>" + "Msg#" + strconv.Itoa(cnt) + "</p>" + "<p>" + strWords[2] + "</p>" + "</td><td><b>" + strWords[1] + "</b></td></tr>")
		cnt++
	}
	f1.WriteString("</body></html>")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	decoder := nativemessaging.NewNativeJSONDecoder(os.Stdin)
	encoder := nativemessaging.NewNativeJSONEncoder(os.Stdout)
	f2, fErr := os.OpenFile("UserSelection.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fErr != nil {
		log.Println(fErr)
	}
	defer f2.Close()
	strNoteFile, _ = os.Getwd()
	strNoteFile = strNoteFile + "\\MySelNotes.html"

	formThePage()
	for {
		var rsp response
		var msg message
		err := decoder.Decode(&msg)

		if err != nil {
			if err == io.EOF {
				// exit
				return
			}
			rsp.Text = err.Error()
		} else {
			if msg.Text == "getpath" {
				rsp.Text = strNoteFile
				rsp.Success = true
			} else if strings.Contains(msg.Text, "Delete#") {
				strTxts := strings.Split(msg.Text, "#")
				i, _ := strconv.Atoi(strTxts[1])
				deleteLine(i)
				formThePage()
				rsp.Text = "Deleted"
				rsp.Success = true

			} else if strings.Contains(msg.Text, "VisitedURL") {
				strTxts := strings.Split(msg.Text, "#")
				i, _ := strconv.Atoi(strTxts[1])
				deleteLine(i)
				formThePage()
				rsp.Text = "Visited"
				rsp.Success = true

			}else {
				nowTime = time.Now()
				f2, fErr := os.OpenFile("UserSelection.log",
					os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if fErr != nil {
					log.Println(fErr)
				}				
				_, err := f2.WriteString(msg.Text + "DEL" + nowTime.Format(time.Stamp) + "\n")
				if err != nil {
				} else {
					f2.Close()
					formThePage()
					rsp.Text = "YourTextisSaved"
					rsp.Success = true
				}
			}
		}
		if err := encoder.Encode(rsp); err != nil {
			return
		}
	}
}
