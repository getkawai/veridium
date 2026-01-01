package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/kawai-network/veridium/pkg/obfuscator"
)

// ConfigVar represents a configuration variable with metadata
type ConfigVar struct {
	EnvKey      string // Original env key (e.g., CF_ACCOUNT_ID)
	GoName      string // Go constant/function name (e.g., CfAccountId)
	Comment     string // Optional comment
	IsNamespace bool   // Is this a KV namespace?
	Order       int    // Sort order
}

// Define all supported Cloudflare variables
var supportedVars = []ConfigVar{
	// Core Cloudflare credentials
	{EnvKey: "CF_ACCOUNT_ID", GoName: "CfAccountId", Comment: "Cloudflare Account ID", Order: 1},
	{EnvKey: "CF_API_TOKEN", GoName: "CfApiToken", Comment: "Cloudflare API Token", Order: 2},

	// Legacy namespace (deprecated)
	{EnvKey: "CF_KV_NAMESPACE_ID", GoName: "CfKvNamespaceId", Comment: "Legacy namespace (deprecated)", IsNamespace: true, Order: 10},

	// Multi-namespace architecture
	{EnvKey: "CF_KV_CONTRIBUTORS_NAMESPACE_ID", GoName: "CfKvContributorsNamespaceId", Comment: "Namespace for contributor data", IsNamespace: true, Order: 20},
	{EnvKey: "CF_KV_PROOFS_NAMESPACE_ID", GoName: "CfKvProofsNamespaceId", Comment: "Namespace for Merkle proofs", IsNamespace: true, Order: 21},
	{EnvKey: "CF_KV_SETTLEMENTS_NAMESPACE_ID", GoName: "CfKvSettlementsNamespaceId", Comment: "Namespace for settlement metadata", IsNamespace: true, Order: 22},
	{EnvKey: "CF_KV_APIKEY_NAMESPACE_ID", GoName: "CfKvApikeyNamespaceId", Comment: "Namespace for API key management", IsNamespace: true, Order: 23},
	{EnvKey: "CF_KV_AUTHZ_NAMESPACE_ID", GoName: "CfKvAuthzNamespaceId", Comment: "Namespace for authorization reverse index (address -> apikey)", IsNamespace: true, Order: 24},
	{EnvKey: "CF_KV_P2PMARKETPLACE_NAMESPACE_ID", GoName: "CfKvP2pMarketplaceNamespaceId", Comment: "Namespace for P2P marketplace data", IsNamespace: true, Order: 25},
}

// Define all supported Telegram variables
var telegramVars = []ConfigVar{
	{EnvKey: "TELEGRAM_BOT_TOKEN", GoName: "TelegramBotToken", Comment: "Telegram Bot Token", Order: 1},
	{EnvKey: "TELEGRAM_CHAT_ID", GoName: "TelegramChatId", Comment: "Telegram Chat ID", Order: 2},
}

func main() {
	envFile := ".env"

	// Parse command line args
	if len(os.Args) > 1 {
		envFile = os.Args[1]
	}

	// Read .env file
	configs, err := readEnvFile(envFile)
	if err != nil {
		log.Fatalf("failed to read %s: %v", envFile, err)
	}

	// 1. Generate Cloudflare constants
	generateCloudflare(configs)

	// 2. Generate Etherscan constants
	generateEtherscan(configs)

	// 3. Generate LLM API keys constants
	generateLLM(configs)

	// 4. Generate Treasury constants
	generateTreasury(configs)

	// 5. Generate Telegram constants
	generateTelegram(configs)
}

func generateCloudflare(configs map[string]string) {
	outputFile := "internal/constant/cloudflare.go"
	foundVars := matchConfigs(configs, supportedVars)

	if len(foundVars) == 0 {
		fmt.Printf("⚠️ No Cloudflare variables found in .env, skipping %s\n", outputFile)
		return
	}

	// Sort by order
	sort.Slice(foundVars, func(i, j int) bool {
		return foundVars[i].Order < foundVars[j].Order
	})

	// Generate Go file
	content := generateCloudflareGoFile(foundVars, configs)

	// Write output
	err := os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		log.Fatalf("failed to write %s: %v", outputFile, err)
	}

	fmt.Printf("✅ Generated %s with %d variables\n", outputFile, len(foundVars))
}

