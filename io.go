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
	}

	if getColor(piece) == White {
		code = strings.ToUpper(code)
	}

	return code
}

func showPiece(piece byte) {
	switch piece & Piece {
	case King:
		fmt.Print("K")
	case Queen:
		fmt.Print("Q")
	case Rook:
		fmt.Print("R")
	case Bishop:
		fmt.Print("B")
	case Knight:
		fmt.Print("H")
	case Pawn:
		fmt.Print("P")
	default:
		fmt.Print("_")
	}
}
func showPosition(position position) {
	for i := 0; i < 128; i++ {
		if i&OffBoard == 0 {
			showPiece(position.board[i])
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
	fen := ""

	// piece placement

	for rank := 7; rank >= 0; rank-- {
		startIndex := rank * 16

		file := 0
		empty := 0

		for {
			if file > 7 {
				if empty != 0 {
					fen += strconv.FormatInt(int64(empty), 10)
				}
				break
			}

			if !piecePresent(position, startIndex+file) {
				empty++
			} else {
				if empty != 0 {
					fen += strconv.FormatInt(int64(empty), 10)
				}
				empty = 0

				fen += pieceToString(position.board[startIndex+file])
			}

			file++
		}

		if rank != 0 {
			fen += "/"
		} else {
			fen += " "
		}

	}

	if position.toMove == White {
		fen += "w "
	} else {
		fen += "b "
	}

	if position.blackCanCastleKingside || position.blackCanCastleQueenside || position.whiteCanCastleKingside || position.whiteCanCastleQueenside {
		if position.whiteCanCastleKingside {
			fen += "K"
		}
		if position.whiteCanCastleQueenside {
			fen += "Q"
		}
		if position.blackCanCastleKingside {
			fen += "k"
		}
		if position.blackCanCastleQueenside {
			fen += "q"
		}
	} else {
		fen += "-"
	}

	fen += " "

	if position.enPassantTarget == -1 {
		fen += "-"
	} else {
		fileLetter := string(position.enPassantTarget%16 + 'a')
		rankNumber := strconv.FormatInt(int64(position.enPassantTarget/16), 10)
		fen += fileLetter
		fen += rankNumber
	}
	fen += " "

	fen += strconv.FormatInt(int64(position.halfmove), 10)
	fen += " "
	fen += strconv.FormatInt(int64(position.fullmove), 10)
	return fen
}
