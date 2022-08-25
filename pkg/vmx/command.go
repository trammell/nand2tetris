package vmx

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

func NewCommand(src string) Command {
	return Command{vmCommand: src, fields: strings.Fields(src)}
}

func (c *Command) GetAsm() string {

	// weed out comments & blank lines
	if regexp.MustCompile(`^\s*(//.*)?$`).MatchString(c.vmCommand) {
		return ``
	}

	// arithmetic commands
	if regexp.MustCompile(`^\s*(add|sub|neg|eq|gt)\s*$`).MatchString(c.vmCommand) {
		return c.Math()
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

func (c *Command) Math() string {
	var asm string
	switch c.fields[0] {
	case "add":
		asm = `
			// add
			@SP        // SP--
			M=M-1
			A=M        // D = *SP
			D=M
			A=A-1      // *(SP-1) = *(SP-1) + D 		
			M=D+M`
	case "sub":
		asm = `
			// sub
			@SP        // SP--
			M=M-1
			A=M        // D = *SP
			D=M
			A=A-1      // *(SP-1) = *(SP-1) - D
			M=M-D`
	case "neg":
		asm = `
			// neg
			@SP        // *(SP-1) = -(*(SP-1))
			A=M
			A=A-1
			M=-M
		`
	default:
		asm = "ERROR"
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
