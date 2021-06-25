package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	hollow "go.metalkube.net/hollow/pkg/api/v1"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "test client",
	Run: func(cmd *cobra.Command, args []string) {
		client(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().String("api", "http://localhost:8000", "address to connect to hollow on")
	viperBindFlag("api", clientCmd.Flags().Lookup("api"))
}

func client(ctx context.Context) {
	uuid := uuid.New()

	client, err := hollow.NewClient("superSecret", viper.GetString("api"), nil)
	if err != nil {
		return
	}

	exampleBiosResults := `{
    "dell": {
      "boot_mode": "Bios",
      "cpu_min_sev_asid": 1,
      "logical_proc": "Enabled",
      "sriov_global_enable": "Enabled",
      "tpm_security": "On"
    }
  }`

	jsonBios, err := json.Marshal(exampleBiosResults)
	if err != nil {
		fmt.Println("failed to convert example bios to json")
		log.Fatal(err)
	}

	bc := hollow.BIOSConfig{
		HardwareUUID: uuid,
		ConfigValues: jsonBios,
	}

	if err := client.BIOSConfig.CreateBIOSConfig(ctx, bc); err != nil {
		fmt.Println("failed to create bios config")
		log.Fatal(err)
	}

	lbc, err := client.Hardware.ListBIOSConfigs(ctx, uuid)
	if err != nil {
		fmt.Println("Failed to list bios configs")
		log.Fatal(err)
	}

	fmt.Printf("Found %d BIOS Configs\n", len(lbc))
}
