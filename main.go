package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

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

	var monitors []w32.HMONITOR
	cb := syscall.NewCallback(func(m w32.HMONITOR, hdc w32.HDC, r *w32.RECT, l w32.LPARAM) uintptr {
		monitors = append(monitors, m)
		return 1
	})
	w32.EnumDisplayMonitors(0, nil, cb, 0)
	for _, monitor := range monitors {
		ok, n := w32.GetNumberOfPhysicalMonitorsFromHMONITOR(monitor)
		if !ok {
			fmt.Println("GetNumberOfPhysicalMonitorsFromHMONITOR failed")
			continue
		}
		fmt.Println(n, "physical monitor(s) found")
		monitors := make([]w32.PHYSICAL_MONITOR, n)
		if !w32.GetPhysicalMonitorsFromHMONITOR(monitor, monitors) {
			fmt.Println("GetPhysicalMonitorsFromHMONITOR failed")
			continue
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
}
