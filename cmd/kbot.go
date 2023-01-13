/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/stianeikeland/go-rpio"
	telebot "gopkg.in/telebot.v3"
)

var (
	// TeleToken bot
	TeleToken = os.Getenv("TELE_TOKEN")
)

// kbotCmd represents the kbot command
var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("kbot %s started", appVersion)

		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Plaese check TELE_TOKEN env variable. %s", err)
			return
		}

		err = rpio.Open()
		if err != nil {
			log.Printf("Unable to open gpio: %s", err.Error())
		}

		defer rpio.Close()

		trafficSignal := make(map[string]map[string]int8)

		trafficSignal["red"] = make(map[string]int8)
		trafficSignal["amber"] = make(map[string]int8)
		trafficSignal["green"] = make(map[string]int8)

		trafficSignal["red"]["pin"] = 12
		//default on/off
		//trafficSignal["red"]["on"]=0
		trafficSignal["amber"]["pin"] = 27
		trafficSignal["green"]["pin"] = 22

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {

			var (
				err error
				pin = rpio.Pin(0)
			)
			log.Print(m.Message().Payload, m.Text())
			payload := m.Message().Payload

			switch payload {
			case "hello":
				err = m.Send(fmt.Sprintf("Hello I'm Kbot %s!", appVersion))

			case "red", "amber", "green":
				pin = rpio.Pin(trafficSignal[payload]["pin"])
				if trafficSignal[payload]["on"] == 0 {
					pin.Output()
					trafficSignal[payload]["on"] = 1
				} else {
					pin.Input()
					trafficSignal[payload]["on"] = 0
				}

				err = m.Send(fmt.Sprintf("Switch %s light signal to %d", payload, trafficSignal[payload]["on"]))

			default:
				err = m.Send("Usage: /s red|amber|green")

			}

			return err

		})

		kbot.Start()
	},
}

func init() {
	rootCmd.AddCommand(kbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
