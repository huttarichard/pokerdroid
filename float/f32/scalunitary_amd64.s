// +build !noasm,!gccgo,!safe

#include "textflag.h"

#define X_PTR SI
#define IDX AX
#define LEN CX
#define TAIL BX
#define ALPHA X0
#define ALPHA_2 X1

// func ScalUnitary(alpha float32, x []float32)
TEXT ·ScalUnitary(SB), NOSPLIT, $0
	MOVSS   alpha+0(FP), ALPHA    // ALPHA = alpha
	SHUFPS  $0, ALPHA, ALPHA      // ALPHA = { alpha, alpha, alpha, alpha }
	MOVQ    x_base+8(FP), X_PTR   // X_PTR = &x
	MOVQ    x_len+16(FP), LEN     // LEN = len(x)
	CMPQ    LEN, $0
	JE      end                   // if LEN == 0 { return }
	XORQ    IDX, IDX              // IDX = 0

	MOVQ    LEN, TAIL
	ANDQ    $15, TAIL             // TAIL = LEN % 16
	SHRQ    $4, LEN               // LEN = floor( LEN / 16 )
	JZ      tail_start            // if LEN == 0 { goto tail_start }

	MOVUPS  ALPHA, ALPHA_2        // ALPHA_2 = ALPHA for pipelining

loop:  // do {  // x[i] *= alpha unrolled 16x.
	MOVUPS  (X_PTR)(IDX*4), X2     // X_i = x[i]
	MOVUPS  16(X_PTR)(IDX*4), X3
	MOVUPS  32(X_PTR)(IDX*4), X4
	MOVUPS  48(X_PTR)(IDX*4), X5

	MULPS   ALPHA, X2             // X_i *= ALPHA
	MULPS   ALPHA_2, X3
	MULPS   ALPHA, X4
	MULPS   ALPHA_2, X5

	MOVUPS  X2, (X_PTR)(IDX*4)     // x[i] = X_i
	MOVUPS  X3, 16(X_PTR)(IDX*4)
	MOVUPS  X4, 32(X_PTR)(IDX*4)
	MOVUPS  X5, 48(X_PTR)(IDX*4)

	ADDQ    $16, IDX              // i += 16
	DECQ    LEN
	JNZ     loop                  // while --LEN > 0
	CMPQ    TAIL, $0
	JE      end                   // if TAIL == 0 { return }

tail_start: // Reset loop registers
	MOVQ    TAIL, LEN             // Loop counter: LEN = TAIL
	SHRQ    $2, LEN               // LEN = floor( TAIL / 4 )
	JZ      tail_one              // if LEN == 0 { goto tail_one }

tail_four: // do {
	MOVUPS  (X_PTR)(IDX*4), X2    // X_i = x[i]
	MULPS   ALPHA, X2             // X_i *= ALPHA
	MOVUPS  X2, (X_PTR)(IDX*4)    // x[i] = X_i
	ADDQ    $4, IDX               // i += 4
	DECQ    LEN
	JNZ     tail_four             // while --LEN > 0

	ANDQ    $3, TAIL
	JZ      end                   // if TAIL == 0 { return }

tail_one:
	// x[i] *= alpha for the remaining elements
	MOVSS   (X_PTR)(IDX*4), X2
	MULSS   ALPHA, X2
	MOVSS   X2, (X_PTR)(IDX*4)
	INCQ    IDX
	DECQ    TAIL
	JNZ     tail_one              // Process one element at a time until done

end:
	RET
