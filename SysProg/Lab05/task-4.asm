;# Бегущий огонь
org 0x8000

; Настройка сегментов (всё в одном сегменте)
.model tiny

; Огни
mov al, 10000000b

main:
	; Задержка
	mov dl, 0x07 ; 

	loop:
		dec dl
		jnz loop
		
		ror al, 1
	
		out 2, al
		jmp main