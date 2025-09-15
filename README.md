# PowerShell Installer for Linux

This script is designed to simplify the process of installing PowerShell on Linux systems. It will download the latest version of PowerShell from the official GitHub repository, extract the archive, move the files to the correct location, and configure the system to run the pwsh command.

The script is designed to be run as a sudo user, as it requires root access to complete the installation. Once the script is complete, you will be able to run PowerShell by simply typing `pwsh` into your terminal.

Please note that this script is provided as-is, and you should use caution when running scripts from unknown sources. Be sure to review the script before running it, and ensure that you trust the source.
## Usage

You can install or update PowerShell on Linux using Make or by running the script directly.

### Prerequisites
- sudo privileges
- wget and tar installed

### Using Make
```bash
make i v=7.5.1
```
- Replace `7.5.1` with the desired PowerShell version number.

### Direct script usage
```bash
chmod +x ./pwsh-i.sh
./pwsh-i.sh 7.5.1
```

### What the installer does
- Removes any existing PowerShell installation at `/opt/microsoft/powershell/7`
- Downloads the specified version from the official GitHub releases
- Extracts and copies files to `/opt/microsoft/powershell/7`
- Makes `pwsh` executable and creates/updates the `/usr/bin/pwsh` symlink
- Cleans up temporary files

### Verify installation
```bash
pwsh -v
```
