package omp

/*
#cgo CXXFLAGS: -std=c++11 -fPIC -pthread -Wno-implicit-const-int-float-conversion

#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include "omp.h"

*/
import "C"
import (
	"fmt"

	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/equity"
)

func Equity(a, b card.Card, board card.Cards, players int) equity.Equity {
	cardsS := fmt.Sprintf("%s%s", a.String(), b.String())

	boardS := ""
	for _, b := range board {
		boardS += b.String()
	}

	var win, draw float64
	calcEquityWithDraw(cardsS, boardS, "", players, &win, &draw, 5, 0.0001)

	return equity.NewEquity(float32(win), float32(draw))
}

func calcEquityWithDraw(cards string, board string, dead string, players int,
	win *float64, draw *float64, boardCards int, stdDev float64) {
	C.hand_equity_with_draw(C.CString(cards), C.CString(board), C.CString(dead),
		C.int(players), (*C.double)(win), (*C.double)(draw),
		C.int(boardCards), C.double(stdDev))
}

// GOBUILD_AMD64_LINUX   = CGO_ENABLED=1 GOARCH=amd64 GOOS=linux   CC="zig cc -target x86_64-linux-gnu"   CXX="zig c++ -target x86_64-linux-gnu"   $(GOCMD) build
// GOBUILD_ARM64_LINUX   = CGO_ENABLED=1 GOARCH=arm64 GOOS=linux   CC="zig cc -target aarch64-linux-gnu"  CXX="zig c++ -target aarch64-linux-gnu"  $(GOCMD) build
// GOBUILD_AMD64_WINDOWS = CGO_ENABLED=1 GOARCH=amd64 GOOS=windows CC="zig cc -target x86_64-windows-gnu" CXX="zig c++ -target x86_64-windows-gnu" $(GOCMD) build
// GOBUILD_ARM64_DARWIN  = CGO_ENABLED=1 GOARCH=arm64 GOOS=darwin  $(GOCMD) build
