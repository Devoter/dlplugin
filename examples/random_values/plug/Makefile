NAME=plug
GOFLAGS=-buildmode=c-shared
GOOS=
GOARCH=
GOARM=
CGO_ENABLED=1
CC=
LIB_EXT=

ifeq ($(GOOS), windows)
	LIB_EXT=dll
else
	LIB_EXT=so
endif

all: $(NAME).$(LIB_EXT)

clean:
	rm -f $(NAME).$(LIB_EXT) $(NAME).h

$(NAME).$(LIB_EXT): main.go random_generator.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) CGO_ENABLED=$(CGO_ENABLED) CC=$(CC) go build $(GOFLAGS) -o $(NAME).$(LIB_EXT) ./...