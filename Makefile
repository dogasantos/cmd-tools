# Define the source files
SRC := uniquelines.go toupper.go tolower.go numbers.go noschars.go letters.go filterlines.go filterips.go

# Define the output directory
OUTDIR := $(HOME)/go/bin

# Define the binaries to be created
BINARIES := $(patsubst %.go,$(OUTDIR)/%,$(SRC))

# Default target to build all binaries
all: $(BINARIES)

# Rule to build each binary
$(OUTDIR)/%: %.go
	@mkdir -p $(OUTDIR)
  go mod init github.com/dogasantos/cmd-tools
  go mod tidy
	go build -o $@ $<

# Install target to copy binaries to the output directory
install: all
	@cp $(BINARIES) $(OUTDIR)

# Clean target to remove built binaries
clean:
	rm -f $(OUTDIR)/*

.PHONY: all install clean