func generateEtherscan(configs map[string]string) {
	outputFile := "internal/constant/etherscan.go"
	apiKeyStr, exists := configs["ETHERSCAN_API_KEYS"]
	if !exists || apiKeyStr == "" {
		fmt.Printf("⚠️ ETHERSCAN_API_KEYS not found in .env, skipping %s\n", outputFile)
		return
	}

	keys := strings.Split(apiKeyStr, ",")
	var cleanedKeys []string
	for _, k := range keys {
		trimmed := strings.TrimSpace(k)
		if trimmed != "" {
			cleanedKeys = append(cleanedKeys, trimmed)
		}
	}

	if len(cleanedKeys) == 0 {
		fmt.Printf("⚠️ No valid Etherscan API keys found, skipping %s\n", outputFile)
		return
	}

	content := generateEtherscanGoFile(cleanedKeys)

	// Write output
	err := os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		log.Fatalf("failed to write %s: %v", outputFile, err)
	}

	fmt.Printf("✅ Generated %s with %d keys\n", outputFile, len(cleanedKeys))
}

func generateTreasury(configs map[string]string) {
	outputFile := "internal/constant/treasury.go"
	addressesStr, exists := configs["TREASURY_ADDRESSES"]
	if !exists || addressesStr == "" {
		fmt.Printf("⚠️ TREASURY_ADDRESSES not found in .env, skipping %s\n", outputFile)
		return
	}

	keys := strings.Split(addressesStr, ",")
	var cleanedKeys []string
	for _, k := range keys {
		trimmed := strings.TrimSpace(k)
		if trimmed != "" {
			if !strings.HasPrefix(trimmed, "0x") {
				trimmed = "0x" + trimmed
			}
			cleanedKeys = append(cleanedKeys, trimmed)
		}
	}

	if len(cleanedKeys) == 0 {
		fmt.Printf("⚠️ No valid Treasury addresses found, skipping %s\n", outputFile)
		return
	}

	content := generateTreasuryGoFile(cleanedKeys)

	// Write output
	err := os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		log.Fatalf("failed to write %s: %v", outputFile, err)
	}

	fmt.Printf("✅ Generated %s with %d addresses\n", outputFile, len(cleanedKeys))
}

func generateTelegram(configs map[string]string) {
	outputFile := "internal/constant/telegram.go"
	foundVars := matchConfigs(configs, telegramVars)

	if len(foundVars) == 0 {
		fmt.Printf("⚠️ No Telegram variables found in .env, skipping %s\n", outputFile)
		return
	}

	// Sort by order
	sort.Slice(foundVars, func(i, j int) bool {
		return foundVars[i].Order < foundVars[j].Order
	})

	// Generate Go file
	content := generateTelegramGoFile(foundVars, configs)

	// Write output
	err := os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		log.Fatalf("failed to write %s: %v", outputFile, err)
	}

	fmt.Printf("✅ Generated %s with %d variables\n", outputFile, len(foundVars))
}

func generateLLM(configs map[string]string) {
	outputFile := "internal/constant/llm.go"

	// Parse OpenRouter API keys
	openrouterKeys := parseCommaSeparatedKeys(configs["OPENROUTER_API_KEYS"])

	// Parse ZAI API keys
	zaiKeys := parseCommaSeparatedKeys(configs["ZAI_API_KEYS"])

	if len(openrouterKeys) == 0 && len(zaiKeys) == 0 {
		fmt.Printf("⚠️ No LLM API keys found in .env, skipping %s\n", outputFile)
		return
	}

	content := generateLLMGoFile(openrouterKeys, zaiKeys)

	// Write output
	err := os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		log.Fatalf("failed to write %s: %v", outputFile, err)
	}

	fmt.Printf("✅ Generated %s with %d OpenRouter keys and %d ZAI keys\n",
		outputFile, len(openrouterKeys), len(zaiKeys))
}

// parseCommaSeparatedKeys parses a comma-separated string of API keys
func parseCommaSeparatedKeys(apiKeyStr string) []string {
	if apiKeyStr == "" {
		return nil
	}

	keys := strings.Split(apiKeyStr, ",")
	var cleanedKeys []string
	for _, k := range keys {
		trimmed := strings.TrimSpace(k)
		if trimmed != "" {
			cleanedKeys = append(cleanedKeys, trimmed)
		}
	}

	return cleanedKeys
}

// readEnvFile reads all KEY=VALUE pairs from an env file
func readEnvFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	configs := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Remove quotes if present
			value = strings.Trim(value, `"'`)

			configs[key] = value
		}
	}

	return configs, scanner.Err()
}

// matchConfigs finds which supported vars have values in configs
func matchConfigs(configs map[string]string, vars []ConfigVar) []ConfigVar {
	var found []ConfigVar

	for _, v := range vars {
		if _, exists := configs[v.EnvKey]; exists {
			found = append(found, v)
		}
	}

	return found
}

