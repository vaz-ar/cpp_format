/*
C++ code fomatting tool
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func indent_connects(lines []string) {

	re_start := regexp.MustCompile(`^(\s*connect\(\S+,)\s*(&\w+::\w+,)\s*$`)
	re_end := regexp.MustCompile(`^(\s*\S+,)\s*(&\w+::\w+\);)\s*$`)

	var (
		ix_start, ix_end          []int
		len_start, len_end, pad   int
		format_connect, prev_line string
		found_connect             bool = false
	)

	for i, line := range lines {

		if !found_connect {
			if ix_start = re_start.FindStringSubmatchIndex(line); ix_start != nil {
				found_connect = true
				prev_line = line
			}
		} else {
			if ix_end = re_end.FindStringSubmatchIndex(line); ix_end != nil {
				len_start = ix_start[3]
				len_end = ix_end[3]

				if len_start > len_end {
					pad = len_start
				} else {
					pad = len_end
				}
				// Use 4 spaces to align
				if mod := pad % 4; mod != 0 {
					pad += mod
				}

				format_connect = fmt.Sprint("%-", pad, "s%s")
				lines[i-1] = fmt.Sprintf(format_connect, prev_line[ix_start[2]:ix_start[3]], prev_line[ix_start[4]:ix_start[5]])
				lines[i] = fmt.Sprintf(format_connect, line[ix_end[2]:ix_end[3]], line[ix_end[4]:ix_end[5]])
			}
			found_connect = false
		}
	}
}

func format(lines []string) {

	re_comment_caps := regexp.MustCompile(`//\s{0,3}(this|if|for|while|delete|void|return|qDebug|m\S+?_\S+)?(\S)`)
	re_dox_backslash := regexp.MustCompile(`(\* )@(brief|param(?:\[(?:in|out)\])?|return)`)
	re_dox_caps := regexp.MustCompile(`(\* \\(?:brief|return))\s+([a-zA-Z])([a-zA-Z]+(?:_|::)[a-zA-Z]+)?`)
	re_dox_colon := regexp.MustCompile(`(\* \\param(?:\[(?:in|out)\])?\s\S+\s)(\S{2,})`)

	var indexes []int
	for i, line := range lines {

		// First letter of comments in capital
		if indexes = re_comment_caps.FindStringSubmatchIndex(line); indexes != nil && indexes[2] == -1 {
			// => the first capture group didn't return anything, which is what we want (nb: No negative lookahead for regexp in Go)
			// Capitalisation of the first letter
			line = fmt.Sprint(line[:indexes[4]], strings.ToUpper(line[indexes[4]:indexes[5]]), line[indexes[5]:])
		}

		// Replace "@" by "\" in doxygen comment blocks
		if indexes = re_dox_backslash.FindStringSubmatchIndex(line); indexes != nil {
			line = fmt.Sprint(line[:indexes[3]], "\\", line[indexes[4]:])
		}

		//  Capitalisation of first letter after brief/return in doxygen comment blocks
		if indexes = re_dox_caps.FindStringSubmatchIndex(line); indexes != nil && indexes[6] == -1 {
			// the third capture group didn't return anything, which is what we want (nb: No negative lookahead for regexp in Go)
			line = fmt.Sprint(line[:indexes[3]], " ", strings.ToUpper(line[indexes[4]:indexes[5]]), line[indexes[5]:])
		}

		//  Add colon after parameter name and capitalise the first letter of the parameter detail in doxygen comment blocks
		if indexes = re_dox_colon.FindStringSubmatchIndex(line); indexes != nil {
			line = fmt.Sprint(line[:indexes[3]], ": ", strings.ToUpper(line[indexes[4]:indexes[4]+1]), line[indexes[4]+1:])
		}

		// If the line used in the loop is modified we replace it in the slice
		if line != lines[i] {
			lines[i] = line
		}
	}
}

func get_file_list(target *string, ignore *string) []string {
	var ext string
	var walkFunc filepath.WalkFunc
	var files []string

	if *ignore != "" {

		content, err := ioutil.ReadFile(*ignore)
		if err != nil {
			log.Fatalln(err)
		}
		ignore_list := strings.Split(string(content), "\n")

		walkFunc = func(path string, f os.FileInfo, err error) error {
			for _, line := range ignore_list {
				if line = strings.Trim(line, " \t\n\r"); line == "" {
					continue
				}

				if matched, _ := regexp.MatchString(line, path); matched {
					return nil
				}
			}
			if ext = filepath.Ext(path); ext == ".cpp" || ext == ".h" {
				files = append(files, path)
			}
			return nil
		}

	} else {
		walkFunc = func(path string, f os.FileInfo, err error) error {
			if ext = filepath.Ext(path); ext == ".cpp" || ext == ".h" {
				files = append(files, path)
			}
			return nil
		}
	}

	// Iterate recursively over each file in the folder, calling walkFunc for each file
	err := filepath.Walk(*target, walkFunc)
	if err != nil {
		log.Fatalln(err)
	}
	return files
}

func main() {

	flag.Usage = func() {
		var sep string
		if runtime.GOOS == "windows" {
			sep = "\\"
		} else {
			sep = "/"
		}
		prog := strings.Split(os.Args[0], sep)
		fmt.Printf("\nUsage: %s target\n", prog[len(prog)-1])
		fmt.Println("Arguments description:")
		flag.PrintDefaults()
		fmt.Println()
	}

	form := flag.Bool("f", false, "Format code")
	indent := flag.Bool("ic", false, "Indent \"connect\" statements")
	ignore := flag.String("ignore", "", "Path to a file with a list of patterns to ignore (One by line), used only if target is a directory")
	// git := flag.Bool("git", false, "Use the program as a git hook")
	flag.Parse()

	// Arguments verification (Needs at least one flag + path to a file/folder)
	if flag.NFlag() <= 0 || len(flag.Args()) != 1 {
		flag.Usage()
		return
	}

	target := flag.Arg(0)
	var files []string
	fi, err := os.Stat(target)
	if err != nil {
		log.Fatalln(err)
	}

	if fi.Mode().IsDir() { // If the target is a folder
		files = get_file_list(&target, ignore)
		if len(files) == 0 {
			fmt.Println("\nThe target folder must contains c++ files")
		}
	} else if ext := filepath.Ext(target); fi.Mode().IsRegular() && (ext == ".cpp" || ext == ".h") { // Else if the target is a cpp file
		files = append(files, target)
	} else {
		fmt.Println("\nThe target must be a regular c++ file or a folder")
		return
	}

	for _, file := range files {
		// Open the file
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		// Split file content into a slice
		lines := strings.Split(string(content), "\n")

		if *indent {
			indent_connects(lines)
		}
		if *form {
			format(lines)
		}

		ioutil.WriteFile(file, []byte(strings.Join(lines, "\n")), os.ModeAppend|os.ModeExclusive)
	}
}
