package validation

import (
	"fmt"
	"net"
	"regexp"
	"time"
)

var (
	addressValidationRegexp = regexp.MustCompile(`^(one[a-zA-Z0-9]{39})|(0x[a-fA-F0-9]{40})`)
)

// ValidateAddress validates that an address is a valid bech32 address (one...) or a valid base16 address (0x...)
func ValidateAddress(address string) error {
	matches := addressValidationRegexp.FindAllStringSubmatch(address, -1)
	if len(matches) == 0 {
		return fmt.Errorf("address supplied (%s) is not valid", address)
	}

	return nil
}

// ValidShardIDs validates senderShard and receiverShard against the shardCount
func ValidShardIDs(senderShard uint32, receiverShard uint32, shardCount uint32) error {
	if !ValidShardID(senderShard, shardCount) {
		return fmt.Errorf(`invalid argument "%d" for "--from-shard"`, senderShard)
	}

	if !ValidShardID(receiverShard, shardCount) {
		return fmt.Errorf(`invalid argument "%d" for "--to-shard"`, receiverShard)
	}

	return nil
}

// ValidShardID validates that a shardID is within the bounds of the shardCount
func ValidShardID(shardID uint32, shardCount uint32) bool {
	if shardID > (shardCount - 1) {
		return false
	}

	return true
}

// ValidateNodeConnection validates that the node can be connected to
func ValidateNodeConnection(node string) error {
	timeout := time.Duration(1 * time.Second)
	_, err := net.DialTimeout("tcp", node, timeout)
	return err
}
