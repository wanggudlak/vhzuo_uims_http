package cmd_template

import (
	"uims/command"
	_ "uims/command/commands/app_key_generator"
	_ "uims/command/commands/app_rsakey_generator"
	_ "uims/command/commands/make"
	_ "uims/command/commands/makemigration"
	_ "uims/command/commands/migrate_data"
	_ "uims/command/commands/migrator"
	_ "uims/command/commands/server"
	_ "uims/command/commands/thrift_RPC_server"
	_ "uims/command/commands/version"
	"uims/pkg/tool"
)

var usageTemplate = `UIMS is a user and authority management system based on THE GIN framework.

{{"用法(注意：如果uims可执行文件不在系统PATH中，请使用<uims所在的目录>/uims)" | headline}}
    {{"uims command [arguments]" | bold}}

{{"可使用的命令" | headline}}
{{range .}}{{if .Runnable}}
    {{.Name | printf "%-75s" | greenbold }} {{.Short | green }}{{ end }}{{ end }}

Use {{"uims help [command]" | bold}} for more information about a command.

{{"额外的帮助" | headline}}
{{range .}}{{if not .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short | green }}{{end}}{{end}}

Use {{"uims help [topic]" | bold}} for more information about that topic.
`

var helpTemplate = `{{"用法" | headline}}
  {{.UsageLine | printf "uims %s" | bold}}
{{if .Options}}{{endline}}{{"OPTIONS" | headline}}{{range $k,$v := .Options}}
  {{$k | printf "-%s" | bold}}
      {{$v}}
  {{end}}{{end}}
{{"DESCRIPTION" | headline}}
  {{tmpltostr .Long . | trim}}
`

var ErrorTemplate = `uims: %s.
Use {{"uims help" | bold}} for more information.
`

func Usage() {
	tool.TmplTextParseAndOutput(usageTemplate, command.CMD.Commands)
}

func Help(args []string) {
	if len(args) == 0 {
		Usage()
	}
	if len(args) != 1 {
		tool.PrintErrorTmplAndExit("Too many arguments", ErrorTemplate)
	}

	arg := args[0]

	if cmd, ok := command.CMD.Get(arg); ok {
		tool.TmplTextParseAndOutput(helpTemplate, cmd)
	} else {
		tool.PrintErrorTmplAndExit("Unknown help topic", ErrorTemplate)
	}
}
