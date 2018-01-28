# ChessSlayer

ChessSlayer is a UCI-compatible chess engine, written in Go. While it doesn't
support the whole UCI standard, it is perfectly capable of playing games when
used with a client like [Scid vs PC](http://scidvspc.sourceforge.net/).

## Why ChessSlayer? Is it actually good?

Not even close. But it works. The aim of this project was to create a chess
engine that could play games correctly and provide some challenge to a novice
player.

The engine uses a negamax search framework with alpha-beta pruning to speed up
the search. While it's very slow, it gives correct results on the [Perft
test](https://chessprogramming.wikispaces.com/Perft) and could therefore
theoretically beat any human player given enough time. In this case, the
required time would be very long indeed.

## Installation

Run `go get github.com/kdelwat/chess` to install the binary, which will be called `chess`.

This binary can then be used as an engine in UCI-compatible programs. My
recommendation is [Scid vs PC](http://scidvspc.sourceforge.net/). Add the
program to the analysis engines list, following the instructions
[here](http://www.watfordchessclub.org/index.php/chess-freeware/54-scid-vs-pc-getting-started),
then choose the `Play -> Computer - UCI Engine` menu item to start a game. The
available analysis modes are depth-based and time-based. I recommend using a
depth of four or five, or a move time of 5 seconds, for the best
performance/time tradeoff.

![A game in progress](https://imgur.com/1ed6fd9e-4e5c-41f2-8017-63b3a0426042)

## Further improvements

While ChessSlayer works as a minimal chess engine, there are a number of
improvements I would like to add in the future:

- Complete implementation of the UCI protocol
- Legal move generator to improve performance
- Better move ordering to improve search speed
- Aspiration windows to increase search efficiency
- A Zobrist transposition table to increase search efficiency

## Credits

I couldn't have made this engine without the invaluable help of the
[Chess programming wiki](https://chessprogramming.wikispaces.com/) and the
[Mediocre chess blog](https://mediocrechess.blogspot.com.au/) by Jonatan Pettersson.