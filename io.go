package main

import (
	"fmt"
	"strconv"
	"strings"
)

var fenCodes = map[byte]piece{
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

func pieceToString(p piece) string {
	var code string

	switch p & Piece {
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

	if p.color() == White {
		code = strings.ToUpper(code)
	}

	return code
}

func showPosition(position position) {
	fmt.Print("====BOARD====\n")
	for i := 0; i < 128; i++ {
		if i&OffBoard == 0 {
			fmt.Print(pieceToString(position.board[i]))
		}
		if (i+1)%16 == 0 {
			fmt.Print("\n")
		}
	}

	var nextMove string

	if position.toMove == White {
		nextMove = "white"
	} else {
		nextMove = "black"
	}

	castling := castleString(position)
	enPassant := enPassantString(position)

	halfMove := strconv.FormatInt(int64(position.halfmove), 10)
	fullMove := strconv.FormatInt(int64(position.fullmove), 10)

	fmt.Printf("====DETAILS====\nNext move: %v\nCastling: %v\nEn passant target: %v\nHalfmove: %v\nFullmove: %v\n\n", nextMove, castling, enPassant, halfMove, fullMove)
}

func showSliding(position position) {
	for i := 0; i < 128; i++ {
		if i&OffBoard == 0 {
			if position.board[i].isSliding() {
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

func showAttackMap(attackMap [128]byte) {
	for i := 0; i < 128; i++ {
		if i&OffBoard == 0 {
			if attackMap[i] == 1 {
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

func toMoveString(move move) string {
	if move.isQuiet() {
		return "Quiet"
	} else if move.isQueenCastle() {
		return "Castle queenside"
	} else if move.isKingCastle() {
		return "Castle kingside"
	} else if move.isPromotionCapture() {
		return "Promotion capture"
	} else if move.isPromotion() {
		return "Promotion"
	} else if move.isEnPassantCapture() {
		return "En passant capture"
	} else if move.isDoublePawnPush() {
		return "Double pawn push"
	} else if move.isCapture() {
		return "Capture"
	}
	return "Invalid"
}

func showMove(move move) {
	var formatString string

	if move == KingCastle {
		formatString = "Castle to the kingside (%v%v)\n"
	} else if move == QueenCastle {
		formatString = "Castle to the queenside (%v%v)\n"
	} else if move&Capture != 0 {
		formatString = "Capture from %v to %v\n"
	} else if move&Promotion != 0 {
		formatString = "Promotion from %v to %v\n"
	} else if move&DoublePawnPush != 0 {
		formatString = "Double pawn push from %v to %v\n"
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

func showBitboard(board uint64) {
	var i int
	var j int
	fmt.Print("\n")

	for i = 56; i >= 0; i -= 8 {
		for j = 0; j < 8; j++ {
			if (board & (1 << (uint(i) + uint(j)))) != 0 {
				fmt.Print("x")
			} else {
				fmt.Print("-")
			}
		}
		fmt.Print("\n")
	}
}

func fromFEN(fen string) position {
	sections := strings.Split(fen, " ")

	boardString := sections[0]
	playerString := sections[1]
	castleString := sections[2]
	enPassantString := sections[3]
	halfmoveString := sections[4]
	fullmoveString := sections[5]

	var startBoard [128]piece

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

	var castling byte
	castling = setCastle(castling, KingCastle, White, strings.Contains(castleString, "K"))
	castling = setCastle(castling, QueenCastle, White, strings.Contains(castleString, "Q"))
	castling = setCastle(castling, KingCastle, Black, strings.Contains(castleString, "k"))
	castling = setCastle(castling, QueenCastle, Black, strings.Contains(castleString, "q"))

	// en passant squares
	var enPassantTarget byte

	if enPassantString == "-" {
		enPassantTarget = NoEnPassant
	} else {
		fileLetter := enPassantString[0]
		rankNumber := int(enPassantString[1] - '0')

		enPassantTarget = byte(((rankNumber - 1) * 16) + int(fileLetter-'a'))
	}

	halfmoveInt, _ := strconv.Atoi(halfmoveString)
	halfmove := byte(halfmoveInt)
	fullmove, _ := strconv.Atoi(fullmoveString)

	startPosition := position{board: startBoard, toMove: toMove, castling: castling, enPassantTarget: enPassantTarget, halfmove: halfmove, fullmove: fullmove}

	return startPosition
}

func castleString(position position) string {
	var castling string

	if getCastle(position.castling, KingCastle, White) {
		castling += "K"
	}
	if getCastle(position.castling, QueenCastle, White) {
		castling += "Q"
	}
	if getCastle(position.castling, KingCastle, Black) {
		castling += "k"
	}
	if getCastle(position.castling, QueenCastle, Black) {
		castling += "q"
	}

	if len(castling) == 0 {
		castling = "-"
	}

	return castling
}

func enPassantString(position position) string {
	var enPassant string

	if position.enPassantTarget == NoEnPassant {
		enPassant = "-"
	} else {
		fileLetter := string(position.enPassantTarget%16 + 'a')
		rankNumber := (position.enPassantTarget / 16) + 1
		enPassant = fmt.Sprintf("%v%v", fileLetter, rankNumber)
	}

	return enPassant
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

	castling = castleString(position)

	enPassant = enPassantString(position)

	fen := fmt.Sprintf("%v %v %v %v %v %v", pieces, player, castling, enPassant, position.halfmove, position.fullmove)

	return fen
}
