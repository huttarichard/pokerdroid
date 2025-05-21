// +build !noasm,!gccgo,!safe

#include "textflag.h"

// PreventZero(v []float32, tol float32) replaces any 0.0 in v with tol:
//   for i := 0; i < len(v); i++ {
//       if v[i] == 0 {
//           v[i] = tol
//       }
//   }
//
// This implementation uses a simple loop with UCOMISS to check equality with 0.

TEXT ·PreventZero(SB), NOSPLIT, $0
    // Go AMD64 calling convention for:
    //   func PreventZero(v []float32, tol float32)
    // Parameters:
    //   v is a slice at v_base+0(FP), v_len+8(FP), v_cap+16(FP)
    //   tol at tol+24(FP)

    // Load v.ptr into DI, v.len into CX
    MOVQ  v_base+0(FP), DI  // DI = &v[0]
    MOVQ  v_len+8(FP), CX   // CX = len(v)
    
    // If len == 0, no work to do
    CMPQ  CX, $0
    JE    done

    // Load tol into X2
    MOVSS tol+24(FP), X2    // X2 = tol
    
    // X0 = 0.0 for comparison
    XORPS X0, X0            // X0 = 0.0
    
    // We'll use AX as our loop index
    XORQ  AX, AX            // AX = 0

loop:
    // If we've processed all elements, we're done
    CMPQ  CX, $0
    JE    done

    // Load v[i] into X1
    MOVSS (DI), X1          // X1 = v[i]

    // Compare X1 with X0
    UCOMISS X0, X1
    JNE   not_zero          // Jump if unequal or unordered (NaN)

    // Store tol at v[i] if v[i] == 0
    MOVSS X2, (DI)

not_zero:
    // Advance to next element
    ADDQ $4, DI             // move pointer by 4 bytes (1 float32)
    DECQ CX                 // decrement length count
    JMP  loop

done:
    RET
