// +build !noasm,!gccgo,!safe

#include "textflag.h"

/*
MakePositive(v []float64) clears negative values:
  for i := 0; i < len(v); i++ {
      if v[i] < 0 {
          v[i] = 0
      }
  }
This version uses SSE2 unrolling to handle 8 elements per loop iteration.
We avoid recursive macro invocation by not redefining X registers.
*/

#define PTRv    DI
#define LENv    CX
#define IDXv    AX
#define TAILv   BX

// func MakePositive(v []float64)
TEXT ·MakePositive(SB), NOSPLIT, $0
    // Load pointer (v_base) and length (v_len)
    MOVQ    v_base+0(FP), PTRv    // &v[0]
    MOVQ    v_len+8(FP),  LENv    // len(v)
    CMPQ    LENv, $0
    JE      done

    // XORPS X0, X0 -> X0 := 0.0...0.0 in double
    XORPS   X0, X0

    // We process 8 doubles at a time.
    XORQ    IDXv, IDXv
    MOVQ    LENv, TAILv
    ANDQ    $7, TAILv        // leftover count
    SHRQ    $3, LENv         // main loop count (len/8)
    JZ      tail_start

loop_main:
    // Load 8 doubles in 4 SSE regs.
    MOVUPS  (PTRv)(IDXv*8), X1
    MOVUPS  16(PTRv)(IDXv*8), X2
    MOVUPS  32(PTRv)(IDXv*8), X3
    MOVUPS  48(PTRv)(IDXv*8), X4

    // maxpd X0, Xn => Xn = max(Xn, X0)
    // This replaces negatives in Xn with zero.
    MAXPD   X0, X1
    MAXPD   X0, X2
    MAXPD   X0, X3
    MAXPD   X0, X4

    // Store the results back
    MOVUPS  X1, (PTRv)(IDXv*8)
    MOVUPS  X2, 16(PTRv)(IDXv*8)
    MOVUPS  X3, 32(PTRv)(IDXv*8)
    MOVUPS  X4, 48(PTRv)(IDXv*8)

    ADDQ    $8, IDXv
    DECQ    LENv
    JNZ     loop_main

    CMPQ    TAILv, $0
    JE      done

tail_start:
    // Handle leftover (1..7) doubles
    MOVQ    TAILv, LENv
    SHRQ    $1, LENv          // process pairs of 2
    JZ      tail_one

tail_two:
    MOVUPS  (PTRv)(IDXv*8), X1
    MAXPD   X0, X1
    MOVUPS  X1, (PTRv)(IDXv*8)
    ADDQ    $2, IDXv
    DECQ    LENv
    JNZ     tail_two

    ANDQ    $1, TAILv
    JZ      done

tail_one:
    // handle single leftover double
    MOVSD   (PTRv)(IDXv*8), X1
    MAXSD   X0, X1
    MOVSD   X1, (PTRv)(IDXv*8)

done:
    RET
