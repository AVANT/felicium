package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/AVANT/felicium/model"
	"github.com/AVANT/felicium/moonrakr/app/lib/boot"
	"github.com/AVANT/felicium/moonrakr/app/models"
	"github.com/robfig/revel"
)

func main() {
	env := flag.String("env", "dev", "allows you to set the env flag")
	flag.Parse()
	revel.ERROR.Println(*env)

	boot.ToolBoot(*env)

	posts, err := models.GetPostsByCreatedAt()
	if err != nil {
		revel.ERROR.Fatal(err)
	}

	for _, i := range *posts {
		revel.ERROR.Printf("%+v", i.Model)
		revel.ERROR.Println(i.GetPublishedAt())
		revel.ERROR.Println(i.GetId())
		revel.ERROR.Println(i.GetRev())
		if fmt.Sprint(i.GetPublishedAt()) == "0001-01-01 00:00:00 +0000 UTC" {
			i.SetPublishedAt(time.Now())
		}
		err := model.Save(i)
		if err != nil {
			revel.ERROR.Fatal(err)
		}
	}
}
