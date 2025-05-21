// +build !noasm,!gccgo,!safe

#include "textflag.h"

// ScalUnitaryToUP implements:
//
//   for i, v := range x {
//       if v > 0 {
//           dst[i] = alphaUp * v
//       } else if v < 0 {
//           dst[i] = alphaDown * v
//       } else {
//           dst[i] = 0
//       }
//   }
//
// It processes two float32 values at a time with SSE instructions.

TEXT ·ScalUnitaryToUP(SB), NOSPLIT, $0
    // 1) Load slice dst:
    //    param dst @ [0..23], pointer => 0(FP), length => 8(FP)
    MOVQ  dst_base+0(FP), DI      // DI = &dst[0]
    MOVQ  dst_len+8(FP), R8       // R8 = len(dst)

    // 2) Load alphaUp, alphaDown (float32):
    //    param #2 => alphaUp at 24(FP)
    //    param #3 => alphaDown at 28(FP)
    MOVSS alphaUp+24(FP), X0      // X0 = alphaUp
    MOVSS alphaDown+28(FP), X1    // X1 = alphaDown

    // 3) Load slice x:
    //    param x @ [32..(32+24)=56], pointer => 32(FP), length => 40(FP)
    MOVQ  x_base+32(FP), SI       // SI = &x[0]
    MOVQ  x_len+40(FP), CX        // CX = len(x)

    // 4) We want to process up to min( len(dst), len(x) ) 
    //    to avoid out-of-bounds. Compare CX, R8:
    CMPQ  CX, R8
    CMOVQGT R8, CX                // CX = min(CX, R8)

    // If there's nothing to process, return.
    TESTQ CX, CX
    JE end

    // We'll process two floats at a time => CX >> 1
    // Keep original length in BX for leftover check:
    MOVQ CX, BX
    SHRQ $1, CX
    JZ  remain

    // X7 will hold zero for sign checking:
    XORPS X7, X7   // single-precision zero vector

loop:
    // --- Load 2 x-values (8 bytes total) in SSE register X2 ---
    //   We must assume unaligned, so MOVUPS or MOVLPS:
    MOVUPS (SI), X2

    // Copy them:
    MOVAPS X2, X3
    MOVAPS X2, X4

    // 1) first float => X3
    UCOMISS X7, X3    // compare X3 to 0
    JA  pos1          // X3 > 0
    JB  neg1          // X3 < 0
    // else => zero
    XORPS X3, X3
    JMP next1

pos1:
    MULSS X0, X3      // X3 *= alphaUp
    JMP next1
neg1:
    MULSS X1, X3      // X3 *= alphaDown
next1:

    // 2) second float => shift X4 by 4
    PSRLDQ $4, X4     // move top float down
    UCOMISS X7, X4
    JA  pos2
    JB  neg2
    // else => zero
    XORPS X4, X4
    JMP next2

pos2:
    MULSS X0, X4
    JMP next2
neg2:
    MULSS X1, X4
next2:

    // Store the two results to dst (DI):
    MOVSS X3, (DI)
    MOVSS X4, 4(DI)

    // Advance:
    ADDQ $8, SI
    ADDQ $8, DI

    DECQ CX
    JNZ loop

remain:
    // If original length was odd, do 1 leftover:
    ANDQ $1, BX
    JZ end  // no leftover => done

    // leftover:
    MOVSS (SI), X2  // load x
    UCOMISS X7, X2
    JA pos_last
    JB neg_last
    // zero:
    XORPS X2, X2
    JMP store_last

pos_last:
    MULSS X0, X2
    JMP store_last
neg_last:
    MULSS X1, X2
store_last:
    MOVSS X2, (DI)

end:
    RET
