package remoteservices

import (
	"github.com/spf13/cobra"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// GenerateSubCommands is the list of commands the remote-services module exposes
func GenerateSubCommands() []*cobra.Command {

	var generateService = &cobra.Command{
		Use:     "remote-service [path to config file]",
		RunE:    actionGenerateService,
		Example: "space-cli generate remote-services config.yaml --project myproject --log-level info",
		Aliases: []string{"remote-services"},
	}
	return []*cobra.Command{generateService}
}

// GetSubCommands is the list of commands the remote-services module exposes
func GetSubCommands() []*cobra.Command {

	var getServices = &cobra.Command{
		Use:               "remote-services",
		Aliases:           []string{"remote-service"},
		RunE:              actionGetRemoteServices,
		ValidArgsFunction: remoteServicesAutoCompleteFun,
	}
	return []*cobra.Command{getServices}
}

func actionGetRemoteServices(cmd *cobra.Command, args []string) error {
	// Get the project and url parameters
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}
	commandName := "remote-service"

	params := map[string]string{}
	if len(args) != 0 {
		params["id"] = args[0]
	}

	objs, err := GetRemoteServices(project, commandName, params)
	if err != nil {
		return err
	}

	if err := utils.PrintYaml(objs); err != nil {
		return err
	}
	return nil
}

func actionGenerateService(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("incorrect number of arguments. Use -h to check usage instructions", nil)
	}
	dbruleConfigFile := args[0]
	dbrule, err := generateService()
	if err != nil {
		return err
	}

	return utils.AppendConfigToDisk(dbrule, dbruleConfigFile)
}

// DeleteSubCommands is the list of commands the remote-services module exposes
func DeleteSubCommands() []*cobra.Command {

	var getServices = &cobra.Command{
		Use:               "remote-services",
		Aliases:           []string{"remote-service"},
		RunE:              actionDeleteRemoteServices,
		ValidArgsFunction: remoteServicesAutoCompleteFun,
		Example:           "space-cli delete remote-services serviceID --project myproject",
	}
	return []*cobra.Command{getServices}
}

func actionDeleteRemoteServices(cmd *cobra.Command, args []string) error {
	// Get the project
	project, check := utils.GetProjectID()
	if !check {
		return utils.LogError("Project not specified in flag", nil)
	}

	prefix := ""
	if len(args) != 0 {
		prefix = args[0]
	}

	return deleteRemoteService(project, prefix)
}
