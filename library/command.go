package library

import (
	"context"
	"errors"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
	"sync"
)

type CliCommand struct {
	name        string
	description string
	usage       string
	configs     []*CommandConfig
	cliApp      *cli.App
	mux         sync.Mutex
	signatures  map[string]*signature
	execCommand *ExecCommand
	initCommand *cli.Command
}

type ExecCommand struct {
	currentCmdName  string
	isCliScriptExec bool
	arguments       map[string][]string
	options         map[string]interface{}
	isEndless       bool
	mu              sync.RWMutex
}

type signature struct {
	description string
	name        string
	category    string
	flags       []*flag
	handleFunc  func(ctx context.Context, cmd *ExecCommand) error
	isEndless   bool
}

type flag struct {
	signType     string
	name         string
	defaultValue string
	description  string
	required     bool
	isArrayValue bool
}

func NewExecCommand() *ExecCommand {
	return &ExecCommand{
		options:   map[string]interface{}{},
		arguments: map[string][]string{},
	}
}

func (cmd *ExecCommand) SetCurrentCmdName(name string) {
	cmd.mu.Lock()
	cmd.currentCmdName = name
	cmd.mu.Unlock()
}

func (cmd *ExecCommand) SetIsEndless(endless bool) {
	cmd.mu.Lock()
	cmd.isEndless = endless
	cmd.mu.Unlock()
}

func (cmd *ExecCommand) SetOption(key string, value interface{}) {
	cmd.mu.Lock()
	cmd.options[key] = value
	cmd.mu.Unlock()
}

func (cmd *ExecCommand) Options() map[string]interface{} {
	cmd.mu.RLock()
	defer cmd.mu.RUnlock()
	return cmd.options
}

func (cmd *ExecCommand) Arguments() map[string][]string {
	cmd.mu.RLock()
	defer cmd.mu.RUnlock()
	return cmd.arguments
}

func (cmd *CliCommand) CustomInitCommand(cliCmd *cli.Command, execCmd *ExecCommand) {
	cmd.initCommand = cliCmd
	cmd.execCommand = execCmd
}

func (cmd *CliCommand) AddConfig(config ...*CommandConfig) {
	cmd.mux.Lock()
	defer cmd.mux.Unlock()
	cmd.configs = append(cmd.configs, config...)
}

func (cmd *ExecCommand) IsScriptExec() bool {
	return cmd.isCliScriptExec
}

func (cmd *ExecCommand) IsEndless() bool {
	return cmd.isEndless
}

func (cmd *CliCommand) Setup() (err error) {
	//TODO 是否已经初始化的判断
	cmd.mux.Lock()
	defer cmd.mux.Unlock()
	for _, config := range cmd.configs {
		s, e := parseSignature(config.Signature)
		s.handleFunc = config.HandleFunc
		s.description = config.Description
		s.isEndless = config.IsEndless
		if e != nil {
			return e
		}
		cmd.signatures[s.name] = s
	}

	var subCommands []*cli.Command
	for name, sign := range cmd.signatures {
		command := &cli.Command{
			Name:        name,
			Usage:       sign.description,
			Description: sign.description,
			Category:    sign.category,
		}

		for _, sflag := range sign.flags {
			if sflag.signType == "option" {
				if sflag.isArrayValue {
					f := &cli.StringSliceFlag{
						Name:     sflag.name,
						Usage:    sflag.description,
						Required: sflag.required,
					}
					command.Flags = append(command.Flags, f)
				} else {
					f := &cli.StringFlag{
						Name:     sflag.name,
						Usage:    sflag.description,
						Value:    sflag.defaultValue,
						Required: sflag.required,
					}
					command.Flags = append(command.Flags, f)
				}

			}
		}
		command.Action = func(ctx *cli.Context) error {
			cmd.execCommand.currentCmdName = ctx.Command.Name
			cmd.execCommand.isCliScriptExec = true
			s := cmd.signatures[ctx.Command.Name]
			cmd.execCommand.options["config"] = ctx.String("config")
			cmd.execCommand.options["skip-schedule"] = ctx.Bool("skip-schedule")
			for _, sflag := range s.flags {
				if sflag.signType == "option" {
					cmd.execCommand.options[sflag.name] = ctx.String(sflag.name)
				} else {
					cmd.execCommand.arguments[sflag.name] = ctx.Args().Slice()
				}
			}
			return nil
		}
		subCommands = append(subCommands, command)
	}

	if cmd.initCommand == nil {
		cmd.initCommand = &cli.Command{
			Name:  "run",
			Usage: "run http service",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "config",
					Value:   ".env.local",
					Aliases: []string{"c"},
				},
				&cli.IntFlag{
					Name:    "port",
					Value:   2345,
					Aliases: []string{"p"},
				},
				&cli.StringFlag{
					Name:    "script",
					Value:   "",
					Aliases: []string{"s"},
				},
			},
			Action: func(context *cli.Context) error {
				cmd.execCommand.currentCmdName = context.Command.Name
				cmd.execCommand.isEndless = true
				cmd.execCommand.options["config"] = context.String("config")
				cmd.execCommand.options["port"] = context.Int("port")
				cmd.execCommand.options["script"] = context.String("script")
				return nil
			},
		}
	}

	app := &cli.App{
		Name:        cmd.name,
		Usage:       cmd.usage,
		Description: cmd.description,
		Commands: []*cli.Command{
			cmd.initCommand,
			{
				Name:  "script",
				Usage: "run script",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Value:   ".env.local",
						Aliases: []string{"c"},
					},
					&cli.BoolFlag{
						Name:    "skip-schedule",
						Value:   true,
						Aliases: []string{"k"},
					},
				},
				Subcommands: subCommands,
			},
		},
	}
	cmd.cliApp = app
	err = app.Run(os.Args)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *CliCommand) ExecCommand() *ExecCommand {
	return cmd.execCommand
}

