package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

//time format RFC3339 to 2020-03-18 08:04:17
func SeqAdd(index int, seq int) int {
	return index + seq
}
func TimeFormat(s string) (t string) {
	tmpT, _ := time.Parse(time.RFC3339, s)
	t = tmpT.Format(time_layout)
	return
}

//dea args to map
func ArgsDeploy(fargs *[]string) (args map[string]string) {
	tmpArg := make(map[string]string)
	for _, v := range *fargs {
		tmpv := strings.Split(v, "=")
		if len(tmpv) == 2 {
			tmpArg[tmpv[0]] = tmpv[1]
		}
	}
	args = tmpArg
	return
}

//dump data to flaybook msgData {"title":"msg":,"text":"msg"}
func DataToFlyBook(msg Message) (flybook Flybook, err error) {
	var result bytes.Buffer
	var tmpl *template.Template
	var funcMap template.FuncMap

	msg.StartsAt = msg.Alerts[0].StartsAt
	msg.EndsAt = msg.Alerts[0].EndsAt
	for index, _ := range msg.Alerts {

		delete(msg.Alerts[index].Labels, "alertname")
		delete(msg.Alerts[index].Labels, "severity")
	}
	funcMap = template.FuncMap{
		"ToUpper":      strings.ToUpper,
		"ToTimeFormat": TimeFormat,
		"Add":          SeqAdd,
	}

	if msg.Status == "firing" {
		tmpl = template.Must(template.New("").Funcs(funcMap).Parse(ftemplJson))
	} else if msg.Status == "resolved" {
		tmpl = template.Must(template.New("").Funcs(funcMap).Parse(rtemplJson))
	}

	if err = tmpl.Execute(&result, msg); err != nil {
		log.Println(err.Error())
	}
	err = json.Unmarshal([]byte(result.String()), &flybook)
	if err != nil {
		log.Println(err.Error())
	}
	return
}

//search key fron map ,if true,return assin to url,else url="default"
func KeySearch(m map[string]string, s string, ss *string) (url *string) {
	url = ss
	tmpurl := m["default"]
	if v, ok := m[s]; ok {
		url = &v
	} else {
		url = &tmpurl
	}
	return
}
func MessageDeploy(w http.ResponseWriter, r *http.Request) {
	var message Message
	var url *string
	var args map[string]string

	if r.Body == nil {
		fmt.Println("No Body")
	}
	s, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal([]byte(s), &message)
	if err != nil {
		log.Println(err.Error())
	}
	flydata, err := DataToFlyBook(message)
	if err != nil {
		log.Println(err.Error())
	}
	args = ArgsDeploy(flyBookHook)
	switch message.CommonLabels["severity"] {
	case message.CommonLabels["severity"]:
		url = KeySearch(args, message.CommonLabels["severity"], url)
	default:
		url = KeySearch(args, "default", url)
	}

	statusCode, err := SendMessage(&flydata, *url)
	if err != nil {
		log.Println(err.Error())

	}
	rourceIp := strings.Split(r.RemoteAddr, ":")[0]
	fmt.Fprintf(w, "send")
	log.Println(fmt.Sprintf("- %s - %d -  Send %s Monitor messge - %s", rourceIp, statusCode, message.GroupLabels["alertname"], message.Status))
}
func SendMessage(msg *Flybook, url string) (respcode int, err error) {
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(msg)
	if err != nil {
		log.Println(err.Error())
	}
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		log.Println(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	//log.Println(b)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	respcode = resp.StatusCode
	return

}
func main() {

	kingpin.Parse()
	log.Println(fmt.Sprintf("Starting Server At %s", boundPort))
	http.HandleFunc("/", MessageDeploy)

	log.Fatal(http.ListenAndServe(boundPort, nil))

}
