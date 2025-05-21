// +build !noasm,!gccgo,!safe

#include "textflag.h"

// PreventZero(v []float64, tol float64) replaces any 0.0 in v with tol:
//   for i := 0; i < len(v); i++ {
//       if v[i] == 0 {
//           v[i] = tol
//       }
//   }
//
// This implementation uses a simple loop with UCOMISD to check equality with 0.

TEXT ·PreventZero(SB), NOSPLIT, $0
    // Go AMD64 calling convention for:
    //   func PreventZero(v []float64, tol float64)
    // On stack (starting at 0(SP)):
    //   0(SP)   : return address
    //   8(SP)   : v.ptr
    //   16(SP)  : v.len
    //   24(SP)  : v.cap
    //   32(SP)  : tol (float64)

    // Load v.ptr into DI, v.len into CX
    MOVQ  8(SP), DI
    MOVQ 16(SP), CX
    // We do not need v.cap, so we skip the 24(SP) slot.

    // Correctly load tol from 32(SP) into X2
    MOVSD 32(SP), X2

    // If len == 0, no work to do.
    CMPQ  CX, $0
    JE    done

    // X0 = 0.0 for comparison
    XORPS X0, X0

loop:
    // If we've processed all elements, we're done.
    CMPQ  CX, $0
    JE    done

    // Load v[i] into X1
    MOVSD (DI), X1

    // Compare X1 with X0. If equal => store tol.
    UCOMISD X0, X1
    JNE   not_zero

    // Store tol at v[i]
    MOVSD X2, (DI)

not_zero:
    // Advance to next element
    ADDQ $8, DI   // move pointer by 8 bytes (1 float64)
    DECQ CX       // decrement length count
    JMP  loop

done:
    RET
