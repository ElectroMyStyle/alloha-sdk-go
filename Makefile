# Makefile for Alloha SDK

# Test the SDK
.PHONY: test
test:
	go test -v -timeout 30s ./...

# Default command when running 'make' without specifying explicit commands
.DEFAULT_GOAL := test
