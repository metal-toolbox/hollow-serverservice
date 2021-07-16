package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

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

var alpha = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = alpha[rand.Intn(len(alpha))]
	}

	return string(b)
}

func client(ctx context.Context) {
	rand.Seed(time.Now().UnixNano())

	hwUUID := uuid.New()
	hwPlan := fmt.Sprintf("plan_%s", randSeq(6)) //nolint

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

	newJSONBios := json.RawMessage([]byte(`{
    "dell": {
      "boot_mode": "UEFI",
      "cpu_min_sev_asid": 1,
      "logical_proc": "Disabled",
      "sriov_global_enable": "Enabled",
      "tpm_security": "Off"
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

		if _, err := client.HardwareComponentType.Create(ctx, t); err != nil {
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
				Values:    json.RawMessage([]byte(fmt.Sprintf(`{"plan_type": "%s"}`, hwPlan))),
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

	if _, err := client.Hardware.Create(ctx, hw); err != nil {
		fmt.Println("failed to create hardware")
		log.Fatal(err)
	}

	bc := hollow.VersionedAttributes{
		Namespace: "net.equinixplatform.bios",
		Values:    jsonBios,
	}

	if _, err := client.Hardware.CreateVersionedAttributes(ctx, hwUUID, bc); err != nil {
		fmt.Println("failed to create bios config")
		log.Fatal(err)
	}

	bcNew := hollow.VersionedAttributes{
		Namespace: "net.equinixplatform.bios",
		Values:    newJSONBios,
	}

	if _, err := client.Hardware.CreateVersionedAttributes(ctx, hwUUID, bcNew); err != nil {
		fmt.Println("failed to create bios config")
		log.Fatal(err)
	}

	bcNew2 := hollow.VersionedAttributes{
		Namespace: "net.hollow.test",
		Values:    jsonBios,
	}

	if _, err := client.Hardware.CreateVersionedAttributes(ctx, hwUUID, bcNew2); err != nil {
		fmt.Println("failed to create bios config")
		log.Fatal(err)
	}

	lbc, err := client.Hardware.GetVersionedAttributes(ctx, hwUUID)
	if err != nil {
		fmt.Println("Failed to get bios configs")
		log.Fatal(err)
	}

	fmt.Printf("Found %d BIOS Configs\n", len(lbc))

	lhw, err := client.Hardware.List(ctx, &hollow.HardwareListParams{
		FacilityCode: "TEST1",
		AttributeListParams: []hollow.AttributeListParams{
			{
				Namespace:  "hollow.client.test.api",
				Keys:       []string{"plan_type"},
				EqualValue: hwPlan,
			},
		},
	})
	if err != nil {
		fmt.Println("Failed to list hardware")
		log.Fatal(err)
	}

	fmt.Printf("Found %d pieces of hardware filtered by plan type (%s)\n", len(lhw), hwPlan)

	lhw, err = client.Hardware.List(ctx, &hollow.HardwareListParams{
		FacilityCode: "TEST1",
		AttributeListParams: []hollow.AttributeListParams{
			{
				Namespace:  "hollow.client.test.api",
				Keys:       []string{"plan_type"},
				EqualValue: hwPlan,
			},
		},
		VersionedAttributeListParams: []hollow.AttributeListParams{
			{
				Namespace:  "net.equinixplatform.bios",
				Keys:       []string{"dell", "tpm_security"},
				EqualValue: "On",
			},
		},
	})
	if err != nil {
		fmt.Println("Failed to list hardware")
		log.Fatal(err)
	}

	fmt.Printf("Found %d pieces of hardware filtered by plan type (%s) AND bios: tpm_security On\n", len(lhw), hwPlan)

	lhw, err = client.Hardware.List(ctx, &hollow.HardwareListParams{
		FacilityCode: "TEST1",
		AttributeListParams: []hollow.AttributeListParams{
			{
				Namespace:  "hollow.client.test.api",
				Keys:       []string{"plan_type"},
				EqualValue: hwPlan,
			},
		},
		VersionedAttributeListParams: []hollow.AttributeListParams{
			{
				Namespace:  "net.equinixplatform.bios",
				Keys:       []string{"dell", "tpm_security"},
				EqualValue: "Off",
			},
		},
	})
	if err != nil {
		fmt.Println("Failed to list hardware")
		log.Fatal(err)
	}

	fmt.Printf("Found %d pieces of hardware filtered by plan type (%s) AND bios: tpm_security Off\n", len(lhw), hwPlan)

	lhw, err = client.Hardware.List(ctx, nil)
	if err != nil {
		fmt.Println("Failed to list hardware")
		log.Fatal(err)
	}

	fmt.Printf("Found %d pieces of hardware in total\n", len(lhw))
}
