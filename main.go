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
	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
)

var app = cli.NewApp()

type LogRecord struct {
    remote_addr string
    time_local  string
    rest string
    request string
    status int
    body_bytes_sent int
    http_referer string
    http_user_agent string
}

var records [] LogRecord

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
		},
	}

	// app.Flags = []cli.Flag {
	// 	cli.StringFlag{
	// 		Name: "format, f",
	// 		Usage: "Log file format ",
	// 		Value: "$remote_addr - - [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\"",
	// 	},
	// }

	
	app.Action = func(c *cli.Context) error {
		name := "/var/log/nginx/acccess.log"
		format:="$remote_addr - - [$time_local] \"$rest $request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\""
		
		// 104.236.216.165 - - [01/Jun/2019:06:25:18 -0400] "POST /wp-cron.php?doing_wp_cron=1559384718.5513510704040527343750 HTTP/1.1" 200 31 "https://www.obedira.com.py/wp-cron.php?doing_wp_cron=1559384718.5513510704040527343750" "WordPress/5.2.1; https://www.obedira.com.py"


		if c.NArg() > 0 {
			
			if(c.String("log")!=""){
				name = c.String("log")
			}


		}

		if(c.String("log")!=""){
			name = c.String("log")
		}

		// if(c.String("format")!=""){
		// 	format = c.String("format")
		// }

		fmt.Println("Hola", format)
		fmt.Println("Hola", name)
		read_log_file(name,format)
		// fmt.Println("Hola", c.String("log"))
		return nil
	}
}

func draw_interface(){

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	termWidth, termHeight := ui.TerminalDimensions()

	
	header := widgets.NewParagraph()
	header.Text = "Apreta ctrl+c para salir"
	header.SetRect(0, 0, termWidth, 1)
	header.Border = true
	header.TextStyle.Bg = ui.ColorBlue

	footer := widgets.NewParagraph()
	footer.Title = "Apreta q para salir ----------------------------------------"
	footer.SetRect(0, 0, termWidth, 1)
	footer.Border = true
	footer.TextStyle.Bg = ui.ColorMagenta

	logs := widgets.NewParagraph()
	logs.Title = "LOGS CRUDOS DEL NGINX ACCESS LOG"
	logs.SetRect(0, 0, termWidth,termHeight-200)
	logs.Border = true
	


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
			ui.NewCol(1.0, logs),
		),
		ui.NewRow(0.1,
			ui.NewCol(1.0, footer),
		),
	)
	

	ui.Render(grid)

	for e := range ui.PollEvents() {

		switch e.ID {
			case "q", "<C-c>":
				return
		}

		// if e.Type == ui.KeyboardEvent {
		// 	break
		// }
	}
}


func main() {

	info()
	flags()
	draw_interface()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func read_log_file(log_file string, format string){
	
	var logReader io.Reader
	file, err := os.Open(log_file)
	if err != nil {
		panic(err)
	}
	logReader = file
	defer file.Close()

	// var format="$remote_addr - - [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\""
	// 104.236.216.165 - - [01/Jun/2019:06:25:18 -0400] "POST /wp-cron.php?doing_wp_cron=1559384718.5513510704040527343750 HTTP/1.1" 200 31 "https://www.obedira.com.py/wp-cron.php?doing_wp_cron=1559384718.5513510704040527343750" "WordPress/5.2.1; https://www.obedira.com.py"


	reader := gonx.NewReader(logReader, format)
	
	var cont int
	cont = 0
	start := time.Now()
	for {
		// var rec *gonx.Entry
		rec, err := reader.Read()
		if err == io.EOF {
			break
		} 
		cont++

		var remote_addr, _ = rec.Field("remote_addr")
		if remote_addr=="nil"{
			fmt.Printf("%+v ", remote_addr )

		}
		// var time_local, _ = rec.Field("time_local")
		// var rest, _ = rec.Field("rest")
		// var request, _ = rec.Field("request")
		// var status, _ = rec.Field("status")
		// var body_bytes_sent, _ = rec.Field("body_bytes_sent")

		// var http_referer, _ = rec.Field("remote_addr")
		// var http_user_agent, _ = rec.Field("remote_addr")

		// fmt.Printf("%+v ", remote_addr )
		// fmt.Printf("%+v ", status )
		// fmt.Printf("%+v ", rest )
		// fmt.Printf("%+v ", time_local )
		// fmt.Printf("%+v ", body_bytes_sent )
		// fmt.Printf("%+v\n", request )
		
		// var _record LogRecord
		// _record.remote_addr
		// _record.time_local 
		// _record.rest
		// _record.request
		// _record.status
		// _record.body_bytes_sent
		// _record.http_referer
		// _record.http_user_agent


		// LogRecord
	}
	
	duration := time.Since(start)
	fmt.Printf("%v lines readed, it takes %v\n", cont, duration)
	// fmt.Println(log_file)
	// fmt.Println(logReader)

}