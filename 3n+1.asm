; 3n+1 problem in stack S16 assembly
; given a number N, repeat until you reach 1:
; if N is even, replace with N/2
; if N is odd, replace with 3N + 1
; an interrupt occurs ebfore every number in the sequence

            PUSH #7     ; the original N
start:      INT
            DUP         ; make a copy for testing if it's 1
            PUSH #1
            CEQ
            JT end
            DUP         ; make a copy for the even test
            PUSH #1
            AND
            JT odd
            PUSH #1
            SHR         ; it's even, divide by 2
            J start        
odd:        PUSH #3
            MUL
            PUSH #1
            ADD
            J start
end:        HALT
