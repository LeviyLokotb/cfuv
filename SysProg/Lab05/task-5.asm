;# Вложенная задержка
org 0x8000

; Настройка сегментов (всё в одном сегменте)
.model tiny

_outer_delay equ 0x0F
_inner_delay equ 0x0F

mov al, 0x01
out 2, al

; 18N_i - 12             ~ тактов за внутренний цикл
; 18N_i*N_o + 10N_o - 8  ~ тактов за внешний цикл 
mov dh, _outer_delay
outer:
    mov dl, _inner_delay     ; 4 такта
    inner:
        dec dl               ; 2 такта
        jnz inner            ; 16 тактов при переходе, иначе 4
    dec dh                   ; 2 такта
    jnz outer                ; 16 тактов при переходе, иначе 4

mov al, 0xFF
out 2, al