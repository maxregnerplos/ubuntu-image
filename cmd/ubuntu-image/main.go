package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/ubuntu/image"
	"github.com/ubuntu/image/partition"
	"github.com/ubuntu/image/types"
)

var (
	cmd = &cobra.Command{
		Use:   "mros-image",
		Short: "Create an mros image",
		RunE:  runCreate,
	}

	outputFile string
	inputDir   string
	size       int
)

func init() {
	cmd.Flags().StringVarP(&outputFile, "output", "o", "mros.img", "output image path")
	cmd.Flags().StringVarP(&inputDir, "input", "i", ".", "input directory path")
	cmd.Flags().IntVar(&size, "size", 2048, "size of the image in megabytes")
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runCreate(cmd *cobra.Command, args []string) error {
	inputDir, err := filepath.Abs(inputDir)
	if err != nil {
		return err
	}

	outputFile, err := filepath.Abs(outputFile)
	if err != nil {
		return err
	}

	fmt.Printf("Creating mros image from %s, size %d MB\n", inputDir, size)

	config := &types.Config{
		TargetSizeMB: size,

		// Add mros-specific configuration options here
		// (e.g. kernel version, partition layout, etc.)
	}

	builder, err := image.NewBuilder(config)
	if err != nil {
		return err
	}

	// Add mros-specific partitions here
	parts := []partition.Partition{
		partition.NewPartition("boot", partition.PartitionTypeEFI, 0, "512M"),
		partition.NewPartition("root", partition.PartitionTypeLinux, 1, ""),
	}

	builder.SetPartitions(parts)

	if err := builder.Build(inputDir, outputFile); err != nil {
		return err
	}

	fmt.Printf("Created mros image at %s\n", outputFile)

	return nil
}
