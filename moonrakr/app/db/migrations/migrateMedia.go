package main

import (
	"flag"
	"strings"

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

	media, err := models.GetMediumByUpdatedAt()
	if err != nil {
		revel.ERROR.Fatal(err)
	}

	for _, i := range *media {
		i.SetUrl(strings.Replace(i.GetUrl(), "vvvnt.s3.amazonaws.com", "jeroboam.s3.amazonaws.com", -1))
		err = model.Save(i)
		if err != nil {
			revel.ERROR.Fatal(err)
		}
	}
	posts, err := models.GetPostsByCreatedAt()
	if err != nil {
		revel.ERROR.Fatal(err)
	}

	for _, i := range *posts {
		i.SetBody(strings.Replace(i.GetBody(), "vvvnt.s3.amazonaws.com", "jeroboam.s3.amazonaws.com", -1))
		err := model.Save(i)
		if err != nil {
			revel.ERROR.Fatal(err)
		}
	}
}
