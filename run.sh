#!/bin/bash

export DYLD_LIBRARY_PATH=$(brew --prefix opencv)/lib:$DYLD_LIBRARY_PATH
export PKG_CONFIG_PATH=$(brew --prefix opencv)/lib/pkgconfig:$PKG_CONFIG_PATH

LOG_LEVEL=debug go run main.go