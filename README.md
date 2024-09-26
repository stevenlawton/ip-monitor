# IP Monitor Discord Bot

This Go application monitors your external IP address and updates a pinned message in a specified Discord channel whenever the IP changes. It also sends a broadcast message notifying about the IP change.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Setup Instructions](#setup-instructions)
    - [1. Clone the Repository](#1-clone-the-repository)
    - [2. Set Up Environment Variables](#2-set-up-environment-variables)
    - [3. Build the Docker Image](#3-build-the-docker-image)
    - [4. Run the Docker Container](#4-run-the-docker-container)
- [Environment Variables](#environment-variables)
- [Stopping the Bot](#stopping-the-bot)
- [Updating the Bot](#updating-the-bot)
- [Notes](#notes)
- [Troubleshooting](#troubleshooting)
- [License](#license)

## Prerequisites

- **Docker**: Ensure Docker is installed on your system.
- **Discord Bot**: You need a Discord bot token and the ID of the channel where the bot will operate.

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/ip-monitor.git
cd ip-monitor
```

### 2. Create a Discord Bot

- Go to the [Discord Developer Portal](https://discord.com/developers/applications).
- Create a new application and add a bot to it.
- Copy the bot's token; you'll need it for authentication.
- Invite the bot to your server with appropriate permissions (send messages, manage messages).