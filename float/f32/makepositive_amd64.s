// +build !noasm,!gccgo,!safe

#include "textflag.h"

/*
MakePositive(v []float32) clears negative values:
  for i := 0; i < len(v); i++ {
      if v[i] < 0 {
          v[i] = 0
      }
  }
This version uses SSE2 unrolling to handle 16 elements per loop iteration.
*/

#define PTRv    DI
#define LENv    CX
#define IDXv    AX
#define TAILv   BX

// func MakePositive(v []float32)
TEXT ·MakePositive(SB), NOSPLIT, $0
    // Load pointer (v_base) and length (v_len)
    MOVQ    v_base+0(FP), PTRv    // &v[0]
    MOVQ    v_len+8(FP),  LENv    // len(v)
    CMPQ    LENv, $0
    JE      done

    // XORPS X0, X0 -> X0 := 0.0...0.0 in float32
    XORPS   X0, X0

    // We process 16 floats at a time (4 floats per XMM register * 4 registers)
    XORQ    IDXv, IDXv
    MOVQ    LENv, TAILv
    ANDQ    $15, TAILv        // leftover count
    SHRQ    $4, LENv          // main loop count (len/16)
    JZ      tail_start        // if LENv == 0 { goto tail_start }

loop_main:
    // Load 16 floats in 4 SSE regs
    MOVUPS  (PTRv)(IDXv*4), X1
    MOVUPS  16(PTRv)(IDXv*4), X2
    MOVUPS  32(PTRv)(IDXv*4), X3
    MOVUPS  48(PTRv)(IDXv*4), X4

    // maxps X0, Xn => Xn = max(Xn, X0)
    // This replaces negatives in Xn with zero
    MAXPS   X0, X1
    MAXPS   X0, X2
    MAXPS   X0, X3
    MAXPS   X0, X4

    // Store the results back
    MOVUPS  X1, (PTRv)(IDXv*4)
    MOVUPS  X2, 16(PTRv)(IDXv*4)
    MOVUPS  X3, 32(PTRv)(IDXv*4)
    MOVUPS  X4, 48(PTRv)(IDXv*4)

    ADDQ    $16, IDXv
    DECQ    LENv
    JNZ     loop_main

    CMPQ    TAILv, $0
    JE      done

tail_start:
    // Handle leftover (1..15) floats
    MOVQ    TAILv, LENv
    SHRQ    $2, LENv          // process groups of 4
    JZ      tail_one

tail_four:
    MOVUPS  (PTRv)(IDXv*4), X1
    MAXPS   X0, X1
    MOVUPS  X1, (PTRv)(IDXv*4)
    ADDQ    $4, IDXv
    DECQ    LENv
    JNZ     tail_four

    ANDQ    $3, TAILv
    JZ      done

tail_one:
    // handle single leftover floats one by one
    MOVSS   (PTRv)(IDXv*4), X1
    MAXSS   X0, X1
    MOVSS   X1, (PTRv)(IDXv*4)
    INCQ    IDXv
    DECQ    TAILv
    JNZ     tail_one

done:
    RET
