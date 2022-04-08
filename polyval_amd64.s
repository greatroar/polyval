// Code generated by command: go run asm.go -out out/polyval_amd64.s -stubs out/stub_amd64.go -pkg polyval. DO NOT EDIT.

//go:build gc && !purego

#include "textflag.h"

DATA polymask<>+0(SB)/8, $0xc200000000000000
DATA polymask<>+8(SB)/8, $0xc200000000000000
GLOBL polymask<>(SB), RODATA|NOPTR, $16

// func polymulAsm(acc *fieldElement, key *fieldElement)
// Requires: PCLMULQDQ, SSE, SSE2
TEXT ·polymulAsm(SB), NOSPLIT, $0-16
	MOVQ  acc+0(FP), AX
	MOVQ  key+8(FP), CX
	MOVOU (AX), X0
	MOVOU (CX), X1
	MOVOU polymask<>+0(SB), X2

	// Karatsuba 1
	PSHUFD    $0xee, X0, X3
	PXOR      X0, X3
	PSHUFD    $0xee, X1, X4
	PXOR      X1, X4
	PCLMULQDQ $0x00, X3, X4
	MOVOU     X0, X3
	PCLMULQDQ $0x11, X1, X3
	PCLMULQDQ $0x00, X1, X0

	// Karatsuba 2
	MOVOU      X0, X1
	SHUFPS     $0x4e, X3, X1
	MOVOU      X3, X5
	PXOR       X0, X5
	PXOR       X1, X5
	PXOR       X4, X5
	MOVHLPS    X5, X3
	PUNPCKLQDQ X5, X0

	// Montgomery reduce
	MOVOU     X2, X1
	PCLMULQDQ $0x00, X0, X1
	PSHUFD    $0x4e, X1, X1
	PXOR      X0, X1
	XORPS     X1, X3
	PCLMULQDQ $0x11, X2, X1
	PXOR      X3, X1
	MOVOU     X1, (AX)
	RET

// func polymulBlocksAsm(acc *fieldElement, pow *[8]fieldElement, input *byte, nblocks int)
// Requires: PCLMULQDQ, SSE, SSE2
TEXT ·polymulBlocksAsm(SB), NOSPLIT, $0-32
	MOVQ  acc+0(FP), AX
	MOVQ  pow+8(FP), CX
	MOVQ  input+16(FP), DX
	MOVQ  nblocks+24(FP), BX
	MOVOU polymask<>+0(SB), X0
	MOVOU (AX), X1
	MOVQ  BX, SI
	ANDQ  $0x07, SI
	JZ    initWideLoop
	MOVOU 112(CX), X2

singleLoop:
	MOVOU (DX), X3
	PXOR  X1, X3

	// Karatsuba 1
	PSHUFD    $0xee, X3, X4
	PXOR      X3, X4
	PSHUFD    $0xee, X2, X1
	PXOR      X2, X1
	PCLMULQDQ $0x00, X4, X1
	MOVOU     X3, X4
	PCLMULQDQ $0x11, X2, X4
	PCLMULQDQ $0x00, X2, X3

	// Karatsuba 2
	MOVOU      X3, X5
	SHUFPS     $0x4e, X4, X5
	MOVOU      X4, X6
	PXOR       X3, X6
	PXOR       X5, X6
	PXOR       X1, X6
	MOVHLPS    X6, X4
	PUNPCKLQDQ X6, X3

	// Montgomery reduce
	MOVOU     X0, X1
	PCLMULQDQ $0x00, X3, X1
	PSHUFD    $0x4e, X1, X1
	PXOR      X3, X1
	XORPS     X1, X4
	PCLMULQDQ $0x11, X0, X1
	PXOR      X4, X1
	ADDQ      $0x10, DX
	SUBQ      $0x01, SI
	JNZ       singleLoop

initWideLoop:
	SHRQ $0x03, BX
	JZ   done

