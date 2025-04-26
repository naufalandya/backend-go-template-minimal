APP_NAME = modular_monolith
BUILD_DIR = build
MAIN_FILE = cmd/main.go

.PHONY: run build

run:
	go run $(MAIN_FILE)

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
