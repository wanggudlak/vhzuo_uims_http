package version

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"uims/command"
	"uims/pkg/color"
)

const version = "1.0.0"

const verboseVersionBanner string = `%s%s______
 -    - -    --     -- 
| |  | | |  /  \   /  \ 
| |  | | | / /\ \ / /\ \
| |__| | |/ /  \ | /  \ \ 
\ ____ /_|_/    \_/    \_|
 v{{ .UIMSVersion }}%s
%s%s
├── GinFramework     : {{ .GinVersion }}
├── GoVersion        : {{ .GoVersion }}
├── GOOS             : {{ .GOOS }}
├── GOARCH           : {{ .GOARCH }}
├── NumCPU           : {{ .NumCPU }}
├── Compiler         : {{ .Compiler }}
└── Date             : {{ Now "Monday, 14 June 2020" }}%s
`

const shortVersionBanner = `
 -    - -    --     -- 
| |  | | |  /  \   /  \ 
| |  | | | / /\ \ / /\ \
| |__| | |/ /  \ | /  \ \ 
\ ____ /_|_/    \_/    \_|
 v{{ .UIMSVersion }}
`

var VersionSN string

var CmdVersion = &command.Command{
	UsageLine: "version",
	Short:     "输出当前UIMS的版本信息",
	Long: `
输出当前UIMS，Gin框架，Go的版本信息以及相关的平台信息。
`,
	Run: versionCMD,
}

var outputFormat string

func init() {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	fs.StringVar(&outputFormat, "o", "", "Set the output format. Either json or yaml.")
	CmdVersion.Flag = *fs
	command.CMD.Register(CmdVersion)
}

func versionCMD(cmd *command.Command, args []string) int {
	_ = cmd.Flag.Parse(args)
	stdout := cmd.Out()

	if outputFormat != "" {
		runtimeInfo := RuntimeInfo{
			GetGoVersion(),
			runtime.GOOS,
			runtime.GOARCH,
			runtime.NumCPU(),
			runtime.Compiler,
			version,
			GetGinVersion(),
		}
		switch outputFormat {
		case "json":
			{
				b, err := json.MarshalIndent(runtimeInfo, "", "    ")
				if err != nil {
					logrus.Error(err.Error())
				}
				fmt.Println(string(b))
				return 0
			}
		case "yaml":
			{
				b, err := yaml.Marshal(&runtimeInfo)
				if err != nil {
					logrus.Error(err.Error())
				}
				fmt.Println(string(b))
				return 0
			}
		}
	}

	coloredBanner := fmt.Sprintf(verboseVersionBanner, "\x1b[35m", "\x1b[1m",
		"\x1b[0m", "\x1b[32m", "\x1b[1m", "\x1b[0m")
	InitBanner(stdout, bytes.NewBufferString(coloredBanner))
	return 0
}

// ShowShortVersionBanner prints the short version banner.
func ShowShortVersionBanner() {
	output := color.NewColorWriter(os.Stdout)
	InitBanner(output, bytes.NewBufferString(color.MagentaBold(shortVersionBanner)))
}

func GetGoVersion() string {
	var (
		cmdOut []byte
		err    error
	)

	if cmdOut, err = exec.Command("go", "version").Output(); err != nil {
		logrus.Fatalf("There was an error running 'go version' command: %s", err)
	}
	return strings.Split(string(cmdOut), " ")[2]
}

func GetGinVersion() string {
	return gin.Version
}