wideLoop:
	PXOR X2, X2
	PXOR X4, X4
	PXOR X3, X3

	// Block 7
	MOVOU 112(DX), X5
	MOVOU 112(CX), X6

	// Karatsuba 1
	PSHUFD    $0xee, X5, X7
	PXOR      X5, X7
	PSHUFD    $0xee, X6, X8
	PXOR      X6, X8
	PCLMULQDQ $0x00, X7, X8
	MOVOU     X5, X7
	PCLMULQDQ $0x11, X6, X7
	PCLMULQDQ $0x00, X6, X5
	PXOR      X7, X2
	PXOR      X5, X4
	PXOR      X8, X3

	// Block 6
	MOVOU 96(DX), X5
	MOVOU 96(CX), X6

	// Karatsuba 1
	PSHUFD    $0xee, X5, X7
	PXOR      X5, X7
	PSHUFD    $0xee, X6, X8
	PXOR      X6, X8
	PCLMULQDQ $0x00, X7, X8
	MOVOU     X5, X7
	PCLMULQDQ $0x11, X6, X7
	PCLMULQDQ $0x00, X6, X5
	PXOR      X7, X2
	PXOR      X5, X4
	PXOR      X8, X3

	// Block 5
	MOVOU 80(DX), X5
	MOVOU 80(CX), X6

	// Karatsuba 1
	PSHUFD    $0xee, X5, X7
	PXOR      X5, X7
	PSHUFD    $0xee, X6, X8
	PXOR      X6, X8
	PCLMULQDQ $0x00, X7, X8
	MOVOU     X5, X7
	PCLMULQDQ $0x11, X6, X7
	PCLMULQDQ $0x00, X6, X5
	PXOR      X7, X2
	PXOR      X5, X4
	PXOR      X8, X3

	// Block 4
	MOVOU 64(DX), X5
	MOVOU 64(CX), X6

	// Karatsuba 1
	PSHUFD    $0xee, X5, X7
	PXOR      X5, X7
	PSHUFD    $0xee, X6, X8
	PXOR      X6, X8
	PCLMULQDQ $0x00, X7, X8
	MOVOU     X5, X7
	PCLMULQDQ $0x11, X6, X7
	PCLMULQDQ $0x00, X6, X5
	PXOR      X7, X2
	PXOR      X5, X4
	PXOR      X8, X3

	// Block 3
	MOVOU 48(DX), X5
	MOVOU 48(CX), X6

	// Karatsuba 1
	PSHUFD    $0xee, X5, X7
	PXOR      X5, X7
	PSHUFD    $0xee, X6, X8
	PXOR      X6, X8
	PCLMULQDQ $0x00, X7, X8
	MOVOU     X5, X7
	PCLMULQDQ $0x11, X6, X7
	PCLMULQDQ $0x00, X6, X5
	PXOR      X7, X2
	PXOR      X5, X4
	PXOR      X8, X3

	// Block 2
	MOVOU 32(DX), X5
	MOVOU 32(CX), X6

	// Karatsuba 1
	PSHUFD    $0xee, X5, X7
	PXOR      X5, X7
	PSHUFD    $0xee, X6, X8
	PXOR      X6, X8
	PCLMULQDQ $0x00, X7, X8
	MOVOU     X5, X7
	PCLMULQDQ $0x11, X6, X7
	PCLMULQDQ $0x00, X6, X5
	PXOR      X7, X2
	PXOR      X5, X4
	PXOR      X8, X3

	// Block 1
	MOVOU 16(DX), X5
	MOVOU 16(CX), X6

	// Karatsuba 1
	PSHUFD    $0xee, X5, X7
	PXOR      X5, X7
	PSHUFD    $0xee, X6, X8
	PXOR      X6, X8
	PCLMULQDQ $0x00, X7, X8
	MOVOU     X5, X7
	PCLMULQDQ $0x11, X6, X7
	PCLMULQDQ $0x00, X6, X5
	PXOR      X7, X2
	PXOR      X5, X4
	PXOR      X8, X3

	// Block 0
	MOVOU (DX), X5
	MOVOU (CX), X6
	PXOR  X1, X5

	// Karatsuba 1
	PSHUFD    $0xee, X5, X1
	PXOR      X5, X1
	PSHUFD    $0xee, X6, X7
	PXOR      X6, X7
	PCLMULQDQ $0x00, X1, X7
	MOVOU     X5, X1
	PCLMULQDQ $0x11, X6, X1
	PCLMULQDQ $0x00, X6, X5
	PXOR      X1, X2
	PXOR      X5, X4
	PXOR      X7, X3

	// Karatsuba 2
	MOVOU      X4, X1
	SHUFPS     $0x4e, X2, X1
	MOVOU      X2, X5
	PXOR       X4, X5
	PXOR       X1, X5
	PXOR       X3, X5
	MOVHLPS    X5, X2
	PUNPCKLQDQ X5, X4

	// Montgomery reduce
	MOVOU     X0, X1
	PCLMULQDQ $0x00, X4, X1
	PSHUFD    $0x4e, X1, X1
	PXOR      X4, X1
	XORPS     X1, X2
	PCLMULQDQ $0x11, X0, X1
	PXOR      X2, X1
	ADDQ      $0x80, DX
	SUBQ      $0x01, BX
	JNZ       wideLoop

done:
	MOVOU X1, (AX)
	RET
