;# Проверка конкретного выключателя
org 0x8000

; Настройка сегментов (всё в одном сегменте)
.model tiny

main:
	in al, 1

	mov bl, 00000010b
	and al,bl
	jz continue
	mov al, 0xFF
	
	continue: out 2, al
	jmp main