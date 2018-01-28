package main

const EngineName = "Ultimate Engine"
const EngineAuthor = "Cadel Watson"

// Codes used to determine the identity of a piece.
const King = 0x03
const Queen = 0x07
const Rook = 0x05
const Bishop = 0x04
const Knight = 0x02
const Pawn = 0x01
const Empty = 0x00

// Player colours
const White = 0x00
const Black = 0x40

// Size of board. This is 128, not 64, because the 0x88 board representation is
// used.
const BoardSize = 128

// Codes used to determine move types.
const Capture = 0x1 << 18
const DoublePawnPush = 0x1 << 16
const Promotion = 0x1 << 19
const EnPassant = 0x1 << 16

const KnightPromotion = 0x0
const BishopPromotion = 0x1 << 16
const RookPromotion = 0x1 << 17
const QueenPromotion = 0x3 << 16

const KingCastle = 0x1 << 17
const QueenCastle = 0x3 << 16

// NoEnPassant is used to represent a lack of an en passant target.
const NoEnPassant = 0
