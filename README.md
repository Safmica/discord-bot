# Discord Bot

The **Discord Bot** is designed to provide an interactive gaming experience on Discord. Currently, it includes features for playing the Undercover game and handling simple message responses. The bot is built using **Golang** and integrates with the Discord API via `discordgo`.

## Features
### Undercover Game
- **Game Management**: Create, start, and manage game sessions.
- **Secret Word Assignment**: Assign unique secret words to players.
- **Ephemeral Messages**: Ensure private messages are only visible to the intended recipient.
- **Role Assignment**: Randomly distribute roles between Civilians and Undercover players.
- **Secure Interaction Handling**: Prevent unauthorized users from manipulating the game state.
### Jackheart Game
- **Game Management**: Create, start, and manage game sessions efficiently.
- **Hidden Symbol System**: Players must spend points to reveal others' symbols.
- **Strategic Lying Mechanic**: Players earn points for providing false information.
- **Point System**: Players start with a calculated point pool, with elimination at zero points and victory upon reaching the maximum threshold.
- **Voting System**: After each round, players vote, and the most voted player loses points.
- **Multiple Win Conditions**: Jack Heart wins by surviving or reaching max points, while Pions win by eliminating Jack or reaching max points.
- **Secure Interaction Handling**: Ensures game integrity by preventing unauthorized actions.

## Tech Stack
- ![VSCode](https://img.shields.io/badge/VSCode-0078D4?style=for-the-badge&logo=visual%20studio%20code&logoColor=white) **Visual Studio Code** - Used as the primary IDE for developing the bot.
- ![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white) **Golang** - The programming language used for development.
- ![Discord](https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white) **Discord API** - Enables the bot to interact with Discord servers.

## How to Run
This bot can be run locally using the following steps:

### Prerequisites
- Golang installed
- A Discord bot token
- A JSON file containing the words for the Undercover game

### Environment Variables
Set up your environment variables as follows:
```bash
export TOKEN=<your_bot_token>
export BOT_PREFIX=<your_bot_prefix>
export UNDERCOVER_WORDS=<your_path_to_word_json>
```

### JSON File Structure Undercover
Ensure your JSON file follows this structure before running the bot:
```json
{
  "words": [
    {
      "word": ["cat", "dog"],
      "used": false
    },
    {
      "word": ["apple", "orange"],
      "used": false
    }
  ]
}
```

### Steps to Run
1. Clone the repository:
```bash
git clone https://github.com/Safmica/discord-bot.git
cd discord-bot
```
2. Install dependencies:
```bash
go mod tidy
```
3. Run the bot:
```bash
go run main.go
```
4. Invite the bot to your Discord server and start playing!

## List of Commands
### Global
- `{prefix}help` - Display a list of available commands.
### Undercover Game
- `{prefix}undercover` - Open the playing room session.
- `{prefix}undercover config {config} {options}` - Configure the game settings.
### Jackheart Game
- `{prefix}jackheart` - Open the playing room session.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

