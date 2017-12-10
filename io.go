package main

import (
	"fmt"
	"strconv"
	"strings"
)

var fenCodes = map[byte]byte{
	'k': 67,
	'q': 71,
	'b': 68,
	'n': 66,
	'r': 69,
	'p': 65,
	'K': 3,
	'Q': 7,
	'B': 4,
	'N': 2,
	'R': 5,
	'P': 1,
}

func pieceToString(piece byte) string {
	var code string

	switch piece & Piece {
	case King:
		code = "k"
	case Queen:
		code = "q"
	case Rook:
		code = "r"
	case Bishop:
		code = "b"
	case Knight:
		code = "n"
	case Pawn:
		code = "p"
	default:
		code = "_"
	}

	if getColor(piece) == White {
		code = strings.ToUpper(code)
	}

	return code
}

func showPosition(position position) {
	for i := 0; i < 128; i++ {
		if i&OffBoard == 0 {
			fmt.Print(pieceToString(position.board[i]))
		}
		if (i+1)%16 == 0 {
			fmt.Print("\n")
		}
	}
}

func showSliding(position position) {
	for i := 0; i < 128; i++ {
		if i&OffBoard == 0 {
			if isSliding(position.board[i]) {
				fmt.Print("T")
			} else {
				fmt.Print("F")
			}
		}

		if (i+1)%16 == 0 {
			fmt.Print("\n")
		}
	}
}

func showMove(move move) {
	var formatString string

	if move == KingCastle {
		formatString = "Castle to the kingside (%v%v)\n"
	} else if move == QueenCastle {
		formatString = "Castle to the queenside (%v%v)\n"
	} else if move&Capture != 0 {
		formatString = "Capture from %v to %v\n"
	} else if move&DoublePawnPush != 0 {
		formatString = "Double pawn push from %v to %v\n"
	} else if move&Promotion != 0 {
		formatString = "Promotion from %v to %v\n"
	} else {
		formatString = "Quiet move from %v to %v\n"
	}

	fmt.Printf(formatString, move.From(), move.To())
}

func showMoves(moves []move) {
	for i := 0; i < len(moves); i++ {
		showMove(moves[i])
	}
}

func fromFen(fen string) position {
	sections := strings.Split(fen, " ")

	boardString := sections[0]
	playerString := sections[1]
	castleString := sections[2]
	enPassantString := sections[3]
	halfmoveString := sections[4]
	fullmoveString := sections[5]

	var startBoard [128]byte

	fmt.Printf("Converting board position: %v\n", boardString)
	boardIndex := 112
	for _, char := range boardString {
		// skip slashes seperating ranks
		//fmt.Printf("Starting at index %v\n", boardIndex)
		if char == '/' {
			//fmt.Print("Skipping /\n")
			continue
		}

		if char >= '1' && char <= '8' {
			//fmt.Printf("Adding %v blank squares\n", char-'0')
			boardIndex += int(char - '0')
		} else {
			//fmt.Printf("Adding piece %v\n", fenCodes[byte(char)])
			startBoard[boardIndex] = fenCodes[byte(char)]
			boardIndex++
		}

		// skip squares not on the board
		if boardIndex%16 > 7 {
			boardIndex = ((boardIndex / 16) - 1) * 16
			//fmt.Printf("Skipping off the board to %v\n", boardIndex)
		}
	}

	var toMove byte

	if playerString == "w" {
		toMove = White
	} else {
		toMove = Black
	}

	var castling = map[byte]map[int]bool{
		White: map[int]bool{KingCastle: strings.Contains(castleString, "K"), QueenCastle: strings.Contains(castleString, "Q")},
		Black: map[int]bool{KingCastle: strings.Contains(castleString, "k"), QueenCastle: strings.Contains(castleString, "q")},
	}

	// en passant squares
	var enPassantTarget int

	if enPassantString == "-" {
		enPassantTarget = -1
	} else {
		fileLetter := enPassantString[0]
		rankNumber := int(enPassantString[1] - '0')

		enPassantTarget = (rankNumber * 16) + int(fileLetter-'a')
	}

	halfmove, _ := strconv.Atoi(halfmoveString)
	fullmove, _ := strconv.Atoi(fullmoveString)

	startPosition := position{board: startBoard, toMove: toMove, castling: castling, enPassantTarget: enPassantTarget, halfmove: halfmove, fullmove: fullmove}

	return startPosition
}

func toFEN(position position) string {
	var pieces string
	var player string
	var castling string
	var enPassant string

	for rank := 7; rank >= 0; rank-- {
		empty := 0

		for i := rank * 16; i < rank*16+8; i++ {
			if piecePresent(position, i) {
				if empty != 0 {
					pieces += strconv.FormatInt(int64(empty), 10)
				}

				empty = 0

				pieces += pieceToString(position.board[i])
			} else {
				empty++
			}
		}

		if empty != 0 {
			pieces += strconv.FormatInt(int64(empty), 10)
		}

		if rank != 0 {
			pieces += "/"
		}
	}

	if position.toMove == White {
		player = "w"
	} else {
		player = "b"
	}

	if position.castling[White][KingCastle] {
		castling += "K"
	}
	if position.castling[White][QueenCastle] {
		castling += "Q"
	}
	if position.castling[Black][KingCastle] {
		castling += "k"
	}
	if position.castling[Black][QueenCastle] {
		castling += "q"
	}

	if len(castling) == 0 {
		castling = "-"
	}

	if position.enPassantTarget == -1 {
		enPassant = "-"
	} else {
		fileLetter := string(position.enPassantTarget%16 + 'a')
		rankNumber := position.enPassantTarget / 16
		enPassant = fmt.Sprintf("%v%v", fileLetter, rankNumber)
	}

	fen := fmt.Sprintf("%v %v %v %v %v %v", pieces, player, castling, enPassant, position.halfmove, position.fullmove)

	return fen
}
