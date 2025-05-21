// +build !noasm,!gccgo,!safe

#include "textflag.h"

// func AddConst(alpha float32, x []float32)
TEXT ·AddConst(SB), NOSPLIT, $0
	MOVQ   x_base+8(FP), SI  // SI = &x
	MOVQ   x_len+16(FP), CX  // CX = len(x)
	CMPQ   CX, $0            // if len(x) == 0 { return }
	JE     ac_end
	MOVSS  alpha+0(FP), X4   // X4 = { a, a, a, a }
	SHUFPS $0, X4, X4
	MOVUPS X4, X5            // X5 = X4
	XORQ   AX, AX            // i = 0
	MOVQ   CX, BX
	ANDQ   $15, BX           // BX = len(x) % 16
	SHRQ   $4, CX            // CX = floor( len(x) / 16 )
	JZ     ac_tail_start     // if CX == 0 { goto ac_tail_start }

ac_loop: // Loop unrolled 16x   do {
	MOVUPS (SI)(AX*4), X0    // X_i = s[i:i+4]
	MOVUPS 16(SI)(AX*4), X1
	MOVUPS 32(SI)(AX*4), X2
	MOVUPS 48(SI)(AX*4), X3
	ADDPS  X4, X0            // X_i += a
	ADDPS  X5, X1
	ADDPS  X4, X2
	ADDPS  X5, X3
	MOVUPS X0, (SI)(AX*4)    // s[i:i+4] = X_i
	MOVUPS X1, 16(SI)(AX*4)
	MOVUPS X2, 32(SI)(AX*4)
	MOVUPS X3, 48(SI)(AX*4)
	ADDQ   $16, AX           // i += 16
	LOOP   ac_loop           // } while --CX > 0
	CMPQ   BX, $0            // if BX == 0 { return }
	JE     ac_end

ac_tail_start: // Reset loop counters
	MOVQ BX, CX // Loop counter: CX = BX

ac_tail: // do {
	MOVSS (SI)(AX*4), X0  // X0 = s[i]
	ADDSS X4, X0          // X0 += a
	MOVSS X0, (SI)(AX*4)  // s[i] = X0
	INCQ  AX              // ++i
	LOOP  ac_tail         // } while --CX > 0

ac_end:
	RET
