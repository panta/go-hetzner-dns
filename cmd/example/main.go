package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	hetzner_dns "github.com/panta/go-hetzner-dns"
)

func usage() {
	fmt.Printf("usage: %s SUBCOMMAND [FLAGS] [ARGS]...\n", os.Args[0])
	fmt.Println("SUBCOMMANDS:")
	fmt.Println("  list")
	fmt.Println("  add-record -zone ZONE-ID NAME TYPE VALUE")
}

func main() {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	addRecordCmd := flag.NewFlagSet("add-record", flag.ExitOnError)
	addRecordZone := addRecordCmd.String("zone", "", "zone id")
	// addRecordName := addRecordCmd.String("name", "", "record name")
	// addRecordType := addRecordCmd.String("type", "", "record type")
	// addRecordValue := addRecordCmd.String("value", "", "record value")

	if len(os.Args) < 2 {
		fmt.Println("ERROR: expected a subcommand")
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		_ = listCmd.Parse(os.Args[2:])
		cmdList(listCmd)

	case "add":
		fallthrough
	case "add-record":
		_ = addRecordCmd.Parse(os.Args[2:])
		cmdAddRecord(addRecordCmd, *addRecordZone)

	default:
		fmt.Println("ERROR: expected a subcommand")
		usage()
		os.Exit(1)
	}
}

func cmdList(flagSet *flag.FlagSet) {
	client := hetzner_dns.Client{}

	zonesResponse, err := client.GetZones(context.Background(),
		"", "", 1, 100)
	if err != nil {
		log.Fatal(err)
	}

	const zoneFmt = "%24s %20s %10s %20s %10v\n"
	header := fmt.Sprintf(zoneFmt, "ID", "Name", "Status", "Project", "# Records")
	fmt.Print(header)
	fmt.Println(repeatStr("=", len(header)))
	for _, zone := range zonesResponse.Zones {
		fmt.Printf(zoneFmt,
			zone.ID, zone.Name, zone.Status,
			zone.Project, zone.RecordsCount)
	}

	for _, zone := range zonesResponse.Zones {
		fmt.Printf("\nRecords for zone %s (ID:%v)\n", zone.Name, zone.ID)

		recordsResponse, err := client.GetRecords(context.Background(), zone.ID, 0, 0)
		if err != nil {
			log.Fatal(err)
		}

		const recordFmt = "%32s %20s %10s %32s %10v %32s\n"
		header := fmt.Sprintf(recordFmt, "ID", "Name", "Type", "Value", "TTL", "Created")
		fmt.Print(header)
		fmt.Println(repeatStr("=", len(header)))
		for _, record := range recordsResponse.Records {
			fmt.Printf(recordFmt,
				record.ID,
				record.Name, record.Type, record.Value,
				record.TTL, record.Created.String())
		}
	}
}

func cmdAddRecord(flagSet *flag.FlagSet, zoneId string) {
	client := hetzner_dns.Client{}

	args := flagSet.Args()
	if len(args) < 3 {
		log.Println(args)
		log.Println("ERROR: too few arguments to 'add-record' (expected NAME TYPE VALUE)")
		usage()
		os.Exit(1)
	}

	recordName := args[0]
	recordType := args[1]
	recordValue := args[2]
	recordResponse, err := client.CreateRecord(context.Background(), hetzner_dns.RecordRequest{
		ZoneID: zoneId,
		Type:   recordType,
		Name:   recordName,
		Value:  recordValue,
	})
	if err != nil {
		log.Printf("FAILED: %v", err)
		os.Exit(1)
	}
	fmt.Println("OK.")
	fmt.Printf("Created: %v\n", recordResponse.Record)
}

func repeatStr(str string, count int) string {
	dst := ""
	for i := 0; i < count; i++ {
		dst += str
	}
	return dst
}
