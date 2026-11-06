;# Бегущие огни
org 0x8000

; Настройка сегментов (всё в одном сегменте)
.model tiny

main:
	mov al, 10010010b
	out 2, al
	call delay

	mov al, 01001001b
	out 2, al
	call delay

	mov al, 00100100b
	out 2, al
	call delay
	
	jmp main

delay:
	mov cx, 0x000F
	loop_delay:
		dec cx
		jnz loop_delay
	ret