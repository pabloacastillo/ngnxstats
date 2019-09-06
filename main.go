package main

import (
	"log"
	"fmt"
	"os"
	"io"
	"time"
	"math"
	"strings"
	"sort"
	"gopkg.in/urfave/cli.v1"
	"github.com/satyrius/gonx"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	// "encoding/json"
	// "github.com/hpcloud/tail"
)

type SummaryRecord struct {
	Title string
	Total int
}

type LogRecord struct {
	remote_addr string
	time_local  string
	rest string
	request string
	status string
	body_bytes_sent string
	http_referer string
	http_user_agent string
}
var records [] LogRecord
var file_size int64 = 0

var app = cli.NewApp()

func info(){
	app.Name = "ngnxstats"
	app.Usage = "A realtime nginx log reader"
	app.Version = "0.0.1"	
}

func flags(){
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "log, l",
			Usage: "Read log from `FILE` ",
			Value: "/var/log/nginx/access.log",
		}, cli.StringFlag{
			Name: "format, f",
			Usage: "Log file format ",
			Value: "$remote_addr - - [$time_local] \"$rest $request $http_type\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\"",
		},
	}


	var name   = "/var/log/nginx/acccess.log"
	// var format = "$remote_addr - - [$time_local] \"$rest $request $http_type\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\""
	var format = "$remote_addr - - [$time_local] \"$rest $request $http_type\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\""

	// 190.104.148.41 - - [01/Jul/2019:06:26:45 -0400] "GET /wp-content/uploads/2017/09/if_Application_728900.png HTTP/1.1" 404 209 "-" "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36"
	

	app.Action = func(c *cli.Context) {

		if(c.String("log")!=""){
			name = c.String("log")
		}

		if(c.String("format")!=""){
			format = c.String("format")
		}
	
		// fmt.Println("Hola", name)
		// fmt.Println("Hola", format)
		file_size = check_file_size(name)
		draw_interface(name,format)
	}
}



// func draw_interface(records []LogRecord){
func draw_interface(log_file string, log_format string) {

	var records []LogRecord =read_log_file(log_file,log_format)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()




	termWidth, termHeight := ui.TerminalDimensions()

	
	// footer := widgets.NewParagraph()
	// footer.Title = "Apreta q para salir "
	// footer.SetRect(0, 0, termWidth, 1)
	// footer.Border = true
	// footer.TextStyle.Bg = ui.ColorMagenta

	// WIDGETS
	ips:= widget_common_ips(records)
	rest_chart := widget_rest_chart(records)
	rest_summary := widget_rest_summary(records)
	raw_logs_table := widget_raw_table(records)


	grid := ui.NewGrid()
	grid.SetRect(0, 0, termWidth, termHeight)
	
	grid.Set(
		ui.NewRow(0.2,
			ui.NewCol(0.15, ips),
			ui.NewCol(0.15, rest_summary),
			ui.NewCol(0.7, rest_chart),
		),
		ui.NewRow(0.8,
			ui.NewCol(1.0, raw_logs_table),
		),
		// ui.NewRow(0.1,
		// 	ui.NewCol(1.0, footer),
		// ),
	)
	

	ui.Render(grid)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C

	for {
		select {
			case e := <-uiEvents:
				switch e.ID {
					case "q", "<C-c>":
						return
					case "<Resize>":
						payload := e.Payload.(ui.Resize)
						grid.SetRect(0, 0, payload.Width, payload.Height)
						ui.Clear()
						ui.Render(grid)
				}
			case <- ticker:
				tfile_size := check_file_size(log_file)
				if tfile_size>file_size {
					file_size = tfile_size
					draw_interface(log_file,log_format)
				}
				// records =read_log_file(log_file,log_format)
		}
	}

}



func main() {
	info()
	flags()


	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}	
}


func check_file_size(log_file string) int64 {
	fi, err := os.Stat(log_file) 
	if err != nil {
	   panic(err)
	}
	fs:=fi.Size()
	return fs
}

func read_log_file(log_file string, format string) ([]LogRecord ){
	
	var logReader io.Reader
	file, err := os.Open(log_file)
	if err != nil {
		panic(err)
	}

	logReader = file
	defer file.Close()

	reader := gonx.NewReader(logReader, format)
	
	var cont int
	cont = 0


	var records []LogRecord

	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		} 
		cont++

		var remote_addr, _ = rec.Field("remote_addr")
		if remote_addr=="nil"{
			fmt.Printf("%+v ", remote_addr )
		}
		var time_local, _ = rec.Field("time_local")
		var rest, _ = rec.Field("rest")
		var request, _ = rec.Field("request")
		var status, _ = rec.Field("status")
		var body_bytes_sent, _ = rec.Field("body_bytes_sent")

		var http_referer, _ = rec.Field("http_referer")
		// var http_user_agent, _ = rec.Field("remote_addr")

		var _record LogRecord
		_record.remote_addr = remote_addr
		_record.status = status
		_record.time_local = time_local
		_record.rest = rest
		_record.request = request
		_record.body_bytes_sent = body_bytes_sent
		_record.http_referer = http_referer
		// _record.http_user_agent




		records = append(records, _record)

	}
	
	var a = records

	for i := len(a)/2-1; i >= 0; i-- {
		opp := len(a)-1-i
		a[i], a[opp] = a[opp], a[i]
	}


	return a
}

