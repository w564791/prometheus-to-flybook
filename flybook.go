package main

import "gopkg.in/alecthomas/kingpin.v2"

//var (
//	//flyBookHook = flag.String("url","","FlyBook webhook url")
//	flyBookHook = kingpin.Flag("url","FlyBook webhook urls,defaule Name must provide").Required().Strings()
//
//)
//
//type Message struct {
//	Status string`json:"status"`
//	Alerts []Alert`json:"alerts"`
//	Receiver string`json:"receiver"`
//	GroupLabels map[string]string`json:"groupLabels"`
//	CommonAnnotations map[string]string`json:"commonAnnotations"`
//	CommonLabels map[string]string`json:"commonLabels"`
//	StartsAt string`json:"startsAt"`
//	EndsAt string`json:"endsAt"`
//}
//type Alert struct {
//	Status string `json:"status"`
//	Labels map[string]string`json:"labels"`
//	Annotations map[string]string`json:"annotations"`
//	StartsAt string`json:"startsAt"`
//	EndsAt string`json:"endsAt"`
//	GeneratorURL string`json:"generatorURL"`
//
//}
//
//type Flybook struct {
//	Title string`json:"title"`
//	Text string`json:"text"`
//}
//
////type Servers struct {
////	Router string
////	Function func
////}
//const (
//	ftemplJson  = `
//{"title": "[{{.CommonLabels.severity|ToUpper}} {{.Status|ToUpper}}] {{.GroupLabels.alertname}} [{{len .Alerts}}]",
//"text": "{{ range $index,$item := .Alerts}}Lebels:\n{{range $key,$value:= .Labels}}\t{{$key}} : {{$value}}\n{{end}}------\n{{end}}StartAt: {{ .StartsAt|ToTimeFormat}}\nAnnotations:\n\tmessage: {{.CommonAnnotations.message}}\n\trunbook_url: {{.CommonAnnotations.runbook_url}}"}
//`
//	rtemplJson  = `
//{"title": "[{{.CommonLabels.severity|ToUpper}} {{.Status|ToUpper}}] {{.GroupLabels.alertname}} [{{len .Alerts}}]",
//"text": "{{ range $index,$item := .Alerts}}Lebels:\n{{range $key,$value:= .Labels}}\t{{$key}} : {{$value}}\n{{end}}------\n{{end}}StartAt: {{ .StartsAt|ToTimeFormat}}\nEndsAt: {{.EndsAt|ToTimeFormat}}\nAnnotations:\n\tmessage: {{.CommonAnnotations.message}}\n\trunbook_url: {{.CommonAnnotations.runbook_url}}"}
//`
//	//	rtemplJson  = `
//	//{"title": "[{{.CommonLabels.severity|ToUpper}} {{.Status|ToUpper}}] {{.GroupLabels.alertname}}] [{{len .Alerts}}]",
//	//"text": "Lebels: \n{{ range $index,$item := .Alerts}} {{range $key,$value:= .Labels}}{{if eq $key "alertname" }}{{end}}\t{{$key}} : {{$value}}\n{{end}}\t---\nStartAt: {{ .StartsAt|ToTimeFormat}}\nEndsAt: {{.EndsAt|ToTimeFormat}}{{end}}\nAnnotations:\n\tmessage: {{.CommonAnnotations.message}}\n\trunbook_url: {{.CommonAnnotations.runbook_url}}"}
//	//`
//	//	url="https://open.feishu.cn/open-apis/bot/hook/bcd7176ae53d444e8d042d9231d0a778"
//	time_layout= "2006-01-02 15:04:05"
//	boundPort=":9090"
//)
//func TimeFormat(s string)(t string){
//	tmpT,_:=time.Parse(time.RFC3339,s)
//	t=tmpT.Format(time_layout)
//	return
//}

var (
	//flyBookHook = flag.String("url","","FlyBook webhook url")
	flyBookHook = kingpin.Flag("url", "FlyBook webhook urls,defaule Name must provide").Required().Strings()
)

type Message struct {
	Status            string            `json:"status"`
	Alerts            []Alert           `json:"alerts"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	CommonLabels      map[string]string `json:"commonLabels"`
	StartsAt          string            `json:"startsAt"`
	EndsAt            string            `json:"endsAt"`
}
type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     string            `json:"startsAt"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
}

type Flybook struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

//type Servers struct {
//	Router string
//	Function func
//}
const (
	ftemplJson = `
{"title": "[{{.CommonLabels.severity|ToUpper}} {{.Status|ToUpper}}] {{.GroupLabels.alertname}} [{{len .Alerts}}]",
"text": "{{ range $index,$item := .Alerts}}{{Add $index 1}}.Lebels:\n{{range $key,$value:= .Labels}}    - {{$key}} : {{$value}}\n{{end}}  StartAt: {{ .StartsAt|ToTimeFormat}}\n  Annotations:\n    message: {{.Annotations.message}}\n{{end}}"}
`

	rtemplJson = `
{"title": "[{{.CommonLabels.severity|ToUpper}} {{.Status|ToUpper}}] {{.GroupLabels.alertname}} [{{len .Alerts}}]",
"text": "{{ range $index,$item := .Alerts}}{{Add $index 1}}.Lebels:\n{{range $key,$value:= .Labels}}    - {{$key}} : {{$value}}\n{{end}}  StartAt: {{ .StartsAt|ToTimeFormat}}\n  EndsAt: {{.EndsAt|ToTimeFormat}}\n  Annotations:\n    message: {{.Annotations.message}}\n{{end}}"}
`

	//	url="https://open.feishu.cn/open-apis/bot/hook/bcd7176ae53d444e8d042d9231d0a778"
	time_layout = "2006-01-02 15:04:05"
	boundPort   = ":9090"
)
