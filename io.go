package main

import (
	"fmt"
	"strconv"
	"strings"
)

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

	if position.whiteCanCastleKingside {
		castling += "K"
	}
	if position.whiteCanCastleQueenside {
		castling += "Q"
	}
	if position.blackCanCastleKingside {
		castling += "k"
	}
	if position.blackCanCastleQueenside {
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