/**

	888       888     8888888     8888888b.       .d8888b.      8888888888     88888888888      .d8888b.
	888   o   888       888       888  "Y88b     d88P  Y88b     888                888         d88P  Y88b
	888  d8b  888       888       888    888     888    888     888                888         Y88b.
	888 d888b 888       888       888    888     888            8888888            888          "Y888b.
	888d88888b888       888       888    888     888  88888     888                888             "Y88b.
	88888P Y88888       888       888    888     888    888     888                888               "888
	8888P   Y8888       888       888  .d88P     Y88b  d88P     888                888         Y88b  d88P
	888P     Y888     8888888     8888888P"       "Y8888P88     8888888888         888          "Y8888P"

*/


// RAW DATA TABLE
func widget_raw_table(rawlogs []LogRecord) (*widgets.Table) {

	termWidth, termHeight := ui.TerminalDimensions()

	logs := widgets.NewTable()
	logs.Title = "LOGS CRUDOS DEL NGINX ACCESS LOG"
	logs.SetRect(0, 0, termWidth,termHeight)
	logs.Border = true
	logs.RowSeparator = false
	logs.FillRow=true
	logs.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)
	logs.Rows = append(logs.Rows, []string{ "IP              .", "TIME                       .", "REST   .", "REQUEST", "STATUS", "SIZE", "REFERER"})

	for k, v := range rawlogs {


		if(v.status[:1]=="4"){
			logs.RowStyles[k+1] = ui.NewStyle(ui.ColorRed)
		}
		if(v.status[:1]=="2"){
			logs.RowStyles[k+1] = ui.NewStyle(ui.ColorGreen)
		}
		if(v.status[:1]=="3"){
			logs.RowStyles[k+1] = ui.NewStyle(ui.ColorCyan)
		}


		logs.Rows = append(logs.Rows, []string{ v.remote_addr, v.time_local, v.rest, v.request, v.status, v.body_bytes_sent, v.http_referer})
	}

	return logs;

}

// COMMON IPs SUMMARY BY HITS
func widget_common_ips (rawlogs []LogRecord) (*widgets.List){
		
	
	ips_sum := make(map[string]int)
	for k := range rawlogs {
		ips_sum[ rawlogs[k].remote_addr ]++ 
	}

	var ss []SummaryRecord
	for k, v := range ips_sum {
		ss = append(ss, SummaryRecord{k, v} )
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Total > ss[j].Total
	})

	ips_sum_list :=[]string{}
	ip_s:=""
	total_s:=""
	ipsum_s:=""
	for _, kv := range ss[:7] {

		ip_s=kv.Title + strings.Repeat(" ", 17 - len(kv.Title) )
		
		total_s = fmt.Sprintf("%d",kv.Total)
		total_s = strings.Repeat(" ", 10 - len(total_s) ) + total_s
		ipsum_s = ip_s + " " +  total_s
		ips_sum_list = append(ips_sum_list, ipsum_s )

	}

	ips := widgets.NewList()
	ips.Rows = ips_sum_list
	ips.Title = "IPs"
	ips.Border = true

	return ips
}

// TOTAL OF REQUEST BY REST TYPE
func widget_rest_summary (rawlogs []LogRecord) (*widgets.List){

	rest_sum := make(map[string]int)
	for k := range rawlogs {
		rest_sum[ rawlogs[k].rest ]++ 
	}

	var ss []SummaryRecord
	for k, v := range rest_sum {
		ss = append(ss, SummaryRecord{k, v} )
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Total > ss[j].Total
	})

	rest_sum_list :=[]string{}
	rest_s:=""
	total_s:=""
	restsum_s:=""
	for _, kv := range ss[:7] {

		if len(kv.Title)<8{
			rest_s=kv.Title + strings.Repeat(" ", 8 - len(kv.Title) )
		} else {
			rest_s=kv.Title[0:8]
		}
		
		total_s = fmt.Sprintf("%d",kv.Total)
		total_s = strings.Repeat(" ", 15 - len(total_s) ) + total_s
		restsum_s = rest_s + " " +  total_s
		rest_sum_list = append(rest_sum_list, restsum_s )

	}


	REST := widgets.NewList()
	REST.Rows = rest_sum_list
	// REST.Rows = []string{
	// 	"GET",
	// 	"HEAD",
	// 	"POST",
	// 	"PUT",
	// 	"PATCH",
	// 	"DELETE",
	// 	"CONNECT",
	// 	"OPTIONS",
	// 	"TRACE",
	// }
	REST.Title = "REST"
	REST.Border = true

	return REST

}

// 
func widget_rest_chart(rawlogs []LogRecord) (*widgets.Plot) {
	
	termWidth, termHeight := ui.TerminalDimensions()

	sinData := func() [][]float64 {
		n := 220
		data := make([][]float64, 2)
		data[0] = make([]float64, n)
		data[1] = make([]float64, n)
		for i := 0; i < n; i++ {
			data[0][i] = 1 + math.Sin(float64(i)/5)
			data[1][i] = 1 + math.Cos(float64(i)/5)
		}
		return data
	}()


	rest_chart := widgets.NewPlot()
	rest_chart.Data = sinData
	rest_chart.AxesColor = ui.ColorWhite
	rest_chart.Marker = widgets.MarkerDot
	rest_chart.ShowAxes = false
	rest_chart.SetRect(0, 0, termWidth, termHeight - termHeight + 1)
	rest_chart.Border = false
	return rest_chart
}







