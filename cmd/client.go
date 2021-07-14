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
	hwUUID := uuid.New()

	client, err := hollow.NewClient("superSecret", viper.GetString("api"), nil)
	if err != nil {
		return
	}

	jsonBios := json.RawMessage([]byte(`{
    "dell": {
      "boot_mode": "Bios",
      "cpu_min_sev_asid": 1,
      "logical_proc": "Enabled",
      "sriov_global_enable": "Enabled",
      "tpm_security": "On"
    }
  }`))

	typeMap := make(map[string]uuid.UUID)

	types, err := client.HardwareComponentType.List(ctx, nil)
	if err != nil {
		fmt.Println("failed to list hardware component types")
		log.Fatal(err)
	}

	for _, name := range []string{"CPU", "Hard Drive"} {
		exists := false

		for _, t := range types {
			if t.Name == name {
				exists = true
				typeMap[name] = t.UUID
				break
			}
		}

		if exists {
			continue
		}

		t := hollow.HardwareComponentType{
			UUID: uuid.New(),
			Name: name,
		}

		if err := client.HardwareComponentType.Create(ctx, t); err != nil {
			fmt.Printf("failed to create %s hardware component type\n", name)
			log.Fatal(err)
		}

		typeMap[name] = t.UUID
	}

	hw := hollow.Hardware{
		UUID:         hwUUID,
		FacilityCode: "TEST1",
		Attributes: []hollow.Attributes{
			{
				Namespace: "hollow.client.test.api",
				Values:    json.RawMessage([]byte(`{"plan_type": "plan_a"}`)),
			},
		},
		HardwareComponents: []hollow.HardwareComponent{
			{
				Model:                     "Xeon",
				Vendor:                    "Intel",
				Serial:                    "123456",
				Name:                      "Intel Xeon Processor",
				HardwareComponentTypeUUID: typeMap["CPU"],
				Attributes: []hollow.Attributes{
					{
						Namespace: "hollow.client.test.firmware",
						Values:    json.RawMessage([]byte(`{"firmware_version": "1.2.2"}`)),
					},
					{
						Namespace: "hollow.client.test.api",
						Values:    json.RawMessage([]byte(`{"packet_api_uuid": "123-321"}`)),
					},
				},
			},
		},
	}

	if err := client.Hardware.Create(ctx, hw); err != nil {
		fmt.Println("failed to create hardware")
		log.Fatal(err)
	}

	bc := hollow.BIOSConfig{
		HardwareUUID: hwUUID,
		ConfigValues: jsonBios,
	}

	if err := client.BIOSConfig.Create(ctx, bc); err != nil {
		fmt.Println("failed to create bios config")
		log.Fatal(err)
	}

	lbc, err := client.Hardware.GetBIOSConfigs(ctx, hwUUID)
	if err != nil {
		fmt.Println("Failed to get bios configs")
		log.Fatal(err)
	}

	fmt.Printf("Found %d BIOS Configs\n", len(lbc))
}
