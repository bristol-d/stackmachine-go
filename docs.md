# The S16 Architecture

## Overview 

The S16 processor is a fictional 16-bit Harvard stack machine for teaching purposes, which can be accessed through a browser-based simulator. It operates on 16-bit words (assumed to be unsigned for arithmetic purposes) and has the following components:

  * A 16-bit program counter (only the low 12 bits are used).
  * An operation stack that holds up to 256 words.
  * A return stack that holds up to 256 words.
  * An ALU that can perform many of the standard arithmetic, logic and comparison operations that you would find in a higher-level programming language.
  * 4K words of code memory (this means 4096).
  * 4K words of data memory.

Each machine instruction is one machine word long, except for ones that take an additional 16-bit argument such as PUSH.

The execution cycle is as follows: 

  1. Fetch the instruction in code memory at the location indicated by the program counter.
  2. Increment the program counter.
  3. Decode the instruction. If it requires an operand, fetch that from code memory and in-crement the program counter again.
  4. Execute the instruction. This may change the program counter.

The all-zero word is the HALT instruction, which decrements the program counter again during execution so the processor will not proceed beyond an all-zero word, even in single-step mode in the simulator.

## Instruction Set

In the following tables, 

  * Op is the opcode in hexadecimal.
  * Mnemonic is the mnemonic to use in the assembler.
  * Arg indicates if this instruction takes an argument (operand).
  * Pre indicates the minimum number of words that must be on the stack before executing this instruction. If the stack underflows, the processor halts with a fault.
  * Post indicates the change in stack size after executing this instruction. If it is negative, the words to remove will always be covered by the precondition, however if it is positive then this implies an additional condition that the stack must not overflow – doing so causes a processor fault.
  * Description is a brief description of the instruction.

In assembly, anything from a semicolon ( ; ) to the end of a line is a comment. Labels are indicated by a trailing colon ( : ). Mnemonics must be capitalised.

Where an operation takes an operand, it may either be a constant in hexadecimal prefixed by a # sign, or a label. Thus `PUSH #2` pushes the constant 2, whereas `PUSH start` pushes the value of the label 'start'. The assembler is a standard two-pass assembler.

### Halt Operation

| Op   | Mnemonic | Arg | Pre | Post | Description |
|------|----------|-----|-----|------|-------------|
| 0000 | HALT     |     |     | 0    | Sets the program counter back to the current instruction, then halts the processor. |

### Stack Operations

| Op   | Mnemonic | Arg | Pre | Post | Description |
|------|----------|-----|-----|------|-------------|
| 0001 | PUSH     | Yes |     | +1   | Pushes a constant onto the stack. |
| 0002 | POP      |     | 1   | –1   | Removes the top element from the stack. |
| 0003 | DUP      |     | 1   | +1   | Duplicates the top stack element. |
| 0004 | SWAP     |     | 2   | 0    | Swaps the top two words on the stack. |

### Arithmetic Operations

| Op   | Mnemonic | Arg | Pre | Post | Description |
|------|----------|-----|-----|------|-------------|
| 0101 | ADD      |     | 2   | –1   | Pops two words off the stack, adds them and pushes the result back on the stack. |
| 0102 | ADDC     |     | 3   | –1   | Pops three words from the stack, adds the two that were not the top word and then adds 1 if the top word was nonzero (the top word represents the carry in). Then pushes first the addition result, then an additional word that was 1 if there was a carry out, else 0. |
| 0103 | SUB      |     | 2   | –1   | Pops two words off the stack, subtracts them and pushes the result back on the stack. The operation is X–Y where Y was the word on top of the stack. |
| 0104 | MUL      |     | 2   | –1   | Same as above, but multiplies the words. |
| 0105 | MULC     |     | 2   | 0    | Pops two words, multiplies them, then pushes the result as a pair of words, with the most significant 16 bits on the stack top. |
| 0106 | DIV      |     | 2   | –1   | Pops two words off the stack, performs integer division and pushes the result back on the stack. The divisor is the word on top of the stack. |
| 0107 | MOD      |     | 2   | –1   | Same as above, but computes the modulus. |

