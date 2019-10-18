package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// SelfPath gets compiled executable file absolute path
func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// SelfDir gets compiled executable file directory
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CreateDirectory(dir string) error {
	if FileExists(dir) {
		return nil
	}
	err := os.MkdirAll(dir, os.ModePerm)
	return err
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// get file name with suffix,eg: "/Users/root/Documents/test.txt"=>test.txt
func FileNameWithSuffix(fullFilename string) string {
	return path.Base(fullFilename)
}

//get file name without suffix,eg:  "/Users/itfanr/Documents/test.txt"=>test
func FileName(fullFilename string) string {
	filenameWithSuffix := path.Base(fullFilename)
	fileSuffix := path.Ext(filenameWithSuffix)
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	return filenameOnly
}

// SearchFile Search a file in paths.
// this is often used in search config file in /etc ~/
func SearchFile(filename string, paths ...string) (fullpath string, err error) {
	for _, path := range paths {
		if fullpath = filepath.Join(path, filename); FileExists(fullpath) {
			return
		}
	}
	err = errors.New(fullpath + " not found in paths")
	return
}

// GrepFile like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
func GrepFile(patten string, filename string) (lines []string, err error) {
	re, err := regexp.Compile(patten)
	if err != nil {
		return
	}

	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	lines = make([]string, 0)
	reader := bufio.NewReader(fd)
	prefix := ""
	var isLongLine bool
	for {
		byteLine, isPrefix, er := reader.ReadLine()
		if er != nil && er != io.EOF {
			return nil, er
		}
		if er == io.EOF {
			break
		}
		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		} else {
			isLongLine = true
		}

		line = prefix + line
		if isLongLine {
			prefix = ""
		}
		if re.MatchString(line) {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

func SaveFile(content string, fileName string, openFileFlag int) (string, error) {
	path := filepath.Dir(fileName)
	CreateDirectory(path)
	file, err := os.OpenFile(fileName, openFileFlag, 0666)
	if err != nil {
		return "", err
	}

	defer file.Close()
	writer := bufio.NewWriter(file)
	writer.WriteString(content)

	writer.Flush()
	return fileName, nil
}
