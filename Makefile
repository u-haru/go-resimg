GOOS	  = windows
GOARCH	  = amd64
TARGET    = $(NAME).exe

GO		  = go
GO_BUILD  = $(GO) build
GO_CLEAN  = $(GO) clean
LDFLAGS   = -w -s
NAME	  = go-resimg
TARGETDIR = .
ENTRY	  = ./src


.PHONY: build clean

build:
	$(GO_BUILD) -ldflags='$(LDFLAGS)' -o $(TARGETDIR)/$(TARGET) $(ENTRY)
	@echo FINISHED!

clean:
	$(GO_CLEAN)
	rm -f $(TARGET)

run: $(TARGET)
	./$(TARGET)

