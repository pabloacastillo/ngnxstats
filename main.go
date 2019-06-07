package main

import (
	"log"
	"fmt"
	"os"
	"gopkg.in/urfave/cli.v1"
	// "github.com/satyrius/gonx"
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

	
	app.Action = func(c *cli.Context) error {
		name := "/var/log/nginx/acccess.log"
		if c.NArg() > 0 {
			// name = c.Args().Get(0)
			name = c.String("log")
		}

		if(c.String("log")!=""){
			name = c.String("log")
		}
		read_log_file(name)
		// fmt.Println("Hola", name)
		// fmt.Println("Hola", c.String("log"))
		return nil
	}
}

func draw_interface(){

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Title = " Hello World! "
	p.SetRect(0, 0, 25, 5)

	ui.Render(p)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
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

func read_log_file(log_file string){
	fmt.Println("Hola", log_file)
	// return nil
}