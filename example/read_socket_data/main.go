package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	bngsdk "github.com/ESilva15/gobngsdk"
)

func main() {
	fmt.Println("Example of how to read from the socket with the SDK")

	output, err := os.OpenFile("./output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o755)
	if err != nil {
		log.Fatalf("Failed to open log file: %+v", err)
	}

	logger := slog.New(
		slog.NewTextHandler(output, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)

	sdk, err := bngsdk.NewBngSDK(bngsdk.Options{
		Logger:           logger.With("service", "bngsdk"),
		SourceType:       bngsdk.UDPData,
		ImportUDPAddress: "127.0.0.1",
		ImportUDPPort:    4444,
	})
	if err != nil {
		log.Fatalf("failed to open the sdk: %+v", err)
	}
	defer sdk.Close()

	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := sdk.Update()
			if err != nil {
				log.Fatalf("failed to update data: %+v", err)
			}

			fmt.Printf("\033[?25l\033[2J\033[H")
			fmt.Printf(
				"Gear: %d, RPM: %f, Speed: %f\n",
				sdk.Data.Gear, sdk.Data.RPM, sdk.Data.Speed,
			)
		}
	}
}
