package vmx

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// Translate all .vm source code into .asm code.
func Translate(file string) {
	// get the list of source files to translate
	srcfiles, err := GetSourceFiles(file)
	if err != nil {
		fmt.Println(err)
	}

	// loop over all source files, translating each
	for _, s := range srcfiles {
		TranslateSingleFile(s)
	}
}

func GetSourceFiles(src string) ([]string, error) {
	log.Info().Msgf(`source file/dir is "%s"`, src)

	// get information about the input path
	stat, err := os.Stat(src)
	if err != nil {
		log.Fatal().Err(err)
	}

	// if `input` is a directory, glob out all the .vm files in it
	if stat.IsDir() {
		// glob out all the *.vm files
		log.Info().Msgf(`"%s" is a dir`, src)
		vmfiles, err := filepath.Glob(src + "/*.vm")
		log.Info().Msgf(`*.vm glob matches: "%v"`, vmfiles)
		if err != nil {
			log.Fatal().Err(err)
		}
		return vmfiles, nil
	}

	// it's not a directory, so return a slice with a single filename
	log.Info().Msgf(`"%s" is a file`, src)
	return []string{src}, nil
}

// Translate VM instructions in a single source file to equivalent .asm.
// Note: I've isolated this functionality in its own function to make `defer`
//       work better.
func TranslateSingleFile(srcfile string) error {
	outfile := strings.TrimSuffix(srcfile, filepath.Ext(srcfile)) + ".asm"
	log.Info().Msgf(`Translating source file "%s" to "%s"`, srcfile, outfile)

	// open the input file for reading
	src, err := os.Open(srcfile)
	if err != nil {
		fmt.Println(err)
	}
	defer src.Close()

	// open the output file for writing
	out, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()

	// read & translate the input file line by line
	scanner := bufio.NewScanner(src)
	lineno := 0
	module := filepath.Base(srcfile)
	for scanner.Scan() {
		lineno++
		cmd := NewCommand(scanner.Text(), fmt.Sprintf("%s.%d", module, lineno))
		if asm := cmd.GetAsm(); len(asm) > 0 {
			out.WriteString(cmd.GetAsm() + "\n")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return err
}
