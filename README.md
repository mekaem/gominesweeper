# gominesweeper

A Minesweeper implementation in Go (interview question)

## Prompt

Minesweeper is a classic single-player puzzle game. The objective of the game is to clear a rectangular board containing hidden mines without detonating any of them.

## Basic Rules

1. **Game Setup**: The game board is initially covered with tiles. Some of these tiles contain mines, while others are safe. The player's objective is to uncover all the safe tiles without triggering any mines.

2. **Game Board**: The game board consists of a grid of square tiles. Each tile can be uncovered or covered. Initially, all tiles are covered.

3. **Mines**: Mines are hidden beneath some of the tiles on the game board. The number of mines is predetermined.

4. **Numbers**: Tiles adjacent to mines contain a number indicating the number of mines in the surrounding tiles. For example, if a tile contains the number "3," it means that there are three mines in the adjacent tiles.

5. **Uncovering Tiles**: The player uncovers tiles by clicking on them.
   - If a player uncovers a tile containing a mine, the game ends, and the player loses.
   - If the tile does not contain a mine, it will reveal a number indicating the number of adjacent mines, or it will be blank if there are no adjacent mines.
   - When a tile with 0 adjacent mines is revealed, it and all adjacent tiles are automatically revealed.

6. **Winning the Game**: The player wins the game when all safe tiles are uncovered and only mines remain uncovered.

## Task

Build a program that allows a person to play the game Minesweeper. Please focus on the internal game logic rather than the presentation / UI layer. First create the board and place mines randomly, and then focus on the game play.

## Input Variables

- `h` - board height in # of cells
- `w` - board width in # of cells
- `m` - number of mines
