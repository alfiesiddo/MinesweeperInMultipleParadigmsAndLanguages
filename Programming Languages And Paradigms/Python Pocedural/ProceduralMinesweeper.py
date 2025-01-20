import os
import random
from colorama import init, Fore, Style
import time

size = 5
mines = 7
unTouchedSpaces = size * size
colWidth = 2
shield = False

def createBlankBoard(size):
    board = [['+' for _ in range(size)] for _ in range(size)]
    return board

def createComparisonBoard(size, mines):
    board = [[0 for _ in range(size)] for _ in range(size)]

    mineCount = 0

    while mineCount < mines:
        row = random.randint(0, size - 1)
        column = random.randint(0, size - 1)

        if board[row][column] == 0:  
            board[row][column] = 1  
            mineCount = mineCount + 1
    return board
    
compareBoard = createComparisonBoard(size, mines)
userBoard = createBlankBoard(size)

def DisplayBoard(board): 
    global mines
    global colWidth
    global unTouchedSpaces
    size = len(board)
    # Print each column with column number
    print(f"There are {mines} mines to be found! \n")
    print(unTouchedSpaces)
    print('  ' * colWidth + ' '.join(f'{j:>{colWidth}}' for j in range(size)) + "\n")
    
    # Print each row with row number
    for i, row in enumerate(board):
        print(f'{i:>{colWidth}}  ' + ' '.join(f'{cell:>{colWidth}}' for cell in row))  
    return


def updateBoard():
    global compareBoard
    global userBoard
    global shield
    global unTouchedSpaces

    if unTouchedSpaces == mines:
                os.system('cls')
                print("You Have Won!")
                input("Press any key to try again!")
                unTouchedSpaces = size * size
                compareBoard = createComparisonBoard(size, mines)
                userBoard = createBlankBoard(size)

    userInput = input("Enter row (or 'S' for shield, 'H' for hint): ")
    
    if userInput.upper() == "S":
        shield = True
        print("Shield activated! Avoiding the next mine.")
        return  # Exit early after activating shield

    if userInput.upper() == "H":
        giveHint()
        return  # Exit after revealing hint

    try:
        inputRow = int(userInput)
        inputCol = int(input("Enter column: "))
    except ValueError:
        # Prevents user entering anything other than an int
        return 

    # Makes sure user can't enter a coordinate that doesn't exist
    if 0 <= inputRow < size and 0 <= inputCol < size:
        if compareBoard[inputRow][inputCol] == 1:
            os.system('cls')
            if shield:
                shield = False
                print("Shield used! You avoided the mine!")
                userBoard[inputRow][inputCol] = "S"  # Mark the shielded tile
            else:
                print("You Hit a Mine, Game Over")
                input("Press any key to try again!")
                unTouchedSpaces = size * size
                compareBoard = createComparisonBoard(size, mines)
                userBoard = createBlankBoard(size)
        else:
            revealTiles(inputRow, inputCol)

def checkSurround(x, y):
    count = 0 
    # directions to check for mines in
    directions = [(-1,-1), (-1,0), (-1, 1),
                  (0, -1),         (0, 1),
                  (1, -1), (1, 0), (1, 1)]
    
    for directionX, directionY in directions:
        neighbourX, neighbourY = x + directionX, y + directionY

        # check if neighbour cell is within board
        if 0 <= neighbourX < size and 0 <= neighbourY < size:
            # amount of mines
            count = count + compareBoard[neighbourX][neighbourY]
    return count


def revealTiles(x, y):
    global unTouchedSpaces
    if userBoard[x][y] != '+': 
        return
    unTouchedSpaces -= 1
    surroundingMines = checkSurround(x, y)  # Check for surrounding mines
    
    if surroundingMines == 0:
        userBoard[x][y] = ' '  # if none surrounding then show selected tile as blank
        directions = [(-1,-1), (-1,0), (-1,1), 
                      (0,-1),           (0,1), 
                      (1,-1), (1,0),    (1,1)]
        
        for directionX, directionY in directions:
            neighbourX, neighbourY = x + directionX, y + directionY
            if 0 <= neighbourX < size and 0 <= neighbourY < size:
                revealTiles(neighbourX, neighbourY)  # Recursively reveal adjacent tiles
    else:
        stringifiedSurroundingMines = str(surroundingMines)
        formattedMines = f" {stringifiedSurroundingMines}"
        match(surroundingMines):
            case 1:
                userBoard[x][y] = Fore.RED + formattedMines + Style.RESET_ALL
            case 2:
                userBoard[x][y] = Fore.LIGHTRED_EX + formattedMines + Style.RESET_ALL
            case 3:
                userBoard[x][y] = Fore.YELLOW + formattedMines + Style.RESET_ALL
            case 4:
                userBoard[x][y] = Fore.LIGHTYELLOW_EX + formattedMines + Style.RESET_ALL
            case 5:
                userBoard[x][y] = Fore.GREEN + formattedMines + Style.RESET_ALL
            case 6:
                userBoard[x][y] = Fore.LIGHTGREEN_EX + formattedMines + Style.RESET_ALL
            case 7:
                userBoard[x][y] = Fore.CYAN + formattedMines + Style.RESET_ALL
            case 8:
                userBoard[x][y] = Fore.BLUE + formattedMines + Style.RESET_ALL

def giveHint():
    os.system('cls')
    print("Hint: Revealing some mines for 3 seconds!")
    revealed = []
    for _ in range(random.randint(1, 3)):  # Reveal between 1 and 3 mines
        while True:
            row = random.randint(0, size - 1)
            col = random.randint(0, size - 1)
            if compareBoard[row][col] == 1 and (row, col) not in revealed:
                revealed.append((row, col))
                userBoard[row][col] = "M"  # Temporarily mark the mine
                break
    DisplayBoard(userBoard)
    time.sleep(3)

    # Hide the revealed mines again
    for row, col in revealed:
        userBoard[row][col] = '+'

while True:
    os.system('cls')
    DisplayBoard(userBoard)
    updateBoard()
