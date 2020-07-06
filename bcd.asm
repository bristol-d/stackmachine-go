; BCD addition of 8-word numbers


start:

; get the two digits, add them, then add the carry
  PUSH x
  CALL load_digit
  PUSH y
  CALL load_digit
  ADD
  LOAD carry
  ADD
  PUSH #0
  STOR carry

; check if we need to carry out
  DUP
  PUSH #0A
  CGE
  JT carry_one

; when we're here, carry has been handled
; write the digit then loop back
write:
  PUSH z
  CALL store_digit
  J increment

carry_one: ; record the carry one, then subtract 10
  PUSH #1
  STOR carry
  PUSH #0A
  SUB
  J write

; increment pos; if it's 8 we're done
increment:
  LOAD pos
  PUSH #1
  ADD
  DUP
  STOR pos
  PUSH #8
  CGE
  JT end
  ; INT
  J start

end:
  HALT ; end of main

; load_digit: loads mem[s + len - 1 - pos] onto the stack, where
; s is the value on the stack top

load_digit:
    LOAD len
    ADD
    PUSH #1
    SUB
    LOAD pos
    SUB
    LODS
    RET

; stack: top = address of destination, below = digit
store_digit:
    LOAD len
    ADD
    PUSH #1
    SUB
    LOAD pos
    SUB
    STRS
    RET

.DATA

; the numbers to add, we want z = x + y

x: 8 = 0 0 0 0 9 9 0 1
y: 8 = 0 0 0 0 0 1 0 5
z: 8

; helpers

pos: 1 = 0
carry: 1 = 0
len: 1 = 8