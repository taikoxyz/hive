#!/bin/sh

RED='\033[1;31m'
NC='\033[0m' # No Color

print_error() {
  local msg="$1"
  echo -e "${RED}$msg${NC}"
}

check_env() {
  local name="$1"
  local value=$(eval echo \$$name)

  if [ -z "$value" ]; then
    print_error "error: $name must be non-empty"
    exit 1
  fi
}