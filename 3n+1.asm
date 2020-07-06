; 3n+1 problem in stack S16 assembly
; given a number N, repeat until you reach 1:
; if N is even, replace with N/2
; if N is odd, replace with 3N + 1
; put all the numbers encountered on the way in data memory
; starting at location 1, with location 0 used as a pointer
; an optional interrupt occurs before every number in the sequence

            PUSH #1     ; initialise data pointer to 1
            STOR #0
            PUSH #7     ; the original N
start:      DUP         ; store in data memory
            LOAD #0
            STRS
            LOAD #0     ; increase pointer in data[0]
            PUSH #1
            ADD
            STOR #0
            ; INT         ; interrupt - enable if needed
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
