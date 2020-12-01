package operations

import (
	"context"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

// Commands is the list of commands the operations module exposes
func Commands() []*cobra.Command {
	clusterNameAutoComplete := func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			utils.LogDebug("Unable to initialize docker client ", nil)
			return nil, cobra.ShellCompDirectiveDefault
		}
		connArr, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(filters.Arg("name", "space-cloud"), filters.Arg("label", "service=gateway"))})
		if err != nil {
			utils.LogDebug("Unable to list space cloud containers ", nil)
			return nil, cobra.ShellCompDirectiveDefault
		}
		accountIDs := []string{}
		for _, v := range connArr {
			arr := strings.Split(strings.Split(v.Names[0], "--")[0], "-")
			if len(arr) != 4 {
				// default gateway container
				continue
			}
			accountIDs = append(accountIDs, arr[2])
		}
		return accountIDs, cobra.ShellCompDirectiveDefault
	}

	var setup = &cobra.Command{
		Use:   "setup",
		Short: "setup development environment",
		PreRun: func(cmd *cobra.Command, args []string) {
			err := viper.BindPFlag("local-chart-dir", cmd.Flags().Lookup("local-chart-dir"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('local-chart-dir')", nil)
			}
			if err := viper.BindPFlag("file", cmd.Flags().Lookup("file")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('file')", nil)
			}
			if err := viper.BindPFlag("set", cmd.Flags().Lookup("set")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('set')", nil)
			}
		},
		RunE: actionSetup,
	}

	setup.Flags().StringP("local-chart-dir", "c", "", "Path to the space cloud helm chart directory")
	err := viper.BindEnv("local-chart-dir", "LOCAL_CHART_DIR")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('local-chart-dir') to environment variables", nil)
	}

	setup.Flags().StringP("file", "f", "", "Path to the config yaml file")
	err = viper.BindEnv("file", "FILE")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('file' to environment variables", nil)
	}

	setup.Flags().StringP("set", "", "", "Set root string values of chart in format foo1=bar1,foo2=bar2")
	err = viper.BindEnv("`set`", "SET")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('`SET`' to environment variables", nil)
	}

	var upgrade = &cobra.Command{
		Use:   "upgrade",
		Short: "upgrades the existing space cloud cluster",
		PreRun: func(cmd *cobra.Command, args []string) {
			err := viper.BindPFlag("local-chart-dir", cmd.Flags().Lookup("local-chart-dir"))
			if err != nil {
				_ = utils.LogError("Unable to bind the flag ('local-chart-dir')", nil)
			}
			if err := viper.BindPFlag("file", cmd.Flags().Lookup("file")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('file')", nil)
			}
			if err := viper.BindPFlag("set", cmd.Flags().Lookup("set")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('set')", nil)
			}
		},
		RunE: actionUpgrade,
	}

	upgrade.Flags().StringP("local-chart-dir", "c", "", "Path to the space cloud helm chart directory")
	err = viper.BindEnv("local-chart-dir", "LOCAL_CHART_DIR")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('local-chart-dir') to environment variables", nil)
	}

	upgrade.Flags().StringP("file", "f", "", "Path to the config yaml file")
	err = viper.BindEnv("file", "FILE")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('file' to environment variables", nil)
	}

	upgrade.Flags().StringP("set", "", "", "Set root string values of chart in format foo1=bar1,foo2=bar2")
	err = viper.BindEnv("`set`", "SET")
	if err != nil {
		_ = utils.LogError("Unable to bind flag ('`SET`' to environment variables", nil)
	}

	var destroy = &cobra.Command{
		Use:   "destroy",
		Short: "Remove the space cloud cluster from kubernetes",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlag("cluster-name", cmd.Flags().Lookup("cluster-name")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('cluster-name')", nil)
			}
		},
		RunE: actionDestroy,
	}

	var apply = &cobra.Command{
		Use:   "apply",
		Short: "Applies a config file or directory",
		RunE:  actionApply,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlag("delay", cmd.Flags().Lookup("delay")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('delay')", err)
			}
		},
	}
	apply.Flags().DurationP("delay", "", time.Duration(0), "Adds a delay between 2 subsequent request made by space cli to space cloud")

	var start = &cobra.Command{
		Use:   "start",
		Short: "Resumes the space-cloud docker environment",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlag("cluster-name", cmd.Flags().Lookup("cluster-name")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('cluster-name')", nil)
			}
		},
		RunE: actionStart,
	}
	start.Flags().StringP("cluster-name", "", "default", "The name of space-cloud cluster")
	if err = viper.BindEnv("cluster-name", "CLUSTER_NAME"); err != nil {
		_ = utils.LogError("Unable to bind lag ('cluster-name') to environment variables", nil)
	}

	if err := start.RegisterFlagCompletionFunc("cluster-name", clusterNameAutoComplete); err != nil {
		utils.LogDebug("Unable to provide suggetion for flag ('project')", nil)
	}

	var stop = &cobra.Command{
		Use:   "stop",
		Short: "Stops the space-cloud docker environment",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.BindPFlag("cluster-name", cmd.Flags().Lookup("cluster-name")); err != nil {
				_ = utils.LogError("Unable to bind the flag ('cluster-name')", nil)
			}
		},
		RunE: actionStop,
	}
	stop.Flags().StringP("cluster-name", "", "default", "The name of space-cloud cluster")
	if err = viper.BindEnv("cluster-name", "CLUSTER_NAME"); err != nil {
		_ = utils.LogError("Unable to bind lag ('cluster-name') to environment variables", nil)
	}

	if err := stop.RegisterFlagCompletionFunc("cluster-name", clusterNameAutoComplete); err != nil {
		utils.LogDebug("Unable to provide suggetion for flag ('project')", nil)
	}
	return []*cobra.Command{setup, upgrade, destroy, apply, start, stop}

}

func actionUpgrade(cmd *cobra.Command, args []string) error {
	chartDir := viper.GetString("local-chart-dir")
	valuesYamlFile := viper.GetString("file")
	setValue := viper.GetString("set")

	return Upgrade(setValue, valuesYamlFile, chartDir)
}

func actionSetup(cmd *cobra.Command, args []string) error {
	chartDir := viper.GetString("local-chart-dir")
	valuesYamlFile := viper.GetString("file")
	setValue := viper.GetString("set")

	return Setup(setValue, valuesYamlFile, chartDir)
}

func actionDestroy(cmd *cobra.Command, args []string) error {
	return Destroy()
}

func actionApply(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return utils.LogError("error while applying service incorrect number of arguments provided", nil)
	}
	delay := viper.GetDuration("delay")
	dirName := args[0]
	return Apply(dirName, delay)
}

func actionStart(cmd *cobra.Command, args []string) error {
	clusterName := viper.GetString("cluster-name")
	return DockerStart(clusterName)
}

func actionStop(cmd *cobra.Command, args []string) error {
	clusterName := viper.GetString("cluster-name")
	return DockerStop(clusterName)
}
