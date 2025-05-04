# Makefile for PowerShell update scripts (make pu v=1.0.0)
pu:
	@chmod +x ./pwsh-i.sh && ./pwsh-i.sh $(v)
	