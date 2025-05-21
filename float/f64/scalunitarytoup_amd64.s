// +build !noasm,!gccgo,!safe

#include "textflag.h"

/*
ScalUnitaryToUP(dst []float64, alphaUp, alphaDown float64, x []float64) implements:

	for i, v := range x {
		if v > 0 {
			dst[i] = alphaUp * v
		} else if v < 0 {
			dst[i] = alphaDown * v
		} else {
			dst[i] = 0
		}
	}
*/


TEXT ·ScalUnitaryToUP(SB), NOSPLIT, $0
    MOVQ    dst_base+0(FP), DI
    MOVSD   alphaUp+24(FP), X0
    MOVSD   alphaDown+32(FP), X1
    MOVQ    x_base+40(FP), SI
    MOVQ    x_len+48(FP), CX
    
    // Early return if length is 0
    TESTQ   CX, CX
    JE      end
    
    // Setup X0 (alphaUp) and X1 (alphaDown) as vectors
    UNPCKLPD X0, X0
    UNPCKLPD X1, X1
    
    // Zero register
    XORPD   X7, X7
    
    // Process 2 doubles at a time
    SHRQ    $1, CX
    JZ      remain
    
loop:
    // Load 2 doubles
    MOVUPD  (SI), X2
    
    // Copy values for comparison
    MOVAPD  X2, X3
    MOVAPD  X2, X4
    
    // Compare with zero
    UCOMISD X7, X3
    JA      pos1
    JB      neg1
    // X3 is zero, keep it zero
    XORPD   X3, X3
    JMP     next1
pos1:
    MULSD   X0, X3
    JMP     next1
neg1:
    MULSD   X1, X3
next1:
    
    // Second value
    PSRLDQ  $8, X4
    UCOMISD X7, X4
    JA      pos2
    JB      neg2
    // X4 is zero, keep it zero
    XORPD   X4, X4
    JMP     next2
pos2:
    MULSD   X0, X4
    JMP     next2
neg2:
    MULSD   X1, X4
next2:
    
    // Combine results
    UNPCKLPD X4, X3
    
    // Store result
    MOVUPD  X3, (DI)
    
    ADDQ    $16, SI
    ADDQ    $16, DI
    DECQ    CX
    JNZ     loop

remain:
    // Check if we have one more element
    MOVQ    x_len+48(FP), CX
    ANDQ    $1, CX
    JZ      end
    
    // Process last element
    MOVSD   (SI), X2
    UCOMISD X7, X2
    JA      pos_last
    JB      neg_last
    // X2 is zero, store zero
    XORPD   X2, X2
    JMP     store_last
pos_last:
    MULSD   X0, X2
    JMP     store_last
neg_last:
    MULSD   X1, X2
store_last:
    MOVSD   X2, (DI)

end:
    RET
