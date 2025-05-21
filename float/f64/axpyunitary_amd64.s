// +build !noasm,!gccgo,!safe

#include "textflag.h"

#define X_PTR SI
#define Y_PTR DI
#define DST_PTR DI
#define IDX AX
#define LEN CX
#define TAIL BX
#define ALPHA X0
#define ALPHA_2 X1

// func AxpyUnitary(alpha float64, x, y []float64)
TEXT ·AxpyUnitary(SB), NOSPLIT, $0
	MOVQ    x_base+8(FP), X_PTR  // X_PTR := &x
	MOVQ    y_base+32(FP), Y_PTR // Y_PTR := &y
	MOVQ    x_len+16(FP), LEN    // LEN = min( len(x), len(y) )
	CMPQ    y_len+40(FP), LEN
	CMOVQLE y_len+40(FP), LEN
	CMPQ    LEN, $0              // if LEN == 0 { return }
	JE      end
	XORQ    IDX, IDX
	MOVSD   alpha+0(FP), ALPHA   // ALPHA := { alpha, alpha }
	SHUFPD  $0, ALPHA, ALPHA
	MOVUPS  ALPHA, ALPHA_2       // ALPHA_2 := ALPHA   for pipelining
	MOVQ    Y_PTR, TAIL          // Check memory alignment
	ANDQ    $15, TAIL            // TAIL = &y % 16
	JZ      no_trim              // if TAIL == 0 { goto no_trim }

	// Align on 16-byte boundary
	MOVSD (X_PTR), X2   // X2 := x[0]
	MULSD ALPHA, X2     // X2 *= a
	ADDSD (Y_PTR), X2   // X2 += y[0]
	MOVSD X2, (DST_PTR) // y[0] = X2
	INCQ  IDX           // i++
	DECQ  LEN           // LEN--
	JZ    end           // if LEN == 0 { return }

no_trim:
	MOVQ LEN, TAIL
	ANDQ $7, TAIL   // TAIL := n % 8
	SHRQ $3, LEN    // LEN = floor( n / 8 )
	JZ   tail_start // if LEN == 0 { goto tail2_start }

loop:  // do {
	// y[i] += alpha * x[i] unrolled 8x.
	MOVUPS (X_PTR)(IDX*8), X2   // X_i = x[i]
	MOVUPS 16(X_PTR)(IDX*8), X3
	MOVUPS 32(X_PTR)(IDX*8), X4
	MOVUPS 48(X_PTR)(IDX*8), X5

	MULPD ALPHA, X2   // X_i *= a
	MULPD ALPHA_2, X3
	MULPD ALPHA, X4
	MULPD ALPHA_2, X5

	ADDPD (Y_PTR)(IDX*8), X2   // X_i += y[i]
	ADDPD 16(Y_PTR)(IDX*8), X3
	ADDPD 32(Y_PTR)(IDX*8), X4
	ADDPD 48(Y_PTR)(IDX*8), X5

	MOVUPS X2, (DST_PTR)(IDX*8)   // y[i] = X_i
	MOVUPS X3, 16(DST_PTR)(IDX*8)
	MOVUPS X4, 32(DST_PTR)(IDX*8)
	MOVUPS X5, 48(DST_PTR)(IDX*8)

	ADDQ $8, IDX  // i += 8
	DECQ LEN
	JNZ  loop     // } while --LEN > 0
	CMPQ TAIL, $0 // if TAIL == 0 { return }
	JE   end

tail_start: // Reset loop registers
	MOVQ TAIL, LEN // Loop counter: LEN = TAIL
	SHRQ $1, LEN   // LEN = floor( TAIL / 2 )
	JZ   tail_one  // if TAIL == 0 { goto tail }

tail_two: // do {
	MOVUPS (X_PTR)(IDX*8), X2   // X2 = x[i]
	MULPD  ALPHA, X2            // X2 *= a
	ADDPD  (Y_PTR)(IDX*8), X2   // X2 += y[i]
	MOVUPS X2, (DST_PTR)(IDX*8) // y[i] = X2
	ADDQ   $2, IDX              // i += 2
	DECQ   LEN
	JNZ    tail_two             // } while --LEN > 0

	ANDQ $1, TAIL
	JZ   end      // if TAIL == 0 { goto end }

tail_one:
	MOVSD (X_PTR)(IDX*8), X2   // X2 = x[i]
	MULSD ALPHA, X2            // X2 *= a
	ADDSD (Y_PTR)(IDX*8), X2   // X2 += y[i]
	MOVSD X2, (DST_PTR)(IDX*8) // y[i] = X2

end:
	RET
