# OS Platform
ifeq ($(OS),Windows_NT) 
    detected_OS := Windows
else
    detected_OS := $(sh -c 'uname 2>/dev/null || echo Unknown')
endif
# Go Build Environment
GO=go
ifeq ($(detected_OS), Windows)
	GOFLAGS = -v -buildmode=exe -gcflags all=-N 
	EXE_EXT=.exe
else
	GOFLAGS = -v -buildmode=pie
	EXE_EXT=
endif
# Makefile's directory (GNU Make >= v3.81)
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))
# GO Project's BIN directory
GO_PROJ_BIN=${mkfile_dir}bin

# Packagers only
PKG_NAME=caesarx
PKG_REVISION=0
PKG_VERSION=1.1-RC5
PKG_ARCH=amd64
PKG_FULLNAME=${PKG_NAME}_${PKG_VERSION}-${PKG_REVISION}_${PKG_ARCH}
PKG_BUILD_DIR=${HOME}/Develop/Distrib/Build/${PKG_NAME}
PKG_PPA_DIR=${HOME}/Develop/Distrib/PPA
# Application stanza
EXEC_TABULA=tabularecta
MAIN_TABULA=cmd/tabularecta/*go
BIN_OUT_1=$(GO_PROJ_BIN)/$(EXEC_TABULA)$(EXE_EXT)
EXEC_CAESAR=caesarx
MAIN_CAESAR=cmd/caesar/*.go
BIN_OUT_2=$(GO_PROJ_BIN)/$(EXEC_CAESAR)$(EXE_EXT)
EXEC_AFFINE=affine
MAIN_AFFINE=cmd/affine/*go
BIN_OUT_3=$(GO_PROJ_BIN)/$(EXEC_AFFINE)$(EXE_EXT)

# Main Targets
.PHONY: clean build

all: tabula, caesar
	
buildwin:
	$(GO) build -tags logx $(GOFLAGS) -o ${BIN_OUT}.exe ${MAIN}

release:
	$(GO) build $(GOFLAGS) -o ${BIN_OUT} ${MAIN}
	strip --strip-unneeded ${BIN_OUT}

# Application Targets

tabula:
	$(GO) build -tags logx $(GOFLAGS) -o ${BIN_OUT_1} ${MAIN_TABULA}

caesar:
	$(GO) build -tags logx $(GOFLAGS) -o ${BIN_OUT_2} ${MAIN_CAESAR}

affine:
	$(GO) build -tags logx $(GOFLAGS) -o ${BIN_OUT_3} ${MAIN_AFFINE}

# Secondary Targets

clean:
	go clean

run:
	go run -race  $(MAIN)

lint: 
	@gofmt -l . | grep ".*\.go"

test:
	go test tests/*test.go	

testall:
	go test ./...	

update:
	go get -u all

# Package Builders

debian:
	rm -fR ${PKG_BUILD_DIR}
	mkdir -p ${PKG_BUILD_DIR}/DEBIAN
	ln -s ${PKG_BUILD_DIR}/DEBIAN ${PKG_BUILD_DIR}/debian
	cp -R distrib/DEBIAN/* ${PKG_BUILD_DIR}/DEBIAN
	mkdir -p ${PKG_BUILD_DIR}/usr/bin
	mkdir -p ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}/assets
	mkdir -p ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}/data
	mkdir -p ${PKG_BUILD_DIR}/usr/share/man/man1
	gzip -n -9 -c distrib/manpages/man1/$(EXEC_CAESAR).1 > ${PKG_BUILD_DIR}/usr/share/man/man1/$(EXEC_CAESAR).1.gz
	gzip -n -9 -c distrib/manpages/man1/$(EXEC_AFFINE).1 > ${PKG_BUILD_DIR}/usr/share/man/man1/$(EXEC_AFFINE).1.gz
	strip --strip-unneeded ${BIN_OUT_1}
	cp ${BIN_OUT_1} ${PKG_BUILD_DIR}/usr/bin
	strip --strip-unneeded ${BIN_OUT_2}
	cp ${BIN_OUT_2} ${PKG_BUILD_DIR}/usr/bin
	strip --strip-unneeded ${BIN_OUT_3}
	cp ${BIN_OUT_3} ${PKG_BUILD_DIR}/usr/bin	
	cp distrib/DEBIAN/copyright ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
	cp docs/README.md ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
	cp docs/LICENSE.md ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
	cp docs/LANGUAGES.md ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
	cp docs/CIPHER_AFFINE.md ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
	cp docs/CIPHER_BELLASO.md ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
	cp docs/CIPHER_CAESAR.md ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
	cp docs/CIPHER_DIDIMUS.md ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
	cp docs/CIPHER_FIBONACCI.md ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
	cp docs/CIPHER_VIGENERE.md ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
	cp docs/assets/* ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}/assets
	cp docs/data/* ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}/data
	gzip -n -9 -c distrib/DEBIAN/changelog > ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}/changelog.gz
	(cd ${PKG_BUILD_DIR} && dpkg-deb --root-owner-group -b ./ ${PKG_FULLNAME}.deb)
	#(cd ${PKG_BUILD_DIR} && fakeroot /usr/bin/dpkg-buildpackage --build=binary -us -uc -b ./ ${PKG_FULLNAME})
	#@mv /tmp/${PKG_FULLNAME}.deb ${DEST_REPOSITORY}
