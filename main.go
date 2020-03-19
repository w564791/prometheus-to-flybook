package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)
type Message struct {
	Status string`json:"status"`
	Alerts []Alert`json:"alerts"`
	Receiver string`json:"receiver"`
	GroupLabels map[string]string`json:"groupLabels"`
	CommonAnnotations map[string]string`json:"commonAnnotations"`
	CommonLabels map[string]string`json:"commonLabels"`
	StartsAt string`json:"startsAt"`
	EndsAt string`json:"endsAt"`
}
type Alert struct {
	Status string `json:"status"`
	Labels map[string]string`json:"labels"`
	Annotations map[string]string`json:"annotations"`
	StartsAt string`json:"startsAt"`
	EndsAt string`json:"endsAt"`
	GeneratorURL string`json:"generatorURL"`

}

type Flybook struct {
	Title string`json:"title"`
	Text string`json:"text"`
}

var (
	flyBookHook = flag.String("url","","FlyBook webhook url")
)

const (
	ftemplJson  = `
{"title": "[{{.CommonLabels.severity|ToUpper}} {{.Status|ToUpper}}] {{.GroupLabels.alertname}} [{{len .Alerts}}]",
"text": "{{ range $index,$item := .Alerts}}Lebels:\n{{range $key,$value:= .Labels}}\t{{$key}} : {{$value}}\n{{end}}------\n{{end}}StartAt: {{ .StartsAt|ToTimeFormat}}\nAnnotations:\n\tmessage: {{.CommonAnnotations.message}}\n\trunbook_url: {{.CommonAnnotations.runbook_url}}"}
`
	rtemplJson  = `
{"title": "[{{.CommonLabels.severity|ToUpper}} {{.Status|ToUpper}}] {{.GroupLabels.alertname}} [{{len .Alerts}}]",
"text": "{{ range $index,$item := .Alerts}}Lebels:\n{{range $key,$value:= .Labels}}\t{{$key}} : {{$value}}\n{{end}}------\n{{end}}StartAt: {{ .StartsAt|ToTimeFormat}}\nEndsAt: {{.EndsAt|ToTimeFormat}}\nAnnotations:\n\tmessage: {{.CommonAnnotations.message}}\n\trunbook_url: {{.CommonAnnotations.runbook_url}}"}
`
//	rtemplJson  = `
//{"title": "[{{.CommonLabels.severity|ToUpper}} {{.Status|ToUpper}}] {{.GroupLabels.alertname}}] [{{len .Alerts}}]",
//"text": "Lebels: \n{{ range $index,$item := .Alerts}} {{range $key,$value:= .Labels}}{{if eq $key "alertname" }}{{end}}\t{{$key}} : {{$value}}\n{{end}}\t---\nStartAt: {{ .StartsAt|ToTimeFormat}}\nEndsAt: {{.EndsAt|ToTimeFormat}}{{end}}\nAnnotations:\n\tmessage: {{.CommonAnnotations.message}}\n\trunbook_url: {{.CommonAnnotations.runbook_url}}"}
//`
//	url="https://open.feishu.cn/open-apis/bot/hook/bcd7176ae53d444e8d042d9231d0a778"
  time_layout= "2006-01-02 15:04:05"
	boundPort=":9090"
	)

func TimeFormat(s string)(t string){
	tmpT,_:=time.Parse(time.RFC3339,s)
	t=tmpT.Format(time_layout)
	return
}

func DataToFlyBook(msg Message)(flybook Flybook,err error)  {
	var result bytes.Buffer
	var tmpl *template.Template
	var funcMap template.FuncMap
	//t,_:=time.Parse(time.RFC3339,"2020-03-18T08:15:47.64901609Z")
	//log.Println(t.Format(time_layout))
	msg.StartsAt=msg.Alerts[0].StartsAt
	msg.EndsAt=msg.Alerts[0].EndsAt

	funcMap = template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToTimeFormat":TimeFormat,
	}

	if msg.Status == "firing" {

		tmpl= template.Must(template.New("").Funcs(funcMap).Parse(ftemplJson))

	}else if  msg.Status == "resolved" {
		tmpl= template.Must(template.New("").Funcs(funcMap).Parse(rtemplJson))
	}

	if err=tmpl.Execute(&result,msg);err!=nil{
		log.Println(err.Error())
	}
	err=json.Unmarshal([]byte(result.String()),&flybook)
	if err!=nil{
		log.Println(err.Error())
	}
	return
}
func MessageDeploy(w http.ResponseWriter, r *http.Request) {
	var message Message
	//var result bytes.Buffer
	//var flybook Flybook
	//var flaybook *Flaybook

	if r.Body == nil{
		fmt.Println("No Body")
	}
	s, _ := ioutil.ReadAll(r.Body)
	err:=json.Unmarshal([]byte(s),&message)
	if err!=nil{
		log.Println(err.Error())
	}
	flydata,err:=DataToFlyBook(message)
	if err !=nil{
		log.Println(err.Error())
	}

	statusCode,err:=SendMessage(&flydata,*flyBookHook)
	if err != nil{
		log.Println(err.Error())

	}
	log.Println(fmt.Sprintf("%d -  Send %s Monitor messge - %s",statusCode,message.GroupLabels["alertname"],message.Status))
}
func SendMessage(msg *Flybook,url string) (respcode int,err error)  {
	b := new(bytes.Buffer)
	err=json.NewEncoder(b).Encode(msg)
	if err != nil {
		log.Println(err.Error())
	}
	req,err := http.NewRequest("POST",url,b)
	if err != nil {
		log.Println(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	//log.Println(b)
	resp,err:=http.DefaultClient.Do(req)
	if err != nil{
		log.Println(err.Error())

	}
	respcode=resp.StatusCode
	return

}
func main() {
	flag.Parse()
	if *flyBookHook == "" {
		panic("Must Provide Flybook webhook")
	}
	http.HandleFunc("/",MessageDeploy )
	log.Println(fmt.Sprintf("Starting Server At %s",boundPort))
	log.Fatal(http.ListenAndServe(boundPort,nil))

}
