#!/bin/bash

set -e # Exit immediately if a command exits with a non-zero status

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Check if a version was provided
if [ -z "$1" ]; then
  echo -e "${RED}Error: Please provide the PowerShell version to install (e.g. ./pwsh-i.sh 7.5.1)${NC}"
  exit 1
fi

VERSION="$1"
ARCHIVE="powershell-$VERSION-linux-x64.tar.gz"
DOWNLOAD_URL="https://github.com/PowerShell/PowerShell/releases/download/v$VERSION/$ARCHIVE"

# Spinner
spinner() {
  local pid=$!
  local delay=0.1
  local spinstr='|/-\'
  while [ "$(ps a | awk '{print $1}' | grep $pid)" ]; do
    local temp=${spinstr#?}
    printf " [%c]  " "$spinstr"
    local spinstr=$temp${spinstr%"$temp"}
    sleep $delay
    printf "\b\b\b\b\b\b"
  done
  printf "    \b\b\b\b"
}

# Remove old installation
if [ -d /opt/microsoft/powershell/7 ]; then
  echo -e "${YELLOW}Removing previous PowerShell installation...${NC}"
  (sudo rm -rf /opt/microsoft/powershell/7) & spinner
  echo -e "${GREEN}Old PowerShell removed.${NC}"
else
  echo -e "${GREEN}No previous PowerShell installation found.${NC}"
fi

# Download
echo -e "${CYAN}Downloading PowerShell v$VERSION...${NC}"
wget "$DOWNLOAD_URL" -O powershell.tar.gz

# Extract
echo -e "${CYAN}Extracting archive...${NC}"
mkdir -p ~/powershell7 && tar -xvf powershell.tar.gz -C ~/powershell7 & spinner
echo -e "${GREEN}Extraction completed.${NC}"

# Move to /opt
echo -e "${CYAN}Moving files to /opt/microsoft/powershell/7...${NC}"
sudo mkdir -p /opt/microsoft/powershell/7 && sudo cp -r ~/powershell7/* /opt/microsoft/powershell/7/ & spinner
echo -e "${GREEN}Files moved successfully.${NC}"

# Make pwsh executable
echo -e "${CYAN}Making pwsh executable...${NC}"
sudo chmod +x /opt/microsoft/powershell/7/pwsh
echo -e "${GREEN}pwsh is now executable.${NC}"

# Create symlink
echo -e "${CYAN}Creating symbolic link to /usr/bin/pwsh...${NC}"
sudo ln -sf /opt/microsoft/powershell/7/pwsh /usr/bin/pwsh & spinner
echo -e "${GREEN}Symbolic link created.${NC}"

# Cleanup
echo -e "${CYAN}Cleaning up...${NC}"
(rm -f powershell.tar.gz "$ARCHIVE" && rm -rf ~/powershell7) & spinner
echo -e "${GREEN}Cleanup done.${NC}"

# Done
echo -e "${GREEN}PowerShell v$VERSION has been successfully installed.${NC}"
echo -e "${GREEN}You can now run PowerShell using the command: pwsh${NC}"
echo -e "${GREEN}Official documentation: https://learn.microsoft.com/en-us/powershell/${NC}"
echo -e "${GREEN}Thanks for using this PowerShell installer script.${NC}"
echo -e "${GREEN}Have a great day!${NC}"
echo -e "${GREEN}----------------------------------------${NC}"

# Version check
pwsh -v
