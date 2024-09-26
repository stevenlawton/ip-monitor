package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const (
	ipCheckURL    = "https://api.ipify.org" // Service to get external IP
	checkInterval = 1 * time.Hour           // Interval between IP checks
)

var (
	previousIP string
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error reading .env file")
	}

	// Now you can use os.Getenv to get the variables
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	channelID := os.Getenv("DISCORD_CHANNEL_ID")

	if botToken == "" || channelID == "" {
		log.Fatal("Bot token or channel ID not set. Please set DISCORD_BOT_TOKEN and DISCORD_CHANNEL_ID.")
	}

	// Create a new Discord session
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Open a WebSocket connection to Discord
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection to Discord: %v", err)
	}
	defer func(dg *discordgo.Session) {
		err := dg.Close()
		if err != nil {
			log.Fatalf("Error closing connection to Discord: %v", err)
		}
	}(dg)

	log.Println("Bot is now running. Press CTRL+C to exit.")

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start the IP monitoring in a goroutine
	go func() {
		for {
			currentIP, err := getExternalIP()
			if err != nil {
				log.Printf("Error fetching external IP: %v", err)
				time.Sleep(checkInterval)
				continue
			}

			if currentIP != previousIP {
				log.Printf("IP changed from %s to %s", previousIP, currentIP)
				err = updateDiscordMessage(dg, channelID, currentIP)
				if err != nil {
					log.Printf("Error updating Discord message: %v", err)
				}
				previousIP = currentIP
			}

			time.Sleep(checkInterval)
		}
	}()

	// Wait for a termination signal
	<-stop
	log.Println("Shutting down bot.")
}

func getExternalIP() (string, error) {
	resp, err := http.Get(ipCheckURL)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing connection to IP Checker: %v", err)
		}
	}(resp.Body)

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

func updateDiscordMessage(dg *discordgo.Session, channelID, currentIP string) error {
	// Fetch pinned messages
	pinnedMessages, err := dg.ChannelMessagesPinned(channelID)
	if err != nil {
		return fmt.Errorf("error fetching pinned messages: %w", err)
	}

	var pinnedMessageID string
	if len(pinnedMessages) > 0 {
		for _, message := range pinnedMessages {
			// Assume the first pinned message is the one to update
			if strings.Contains(message.Content, "Current IP Address:") {
				pinnedMessageID = message.ID
				_, err = dg.ChannelMessageEdit(channelID, pinnedMessageID, fmt.Sprintf("Current IP Address: `%s`", currentIP))
				if err != nil {
					return fmt.Errorf("error editing pinned message: %w", err)
				}
				// Send a broadcast message
				_, err = dg.ChannelMessageSend(channelID, fmt.Sprintf("IP Address has changed to `%s`", currentIP))
				if err != nil {
					return fmt.Errorf("error sending broadcast message: %w", err)
				}

				return nil
			}
		}
	}
	// Send a new message and pin it
	msg, err := dg.ChannelMessageSend(channelID, fmt.Sprintf("Current IP Address: `%s`", currentIP))
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	err = dg.ChannelMessagePin(channelID, msg.ID)
	if err != nil {
		return fmt.Errorf("error pinning message: %w", err)
	}

	// Send a broadcast message
	_, err = dg.ChannelMessageSend(channelID, fmt.Sprintf("IP Address has changed to `%s`", currentIP))
	if err != nil {
		return fmt.Errorf("error sending broadcast message: %w", err)
	}

	return nil
}
