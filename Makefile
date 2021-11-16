INSTALL_DIR = /usr/local/bin

ifdef prefix
	INSTALL_DIR = $(prefix)
endif

build:
	go build -o bin/sequence-length ./cmd/sequence-length/main.go
	go build -o bin/sequence-random ./cmd/sequence-random/main.go
	go build -o bin/sequence-shuffle ./cmd/sequence-shuffle/main.go
	go build -o bin/fastq-sample ./cmd/fastq-sample/main.go

install:
	cp bin/sequence-length $(INSTALL_DIR)/sequence-length
	cp bin/sequence-random $(INSTALL_DIR)/sequence-random
	cp bin/sequence-shuffle $(INSTALL_DIR)/sequence-shuffle
	cp bin/fastq-sample $(INSTALL_DIR)/fastq-sample

uninstall:
	rm -f $(INSTALL_DIR)/sequence-length
	rm -f $(INSTALL_DIR)/sequence-random
	rm -f $(INSTALL_DIR)/sequence-shuffle
	rm -f $(INSTALL_DIR)/fastq-sample
