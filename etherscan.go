package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// +------------------+
// | ContractCreation |
// +------------------+

type ContractCreation struct {
	ContractAddress string `json:"contractAddress"`
	ContractCreator string `json:"contractCreator"`
	TxHash          string `json:"txHash"`
}

// +----------------+
// | ContractSource |
// +----------------+

//nolint:tagliatelle
type ContractSource struct {
	SourceCode           string `json:"SourceCode"`
	ABI                  string `json:"ABI"`
	ContractName         string `json:"ContractName"`
	CompilerVersion      string `json:"CompilerVersion"`
	OptimizationUsed     string `json:"OptimizationUsed"`
	Runs                 string `json:"Runs"`
	ConstructorArguments string `json:"ConstructorArguments"`
	EVMVersion           string `json:"EVMVersion"`
	Library              string `json:"Library"`
	LicenseType          string `json:"LicenseType"`
	Proxy                string `json:"Proxy"`
	Implementation       string `json:"Implementation"`
	SwarmSource          string `json:"SwarmSource"`
}

func (s ContractSource) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "contractSource:\n")
	fmt.Fprintf(&b, "\tSourceCode           : %s\n", s.SourceCode)
	fmt.Fprintf(&b, "\tABI                  : %s\n", s.ABI)
	fmt.Fprintf(&b, "\tContractName         : %s\n", s.ContractName)
	fmt.Fprintf(&b, "\tCompilerVersion      : %s\n", s.CompilerVersion)
	fmt.Fprintf(&b, "\tOptimizationUsed     : %s\n", s.OptimizationUsed)
	fmt.Fprintf(&b, "\tRuns                 : %s\n", s.Runs)
	fmt.Fprintf(&b, "\tConstructorArguments : %s\n", s.ConstructorArguments)
	fmt.Fprintf(&b, "\tEVMVersion           : %s\n", s.EVMVersion)
	fmt.Fprintf(&b, "\tLibrary              : %s\n", s.Library)
	fmt.Fprintf(&b, "\tLicenseType          : %s\n", s.LicenseType)
	fmt.Fprintf(&b, "\tProxy                : %s\n", s.Proxy)
	fmt.Fprintf(&b, "\tImplementation       : %s\n", s.Implementation)
	fmt.Fprintf(&b, "\tSwarmSource          : %s\n", s.SwarmSource)
	return b.String()
}

// +---------------------+
// | Constants and Types |
// +---------------------+

const (
	envEtherscanKey = "ETHERSCAN_APIKEY"

	urlContract = "https://api.etherscan.io/api" +
		"?module=contract" +
		"&action=%s" +
		"&%s=%s" +
		"&apikey=%s"
)

type etherscanResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  T      `json:"result"`
}

// +---------+
// | Private |
// +---------+

func makeURL(multipleContracts bool, action string, contracts []string) (string, error) {
	// Check action.
	if action != "getabi" &&
		action != "getsourcecode" &&
		action != "getcontractcreation" &&
		action != "checkverifystatus" {
		//
		panic("unrecognized action")
	}

	// Get API key.
	apikey := os.Getenv(envEtherscanKey)
	if apikey == "" {
		return "", fmt.Errorf("%w: %s", ErrUnsetEnvironmentVar, envEtherscanKey)
	}

	// Now bundle the URL together.
	key := "address"
	if multipleContracts {
		key = "contractaddresses"
	}
	joined := strings.Join(contracts, ",")
	return fmt.Sprintf(urlContract, action, key, joined, apikey), nil
}

func etherscanGet[T any](multipleContracts bool, action string, addresses []string, out *T) error {
	// Check addresses.
	for _, contract := range addresses {
		if !ValidateAddress(contract) {
			return makeErrorHex(ErrInvalidContractAddress, contract)
		}
	}

	// Make the GET request.
	url, err := makeURL(multipleContracts, action, addresses)
	if err != nil {
		return err
	}
	response, err := http.Get(url) //nolint:gosec,noctx
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}
	defer response.Body.Close()

	// Check if the request was unsuccessful.
	if http.StatusBadRequest <= response.StatusCode {
		return fmt.Errorf("%w: code=%d", ErrInvalidResponse, response.StatusCode)
	}

	// Read the response content.
	content, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	// Parse the JSON response.
	var body etherscanResponse[T]
	err = json.Unmarshal(content, &body)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	if body.Status != "1" || body.Message != "OK" {
		return fmt.Errorf(
			"%w: code=%d status=%s message=%s result=%v",
			ErrInvalidResponseBody,
			response.StatusCode,
			body.Status,
			body.Message,
			body.Result,
		)
	}

	*out = body.Result
	return nil
}

func etherscanGet1[T any](action, address string, out *T) error {
	return etherscanGet(false, action, []string{address}, out)
}

func etherscanGet2[T any](action string, addresses []string, out *T) error {
	return etherscanGet(true, action, addresses, out)
}

// +--------+
// | Public |
// +--------+

func GetContractABI(address string) (abi.ABI, error) {
	var (
		result string
		parsed abi.ABI
	)

	err := etherscanGet1("getabi", address, &result)
	if err != nil {
		return parsed, err
	}

	parsed, err = abi.JSON(strings.NewReader(result))
	if err != nil {
		return parsed, fmt.Errorf("read json: %w", err)
	}

	return parsed, nil
}

func GetContractEvents(address string) (map[string]abi.Event, error) {
	iface, err := GetContractABI(address)
	if err != nil {
		return nil, err
	}
	return iface.Events, nil
}

func GetContractCreation(contracts []string) ([]ContractCreation, error) {
	var xs []ContractCreation
	err := etherscanGet2("getcontractcreation", contracts, &xs)
	if err != nil {
		return nil, err
	}
	if len(xs) != len(contracts) {
		return xs, fmt.Errorf(
			"%w: want=%d have=%d",
			ErrInvalidResponseBody,
			len(contracts),
			len(xs),
		)
	}
	return xs, nil
}

func GetContractCreation1(contract string) (ContractCreation, error) {
	xs, err := GetContractCreation([]string{contract})
	if err != nil {
		var empty ContractCreation
		return empty, err
	}
	if len(xs) != 1 {
		var empty ContractCreation
		return empty, fmt.Errorf("%w: want=1 have=%d", ErrInvalidResponseBody, len(xs))
	}
	return xs[0], nil
}

func GetContractSource(address string) ([]ContractSource, error) {
	var ans []ContractSource
	err := etherscanGet1("getsourcecode", address, &ans)
	if err != nil {
		return nil, err
	}
	return ans, nil
}