ADDC stands for 'add with carry'.

### Binary Operations

| Op   | Mnemonic | Arg | Pre | Post | Description |
|------|----------|-----|-----|------|-------------|
| 0201 | AND      |     | 2   | –1   | Removes the top two words from the stack, performs the binary AND operation on them and places the result back on the stack. |
| 0202 | OR       |     | 2   | –1   | Like above, but performs the binary OR. |
| 0203 | XOR      |     | 2   | –1   | Like above, but performs the binary XOR. |
| 0204 | NAND     |     | 2   | –1   | Like above, but performs the binary NAND. |
| 0205 | NOT      |     | 1   | 0    | Removes the top word from the stack, performs binary NOT and puts the result back on the stack. |
| 0206 | SHR      |     | 2   | –1   | Shifts the second-from-top element right by N bits, where N is the stack top element. Shifted in bits are all zeroes. |
| 0207 | SSR      |     | 2   | –1   | Signed shift right: same as above but shifts in the bit that used to be the most significant bit. |
| 0208 | SHL      |     | 2   | –1   | Same as SHR, but shifts left. |
| 0209 | SWE      |     | 1   | 0    | Swap endianness: swaps the low and high bytes of the stack top word. |

SHL/R stands for shift left/right; SSR is 'shift signed right'.

### Comparison and Test Operations

| Op   | Mnemonic | Arg | Pre | Post | Description |
|------|----------|-----|-----|------|-------------|
| 0301 | CEQ      |     | 2   | –1   | Remove two words from the stack and push 1 if they are equal, otherwise 0. |
| 0302 | CNE      |     | 2   | –1   | Remove two words and push 1 if X != Y, else 0. Here Y is the word on the stack top and X the word below. |
| 0303 | CGT      |     | 2   | –1   | Same as above for X > Y. |
| 0304 | CGE      |     | 2   | –1   | Same as above for X >= Y. |
| 0305 | CLT      |     | 2   | –1   | Same as above for X < Y. |
| 0306 | CLE      |     | 2   | –1   | Same as above for X <= Y. |
| 0310 | TZ       |     | 1   | 0    | Pop the top word on the stack. Push 1 if the word was zero, else 0. |
| 0311 | TN       |     | 1   | 0    | Pop the top word on the stack. Push 1 if the word was nonzero, else 0. |
| 0312 | TM       |     | 1   | 0    | Pop the top word on the stack. Push 1 if the most significant bit of the word was set, else 0. |
| 0313 | TL       |     | 1   | 0    | Pop the top word on the stack. Push 1 if the least significant bit of the word was set, else 0. |

### Jump (Branch) Instructions

| Op   | Mnemonic | Arg | Pre | Post | Description |
|------|----------|-----|-----|------|-------------|
| 0401 | J        | yes | 0   | 0    | Jump (set PC) to the location indicated by the operand. |
| 0402 | JS       |     | 1   | -1   | Pop the top word off the stack and jump to that location. |
| 0403 | JT       | yes | 1   | -1   | Pop the top word and jump to the location in the operand if the least-significant bit was 1 (=if true). |
| 0404 | JTS      |     | 2   | -2   | Pop first a jump target, then a condition word. Jump to the target if the condition word had its least-significant bit set. |
| 0405 | JF       | yes | 1   | -1   | Like JT but jump if the LSB of the popped word is 0 (jump if false). |
| 0406 | JFS      |     | 2   | -2   | Like JTS, but jump if false. |
| 0410 | CALL     | yes | 0   | R+1  | Push the address of the following instruction onto the return stack, then jump to the argument. |
| 0411 | RET      |     | R1  | R-1  | Pop the top value off the return stack and jump to it. |

### Load/Store Instructions

