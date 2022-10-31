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

func printMaybe(name string, value string) {
    if value != "" {
    	fmt.Printf("%s=%s,", name, value)
    }
}

func printBattery(things map[string]string) {
	idx := strings.Replace(things["BAT_NAME"], "BAT", "", 1)

	// tags
	fmt.Printf("battery_linux,battery=%s,idx=%s ", things["BAT_NAME"], idx)

	// values
	printMaybe("present", things["POWER_SUPPLY_PRESENT"])
	printMaybe("cycle_count", things["POWER_SUPPLY_CYCLE_COUNT"])
	printMaybe("voltage_min_design", things["POWER_SUPPLY_VOLTAGE_MIN_DESIGN"])
	printMaybe("voltage_now", things["POWER_SUPPLY_VOLTAGE_NOW"])
	printMaybe("power_now", things["POWER_SUPPLY_POWER_NOW"])
	printMaybe("energy_full_design", things["POWER_SUPPLY_ENERGY_FULL_DESIGN"])
	printMaybe("energy_full", things["POWER_SUPPLY_ENERGY_FULL"])
	printMaybe("energy_now", things["POWER_SUPPLY_ENERGY_NOW"])
	printMaybe("capacity", things["POWER_SUPPLY_CAPACITY"])
	printMaybe("battery_status", things["POWER_SUPPLY_STATUS"])
	printMaybe("capacity_level", things["POWER_SUPPLY_CAPACITY_LEVEL"])
	fmt.Printf("\n")
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
