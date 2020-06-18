# Stack machine example

This is an example 16-bit stack machine for teaching computer architecture.

The machine has a 16-bit integer data type, 4K each of code and data memory,
a value stack (capacity 256), a return stack (capacity 256) and the following
instructions. The "need" column shows how many values must at least be on the
stack when executing this instruction; if there are fewer, the machine halts
with an error.

Instructions are 16 bits long; instructions that take an argument are a total
of 32 bits long.

name    args    need    diff    desc
DROP    0       1       -1      Removes the top stack element.
PUSH    1       0       +1      Pushes a value onto the stack.
DUP     0       1       +1      Duplicates the stack top.
SWAP    0       2       =0      Swaps the top two stack elements.

ADD     0       2       -1      Adds two elements and pushes the result.
ADDC    0       2       =0      Add and stack the result and a carry word.

SUB     0       2       -2      Subtracts two elements and pushes the result.
MUL     0       2       -1      Multiplies two elements and pushes the low 16 bits of the result.
MULD    0       2       -0      Multiplies and pushes both the low and high words of the result.
MOD     0       2       -1      Computes the modulo of the top two words.
DIV     0       2       -1      Computes integer division of the top two words.
AND     0       2       -1      Computes the logical AND of the top two words.
OR      0       2       -1      Computes the logical OR of the top two words.
XOR     0       2       -1      Computes the logical XOR of the top two words.
NAND    0       2       -1      Computes the logical NAND of the top two words.
NOT     0       1       =0      Logically negates the top word.

ISZE    0       1       =0      Pops a word and pushes 1 if it was zero, else 0.
ISNZ    0       1       =0      Pops a word and pushes 0 if it was zero, else 1.
ISHI    0       1       =0      Pops a word and pushes 1 if the high bit was set, else 0.
ISNH    0       1       =0      Pops a word and pushes 0 if the high bit was set, else 1.
ISGT    0       2       -1      Pushes 1 if x > y, where y is the stack top, else 0.
ISGE    0       2       -1      Pushes 1 if x >= y, where y is the stack top, else 0.
ISLT    0       2       -1      Pushes 1 if x < y, where y is the stack top, else 0.
ISLE    0       2       -1      Pushes 1 if x <= y, where y is the stack top, else 0.
ISEQ    0       2       -1      Pushes 1 if x == y, where y is the stack top, else 0.
ISNE    0       2       -1      Pushes 1 if x != y, where y is the stack top, else 0.

STORE   1       1       -1      Pops a word and stores it at the memory location in the argument.
STORI   0       2       -2      Stores a word at the location on the stack top.
LOAD    1       0       +1      Loads the word at the argument location onto the stack.
LOADI   0       1       =0      Loads the word at the location in the stack top.

JUMP    1       0       =0      Absolute jump to location in argument.
JUMPC   1       1       -1      Conditional jump if top word is nonzero.
JUMPI   0       1       -1      Indirect jump to location on stack.
JUMPCI  0       2       -2      Conditional indirect jump.

CALL    1       0       =0      Function call, push return value on return stack then jump.
RET     0       0       =0      Return from function call - error if return stack empty.


ADDC pops x and y off the stack and pushes z=x+y and w, which is 1 if there was
a carry out and 0 otherwise.