| Op   | Mnemonic | Arg | Pre | Post | Description |
|------|----------|-----|-----|------|-------------|
| 0501 | LOAD     | yes |     | 1    | Load a word from data memory at the location indicated by the operand. |
| 0502 | STOR     | yes | 1   | -1   | Store the top word on the stack in the data memory location indicated by the operand. |
| 0503 | LODS     |     | 1   | 0    | Pop the top word and load a word from data memory at the location indicated by this word. |
| 0504 | STRS     |     | 2   | -2   | Store the second-from-top word on the stack at the location indicated by the top word. Both words are popped off the stack. |

Referencing a memory location greater or equal to 4096 causes the processor to halt with an error.

### Interrupt Instruction

| Op   | Mnemonic | Arg | Pre | Post | Description |
|------|----------|-----|-----|------|-------------|
| FFFF | INT      |     |     | 0    | The processor returns control to the simulation. |

The INT instruction is useful for debugging in the simulator: when the processor encounters and INT in 'run' mode, it returns control to the simulator but you can continue to step/run from the next instruction onwards. This has the same effect as a breakpoint in other debuggers: 'run' means 'run until you encounter an interrupt'. 

## Using the Simulator

Type your assembly code in the box on the left and use the 'Assemble' menu item. If something goes wrong, it will tell you the line where a problem occurred (line numbers are shown in the status bar as you move aroudn the text box). If all is well, the left box will show the assembled machine code. You can use the 'Edit' menu to go back to the assembly code.

The simulator does not save your code, so you should save a copy in a text file yourself.

The menu items in the simulator are:

  * Assemble: shown only in edit mode, attempts to assemble your code. If successful, you see the machine code and the simulator switches to run mode. In case of an error, you get an error message.
  * Edit: shown only in run mode while the processor is halted. Switches back to edit mode.
  * Reset: shown only in run mode, resets the program counter to 0. After assembling your code, you must reset the processor to 'start it up' before running it.
  * Step: in run mode, executes a single instruction. The result (next instruction, stack state) is shown in the simulator. After assembling code, you must reset the processor before this menu item appears.
  * Run: in run mode, starts running the code until a HALT or INT instruction or an error occurs, or you use the interrupt menu item. After assembling code, you must reset the processor before this menu item appears.
  * Faster/Slower: decreases or increases an artificial delay of 1/2 seconds after each inst¬ruction.
  * Interrupt: displayed only when the processor is running. Stops execution after the current instruction, useful for example if you have produced an infinite loop.
  * About: displays information about the simulator.

In run mode, the top right box shows either the stack or the data memory (the first 128 bytes). You can switch between the two by clicking on the 'stack' or 'memory' headings.

## DATA Labels and Constants

The assembler supports giving names to memory locations in the data memory to simulate 'variables'. To use this, at the end of your program, place the command .DATA on a line on its own. After this line, you may declare memory locations starting at 0 with labels and an optional length as a decimal number, for example 'x:' (without quotes, on a line of its own) would create a one-word variable called x at the current location whereas 'y: 3' would allocate the next three words and make the label y point to the first one. For example, the following declarations:

```
.DATA
x:
y:
array: 6
z:
```

would define the labels x = 0, y = 1, array = 2, z = 8 (words 2–7 are allocated to 'array'). You could now write 'STOR y' in your code to store the top value on the stack in memory location 1 or 'PUSH array' to load the base address of the array onto the stack.

You can also write constants into data memory that will be loaded whenever you reset the processor. To do this, write an equals sign after the length specifier in a data declaration, then write the constants as hexadecimal numbers separated by spaces, with optional leading zeroes. Do not write an '0x' specifier, and use uppercase letters for digits beyond 9. For example:

```
.DATA
x: 1 = 0
y: 1 = 100
array: 6 = 0A00 0B00 0C00 0D00 0E00 0F00 
z: 1 = 4
```

This would initialise the data memory to the following values:

    [0x0000, 0x0100, 0x0A00, 0x0B00, 0x0C00, 0x0D00, 0x0E00, 0x0F00, 0x0004] 

with all further memory locations set to the all-zero value.
