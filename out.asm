global	main

main:
	push	%rbp
	mov		%rsp, %rbp
	mov		0, %rax
	mov		%rbp, %rsp
	pop		%rbp
	ret
