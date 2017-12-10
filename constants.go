package main

// Piece codes
const King = 0x03
const Queen = 0x07
const Rook = 0x05
const Bishop = 0x04
const Knight = 0x02
const Pawn = 0x01
const Empty = 0x00

const White = 0x00
const Black = 0x40

const OffBoard = 0x88
const Sliding = 0x04
const Piece = 0x0F
const Color = 0x40

const BoardSize = 128

// moves
const Capture = 0x1 << 18
const DoublePawnPush = 0x1 << 16
const Promotion = 0x1 << 19

const KnightPromotion = 0x0
const BishopPromotion = 0x1 << 16
const RookPromotion = 0x1 << 17
const QueenPromotion = 0x2 << 16
