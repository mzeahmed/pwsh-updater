# Makefile for PowerShell update scripts (make i v=1.0.0)
i:
	@chmod +x ./pwsh-i.sh && ./pwsh-i.sh $(v)
	