// generateCloudflareGoFile creates the Go source file content for Cloudflare
func generateCloudflareGoFile(vars []ConfigVar, configs map[string]string) string {
	var sb strings.Builder

	// Header
	sb.WriteString("// Code generated by obfuscator-gen. DO NOT EDIT.\n")
	sb.WriteString("package constant\n\n")
	sb.WriteString("import \"github.com/kawai-network/veridium/pkg/obfuscator\"\n\n")

	// Constants
	sb.WriteString("const (\n")

	// Group: Core credentials
	hasCredentials := false
	for _, v := range vars {
		if !v.IsNamespace {
			hasCredentials = true
			break
		}
	}

	if hasCredentials {
		sb.WriteString("\t// Core Cloudflare credentials\n")
		for _, v := range vars {
			if !v.IsNamespace {
				value := configs[v.EnvKey]
				encoded := obfuscator.EncodeString(value)
				sb.WriteString(fmt.Sprintf("\tobfuscated%s = \"%s\"\n", v.GoName, encoded))
			}
		}
		sb.WriteString("\n")
	}

	// Group: Namespaces
	hasNamespaces := false
	for _, v := range vars {
		if v.IsNamespace {
			hasNamespaces = true
			break
		}
	}

	if hasNamespaces {
		sb.WriteString("\t// KV Namespaces - Multiple namespaces for data isolation\n")
		sb.WriteString("\t// - Contributors: User data, balances, heartbeat\n")
		sb.WriteString("\t// - Proofs: Merkle proofs for each settlement period\n")
		sb.WriteString("\t// - Settlements: Settlement period metadata\n")
		sb.WriteString("\t// - Apikey: API key management (apikey -> address)\n")
		sb.WriteString("\t// - Authz: Authorization reverse index (address -> apikey)\n")
		sb.WriteString("\t// - P2pMarketplace: P2P marketplace data\n")

		for _, v := range vars {
			if v.IsNamespace {
				value := configs[v.EnvKey]
				encoded := obfuscator.EncodeString(value)
				sb.WriteString(fmt.Sprintf("\tobfuscated%s = \"%s\"\n", v.GoName, encoded))
			}
		}
	}

	sb.WriteString(")\n\n")

	// Functions
	for _, v := range vars {
		sb.WriteString(fmt.Sprintf("// Get%s returns the decoded value of %s\n", v.GoName, v.EnvKey))
		if v.Comment != "" {
			sb.WriteString(fmt.Sprintf("// %s\n", v.Comment))
		}
		sb.WriteString(fmt.Sprintf("func Get%s() string {\n", v.GoName))
		sb.WriteString(fmt.Sprintf("\tval, _ := obfuscator.DecodeString(obfuscated%s)\n", v.GoName))
		sb.WriteString("\treturn val\n")
		sb.WriteString("}\n\n")
	}

	return sb.String()
}

func generateEtherscanGoFile(keys []string) string {
	var sb strings.Builder

	// Header
	sb.WriteString("// Code generated by obfuscator-gen. DO NOT EDIT.\n")
	sb.WriteString("package constant\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"math/rand\"\n")
	sb.WriteString("\t\"time\"\n\n")
	sb.WriteString("\t\"github.com/kawai-network/veridium/pkg/obfuscator\"\n")
	sb.WriteString(")\n\n")

	// Constants
	sb.WriteString("const (\n")
	for i, key := range keys {
		encoded := obfuscator.EncodeString(key)
		sb.WriteString(fmt.Sprintf("\tobfuscatedEtherscanApiKey%d = \"%s\"\n", i, encoded))
	}
	sb.WriteString(")\n\n")

	// Random picker
	sb.WriteString("// GetRandomEtherscanApiKey returns a random decoded Etherscan API key from the pool\n")
	sb.WriteString("func GetRandomEtherscanApiKey() string {\n")
	sb.WriteString("\tkeys := getEtherscanApiKeys()\n")
	sb.WriteString("\tif len(keys) == 0 {\n")
	sb.WriteString("\t\treturn \"\"\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\tr := rand.New(rand.NewSource(time.Now().UnixNano()))\n")
	sb.WriteString("\treturn keys[r.Intn(len(keys))]\n")
	sb.WriteString("}\n\n")

	// Array of all keys
	sb.WriteString("// getEtherscanApiKeys returns a slice of all decoded Etherscan API keys\n")
	sb.WriteString("func getEtherscanApiKeys() []string {\n")
	sb.WriteString("\treturn []string{\n")
	for i := range keys {
		sb.WriteString(fmt.Sprintf("\t\tgetEtherscanApiKey%d(),\n", i))
	}
	sb.WriteString("\t}\n")
	sb.WriteString("}\n\n")

	// Individual getters (private)
	for i := range keys {
		sb.WriteString(fmt.Sprintf("func getEtherscanApiKey%d() string {\n", i))
		sb.WriteString(fmt.Sprintf("\tval, _ := obfuscator.DecodeString(obfuscatedEtherscanApiKey%d)\n", i))
		sb.WriteString("\treturn val\n")
		sb.WriteString("}\n\n")
	}

	return sb.String()
}

