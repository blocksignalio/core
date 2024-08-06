#!/bin/bash

ADDRESS=${1:-0xBB9bc244D798123fDe783fCc1C72d3Bb8C189413}

# out() { printf "$1 $2\n" "${@:3}"; }
out() { printf "%s %s\n" "$1" "$2"; }
error() { out "==> ERROR:" "$@"; } >&2
die() { error "$@"; exit 1; }

[ -z "$ETHERSCAN_APIKEY" ] && die "ETHERSCAN_APIKEY is not set."

# Needs ETHERSCAN_APIKEY.
make_url() {
    local address="$1"
    echo "https://api.etherscan.io/api?module=contract&action=getabi&address=${address}&apikey=${ETHERSCAN_APIKEY}"
}

parse() {
    local json="$1"
    # eval echo "$json" | jq
    echo "$json" | sed 's/\\"/"/g' | sed 's/^"//' | sed 's/"$//' | jq    
}

fetch() {
    local address="$1"

    local url
    url=$(make_url "$address")

    local response
    response=$(curl -s "$url")

    local result
    result=$(echo "$response" | jq .result)

    [ -z "$response" ] && die "Request failed: url='${url}'"

    # TODO: Check for "invalid apikey" or "rate limit" errors.
    # TODO: Check for "unverified source code" errors.

    parse "$result"
}

fetch "$ADDRESS" | jq '[.[] | select(.type == "event")]'

# TODO: Detect proxy patterns:
#   - OpenZeppelin's Unstructured Storage proxy pattern
#   - EIP-1967 ( https://eips.ethereum.org/EIPS/eip-1967 )
