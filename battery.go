package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func batteryThings(battery string) (map[string]string, bool) {
	things := make(map[string]string)

	filename := fmt.Sprintf("/sys/class/power_supply/%s/uevent", battery)
	//fmt.Println(os.Stderr, filename)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s error trying to open the file: %s", filename, err)
		return nil, true
	}

	things["BAT_NAME"] = battery // Add battery name "BATx" to the things

	for _, line := range strings.Split(string(b), "\n") {
		if strings.TrimSpace(line) == "" {
			continue // if empty
		}
		l := strings.Split(line, "=")
		things[l[0]] = l[1]
	}

	return things, false
}

func appendMaybe(list *[]string, name string, value string) {
	if value != "" {
		*list = append(*list, fmt.Sprintf("%s=%s", name, value))
	}
}

func appendQuotedMaybe(list *[]string, name string, value string) {
	if value != "" {
		*list = append(*list, fmt.Sprintf("%s=\"%s\"", name, value))
	}
}

func printBattery(things map[string]string) {
	idx := strings.Replace(things["BAT_NAME"], "BAT", "", 1)

	// tags
	fmt.Printf("battery_linux,battery=%s,idx=%s ", things["BAT_NAME"], idx)

    var fields []string
    
	// values
	appendMaybe(&fields, "present", things["POWER_SUPPLY_PRESENT"])
	appendMaybe(&fields, "cycle_count", things["POWER_SUPPLY_CYCLE_COUNT"])
	appendMaybe(&fields, "voltage_min_design", things["POWER_SUPPLY_VOLTAGE_MIN_DESIGN"])
	appendMaybe(&fields, "voltage_now", things["POWER_SUPPLY_VOLTAGE_NOW"])
	appendMaybe(&fields, "power_now", things["POWER_SUPPLY_POWER_NOW"])
	appendMaybe(&fields, "energy_full_design", things["POWER_SUPPLY_ENERGY_FULL_DESIGN"])
	appendMaybe(&fields, "energy_full", things["POWER_SUPPLY_ENERGY_FULL"])
	appendMaybe(&fields, "energy_now", things["POWER_SUPPLY_ENERGY_NOW"])
	appendMaybe(&fields, "capacity", things["POWER_SUPPLY_CAPACITY"])
	appendQuotedMaybe(&fields, "battery_status", things["POWER_SUPPLY_STATUS"])
	appendQuotedMaybe(&fields, "capacity_level", things["POWER_SUPPLY_CAPACITY_LEVEL"])
	fmt.Printf("%s\n", strings.Join(fields, ","))
	return
}

func main() {
	files, _ := ioutil.ReadDir("/sys/class/power_supply/")
	for _, f := range files {
		if strings.Contains(f.Name(), "BAT") {
			stats, err := batteryThings(f.Name())
			if err {
				fmt.Fprintln(os.Stderr, "BatteryThings error")
				return
			}
			printBattery(stats)
		}
	}
}