func (cmd *CliCommand) Exec(ctx context.Context) (err error) {
	sign, ok := cmd.signatures[cmd.execCommand.currentCmdName]
	if ok && sign.handleFunc != nil {
		err = sign.handleFunc(ctx, cmd.execCommand)
	}
	return err
}

func parseSignature(s string) (sign *signature, err error) {
	sign = &signature{}
	trimString := strings.Trim(s, " ")
	if len(trimString) == 0 {
		return nil, errors.New("command signature is not exist")
	}
	seg := strings.Split(trimString, "{")
	if len(seg) >= 1 {
		sign.name, sign.category, err = parseCommandName(seg[0])
		if err != nil {
			return nil, err
		}
		for i := 1; i < len(seg); i++ {
			f, e := parseFlag("{" + seg[i])
			if e != nil {
				return nil, e
			}
			sign.flags = append(sign.flags, f)
		}
	}
	return sign, nil
}

func parseFlag(s string) (f *flag, err error) {
	trimString := strings.Trim(s, " ")
	if len(trimString) == 0 {
		return nil, errors.New("command signature with error flag:" + trimString)
	}
	if strings.Index(trimString, "{") != 0 || strings.Index(trimString, "}") != len(trimString)-1 {
		return nil, errors.New("command signature with error flag:" + trimString)
	}
	flagName := strings.Trim(trimString[1:len(trimString)-1], " ")
	if len(flagName) == 0 {
		return nil, errors.New("command signature with empty flag:" + trimString)
	}
	f = &flag{}

	// Options type:
	// no value {--sendEmail : with Description}
	// required a value {--password=}
	// withDefault {--queue=default}
	// array options {--queue=*}
	if strings.Index(flagName, "--") != -1 && strings.Index(flagName, "--") == 0 {
		f.signType = "option"
		seg := strings.Split(flagName, ":")
		if len(seg) > 1 {
			//描述信息
			for i, v := range seg {
				if i > 0 {
					f.description = f.description + v
				}
			}
		}
		seg1 := strings.Split(seg[0], "=")
		if strings.Index(trimString, "=") != -1 {
			f.required = true
		}
		if len(seg1) == 2 && seg1[1] != "" {
			//有默认值
			f.required = false
			if strings.Trim(seg1[1], " ") != "*" {
				f.defaultValue = strings.Trim(seg1[1], " ")
			} else {
				f.isArrayValue = true
			}

		}

		f.name = strings.TrimLeft(strings.Trim(seg1[0], " "), "--")
		return f, nil
	}

	//Arguments Type:
	// required {userId : with Description}
	// option {userId?}
	// withDefault {userId=1}
	// array options {userIds*}
	if strings.Index(flagName, "--") == -1 {
		f.signType = "argument"
		seg := strings.Split(flagName, ":")
		if len(seg) > 1 {
			//描述信息
			for i, v := range seg {
				if i > 1 {
					f.description = f.description + v
				}
			}
		}
		seg1 := strings.Split(seg[0], "=")

		if len(seg1) == 2 {
			//有默认值
			f.defaultValue = strings.Trim(seg1[1], " ")
		}

		f.required = true

		if string(seg[0][len(seg[0])-1]) == "?" {
			f.required = false
		}

		//TODO userIds*? 这种情况的处理
		if string(seg[0][len(seg[0])-1]) == "*" {
			f.isArrayValue = true
		}

		f.name = strings.Trim(seg1[0], " ")
		return f, nil
	}
	return nil, errors.New("command signature with error flag:" + s)
}

func parseCommandName(s string) (commandName, category string, err error) {
	trimString := strings.Trim(s, " ")
	if len(trimString) == 0 {
		return "", "", errors.New("command signature is not exist")
	}
	seg := strings.Split(trimString, ":")
	if len(seg) == 1 {
		return trimString, "", nil
	}

	return trimString, seg[0], nil
}

func NewCliCommand(AppName, AppDescription, AppUsage string) *CliCommand {
	return &CliCommand{
		name:        AppName,
		description: AppDescription,
		usage:       AppUsage,
		signatures:  map[string]*signature{},
		execCommand: &ExecCommand{
			options:   map[string]interface{}{},
			arguments: map[string][]string{},
		},
	}
}
