package main

import (
	"fmt"
	"log"

	"github.com/kawai-network/veridium/pkg/nodeos"
)

func main() {
	fmt.Println("=== Node.js os module equivalents ===")

	// Example 1: System information
	fmt.Println("\n1. System Information:")
	fmt.Printf("Architecture: %s\n", nodeos.Arch())
	fmt.Printf("Platform: %s\n", nodeos.Platform())
	fmt.Printf("Type: %s\n", nodeos.Type())
	fmt.Printf("Hostname: %s\n", nodeos.Hostname())
	fmt.Printf("Machine: %s\n", nodeos.Machine())
	fmt.Printf("Version: %s\n", nodeos.Version())

	// Example 2: Directory paths
	fmt.Println("\n2. Directory Paths:")
	fmt.Printf("Home directory: %s\n", nodeos.Homedir())
	fmt.Printf("Temp directory: %s\n", nodeos.Tmpdir())

	// Example 3: CPU information
	fmt.Println("\n3. CPU Information:")
	cpus := nodeos.Cpus()
	fmt.Printf("Number of CPUs: %d\n", len(cpus))
	if len(cpus) > 0 {
		fmt.Printf("First CPU model: %s\n", cpus[0].Model)
		fmt.Printf("First CPU speed: %d MHz\n", cpus[0].Speed)
		fmt.Printf("CPU times - User: %d, Sys: %d, Idle: %d\n",
			cpus[0].Times.User, cpus[0].Times.Sys, cpus[0].Times.Idle)
	}

	// Example 4: Memory information (limited in Go stdlib)
	fmt.Println("\n4. Memory Information:")
	fmt.Printf("Total memory: %d bytes (0 = not available in Go stdlib)\n", nodeos.Totalmem())
	fmt.Printf("Free memory: %d bytes (0 = not available in Go stdlib)\n", nodeos.Freemem())

	// Example 5: Load average (limited in Go stdlib)
	fmt.Println("\n5. Load Average:")
	loadavg := nodeos.Loadavg()
	fmt.Printf("Load average (1, 5, 15 min): %.2f, %.2f, %.2f (0 = not available)\n",
		loadavg[0], loadavg[1], loadavg[2])

	// Example 6: User information
	fmt.Println("\n6. User Information:")
	userInfo, err := nodeos.UserInfo()
	if err != nil {
		log.Printf("Error getting user info: %v", err)
	} else {
		fmt.Printf("Username: %s\n", userInfo.Username)
		fmt.Printf("Home directory: %s\n", userInfo.Homedir)
		fmt.Printf("Shell: %s\n", userInfo.Shell)
		// UID/GID may be 0 (not available on all platforms)
		fmt.Printf("UID: %d, GID: %d (0 = not available on this platform)\n",
			userInfo.Uid, userInfo.Gid)
	}

	// Example 7: Endianness
	fmt.Println("\n7. System Properties:")
	fmt.Printf("Endianness: %s\n", nodeos.Endianness())

	// Example 8: Uptime (limited in Go stdlib)
	fmt.Println("\n8. System Uptime:")
	uptime := nodeos.Uptime()
	if uptime == 0 {
		fmt.Println("Uptime: not available (would require platform-specific code)")
	} else {
		fmt.Printf("Uptime: %d seconds\n", uptime)
	}

	// Example 9: Constants
	fmt.Println("\n9. Constants:")
	fmt.Printf("EOL (End of Line): %q\n", nodeos.EOL)
	fmt.Printf("DevNull: %s\n", nodeos.DevNull)

	// Example 10: Priority (limited support)
	fmt.Println("\n10. Process Priority:")
	currentPriority := nodeos.GetPriority()
	fmt.Printf("Current process priority: %d\n", currentPriority)

	// Try to set priority (may fail on some systems)
	err = nodeos.SetPriority(nodeos.PRIORITY_HIGH)
	if err != nil {
		fmt.Printf("Setting priority failed (expected on some systems): %v\n", err)
	} else {
		fmt.Println("Priority set successfully")
	}

	// Example 11: Available parallelism
	fmt.Println("\n11. Parallelism:")
	fmt.Printf("Available parallelism (CPU cores): %d\n", nodeos.AvailableParallelism())

	// Example 12: Network interfaces (limited in Go stdlib)
	fmt.Println("\n12. Network Interfaces:")
	interfaces := nodeos.NetworkInterfaces()
	if len(interfaces) == 0 {
		fmt.Println("Network interfaces: not available (would require platform-specific code)")
	} else {
		fmt.Printf("Found %d network interfaces\n", len(interfaces))
		for name, ifaceList := range interfaces {
			fmt.Printf("  %s: %d interfaces\n", name, len(ifaceList))
		}
	}

	fmt.Println("\n=== Examples completed ===")
	fmt.Println("\nNote: Some values show 0 or empty strings because Go's standard")
	fmt.Println("library doesn't provide all the information that Node.js does.")
	fmt.Println("Platform-specific code would be needed for full functionality.")
}
