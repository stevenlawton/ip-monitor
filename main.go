package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/docker/docker/client"
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

				err = restartARKServer()
				if err != nil {
					log.Printf("Error restarting ARK: %v", err)
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

func restartARKServer() error {
	// Connect to Docker daemon
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	restartContainers(context.Background(), cli)

	return nil
}

// RestartContainers stops and starts containers based on specific criteria
func restartContainers(ctx context.Context, cli *client.Client) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		log.Fatalf("Error listing containers: %v", err)
	}

	containerIDs := make([]string, 10)

	for _, ct := range containers {
		if strings.Contains(ct.Image, "acekorneya/asa_server") {
			for _, name := range ct.Names {
				if strings.Contains(name, "asa_ARK_Lawton") || slices.Contains(containerIDs, ct.ID) {
					containerIDs = append(containerIDs, ct.ID)
					fmt.Printf("Container %s(%s) stopping...\n", ct.Names, ct.ID)

					err := cli.ContainerStop(ctx, ct.ID, container.StopOptions{})
					if err != nil {
						log.Printf("Error stopping container %s: %v\n", ct.ID, err)
						continue
					}

					// Wait for the container to stop
					if waitForStatus(ctx, cli, ct.ID, "exited") {
						fmt.Printf("Container %s stopped successfully.\n", ct.ID)
					} else {
						log.Printf("Timeout waiting for container %s to stop.\n", ct.ID)
						continue
					}

					fmt.Printf("Container %s starting...\n", ct.ID)
					err = cli.ContainerStart(ctx, ct.ID, container.StartOptions{})
					if err != nil {
						log.Printf("Error starting container %s: %v\n", ct.ID, err)
						continue
					}

					// Wait for the container to start
					if waitForStatus(ctx, cli, ct.ID, "running") {
						fmt.Printf("Container %s started successfully.\n", ct.ID)
					} else {
						log.Printf("Timeout waiting for container %s to start.\n", ct.ID)
					}
				}
			}
		}
	}
}

// waitForStatus waits until the container status matches the desired status
func waitForStatus(ctx context.Context, cli *client.Client, containerID string, desiredStatus string) bool {
	timeout := time.Now().Add(30 * time.Second) // Adjust timeout as needed

	for time.Now().Before(timeout) {
		ct, err := cli.ContainerInspect(ctx, containerID)
		if err != nil {
			log.Printf("Error inspecting container %s: %v\n", containerID, err)
			return false
		}

		if ct.State.Status == desiredStatus {
			return true
		}

		time.Sleep(1 * time.Second) // Avoid busy waiting
	}

	return false
}
