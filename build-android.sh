#!/bin/bash

# Android Mobile Build Script for GoBookMarker

# Ensure CGo is enabled
export CGO_ENABLED=1

# Set Android and architecture-specific configurations
export GOOS=android
export GOARCH=386

# Path to Android NDK (you may need to adjust this)
export NDK_PATH="/path/to/android-ndk"

# Compiler configuration for CGo
export CC="${NDK_PATH}/toolchains/llvm/prebuilt/darwin-x86_64/bin/i686-linux-android21-clang"
export CXX="${NDK_PATH}/toolchains/llvm/prebuilt/darwin-x86_64/bin/i686-linux-android21-clang++"

# Build the mobile application
go build -tags android -ldflags="-s -w" ./cmd/mobile

echo "Android 386 build completed successfully!"
