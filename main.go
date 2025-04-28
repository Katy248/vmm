package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"vmm/manage"
	"vmm/shared"
	"vmm/vm"

	cli "github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "vmm",
		Usage: "Virtual Machine Manager",
		Commands: []*cli.Command{
			{
				Name:    "create",
				Usage:   "Creates and runs new virtual machine",
				Aliases: []string{"c"},
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:      "name",
						UsageText: "NAME ",
					},
					&cli.StringArg{
						Name:      "image",
						UsageText: "IMAGE",
					},
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "size",
						Aliases: []string{"s"},
						Usage:   "Size of virtual machine in GB",
						Value:   20,
					},
					&cli.IntFlag{
						Name:    "ram",
						Aliases: []string{"r"},
						Usage:   "Size of virtual machine RAM in GB",
						Value:   4,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					name := cmd.StringArg("name")
					if err := shared.ValidateName(name); err != nil {
						return err
					}
					image := cmd.StringArg("image")
					diskSize := cmd.Int("size")
					ram := cmd.Int("ram")

					machine := vm.New(name)

					if err := manage.Init(machine, diskSize); err != nil {
						return err
					}
					if err := manage.Start(machine, ram, image); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "status",
				Usage:   "Shows status of a virtual machine",
				Aliases: []string{"s"},
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:      "name",
						UsageText: "NAME",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					name := cmd.StringArg("name")
					if err := shared.ValidateName(name); err != nil {
						return err
					}
					machine := vm.New(name)

					status, err := manage.GetStatus(machine)
					if err != nil {
						return err
					}

					fmt.Printf("Status: %s\n", status.Status)

					return nil
				},
			},
			{
				Name:    "delete",
				Usage:   "Stops and deletes a virtual machine",
				Aliases: []string{"d", "rm"},
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:      "name",
						UsageText: "NAME",
					},
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "force",
						Aliases: []string{"f"},
						Usage:   "Force deletion of a virtual machine",
						Value:   false,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					name := cmd.StringArg("name")
					if err := shared.ValidateName(name); err != nil {
						return err
					}
					machine := vm.New(name)

					if err := manage.Stop(machine); err != nil && !cmd.Bool("force") {
						return err
					}
					if err := manage.Delete(machine); err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("Err: %s", err)
	}
}
