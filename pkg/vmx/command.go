package vmx

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

func NewCommand(src string, unique string) *Command {
	return &Command{
		vmCommand: src,
		unique:    unique,
		fields:    strings.Fields(src),
	}
}

func (c *Command) GetAsm() string {

	// weed out comments & blank lines
	if regexp.MustCompile(`^\s*(//.*)?$`).MatchString(c.vmCommand) {
		return ``
	}

	// arithmetic commands
	if regexp.MustCompile(`^\s*(add|sub|neg|eq|gt|lt|and|or|not)\s*$`).MatchString(c.vmCommand) {
		return c.Arithmetic()
	}

	// push
	if c.fields[0] == "push" {
		return c.Push()
	}

	// pop
	if c.fields[0] == "pop" {
		return c.Pop()
	}

	// catchall
	return fmt.Sprintf("ERROR: I don't recognize VM command '%s'\n", c.vmCommand)
}

func (c *Command) Arithmetic() string {
	var asm string
	switch c.fields[0] {
	case "add":
		// pop top of stack into D, then add it to *(SP-1)
		asm = `
			// add
			@SP        // SP--
			M=M-1
			A=M        // D = *SP
			D=M
			A=A-1      // *(SP-1) = *(SP-1) + D 		
			M=D+M
		`
	case "sub":
		// pop top of stack into D, then subtract it from *(SP-1)
		asm = `
			// sub
			@SP        // SP--
			M=M-1
			A=M        // D = *SP
			D=M
			A=A-1      // *(SP-1) = *(SP-1) - D
			M=M-D
		`
	case "neg":
		// stack size doesn't change, we can modify *(SP-1) in place
		asm = `
			// neg
			@SP        // *(SP-1) = -(*(SP-1))
			A=M
			A=A-1
			M=-M
		`
	case "eq":
		// FIXME: need to make this cleaner.
		labelIf := "eq_if." + c.unique
		labelElse := "eq_else." + c.unique
		asms := []string{`
			// eq
			@SP        // SP--
			AM=M-1     // D = *SP
			D=M
			A=A-1      // *(SP-1) = *(SP-1) - D
			D=M-D
			`,
			"@" + labelIf,
			`
			D;JEQ       // jump to eq_if 
            D=0         // ELSE block D = false
			@labelEnd
			0;JMP
			(@labelIf)
			D=-1        // D = true
			@labelEnd
            @SP         // *SP = D
			A=M-1
			M=D
		`}
		asm = strings.Join(asms, "\n")
	case "gt":
		// pop top of stack into D
		asm = `
			// lt
			@SP        // SP--
			AM=M-1
			D=M        // D = *SP
			A=A-1      // D = *(SP-1) - D
			D=M-D
			// D = D & 0x8000
			M=!M       // *(SP-1) = ! *(SP-1)
		`
	case "lt":
		// a lot like sub but check the MSBit
		asm = `
			// lt
			@SP
			AM=M-1
			D=M
			A=A-1
			D=M-D

		`
	case "and":
		asm = `
			// sub
			@SP        // SP--
			M=M-1
			A=M        // D = *SP
			D=M
			A=A-1      // *(SP-1) = *(SP-1) & D
			M=M&D
		`
	case "or":
		asm = `
			// sub
			@SP        // SP--
			M=M-1
			A=M        // D = *SP
			D=M
			A=A-1      // *(SP-1) = *(SP-1) | D
			M=M|D
		`
	case "not":
		// a lot like neg, change *(SP-1) in place
		asm = `
			// not
			@SP        // *(SP-1) = ! *(SP-1)
			A=M
			A=A-1
			M=!M
		`
	default:
		asm = "ERROR"
		log.Fatal().Msgf(`Unrecognized VM command: %s`, c.fields[0])
	}
	return Trim(asm)
}

func (c *Command) Push() string {
	segment := c.fields[1]
	var asm string
	switch segment {
	case `constant`:
		asm = fmt.Sprintf("// %s\n", c.vmCommand) +
			fmt.Sprintf("@%-8s  // D = %s\n", c.fields[2], c.fields[2]) +
			`D=A
			@SP        // *SP = D
			A=M
			M=D
			@SP        // SP++
			M=M+1
			`
	default:
		asm = `ERROR`
		log.Fatal().Msgf(`Unrecognized segment: %s`, segment)
	}
	return Trim(asm)
}

func (c *Command) Pop() string {
	return `ERROR\n`
}

// trim leading whitespace on each line in a multiline string
// and remove empty lines
func Trim(str string) string {
	var out string
	for _, element := range strings.Split(str, "\n") {
		if len(strings.TrimSpace(element)) != 0 {
			out += strings.TrimSpace(element) + "\n"
		}
	}
	return out
}
