#include <iostream>
#include <vector>
#include <ctime>
#include <cstdlib>
using namespace std;

class Tile {
private:
    bool isMine = false;
    bool isRevealed = false;
    int adjacentMines = 0;

public:
    Tile() {
        isMine = false;
        isRevealed = false;
        adjacentMines = 0;
    }

    void makeMine() {
        isMine = true;
    }

    bool IsMineCheck() const {
        return isMine;
    }

    bool isRevealedCheck() const {
        return isRevealed;
    }

    void setAdjacentMines(int count) {
        adjacentMines = count;
    }

    int getAdjacentMines() const {
        return adjacentMines;
    }

    string revealTile() {
        isRevealed = true;

        if (isMine) {
            return "You Lost, You Hit a Mine!";
        }
        else {
            return "Safe";
        }
    }
};

class Board {
private:
    int size;
    int maxMines;
    int hintCount = 0;
    bool playerHint = false;
    bool shieldUsed = false;
    vector<vector<Tile>> logicBoard;
    vector<vector<char>> drawboard;

    void resetPlayerHint() {
        playerHint = false;
    }

public:
    Board(int Size, int MaxMines) : size(Size), maxMines(MaxMines) {
        logicBoard.resize(size, vector<Tile>(size));
    }

    void drawBoard() {
        drawboard.resize(size, vector<char>(size, ' '));
        // Print column numbers (top row)
        cout << "    "; // Padding for row numbers
        for (int i = 0; i < size; i++) {
            cout << i << "  "; // Print column number
        }
        cout << "\n" << endl;
        // Print the board with row numbers
        for (int i = 0; i < size; i++) {
            cout << i << "   ";
            for (int j = 0; j < size; j++) {
                if (playerHint && logicBoard[i][j].IsMineCheck()) {
                    drawboard[i][j] = 'M';
                }
                else if (logicBoard[i][j].isRevealedCheck()) {
                    if (logicBoard[i][j].IsMineCheck()) {
                        drawboard[i][j] = 'S'; // Revealed mine
                    }
                    else if (logicBoard[i][j].getAdjacentMines() == 0) {
                        drawboard[i][j] = ' ';
                    }
                    else {
                        drawboard[i][j] = '0' + logicBoard[i][j].getAdjacentMines(); // Adjacent mines count
                    }
                }
                else {
                    drawboard[i][j] = '+'; // Not revealed
                }
                cout << drawboard[i][j] << "  ";
            }
            cout << endl;
        }
        resetPlayerHint();
    }

    void setLogicBoard() {
        int placedMines = 0;

        while (placedMines < maxMines) {
            int row = rand() % size;
            int col = rand() % size;

            if (!logicBoard[row][col].IsMineCheck()) {
                logicBoard[row][col].makeMine();
                placedMines++;
            }
        }

        for (int i = 0; i < size; i++) {
            for (int j = 0; j < size; j++) {
                if (!logicBoard[i][j].IsMineCheck()) {
                    logicBoard[i][j].setAdjacentMines(countAdjacentMines(i, j));
                }
            }
        }
    }

    int countAdjacentMines(int row, int col) {
        int count = 0;

        // Check all neighbors
        for (int dr = -1; dr <= 1; dr++) {
            for (int dc = -1; dc <= 1; dc++) {
                if ((dr || dc) && row + dr >= 0 && row + dr < size && col + dc >= 0 && col + dc < size) {
                    count += logicBoard[row + dr][col + dc].IsMineCheck();
                }
            }
        }
        return count;
    }

    string revealTile(int row, int col, bool& shieldActive) {
        if (row < 0 || row >= size || col < 0 || col >= size) {
            return "Invalid coordinates!";
        }

        if (logicBoard[row][col].isRevealedCheck()) {
            return "Tile already revealed!";
        }

        string result = logicBoard[row][col].revealTile();

        if (result == "You Lost, You Hit a Mine!") {
            if (shieldActive && !shieldUsed) {
                shieldUsed = true;
                shieldActive = false;
                return "Shield activated! You are safe this time.";
            }
            return result; // Game over
        }

        if (logicBoard[row][col].getAdjacentMines() == 0) {
            // Recursively reveal adjacent tiles
            for (int dr = -1; dr <= 1; dr++) {
                for (int dc = -1; dc <= 1; dc++) {
                    if ((dr || dc) && row + dr >= 0 && row + dr < size && col + dc >= 0 && col + dc < size) {
                        revealTile(row + dr, col + dc, shieldActive);
                    }
                }
            }
        }

        return result;
    }

    bool allSafeTilesRevealed() {
        for (int i = 0; i < size; i++) {
            for (int j = 0; j < size; j++) {
                if (!logicBoard[i][j].IsMineCheck() && !logicBoard[i][j].isRevealedCheck()) {
                    return false;
                }
            }
        }
        return true;
    }

    bool playerHintActivate(char hint) {
        if ((hint == 'h' || hint == 'H') && hintCount == 0) {
            playerHint = true;  // Activate hint
            hintCount = 1;
            return playerHint;
        }
        else {
            return false;  // No action taken
        }
    }

    bool activateShield(char input) {
        if ((input == 's' || input == 'S') && !shieldUsed) {
            return true;  // Shield activated
        }
        return false;
    }
};

class Game {
private:
    Board board;
    bool gameOver;
    bool shieldActive;

public:
    Game(int size, int maxMines) : board(size, maxMines), gameOver(false), shieldActive(false) {}

    void start() {
        board.setLogicBoard();

        while (!gameOver) {
            system("cls");
            board.drawBoard();

            int row, col;
            char input;

            cout << "Enter the row and column to reveal (e.g., 1 1) \nOr Enter 'H' for a Hint or 'S' to activate the shield: ";
            cin >> input;

            if (isdigit(input)) {
                cin.putback(input);
                cin >> row >> col;
            }
            else if (isalpha(input)) {
                if (input == 'H' || input == 'h') {
                    board.playerHintActivate(input);
                }
                else if (input == 'S' || input == 's') {
                    shieldActive = board.activateShield(input);
                    if (shieldActive) {
                        cout << "Shield activated!" << endl;
                        continue;
                    }
                    else {
                        cout << "Shield already used or unavailable!" << endl;
                        continue;
                    }
                }
                continue;
            }

            string result = board.revealTile(row, col, shieldActive);

            if (result == "You Lost, You Hit a Mine!") {
                gameOver = true;
                cout << result << endl;
                break;
            }
            else if (board.allSafeTilesRevealed()) {
                gameOver = true;
                cout << "You've Won!" << endl;
            }
        }
    }
};

int main() {
    int size, mines;
    cout << "Enter the board size: (Beginner usually is 8x8 and hard, 30x30) ";
    cin >> size;
    mines = size * 2;

    Game game(size, mines);
    game.start();

    return 0;
}
