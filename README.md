# IP Monitor Discord Bot

This Go application monitors your external IP address and updates a pinned message in a specified Discord channel whenever the IP changes. It also sends a broadcast message notifying about the IP change.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Setup Instructions](#setup-instructions)
  - [1. Clone the Repository](#1-clone-the-repository)
  - [2. Create a Discord Bot](#2-create-a-discord-bot)
  - [3. Create a `.env` File](#3-create-a-env-file)
  - [4. Install Dependencies](#4-install-dependencies)
  - [5. Build and Run Locally (Optional)](#5-build-and-run-locally-optional)
  - [6. Build the Docker Image](#6-build-the-docker-image)
  - [7. Run the Docker Container](#7-run-the-docker-container)
- [Environment Variables](#environment-variables)
- [Stopping the Bot](#stopping-the-bot)
- [Updating the Bot](#updating-the-bot)
- [Additional Considerations](#additional-considerations)
  - [Security Notes](#security-notes)
  - [Adjusting the Check Interval](#adjusting-the-check-interval)
  - [Logging and Monitoring](#logging-and-monitoring)
  - [Error Handling](#error-handling)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- **Periodic IP Monitoring**: Checks your external IP address at a configurable interval.
- **Discord Integration**: Updates a pinned message in a specified Discord channel when the IP changes.
- **Broadcast Notifications**: Sends a message to notify channel members of IP changes.
- **Environment Configuration**: Uses a `.env` file for configuration, keeping sensitive data secure.
- **Dockerized**: Easily build and run the application within a Docker container.

---

## Prerequisites

- **Go**: (Optional for local development) [Install Go](https://golang.org/doc/install) if you plan to build and run the application locally.
- **Docker**: Ensure Docker is installed on your system.
- **Discord Account**: You need a Discord account to create a bot and add it to your server.
- **Discord Server**: Access to a Discord server where you can add the bot.
- **Discord Bot Token and Channel ID**: You'll need these to configure the bot.

---

## Setup Instructions

### 1. Clone the Repository

Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/ip-monitor.git
cd ip-monitor
```

### 2. Create a Discord Bot

1. **Go to the Discord Developer Portal**: [Discord Developer Portal](https://discord.com/developers/applications)
2. **Create a New Application**:
  - Click on **"New Application"**.
  - Enter a name for your application and click **"Create"**.
3. **Add a Bot to Your Application**:
  - Navigate to the **"Bot"** tab on the left.
  - Click **"Add Bot"** and confirm.
4. **Copy the Bot Token**:
  - Under the bot's username, click **"Copy"** to copy the bot's token.
  - **Important**: Keep this token secure!
5. **Invite the Bot to Your Server**:
  - Go to the **"OAuth2"** tab and then **"URL Generator"**.
  - Under **"Scopes"**, select **"bot"**.
  - Under **"Bot Permissions"**, select the following permissions:
    - **Send Messages**
    - **Manage Messages**
    - **Read Message History**
  - Copy the generated URL and paste it into your browser to invite the bot to your server.
6. **Get the Channel ID**:
  - Enable **Developer Mode** in Discord settings (User Settings > Advanced > Developer Mode).
  - Right-click the channel where you want the bot to post and select **"Copy ID"**.

### 3. Create a `.env` File

In the project root directory, create a `.env` file to store your environment variables:

```dotenv
DISCORD_BOT_TOKEN=your_bot_token
DISCORD_CHANNEL_ID=your_channel_id
```

- Replace `your_bot_token` with the bot token you copied earlier.
- Replace `your_channel_id` with the ID of the channel where the bot will operate.

**Important:** Do not share this file or commit it to version control.

### 4. Install Dependencies

If you plan to build and run the application locally (without Docker), install the required Go packages:

```bash
go mod download
```

### 5. Build and Run Locally (Optional)

To build and run the application locally without Docker:

```bash
go build -o ip-monitor .
./ip-monitor
```

### 6. Build the Docker Image

Build the Docker image using the provided `Dockerfile`:

```bash
docker build -t ip-monitor .
```

### 7. Run the Docker Container

Run the Docker container, passing the `.env` file:

```bash
docker run -d \
  --name ip-monitor \
  --env-file .env \
  ip-monitor
```

---

## Environment Variables

- **DISCORD_BOT_TOKEN**: Your Discord bot token (required).
- **DISCORD_CHANNEL_ID**: The ID of the Discord channel where the bot will post messages (required).

---

## Stopping the Bot

To stop and remove the Docker container:

```bash
docker stop ip-monitor
docker rm ip-monitor
```

---

## Updating the Bot

If you make changes to the code:

1. **Rebuild the Docker Image**:

   ```bash
   docker build -t ip-monitor .
   ```

2. **Restart the Container**:

   ```bash
   docker stop ip-monitor
   docker rm ip-monitor
   docker run -d \
     --name ip-monitor \
     --env-file .env \
     ip-monitor
   ```

---

## Additional Considerations

### Security Notes

- **Protect Your `.env` File**: Ensure your `.env` file is not committed to version control. It contains sensitive information.
- **Environment Variables in Production**: For production environments, consider using a secrets manager or setting environment variables directly in your deployment environment.
- **Bot Permissions**: Ensure your bot has the correct permissions in the Discord channel:
  - **Send Messages**
  - **Manage Messages** (to pin messages)
  - **Read Message History**

### Adjusting the Check Interval

- The `checkInterval` constant in `main.go` controls how often the application checks for IP changes.
- Default is set to check every 5 minutes:

  ```go
  const checkInterval = 5 * time.Minute
  ```

- Adjust as needed, then rebuild the Docker image.

### Logging and Monitoring

- **View Logs**: To view the application logs, use:

  ```bash
  docker logs -f ip-monitor
  ```

- **Enhanced Logging**: Consider integrating a logging framework like `logrus` for more advanced logging capabilities.

### Error Handling

- The application includes basic error handling.
- For production use, consider implementing retry mechanisms and more sophisticated error handling.

---

## Troubleshooting

- **Bot Not Responding**:
  - Verify the bot is online in your Discord server.
  - Ensure the bot has the necessary permissions.
  - Check the logs for errors:

    ```bash
    docker logs ip-monitor
    ```

- **Environment Variables Not Set**:
  - Ensure the `DISCORD_BOT_TOKEN` and `DISCORD_CHANNEL_ID` are correctly set in your `.env` file.
  - Verify that the `.env` file is in the project root and is correctly formatted.

- **Cannot Connect to Discord**:
  - Ensure your network allows outbound connections to Discord's API.
  - Check for firewall or network restrictions.

---

## Contributing

Contributions are welcome! Please follow these steps:

1. **Fork the Repository**
2. **Create a Feature Branch**

   ```bash
   git checkout -b feature/YourFeature
   ```

3. **Commit Your Changes**

   ```bash
   git commit -m "Add your message here"
   ```

4. **Push to the Branch**

   ```bash
   git push origin feature/YourFeature
   ```

5. **Open a Pull Request**

---

## License

This project is licensed under the [MIT License](LICENSE).

---

## Additional Enhancements (Optional)

If you're interested in extending the functionality of the bot, consider the following enhancements:

### Implement a Dynamic DNS Service

- **Why**: Instead of sharing your raw IP address, use a domain name that updates automatically.
- **How**: Integrate with a dynamic DNS provider like DuckDNS or No-IP.

### Store IP Change History

- **Why**: Keeping a log of IP changes can be useful for auditing.
- **How**: Append changes to a local file or store them in a database.

### Multiple Notification Channels

- **Why**: Notify users across different platforms.
- **How**: Integrate with Slack, Telegram, or send email notifications.

### Improve Error Reporting

- **How**: Use an error tracking service like Sentry to capture and analyze exceptions.

---

## Acknowledgments

- **[discordgo](https://github.com/bwmarrin/discordgo)**: A Go package that provides low-level bindings to the Discord chat client API.

---

## Contact

- **Author**: Your Name
- **Email**: your.email@example.com
- **GitHub**: [yourusername](https://github.com/yourusername)
