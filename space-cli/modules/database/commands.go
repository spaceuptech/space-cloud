package database

import (
	"fmt"

	"github.com/spaceuptech/space-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Commands is the list of commands the database module exposes
func Commands() []*cobra.Command {
	var generateSubCommands = &cobra.Command{}

	var generaterule = &cobra.Command{
		Use:  "db-rules",
		RunE: actionGenerateDBRule,
	}

	var generateconfig = &cobra.Command{
		Use:  "db-config",
		RunE: actionGenerateDBConfig,
	}

	var generateschema = &cobra.Command{
		Use:  "db-schema",
		RunE: actionGenerateDBSchema,
	}

	var getSubCommands = &cobra.Command{}

	var getrule = &cobra.Command{
		Use:  "db-rules",
		RunE: actionGenerateDBRule,
	}

	var getconfig = &cobra.Command{
		Use:  "db-config",
		RunE: actionGetDbConfig,
	}

	var getschema = &cobra.Command{
		Use:  "db-schema",
		RunE: actionGetDbSchema,
	}

	generateSubCommands.AddCommand(generaterule)
	generateSubCommands.AddCommand(generateconfig)
	generateSubCommands.AddCommand(generateschema)
	getSubCommands.AddCommand(getrule)
	getSubCommands.AddCommand(getconfig)
	getSubCommands.AddCommand(getschema)

	command := make([]*cobra.Command, 0)
	command = append(command, generateSubCommands)
	command = append(command, getSubCommands)
	return command
}

// GenerateSubCommands is the list of commands the database module exposes
// var GenerateSubCommands = []cli.Command{
// 	{
// 		Name:   "db-rules",
// 		Action: actionGenerateDBRule,
// 	},
// 	{
// 		Name:   "db-config",
// 		Action: actionGenerateDBConfig,
// 	},
// 	{
// 		Name:   "db-schema",
// 		Action: actionGenerateDBSchema,
// 	},
// }

// // GetSubCommands is the list of commands the operations module exposes
// var GetSubCommands = []cli.Command{
// 	{
// 		Name:   "db-rules",
// 		Action: actionGetDbRules,
// 	},
// 	{
// 		Name:   "db-config",
// 		Action: actionGetDbConfig,
// 	},
// 	{
// 		Name:   "db-schema",
// 		Action: actionGetDbSchema,
// 	},
// }

func actionGetDbRules(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.CalledAs()

	params := map[string]string{}
	switch len(args) {
	case 1:
		params["dbAlias"] = args[0]
	case 2:
		params["dbAlias"] = args[0]
		params["col"] = args[1]
	}
	objs, err := GetDbRule(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetDbConfig(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.CalledAs()

	params := map[string]string{}
	if len(args) != 0 {
		params["dbAlias"] = args[0]
	}
	objs, err := GetDbConfig(project, commandName, params)
	if err != nil {
		return err
	}

	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGetDbSchema(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project := viper.GetString("project")
	commandName := cmd.CalledAs()

	params := map[string]string{}
	switch len(args) {
	case 1:
		params["dbAlias"] = args[0]
	case 2:
		params["dbAlias"] = args[0]
		params["col"] = args[1]
	}

	objs, err := GetDbSchema(project, commandName, params)
	if err != nil {
		return err
	}
	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGenerateDBRule(cmd *cobra.Command, args []string) error {
	argsArr := args
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateDBRule()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionGenerateDBConfig(cmd *cobra.Command, args []string) error {
	argsArr := args
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateDBConfig()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

func actionGenerateDBSchema(cmd *cobra.Command, args []string) error {
	argsArr := args
	if len(argsArr) != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	dbruleConfigFile := argsArr[0]
	dbrule, err := generateDBSchema()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}