func generateTreasuryGoFile(addresses []string) string {
	var sb strings.Builder

	// Header
	sb.WriteString("// Code generated by obfuscator-gen. DO NOT EDIT.\n")
	sb.WriteString("package constant\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"math/rand\"\n")
	sb.WriteString("\t\"time\"\n\n")
	sb.WriteString("\t\"github.com/ethereum/go-ethereum/common\"\n")
	sb.WriteString("\t\"github.com/kawai-network/veridium/pkg/obfuscator\"\n")
	sb.WriteString(")\n\n")

	// Constants
	sb.WriteString("const (\n")
	for i, address := range addresses {
		encoded := obfuscator.EncodeString(address)
		sb.WriteString(fmt.Sprintf("\tobfuscatedTreasuryAddress%d = \"%s\"\n", i, encoded))
	}
	sb.WriteString(")\n\n")

	// Random picker
	sb.WriteString("// GetRandomTreasuryAddress returns a random decoded Treasury Address from the pool\n")
	sb.WriteString("func GetRandomTreasuryAddress() string {\n")
	sb.WriteString("\taddresses := GetTreasuryAddresses()\n")
	sb.WriteString("\tif len(addresses) == 0 {\n")
	sb.WriteString("\t\treturn \"\"\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\tr := rand.New(rand.NewSource(time.Now().UnixNano()))\n")
	sb.WriteString("\taddress := addresses[r.Intn(len(addresses))]\n")
	sb.WriteString("\treturn common.HexToAddress(address).Hex()\n")
	sb.WriteString("}\n\n")

	// Array of all addresses
	sb.WriteString("// GetTreasuryAddresses returns a slice of all decoded Treasury Addresses\n")
	sb.WriteString("func GetTreasuryAddresses() []string {\n")
	sb.WriteString("\treturn []string{\n")
	for i := range addresses {
		sb.WriteString(fmt.Sprintf("\t\tgetTreasuryAddress%d(),\n", i))
	}
	sb.WriteString("\t}\n")
	sb.WriteString("}\n\n")

	// Individual getters (private)
	for i := range addresses {
		sb.WriteString(fmt.Sprintf("func getTreasuryAddress%d() string {\n", i))
		sb.WriteString(fmt.Sprintf("\tval, _ := obfuscator.DecodeString(obfuscatedTreasuryAddress%d)\n", i))
		sb.WriteString("\treturn val\n")
		sb.WriteString("}\n\n")
	}

	return sb.String()
}

