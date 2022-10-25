# From https://github.com/TheNetAdmin/Makefile-Templates

# tool macros
CC ?= gcc
CFLAGS := 

# path macros
BIN_PATH := bin
SRC_PATH := src

# compile macros
TARGET_NAME := jotacoin
TARGET := $(BIN_PATH)/$(TARGET_NAME)

# src files
SRC := 	src/main.c

# clean files list
CLEAN_LIST := $(TARGET)

# default rule
default: makedir all

# non-phony targets
$(TARGET): $(SRC)
	$(CC) -o $@ $(SRC) $(CFLAGS)

# phony rules
.PHONY: makedir
makedir:
	@mkdir -p $(BIN_PATH)

.PHONY: all
all: $(TARGET)

.PHONY: clean
clean:
	@echo CLEAN $(CLEAN_LIST)
	@rm -f $(CLEAN_LIST)
