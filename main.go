package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"

	"telegramalert/helpers"
	"telegramalert/node"
)

const (
	// TimeTicker dur ticker run
	TimeTicker = 30 // seconds
	// MaxTimes is a max try to sending messages
	MaxTimes = 3 //the times try sends message
)

type Client struct {
	TeleCLient *helpers.Telegram
	BlcClient  *helpers.Blockchain
	SendCount  int
}

func main() {
	app := cli.NewApp()
	app.Name = "telegramalert"
	app.Usage = "sends messages to telegram when node dont increase blocks"
	app.Version = "0.0.1"
	app.Commands = healthCheckCommand()

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func healthCheckCommand() []cli.Command {
	healthCheckCmd := cli.Command{
		Action:      start,
		Name:        "start",
		Usage:       "Alert to telegram whem block is stuck",
		Description: `Alert to telegram whem block is stuck`,
	}
	healthCheckCmd.Flags = helpers.NewTeleClientFlag()
	healthCheckCmd.Flags = append(healthCheckCmd.Flags, node.NewEvrynetNodeFlags()...)

	return []cli.Command{healthCheckCmd}
}

func start(ctx *cli.Context) {
	client := &Client{
		SendCount: 0,
	}

	teleClient, err := helpers.NewTeleClientFromFlag(ctx)
	if err != nil {
		log.Printf("can not init telegram bot %s", err.Error())
		return
	}
	log.Print("Connected to telegram")
	client.TeleCLient = teleClient

	blcClient, err := helpers.NewBlcClientFromFlags(ctx)
	if err != nil {
		log.Printf("can not connect to evrynet node %s", err.Error())
		// send SOS
		sendAlert(client, err.Error(), "SOS", false)
		return
	}
	client.BlcClient = blcClient
	log.Print("evrynet client is created")

	ticker := time.NewTicker(TimeTicker * time.Second)
	for range ticker.C {
		client.checkNodeAlive()
	}
}

func (client *Client) checkNodeAlive() {
	lastBlock, err := client.BlcClient.GetLastBlock()
	if err != nil {
		sendAlert(client, err.Error(), "SOS", false)
		return
	}

	if lastBlock == nil {
		// send alert
		sendAlert(client,
			fmt.Sprintf("[%s] Block is stuck, latest block: %d", time.Now().Format(time.RFC3339), client.BlcClient.LatestBlock),
			"SOS", false)
		client.SendCount++
		return
	}
	if lastBlock.Cmp(client.BlcClient.LatestBlock) <= 0 {
		// send alert
		sendAlert(client,
			fmt.Sprintf("[%s] Block is stuck, latest block: %d", time.Now().Format(time.RFC3339), client.BlcClient.LatestBlock),
			"SOS", false)
		return
	}

	// check if node is ok from failed
	if client.SendCount >= MaxTimes {
		sendAlert(client, fmt.Sprintf("[%s] Node is ok", time.Now().Format(time.RFC3339)), "OK", true)
	}
	client.SendCount = 0
	client.BlcClient.LatestBlock = lastBlock
	log.Printf("Current block is: %d\n", lastBlock)
}

func sendAlert(client *Client, msg string, caption string, forceSend bool) {
	if forceSend {
		// send message not increase counter
		log.Printf("================send msg: %s", msg)
		client.TeleCLient.SendMessage(msg, caption)
		return
	}
	if client.SendCount >= MaxTimes {
		return
	}
	log.Printf("================send msg: %s", msg)
	client.TeleCLient.SendMessage(msg, caption)
	client.SendCount++
}
