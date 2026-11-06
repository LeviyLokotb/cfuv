;# Демонстрация действия логических команд
org 0x8000

; Настройка сегментов (всё в одном сегменте)
.model tiny

main:
	in al, 1

	mov bl, 0xFA ; 1111 1010
	and al, bl
	; or al, bl
	; xor al, bl

	out 2, al
	jmp main