package vmx

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// Translate .vm source code into .asm code. Create one CodeWriter object,
// and one Parser object for each source file.
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

func GetSourceFiles(input string) ([]string, error) {
	log.Info().Msgf(`input file/dir is "%s"`, input)

	// get information about the input path
	stat, err := os.Stat(input)
	if err != nil {
		log.Fatal().Err(err)
	}

	// if `path` is a directory, glob out all the .vm files in it
	if stat.IsDir() {
		// glob out all the *.vm files
		log.Info().Msgf(`"%s" is a dir`, input)
		vmfiles, err := filepath.Glob(input + "/*.vm")
		log.Info().Msgf(`*.vm glob matches: "%v"`, vmfiles)
		if err != nil {
			log.Fatal().Err(err)
		}
		return vmfiles, nil
	}

	// it's not a directory, so return a slice with a single filename
	log.Info().Msgf(`"%s" is a file`, input)
	return []string{input}, nil
}

// translate VM instructions in a single file to equivalent .asm
// Note: put this functionality in its own function to make `defer` work better
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
	for scanner.Scan() {
		cmd := NewCommand(scanner.Text()) // FIXME: return pointer?
		if asm := cmd.GetAsm(); len(asm) > 0 {
			out.WriteString(cmd.GetAsm() + "\n")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return err
}
