package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	kz "github.com/kazuya0202/kazuya0202"
	"github.com/spf13/cobra"
)

// DefaultArguments is struct.
type DefaultArguments struct {
	Name      string
	Path      string
	SeqNum    int
	ShowNum   int
	IsForce   bool
	IsAllShow bool
}

var (
	fu   FileUtilify
	defs DefaultArguments

	// log status
	status = kz.NewStatusString()

	rootCmd = &cobra.Command{
		Use:   "seqren [name]",
		Short: "Rename filename in sequence.",
		Long:  "Rename filename in sequence.",
		Run: func(cmd *cobra.Command, args []string) {
			rootFunction(cmd, args)
		},
	}
)

func rootFunction(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		status.DisplayInfo("Please specify rename filename string.")
		fu.Name = kz.GetInput("enter")
	} else {
		fu.Name = args[0]
		s := fmt.Sprintf("Set rename string: %s", fu.Name)
		status.DisplayInfo(s)
	}

	// validation filename
	if fu.Name == "" {
		status.DisplayErrorAndExit("'' is not valid name.")
	} else if strings.ContainsAny(fu.Name, "\\/:*?\"<>|") {
		// unusable character
		s := fmt.Sprintf("'%s' has unusable character in Windows.", fu.Name)
		status.DisplayErrorAndExit(s)
	}

	// clean path
	defs.Path = path.Clean(defs.Path)

	// files in target directory
	fu.Files = kz.GetFiles(defs.Path)
	fu.Len.File = len(fu.Files)

	if fu.Files == nil || fu.Len.File == 0 {
		s := fmt.Sprintf("'%s' may not exist directory name. or file is not exist in directory.", fu.Name)
		status.DisplayErrorAndExit(s)
	}

	if defs.SeqNum <= 0 {
		s := fmt.Sprintf("'%d' is invalid. Please specify 1 or more.", defs.SeqNum)
		status.DisplayErrorAndExit(s)
	}

	// sprintf format - `%s%0Nd%s`
	fu.Format = fmt.Sprintf("%%s%%0%dd%%s", defs.SeqNum)
	fu.determineRenamingTarget()

	if fu.Len.Conf == 0 {
		status.DisplayInfoAndExit("There is no file to rename.")
	}

	// show file num
	var showNum int
	if defs.IsAllShow {
		showNum = fu.Len.Conf
	} else {
		showNum = kz.Min(defs.ShowNum, fu.Len.Conf)
	}
	fu.displayAction(showNum)
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&defs.Path, "path", "p", "", "path of target directory")
	rootCmd.Flags().IntVarP(&defs.SeqNum, "seq", "s", 3, "N-digit 0 filling")
	rootCmd.Flags().IntVarP(&defs.ShowNum, "num", "n", 10, "display lines of renaming target")
	rootCmd.Flags().BoolVarP(&defs.IsForce, "force", "f", false, "execute command without confirmation")
	rootCmd.Flags().BoolVarP(&defs.IsAllShow, "all-show", "a", false, "display all lines of renaming target")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {}
