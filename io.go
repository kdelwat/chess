package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Used to extract the piece identity information.
const pieceIdentity = 0x0F

// Maps the piece code (in byte form) to the correct string representation in
// FEN.
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

func fromFEN(fen string) position {
	// The FEN for the starting position looks like this:
	// 	rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1

	// Split it on spaces
	sections := strings.Split(fen, " ")

	boardString := sections[0]     // rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR
	playerString := sections[1]    // w
	castleString := sections[2]    // KQkq
	enPassantString := sections[3] // -
	halfmoveString := sections[4]  // 0
	fullmoveString := sections[5]  // 1

	// We need to convert the board string into a piece array. Due to the way
	// FEN is structured, the first piece is at a8, which translates to 0x88
	// index 112.
	var startBoard [128]piece
	boardIndex := 112

	// Loop through the board string.
	for _, char := range boardString {
		// Skip the slashes which seperate ranks
		if char == '/' {
			continue
		}

		// Numbers in the board string indicate empty spaces, so we advance the
		// board index by that number of spaces since there aren't any pieces in
		// those positions.
		if char >= '1' && char <= '8' {
			boardIndex += int(char - '0')
		} else {
			// Otherwise, look up the correct piece code and insert it.
			startBoard[boardIndex] = fenCodes[byte(char)]
			boardIndex++
		}

		// Skip squares that aren't on the board
		if boardIndex%16 > 7 {
			boardIndex = ((boardIndex / 16) - 1) * 16
		}
	}

	// Convert the player string to a colour
	var toMove byte
	if playerString == "w" {
		toMove = White
	} else {
		toMove = Black
	}

	// Convert the castling string to a castle byte
	var castling byte
	castling = setCastle(castling, KingCastle, White, strings.Contains(castleString, "K"))
	castling = setCastle(castling, QueenCastle, White, strings.Contains(castleString, "Q"))
	castling = setCastle(castling, KingCastle, Black, strings.Contains(castleString, "k"))
	castling = setCastle(castling, QueenCastle, Black, strings.Contains(castleString, "q"))

	// Convert the en passant target string to the index it represents.
	var enPassantTarget byte

	if enPassantString == "-" {
		enPassantTarget = NoEnPassant
	} else {
		fileLetter := enPassantString[0]
		rankNumber := int(enPassantString[1] - '0')

		enPassantTarget = byte(((rankNumber - 1) * 16) + int(fileLetter-'a'))
	}

	// Convert the half and full move.
	halfmoveInt, _ := strconv.Atoi(halfmoveString)
	halfmove := byte(halfmoveInt)
	fullmove, _ := strconv.Atoi(fullmoveString)

	// Initialise the full position and return it.
	startPosition := position{board: startBoard, toMove: toMove, castling: castling, enPassantTarget: enPassantTarget, halfmove: halfmove, fullmove: fullmove}

	return startPosition
}

// Given an internal position object, convert it to a string in Forsyth–Edwards
// Notation (FEN)
func toFEN(position position) string {
	var pieces string
	var player string
	var castling string
	var enPassant string

	// Loop in reverse through the ranks, from 8 to 1.
	for rank := 7; rank >= 0; rank-- {
		// We need to keep track of the number of empty squares accumulated so
		// far in the rank.
		empty := 0

		// Loop forwards through the files, from a to h.
		for i := rank * 16; i < rank*16+8; i++ {
			// If a piece is present, add the number of empty squares
			// encountered so far (if non-zero) to the output string, then rest
			// the counter and add the current piece to the string. Otherwise,
			// increment the empty square count.
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

		// If no pieces were encountered, add the empty squares count to the
		// string.
		if empty != 0 {
			pieces += strconv.FormatInt(int64(empty), 10)
		}

		// At the end of the rank, add a slash.
		if rank != 0 {
			pieces += "/"
		}
	}

	// Convert the current player to a string.
	if position.toMove == White {
		player = "w"
	} else {
		player = "b"
	}

	castling = castleString(position)

	enPassant = enPassantString(position)

	// Format and return the FEN.
	fen := fmt.Sprintf("%v %v %v %v %v %v", pieces, player, castling, enPassant, position.halfmove, position.fullmove)

	return fen
}

// Converts a piece to a FEN string.
func pieceToString(p piece) string {
	var code string

	switch p & pieceIdentity {
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

// Converts a move to a string representation, used for tests.
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

// Given a string in Forsyth–Edwards Notation (FEN), convert it to the internal
// position object.

// Generate a string representing castling rights from a position.
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

// Generate a string representing the en passant target from a position.
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
