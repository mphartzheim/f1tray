.PHONY: prepare build-all build-all-clean build-all-debug build-all-clean-debug

prepare:
	./builds/prepare-release.sh

build-all:
	./builds/build-all.sh

build-all-clean:
	./builds/build-all.sh --clean

build-all-debug:
	./builds/build-all.sh --debug

build-all-clean-debug:
	./builds/build-all.sh --clean --debug
