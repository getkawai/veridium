package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kawai-network/x/store"
)

func main() {
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	address := "0x0439418dbAdA59FbdCd79a112F9Cf199AeB6c610"

	fmt.Println("🔍 Checking KV for:", address)
	
	rewards, err := kv.GetJobRewardsSinceLastSettlement(ctx, address, "kawai")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found: %d unsettled rewards\n", len(rewards))
	for i, r := range rewards {
		fmt.Printf("%d. %s KAWAI @ %s\n", i+1, r.ContributorAmount, r.Timestamp)
	}
}
