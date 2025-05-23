package main

import (
	"cloudflare-internal-ip-ddns/commands"
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2" // imports as package "cli"
)

// isIPInCIDR checks if an IP address is within a CIDR range
func isIPInCIDR(ipStr, cidrStr string) (bool, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false, fmt.Errorf("invalid CIDR: %s", cidrStr)
	}

	return ipNet.Contains(ip), nil
}

func getLocalIP(subnet string) (string, error) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// Check if the address type is IP address and not loopback
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if ok, _ := isIPInCIDR(ipnet.IP.String(), subnet); ok {
					return ipnet.IP.String(), nil
				}
			}
		}
	}
	return "", fmt.Errorf("no local IP address found")
}

func ddns(
	apiToken string,
	zoneName string,
	recordName string,
	subnet string,
) error {
	ipAddress, err := getLocalIP(subnet)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("local IP address:", ipAddress)

	ctx := context.Background()

	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return err
	}
	zoneId, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return err
	}
	println("zoneId:", zoneId)

	fullRecordName := fmt.Sprintf("%s.%s", recordName, zoneName)

	rs, ret, err := api.ListDNSRecords(ctx,
		&cloudflare.ResourceContainer{
			Identifier: zoneId,
		}, cloudflare.ListDNSRecordsParams{
			Type: "A",
			Name: fullRecordName,
		})
	if err != nil {
		return err
	}
	fmt.Println("ret:", ret.Count)
	for _, r := range rs {
		fmt.Println("r:", r.Name)
	}

	if ret.Count == 0 {
		// Create a new DNS record
		_, err := api.CreateDNSRecord(ctx,
			&cloudflare.ResourceContainer{
				Identifier: zoneId,
			}, cloudflare.CreateDNSRecordParams{
				Type:    "A",
				Name:    fullRecordName,
				Content: ipAddress,
				TTL:     120,
			},
		)
		if err != nil {
			return err
		}
		println("DNS record created successfully")
	} else {
		// Update the existing DNS record
		recordId := rs[0].ID
		if rs[0].Content != ipAddress {
			_, err := api.UpdateDNSRecord(ctx,
				&cloudflare.ResourceContainer{
					Identifier: zoneId,
				}, cloudflare.UpdateDNSRecordParams{
					Type:    "A",
					ID:      recordId,
					Name:    fullRecordName,
					Content: ipAddress,
					TTL:     120,
				},
			)
			if err != nil {
				return err
			}
			println("DNS record updated successfully")
		} else {
			println("DNS record is already up-to-date")
		}
	}
	return nil
}

func main() {
	godotenv.Load()

	app := &cli.App{
		Name:    "cloudflare-internal-ip-ddns",
		Usage:   "Update Cloudflare DNS record with internal IP address",
		Version: commands.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "token",
				Aliases: []string{"t"},
				Usage:   "Cloudflare API token",
				EnvVars: []string{"CF_API_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "domain",
				Aliases: []string{"d"},
				Usage:   "Top-level domain name (Zone name), e.g. 'example.com'",
				EnvVars: []string{"CF_ZONE_NAME"},
			},
			&cli.StringFlag{
				Name:    "subdomain",
				Aliases: []string{"s"},
				Usage:   "Subdomain name (Record name), e.g. 'home' for 'home.example.com'",
				EnvVars: []string{"CF_RECORD_NAME"},
			},
			&cli.StringFlag{
				Name:        "subnet",
				Aliases:     []string{"ip"},
				Usage:       "Look for IP within a subnet",
				EnvVars:     []string{"IP_SUBNET"},
				DefaultText: "0.0.0.0/0",
			},
		},
		Action: func(cCtx *cli.Context) error {
			apiToken := cCtx.String("token")
			zoneName := cCtx.String("domain")
			recordName := cCtx.String("subdomain")
			subnet := cCtx.String("subnet")

			if apiToken == "" {
				log.Fatal("token is required")
			}
			if zoneName == "" {
				log.Fatal("top level domain is required")
			}
			if recordName == "" {
				recordName = "*"
			}

			err := ddns(apiToken, zoneName, recordName, subnet)

			return err
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
