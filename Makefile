APP_NAME=cascading-pemda-service
# DEFAULT TARGET
.PHONY: all
all: build

# Build binary
.PHONY: build
build:
	@echo ">>> Building $(APP_NAME)..."
	@go build -o $(APP_NAME) .
	@echo ">>> SUCCESS..."

# Run with env
.PHONY: run
run: build myenv
	@echo ">>> Running $(APP_NAME)..."
	./$(APP_NAME)

# check env
.PHONY: myenv
myenv:
	@echo "REQUIRED ENV"
	@echo "PERENCANAAN_DB_URL: $(PERENCANAAN_DB_URL)"

# clean
clean:
	@echo "CLEANING UP"
	rm $(APP_NAME)
