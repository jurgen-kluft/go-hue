package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jurgen-kluft/go-hue/groups"
	"github.com/jurgen-kluft/go-hue/lights"
	"github.com/jurgen-kluft/go-hue/portal"
)

var (
	apiKey     string = ""
	blinkState lights.State
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: get-light-state -key=[string]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	on := new(bool)
	*on = true
	blinkState = lights.State{On: on, Alert: "lselect"}
	flag.StringVar(&apiKey, "key", os.Getenv("HUE_USERNAME"), "hue light api key")
	flag.Parse()
	flag.Usage = usage
}

func main() {
	if apiKey != "" {
		pp, err := portal.GetPortal()
		if err != nil {
			fmt.Println("portal.GetPortal() ERROR: ", err)
			os.Exit(1)
		}
		ll := lights.New(pp[0].InternalIPAddress, apiKey)
		allLights, err := ll.GetAllLights()
		if err != nil {
			fmt.Println("lights.GetAllLights() ERROR: ", err)
			os.Exit(1)
		}
		fmt.Println()
		fmt.Println("Lights")
		fmt.Println("------")
		for _, l := range allLights {
			fmt.Printf("ID: %d Name: %s\n", l.ID, l.Name)
		}
		gg := groups.New(pp[0].InternalIPAddress, apiKey)
		allGroups, err := gg.GetAllGroups()
		if err != nil {
			fmt.Println("groups.GetAllGroups() ERROR: ", err)
			os.Exit(1)
		}
		fmt.Println()
		fmt.Println("Groups")
		fmt.Println("------")
		for _, g := range allGroups {
			fmt.Printf("ID: %d Name: %s\n", g.ID, g.Name)
			for _, lll := range g.Lights {
				fmt.Println("\t", lll)
			}
			previousState := g.Action
			_, err := gg.SetGroupState(g.ID, blinkState)
			if err != nil {
				fmt.Println("groups.SetGroupState() ERROR: ", err)
				os.Exit(1)
			}
			time.Sleep(time.Second * time.Duration(10))
			_, err = gg.SetGroupState(g.ID, previousState)
			if err != nil {
				fmt.Println("groups.SetGroupState() ERROR: ", err)
				os.Exit(1)
			}
		}
	} else {
		usage()
	}
}
