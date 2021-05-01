#
# primitivitives for high-speed packet statistics processing
#
PROGRAM := datapath
BUILD_DIR := build

all: $(PROGRAM)

$(PROGRAM):
	@mkdir -p $(BUILD_DIR)
	@GO111MODULE=on go build -i -v -o $(BUILD_DIR)/ ./...
	@echo "Done $@"

test:
	@echo "Starting test.."
	@go clean -testcache
	@go test -v ./...

benchmark:
	@echo "Starting benchmark test.."
	@go clean -testcache
	@go test -v ./... -bench=.

clean:
	@rm -f $(BUILD_DIR)/$(PROGRAM)
	@echo "Cleaned all"

clobber:
	@rm -rf $(BUILD_DIR)
	@echo "Clobber $(BUILD_DIR)"

.PHONY: all
