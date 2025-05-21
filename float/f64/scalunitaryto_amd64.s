// +build !noasm,!gccgo,!safe

#include "textflag.h"

/*
ScalUnitaryTo(dst, alpha, x) implements:

  for i, v := range x {
	dst[i] = alpha * v
  }
*/


#define MOVDDUP_ALPHA    LONG $0x44120FF2; WORD $0x2024 // @ MOVDDUP 32(SP), X0  /*XMM0, 32[RSP]*/

#define X_PTR SI
#define DST_PTR DI
#define IDX AX
#define LEN CX
#define TAIL BX
#define ALPHA X0
#define ALPHA_2 X1

// func ScalUnitaryTo(dst []float64, alpha float64, x []float64)
// This function assumes len(dst) >= len(x).
TEXT ·ScalUnitaryTo(SB), NOSPLIT, $0
	MOVQ x_base+32(FP), X_PTR    // X_PTR = &x
	MOVQ dst_base+0(FP), DST_PTR // DST_PTR = &dst
	MOVDDUP_ALPHA                // ALPHA = { alpha, alpha }
	MOVQ x_len+40(FP), LEN       // LEN = len(x)
	CMPQ LEN, $0
	JE   end                     // if LEN == 0 { return }

	XORQ IDX, IDX   // IDX = 0
	MOVQ LEN, TAIL
	ANDQ $7, TAIL   // TAIL = LEN % 8
	SHRQ $3, LEN    // LEN = floor( LEN / 8 )
	JZ   tail_start // if LEN == 0 { goto tail_start }

	MOVUPS ALPHA, ALPHA_2 // ALPHA_2 = ALPHA for pipelining

loop:  // do { // dst[i] = alpha * x[i] unrolled 8x.
	MOVUPS (X_PTR)(IDX*8), X2   // X_i = x[i]
	MOVUPS 16(X_PTR)(IDX*8), X3
	MOVUPS 32(X_PTR)(IDX*8), X4
	MOVUPS 48(X_PTR)(IDX*8), X5

	MULPD ALPHA, X2   // X_i *= ALPHA
	MULPD ALPHA_2, X3
	MULPD ALPHA, X4
	MULPD ALPHA_2, X5

	MOVUPS X2, (DST_PTR)(IDX*8)   // dst[i] = X_i
	MOVUPS X3, 16(DST_PTR)(IDX*8)
	MOVUPS X4, 32(DST_PTR)(IDX*8)
	MOVUPS X5, 48(DST_PTR)(IDX*8)

	ADDQ $8, IDX  // i += 8
	DECQ LEN
	JNZ  loop     // while --LEN > 0
	CMPQ TAIL, $0
	JE   end      // if TAIL == 0 { return }

tail_start: // Reset loop counters
	MOVQ TAIL, LEN // Loop counter: LEN = TAIL
	SHRQ $1, LEN   // LEN = floor( TAIL / 2 )
	JZ   tail_one  // if LEN == 0 { goto tail_one }

tail_two: // do {
	MOVUPS (X_PTR)(IDX*8), X2   // X_i = x[i]
	MULPD  ALPHA, X2            // X_i *= ALPHA
	MOVUPS X2, (DST_PTR)(IDX*8) // dst[i] = X_i
	ADDQ   $2, IDX              // i += 2
	DECQ   LEN
	JNZ    tail_two             // while --LEN > 0

	ANDQ $1, TAIL
	JZ   end      // if TAIL == 0 { return }

tail_one:
	MOVSD (X_PTR)(IDX*8), X2   // X_i = x[i]
	MULSD ALPHA, X2            // X_i *= ALPHA
	MOVSD X2, (DST_PTR)(IDX*8) // dst[i] = X_i

end:
	RET
