package main

import (
	"log"
	"fmt"
	"os"
	"io"
	"time"
	// "strings"
	"gopkg.in/urfave/cli.v1"
	"github.com/satyrius/gonx"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

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

		draw_interface(name,format)
	}
}



// func draw_interface(records []LogRecord){
func draw_interface(log_file string, log_format string){

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	var records []LogRecord =read_log_file(log_file,log_format)



	termWidth, termHeight := ui.TerminalDimensions()

	
	header := widgets.NewParagraph()
	header.Text = "Apreta ctrl+c para salir"
	header.SetRect(0, 0, termWidth, 1)
	header.Border = true
	header.TextStyle.Bg = ui.ColorBlue

	footer := widgets.NewParagraph()
	footer.Title = "Apreta q para salir "
	footer.SetRect(0, 0, termWidth, 1)
	footer.Border = true
	footer.TextStyle.Bg = ui.ColorMagenta

	raw_logs_table := raw_table(records)
	
	ips := widgets.NewList()
	ips.Rows = []string{
		"1230 999.999.999.000",
		"1230 999.999.999.999",
		"1230 999.999.999.999",
		"1230 999.999.999.999",
		"1230 999.999.999.999",
		"1230 999.999.999.999",
	}
	ips.Title = "IPs"
	ips.Border = true

	REST := widgets.NewList()
	REST.Rows = []string{
		"GET",
		"POST",
		"HEAD",
		"UPDATE",
		"DELETE",
	}
	REST.Title = "REST"
	REST.Border = true

	grid := ui.NewGrid()
	grid.SetRect(0, 0, termWidth, termHeight)
	
	grid.Set(
		ui.NewRow(0.2,
			ui.NewCol(0.2, ips),
			ui.NewCol(0.7, header),
			ui.NewCol(0.1, REST),
		),
		ui.NewRow(0.7,
			ui.NewCol(1.0, raw_logs_table),
		),
		ui.NewRow(0.1,
			ui.NewCol(1.0, footer),
		),
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
				// records =read_log_file(log_file,log_format)
				// draw_interface(log_file,log_format)
				ui.Render(grid)
				// ui.update(grid)
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
	// start := time.Now()
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


// TABLA DE DATOS CRUDOS
func raw_table(rawlogs []LogRecord) (*widgets.Table) {

	termWidth, termHeight := ui.TerminalDimensions()

	logs := widgets.NewTable()
	logs.Title = "LOGS CRUDOS DEL NGINX ACCESS LOG"
	logs.SetRect(0, 0, termWidth,termHeight)
	logs.Border = true
	logs.RowSeparator = false
	logs.FillRow=true
	logs.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)
	logs.Rows = append(logs.Rows, []string{ "IP              .", "TIME                       .", "REST", "REQUEST", "STATUS", "SIZE", "REFERER"})

	for k, v := range rawlogs {


		if(v.status[:1]=="4"){
			logs.RowStyles[k+1] = ui.NewStyle(ui.ColorRed, ui.ColorBlack)
		}
		if(v.status[:1]=="2"){
			logs.RowStyles[k+1] = ui.NewStyle(ui.ColorGreen, ui.ColorBlack)
		}
		if(v.status[:1]=="3"){
			logs.RowStyles[k+1] = ui.NewStyle(ui.ColorCyan, ui.ColorBlack)
		}


		logs.Rows = append(logs.Rows, []string{ v.remote_addr, v.time_local, v.rest, v.request, v.status, v.body_bytes_sent, v.http_referer})
	}

	return logs;

}

