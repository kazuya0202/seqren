package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	kz "github.com/kazuya0202/kazuya0202"
)

// FileUtilify ...
type FileUtilify struct {
	Files    []string
	FileConf []FileNameConfig
	Name     string
	Format   string
	Len      struct {
		File    int
		Conf    int
		MaxNext int
		MaxPrev int
	}
}

// FileNameConfig ...
type FileNameConfig struct {
	Prev string
	Next string
}

func (fu *FileUtilify) determineRenamingTarget() {
	var existConf []FileNameConfig

	seq := 0 // start index
	fileIdx := 0

	for fileIdx < fu.Len.File {
		seq++

		file := fu.Files[fileIdx]
		ext := path.Ext(file)
		name := fmt.Sprintf(fu.Format, fu.Name, seq, ext)

		fnc := FileNameConfig{
			Prev: path.Join(defs.Path, file),
			Next: path.Join(defs.Path, name),
		}

		if kz.Exists(fnc.Next) {
			existConf = append(existConf, fnc)
			fileIdx++
			continue
		}

		fu.appendConfig(fnc)
		fileIdx++
	}

	xs := kz.ToInterfaceArray(existConf)
	ys := kz.ToInterfaceArray(fu.FileConf)
	for _, c := range kz.Cartesian(xs, ys) {
		x := c[0].(FileNameConfig)
		y := c[1].(FileNameConfig)

		// last processing
		//  => After renamed exist name to another name
		if x.Next == y.Prev {
			fu.appendConfig(x)
		}
	}
	// set length
	fu.Len.Conf = len(fu.FileConf)
}

func (fu *FileUtilify) insertConfig(fnc FileNameConfig, i int) {
	fu.FileConf = append(fu.FileConf[:i], append([]FileNameConfig{fnc}, fu.FileConf[i:]...)...)
}

func (fu *FileUtilify) appendConfig(fnc FileNameConfig) {
	fu.FileConf = append(fu.FileConf, fnc)
	fu.updateNameLength(fnc)
}

func (fu *FileUtilify) updateNameLength(fnc FileNameConfig) {
	// update next, prev length
	fu.Len.MaxNext = kz.Max(fu.Len.MaxNext, len(fnc.Next))
	fu.Len.MaxPrev = kz.Max(fu.Len.MaxPrev, len(fnc.Prev))
}

func (fu *FileUtilify) displayAction(num int) {
	// adjust length
	prevTitle := "prev"
	nextTitle := "next"
	prevLen := kz.Max(fu.Len.MaxPrev, len(prevTitle))
	nextLen := kz.Max(fu.Len.MaxNext, len(nextTitle))

	format := fmt.Sprintf("%%-%ds | %%-%ds\n", prevLen, nextLen)

	println()
	fmt.Printf(format, prevTitle, nextTitle)
	fmt.Printf("%s-+-%s\n", strings.Repeat("-", prevLen), strings.Repeat("-", nextLen))

	var i int
	for i = 0; i < num-1; i++ {
		fmt.Printf(format, fu.FileConf[i].Prev, fu.FileConf[i].Next)
	}

	if fu.Len.Conf == i+1 {
		fmt.Printf(format, fu.FileConf[i].Prev, fu.FileConf[i].Next)
	} else if fu.Len.Conf > num-1 {
		fmt.Printf(format, "...", "...")
		fmt.Printf(format, fu.FileConf[fu.Len.File-1].Prev, fu.FileConf[fu.Len.File-1].Next)
	}
	println()

	if !defs.IsForce {
		yn := kz.GetInput("Execute OK? (y/n)")

		regs := []string{"y", "Y", "yes", "Yes", "YES"}
		itf := kz.StringArrayToInterfaceArray(regs)

		// other than `yes`
		if !kz.HasValueInArray(yn, itf) {
			status.DisplayInfo("Interrrupted.")
			os.Exit(0)
		}

	}
	for _, fc := range fu.FileConf {
		if err := os.Rename(fc.Prev, fc.Next); err != nil {
			fmt.Println(err)
		}
	}
	println()
	status.DisplayInfo("Finished renaming.")
}
