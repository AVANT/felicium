package main

import (
	"flag"

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
		err := model.Save(i)
		if err != nil {
			revel.ERROR.Fatal(err)
		}
	}
}
