package main

// сюда писать функцию DirTree
import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func plot(out io.Writer, info os.FileInfo, isEnd bool, printFiles bool) {
	var size string
	if !info.IsDir() {
		if info.Size() == 0 {
			size = " (empty)"
		} else {
			size = " (" + strconv.FormatInt(info.Size(), 10) + "b)"
		}
	}

	if isEnd {
		fmt.Fprintf(out, "└───%s%s\n", info.Name(), size)
		return
	}
	fmt.Fprintf(out, "├───%s%s\n", info.Name(), size)
}

func draw(out io.Writer, dataToDraw []string) {
	for _, value := range dataToDraw {
		fmt.Fprintf(out, value)
	}
}

func plotTree(out io.Writer, path string, printFiles bool, dataToPlot []string) error {
	filesAndDirs, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var AllDirs int
	for _, dataFromDir := range filesAndDirs {
		if dataFromDir.IsDir() == true {
			AllDirs++
		}
	}

	var isEnd bool
	var isAllDirs bool
	var countDirs int
	lengthOfDir := len(filesAndDirs) - 1

	for index, dataFromDir := range filesAndDirs {
		if AllDirs-1 == countDirs {
			isAllDirs = true
		}

		if index == lengthOfDir || isAllDirs && !printFiles {
			isEnd = true
		}

		if printFiles {
			draw(out, dataToPlot)
			plot(out, dataFromDir, isEnd, printFiles)
		}

		if dataFromDir.IsDir() {
			countDirs++
			if !printFiles {
				draw(out, dataToPlot)
				plot(out, dataFromDir, isEnd, printFiles)
			}

			if isEnd {
				dataToPlot = append(dataToPlot, "\t")
				err := plotTree(out, path+"/"+dataFromDir.Name(), printFiles, dataToPlot)
				if err != nil {
					return err
				}
				if !printFiles {
					dataToPlot = dataToPlot[:len(dataToPlot)-1]
				}
				continue
			}
			dataToPlot = append(dataToPlot, "│\t")
			err := plotTree(out, path+"/"+dataFromDir.Name(), printFiles, dataToPlot)
			if err != nil {
				return err
			}
			dataToPlot = dataToPlot[:len(dataToPlot)-1]
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	dataToPlot := make([]string, 0, 1)
	err := plotTree(out, path, printFiles, dataToPlot)
	if err != nil {
		return err
	}
	return nil
}
