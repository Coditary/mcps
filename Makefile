SERVERS := reqpack ast-cli ipmc tempify prebyte beez ycallr

.PHONY: all tidy build clean

all: build

tidy:
	@for dir in pkg $(SERVERS); do \
		echo "==> tidy $$dir"; \
		(cd $$dir && go mod tidy); \
	done

build: tidy
	@mkdir -p bin
	@for dir in $(SERVERS); do \
		echo "==> build mcp-$$dir"; \
		(cd $$dir && go build -o ../bin/mcp-$$dir .); \
	done

clean:
	rm -rf bin
