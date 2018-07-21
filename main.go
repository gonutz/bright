package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gonutz/w32"
)

func usage() {
	fmt.Println(`usage: bright <percent>
  sets the brightness of all monitors to the given percentage
  the percentage must be in the range of 0 to 100`)
}

func main() {
	if len(os.Args) != 2 {
		usage()
		return
	}
	brightness, err := strconv.Atoi(os.Args[1])
	if err != nil || brightness < 0 || brightness > 100 {
		usage()
		return
	}

	monitor := w32.MonitorFromWindow(
		w32.GetDesktopWindow(),
		w32.MONITOR_DEFAULTTOPRIMARY,
	)
	if monitor == 0 {
		fmt.Println("no monitor found")
		return
	}
	ok, n := w32.GetNumberOfPhysicalMonitorsFromHMONITOR(monitor)
	if !ok {
		fmt.Println("GetNumberOfPhysicalMonitorsFromHMONITOR failed")
		return
	}
	fmt.Println(n, "physical monitor(s) found")
	monitors := make([]w32.PHYSICAL_MONITOR, n)
	if !w32.GetPhysicalMonitorsFromHMONITOR(monitor, monitors) {
		fmt.Println("GetPhysicalMonitorsFromHMONITOR failed")
		return
	}
	for i, m := range monitors {
		ok, min, _, max := w32.GetMonitorBrightness(m.Monitor)
		if ok {
			value := min + w32.DWORD(float64(brightness)/100.0*float64(max-min)+0.5)
			if !w32.SetMonitorBrightness(m.Monitor, value) {
				fmt.Println("unable to set brightness for monitor", i)
			}
		} else {
			fmt.Println("unable to query brightness range for monitor", i)
		}
	}
}
