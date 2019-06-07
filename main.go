package main

import (
	"log"
	"fmt"
	"os"
	"io"
	// "strings"
	"gopkg.in/urfave/cli.v1"
	"github.com/satyrius/gonx"
	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
)

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
		format:="$remote_addr - - [$time_local] \"$request\" $status $body_bytes_sent \"$http_referer\" \"$http_user_agent\""

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


	header := widgets.NewParagraph()
	header.Text = "Press q to quit, Press h or l to switch tabs"
	header.SetRect(0, 0, 50, 1)
	header.Border = false
	header.TextStyle.Bg = ui.ColorBlue


	ui.Render(header)

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
	// draw_interface()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func read_log_file(log_file string, format string){
	// fmt.Println("Hola", log_file)
	// return nil

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
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		} 
		cont++
		// rec
		// Process the record... e.g.
		fmt.Printf("Parsed entry: %+v\n", rec)
		// fmt.Println(rec)
	}
	// fmt.Printf("Parsed entry: \n", cont)
	fmt.Println(cont)
	// fmt.Println(log_file)
	// fmt.Println(logReader)

}