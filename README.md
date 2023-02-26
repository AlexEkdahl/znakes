# TCP Snake Game

This is a network game of snakes that allows multiple players to compete against each other over TCP. Players control a snake that moves around the game board, eating fruit and avoiding obstacles. The player with the highest score at the end of the game wins.

## Todos

- [ ] Improve game performance
- [ ] Add support for multiple game modes
- [ ] Create a tutorial for new users
- [ ] Refactor code for better readability
- [ ] Write more unit tests
- [ ] Implement score system
- [ ] Integrate with a database for high scores
- [ ] Figure out how to halt the game with only one player in multiplayer mode (one player leaves)
- [ ] Add collision detection
- [ ] Add fruits

## Getting Started

To build the game, you will need to have Go installed on your system. You can download Go from the official website at [golang.org](https://golang.org).

Once you have Go installed, you can build the server and client binaries by running `make build`. This will compile the server and client binaries and put them in the `bin/` directory.

You can then start the server by running `make run-server`, and start the client by running `make run-client`. Alternatively, you can start both the server and client by running `make run`.

Players can connect to the server using any TCP client that supports the Telnet protocol.

## Gameplay

- `k` or `up arrow`: move the snake up
- `j` or `down arrow`: move the snake down
- `h` or `left arrow`: move the snake left
- `l` or `right arrow`: move the snake right

The player's score increases each time the snake eats a piece of fruit. The game ends when all players have died, or when the maximum game time has been reached. The player with the highest score at the end of the game is declared the winner.

## Contributing

Contributions to the game are welcome! If you have an idea for a new feature, or if you find a bug that needs fixing, please open an issue or submit a pull request.

## File Structure

- `cmd/`: contains the main executable files for the server and client.
- `pkg/game`: contains the core game logic, including the `Game`, `Player`, and `Snake` structs.
- `pkg/network`: contains the code for the server and client network connections, as well as the `Message` and `Print` structs and protobuf message definitions.
- `Makefile`: contains commands for building, running, and testing the game.