func generateLLMGoFile(openrouterKeys, zaiKeys []string) string {
	var sb strings.Builder

	// Header
	sb.WriteString("// Code generated by obfuscator-gen. DO NOT EDIT.\n")
	sb.WriteString("package constant\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"math/rand\"\n")
	sb.WriteString("\t\"time\"\n\n")
	sb.WriteString("\t\"github.com/kawai-network/veridium/pkg/obfuscator\"\n")
	sb.WriteString(")\n\n")

	// Constants
	sb.WriteString("const (\n")

	// OpenRouter keys
	if len(openrouterKeys) > 0 {
		sb.WriteString("\t// OpenRouter API Keys\n")
		for i, key := range openrouterKeys {
			encoded := obfuscator.EncodeString(key)
			sb.WriteString(fmt.Sprintf("\tobfuscatedOpenRouterApiKey%d = \"%s\"\n", i, encoded))
		}
		sb.WriteString("\n")
	}

	// ZAI keys
	if len(zaiKeys) > 0 {
		sb.WriteString("\t// ZAI API Keys\n")
		for i, key := range zaiKeys {
			encoded := obfuscator.EncodeString(key)
			sb.WriteString(fmt.Sprintf("\tobfuscatedZaiApiKey%d = \"%s\"\n", i, encoded))
		}
	}

	sb.WriteString(")\n\n")

	// OpenRouter functions
	if len(openrouterKeys) > 0 {
		// Random picker
		sb.WriteString("// GetRandomOpenRouterApiKey returns a random decoded OpenRouter API key from the pool\n")
		sb.WriteString("func GetRandomOpenRouterApiKey() string {\n")
		sb.WriteString("\tkeys := getOpenRouterApiKeys()\n")
		sb.WriteString("\tif len(keys) == 0 {\n")
		sb.WriteString("\t\treturn \"\"\n")
		sb.WriteString("\t}\n")
		sb.WriteString("\tr := rand.New(rand.NewSource(time.Now().UnixNano()))\n")
		sb.WriteString("\treturn keys[r.Intn(len(keys))]\n")
		sb.WriteString("}\n\n")

		// Array of all keys
		sb.WriteString("// getOpenRouterApiKeys returns a slice of all decoded OpenRouter API keys\n")
		sb.WriteString("func getOpenRouterApiKeys() []string {\n")
		sb.WriteString("\treturn []string{\n")
		for i := range openrouterKeys {
			sb.WriteString(fmt.Sprintf("\t\tgetOpenRouterApiKey%d(),\n", i))
		}
		sb.WriteString("\t}\n")
		sb.WriteString("}\n\n")

		// Individual getters (private)
		for i := range openrouterKeys {
			sb.WriteString(fmt.Sprintf("func getOpenRouterApiKey%d() string {\n", i))
			sb.WriteString(fmt.Sprintf("\tval, _ := obfuscator.DecodeString(obfuscatedOpenRouterApiKey%d)\n", i))
			sb.WriteString("\treturn val\n")
			sb.WriteString("}\n\n")
		}
	}

	// ZAI functions
	if len(zaiKeys) > 0 {
		// Random picker
		sb.WriteString("// GetRandomZaiApiKey returns a random decoded ZAI API key from the pool\n")
		sb.WriteString("func GetRandomZaiApiKey() string {\n")
		sb.WriteString("\tkeys := getZaiApiKeys()\n")
		sb.WriteString("\tif len(keys) == 0 {\n")
		sb.WriteString("\t\treturn \"\"\n")
		sb.WriteString("\t}\n")
		sb.WriteString("\tr := rand.New(rand.NewSource(time.Now().UnixNano()))\n")
		sb.WriteString("\treturn keys[r.Intn(len(keys))]\n")
		sb.WriteString("}\n\n")

		// Array of all keys
		sb.WriteString("// getZaiApiKeys returns a slice of all decoded ZAI API keys\n")
		sb.WriteString("func getZaiApiKeys() []string {\n")
		sb.WriteString("\treturn []string{\n")
		for i := range zaiKeys {
			sb.WriteString(fmt.Sprintf("\t\tgetZaiApiKey%d(),\n", i))
		}
		sb.WriteString("\t}\n")
		sb.WriteString("}\n\n")

		// Individual getters (private)
		for i := range zaiKeys {
			sb.WriteString(fmt.Sprintf("func getZaiApiKey%d() string {\n", i))
			sb.WriteString(fmt.Sprintf("\tval, _ := obfuscator.DecodeString(obfuscatedZaiApiKey%d)\n", i))
			sb.WriteString("\treturn val\n")
			sb.WriteString("}\n\n")
		}
	}

	return sb.String()
}

// generateTelegramGoFile creates the Go source file content for Telegram
func generateTelegramGoFile(vars []ConfigVar, configs map[string]string) string {
	var sb strings.Builder

	// Header
	sb.WriteString("// Code generated by obfuscator-gen. DO NOT EDIT.\n")
	sb.WriteString("package constant\n\n")
	sb.WriteString("import \"github.com/kawai-network/veridium/pkg/obfuscator\"\n\n")

	// Constants
	sb.WriteString("const (\n")

	for _, v := range vars {
		value := configs[v.EnvKey]
		encoded := obfuscator.EncodeString(value)
		sb.WriteString(fmt.Sprintf("\tobfuscated%s = \"%s\"\n", v.GoName, encoded))
	}

	sb.WriteString(")\n\n")

	// Functions
	for _, v := range vars {
		sb.WriteString(fmt.Sprintf("// Get%s returns the decoded value of %s\n", v.GoName, v.EnvKey))
		if v.Comment != "" {
			sb.WriteString(fmt.Sprintf("// %s\n", v.Comment))
		}
		sb.WriteString(fmt.Sprintf("func Get%s() string {\n", v.GoName))
		sb.WriteString(fmt.Sprintf("\tval, _ := obfuscator.DecodeString(obfuscated%s)\n", v.GoName))
		sb.WriteString("\treturn val\n")
		sb.WriteString("}\n\n")
	}

	return sb.String()
}
