package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

var (
	proxyHost = "127.0.0.1"
	proxyPort = "10809"
	proxyType = "http"
)

func getProxyURL() string {
	return fmt.Sprintf("%s://%s:%s", proxyType, proxyHost, proxyPort)
}

func getIDESettingsPath() string {
	appData := os.Getenv("APPDATA")
	return filepath.Join(appData, "Antigravity IDE", "User", "settings.json")
}

func checkIDEStatus() (bool, string) {
	path := getIDESettingsPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return false, ""
	}
	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return false, ""
	}
	val, ok := settings["http.proxy"]
	if ok && val != nil {
		strVal := fmt.Sprintf("%v", val)
		if strVal != "" {
			return true, strVal
		}
	}
	return false, ""
}

func checkV2Status() (bool, string) {
	out, err := exec.Command("reg", "query", `HKCU\Environment`, "/v", "HTTP_PROXY").Output()
	if err != nil {
		return false, ""
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "HTTP_PROXY") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				return true, parts[2]
			}
		}
	}
	return false, ""
}

func setIDEProxy(proxyURL string) error {
	path := getIDESettingsPath()
	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, 0755)

	settings := make(map[string]interface{})
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &settings)
	}

	settings["http.proxy"] = proxyURL
	settings["http.proxySupport"] = "override"
	settings["http.proxyStrictSSL"] = false

	newData, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, newData, 0644)
}

func resetIDEProxy() error {
	path := getIDESettingsPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	settings := make(map[string]interface{})
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil
	}

	delete(settings, "http.proxy")
	delete(settings, "http.proxySupport")
	delete(settings, "http.proxyStrictSSL")

	newData, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, newData, 0644)
}

func setV2Proxy(proxyURL string) error {
	_ = exec.Command("setx", "HTTP_PROXY", proxyURL).Run()
	_ = exec.Command("setx", "HTTPS_PROXY", proxyURL).Run()
	return nil
}

func resetV2Proxy() error {
	_ = exec.Command("reg", "delete", `HKCU\Environment`, "/F", "/V", "HTTP_PROXY").Run()
	_ = exec.Command("reg", "delete", `HKCU\Environment`, "/F", "/V", "HTTPS_PROXY").Run()
	return nil
}

func testConnection(proxyURL string) {
	fmt.Printf("\n[*] Testing connection to https://www.google.com via %s ...\n", proxyURL)
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		fmt.Printf("%s[FAILED] Invalid proxy URL: %v%s\n", colorRed, err, colorReset)
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(parsedURL),
		},
		Timeout: 7 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get("https://www.google.com")
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("%s[FAILED] Connection failed: %v%s\n", colorRed, err, colorReset)
		return
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	fmt.Printf("%s[SUCCESS] Connected!%s Status: %d OK | Latency: %s%d ms%s\n",
		colorGreen, colorReset, resp.StatusCode, colorCyan, duration.Milliseconds(), colorReset)
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		currentProxy := getProxyURL()
		ideOn, ideVal := checkIDEStatus()
		v2On, v2Val := checkV2Status()

		fmt.Println("======================================================================")
		fmt.Printf("           %sAntigravity Proxy Manager v2.0 (EXE)%s\n", colorBold, colorReset)
		fmt.Println("======================================================================")
		fmt.Println()
		fmt.Println("  [ CURRENT STATUS ]")
		if ideOn {
			fmt.Printf("  * Antigravity IDE Proxy : %s[ENABLED]%s  (%s)\n", colorGreen, colorReset, ideVal)
		} else {
			fmt.Printf("  * Antigravity IDE Proxy : %s[DISABLED]%s\n", colorRed, colorReset)
		}

		if v2On {
			fmt.Printf("  * Antigravity 2.0 Proxy : %s[ENABLED]%s  (%s)\n", colorGreen, colorReset, v2Val)
		} else {
			fmt.Printf("  * Antigravity 2.0 Proxy : %s[DISABLED]%s\n", colorRed, colorReset)
		}

		fmt.Println()
		fmt.Println("  [ CONFIGURATION ]")
		fmt.Printf("  * Target Proxy URL      : %s%s%s\n", colorYellow, currentProxy, colorReset)
		fmt.Println()
		fmt.Println("----------------------------------------------------------------------")
		fmt.Println("  [1] Set Proxy for Antigravity IDE")
		fmt.Println("  [2] Set Proxy for Antigravity 2.0")
		fmt.Println("  [3] Set Proxy for BOTH (IDE + 2.0)")
		fmt.Println("  [4] Reset / Disable ALL Proxies")
		fmt.Println("  --------------------------------------------------------------------")
		fmt.Println("  [5] Test Proxy Latency and Connectivity")
		fmt.Println("  [6] Change Proxy Settings (IP / Port / Protocol)")
		fmt.Println("  [0] Exit")
		fmt.Println("======================================================================")
		fmt.Print("Enter your choice [0-6]: ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			_ = setIDEProxy(currentProxy)
			fmt.Printf("\n%s[SUCCESS] Proxy applied to Antigravity IDE!%s\n", colorGreen, colorReset)
			time.Sleep(1500 * time.Millisecond)
		case "2":
			_ = setV2Proxy(currentProxy)
			fmt.Printf("\n%s[SUCCESS] Proxy applied to Antigravity 2.0!%s\n", colorGreen, colorReset)
			time.Sleep(1500 * time.Millisecond)
		case "3":
			_ = setIDEProxy(currentProxy)
			_ = setV2Proxy(currentProxy)
			fmt.Printf("\n%s[SUCCESS] Proxy applied to BOTH Antigravity IDE and 2.0!%s\n", colorGreen, colorReset)
			time.Sleep(1500 * time.Millisecond)
		case "4":
			_ = resetIDEProxy()
			_ = resetV2Proxy()
			fmt.Printf("\n%s[SUCCESS] All proxies disabled and reset!%s\n", colorYellow, colorReset)
			time.Sleep(1500 * time.Millisecond)
		case "5":
			testConnection(currentProxy)
			fmt.Println("\nPress Enter to return to menu...")
			_, _ = reader.ReadString('\n')
		case "6":
			fmt.Print("\nEnter IP/Host (press Enter for " + proxyHost + "): ")
			h, _ := reader.ReadString('\n')
			h = strings.TrimSpace(h)
			if h != "" {
				proxyHost = h
			}

			fmt.Print("Enter Port (press Enter for " + proxyPort + "): ")
			p, _ := reader.ReadString('\n')
			p = strings.TrimSpace(p)
			if p != "" {
				proxyPort = p
			}

			fmt.Print("Enter Protocol [http / socks5] (press Enter for " + proxyType + "): ")
			pt, _ := reader.ReadString('\n')
			pt = strings.TrimSpace(strings.ToLower(pt))
			if pt == "http" || pt == "socks5" {
				proxyType = pt
			}
		case "0":
			return
		}
	}
}
