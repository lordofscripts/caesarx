# - Host OS Platform
ifeq ($(OS),Windows_NT) 
    detected_OS := Windows
else
    detected_OS := $(sh -c 'uname 2>/dev/null || echo Unknown')
endif

# - Go Build Environment
GO=go
GO_TAGS=-tags mlog
ifeq ($(detected_OS), Windows)
	GOFLAGS = -v -buildmode=exe -gcflags all=-N 
	EXE_EXT=.exe
else
	GOFLAGS = -v -buildmode=pie
	EXE_EXT=
endif

# - Source Project Environment
# get the Makefile's directory (GNU Make >= v3.81)
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))
# set the GO Project's BIN directory
GO_PROJ_BIN=${mkfile_dir}bin

# - Packagers only
PKG_FULL_VERSION=$(shell grep -m 1 'MANUAL_VERSION' version.go | sed -E 's/.*"([^"]+)".*/\1/')
PKG_NAME=caesarx
PKG_REVISION=1
PKG_VERSION=1.2
PKG_ARCH=amd64
PKG_FULLNAME=${PKG_NAME}_${PKG_VERSION}-${PKG_REVISION}_${PKG_ARCH}
PKG_BUILD_DIR=${HOME}/Develop/Distrib/Build/${PKG_NAME}
PKG_PPA_DIR=${HOME}/Develop/Distrib/PPA

# - Application stanza
EXEC_TABULA=tabularecta
MAIN_TABULA=cmd/tabularecta/*go
BIN_OUT_1=$(GO_PROJ_BIN)/$(EXEC_TABULA)$(EXE_EXT)
EXEC_CAESAR=caesarx
MAIN_CAESAR=cmd/caesar/*.go
BIN_OUT_2=$(GO_PROJ_BIN)/$(EXEC_CAESAR)$(EXE_EXT)
EXEC_AFFINE=affine
MAIN_AFFINE=cmd/affine/*go
BIN_OUT_3=$(GO_PROJ_BIN)/$(EXEC_AFFINE)$(EXE_EXT)
EXEC_CODEBOOK=codebook
MAIN_CODEBOOK=cmd/codebook/*go
BIN_OUT_4=$(GO_PROJ_BIN)/$(EXEC_CODEBOOK)$(EXE_EXT)
EXEC_LITE=caesar
MAIN_LITE=cmd/lite/*go
BIN_OUT_5=$(GO_PROJ_BIN)/$(EXEC_LITE)$(EXE_EXT)

# - Main Targets
.PHONY: clean build

all: tabula, caesar, affine, lite
	
allwin:
	$(GO) build $(GO_TAGS) $(GOFLAGS) -o ${BIN_OUT_2}.exe ${MAIN_CAESAR}
	$(GO) build $(GO_TAGS) $(GOFLAGS) -o ${BIN_OUT_3}.exe ${MAIN_AFFINE}
	$(GO) build $(GO_TAGS) $(GOFLAGS) -o ${BIN_OUT_1}.exe ${MAIN_TABULA}
	$(GO) build $(GO_TAGS) $(GOFLAGS) -o ${BIN_OUT_4}.exe ${MAIN_CODEBOOK}
	$(GO) build $(GO_TAGS) $(GOFLAGS) -o ${BIN_OUT_5}.exe ${MAIN_LITE}

release:
	strip --strip-unneeded ${BIN_OUT_1}
	strip --strip-unneeded ${BIN_OUT_2}
	strip --strip-unneeded ${BIN_OUT_3}
	strip --strip-unneeded ${BIN_OUT_4}
	strip --strip-unneeded ${BIN_OUT_5}

version:
	@echo $(PKG_FULL_VERSION)

# - Application Targets

tabula:
	$(GO) build $(GO_TAGS) $(GOFLAGS) -o ${BIN_OUT_1} ${MAIN_TABULA}

caesar:
	$(GO) build $(GO_TAGS) $(GOFLAGS) -o ${BIN_OUT_2} ${MAIN_CAESAR}

lite:	
	$(GO) build $(GO_TAGS) $(GOFLAGS) -o ${BIN_OUT_5} ${MAIN_LITE}

affine:
	$(GO) build $(GO_TAGS) $(GOFLAGS) -o ${BIN_OUT_3} ${MAIN_AFFINE}

codebook:
	$(GO) build $(GO_TAGS) $(GOFLAGS) -o ${BIN_OUT_4} ${MAIN_CODEBOOK}

# - Secondary Targets

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

help:
	@echo "· Application related"
	@echo  "\tall - make ALL application targets on/for Linux"
	@echo  "\tallwin - make ALL application targets on/for Windows"
	@echo  "\tversion - print the application version found in the source code"
	@echo "· GO Language targets"
	@echo  "\ttest - Run all the tests"
	@echo  "\tupdate - Update all 3rd party GO package dependencies"
	@echo  "\tlint - Run GO Lint"
	@echo "· Package Building"
	@echo  "\tdebian - Build the Debian (${PKG_FULLNAME}-${PKG_REVISION}.deb) package"
	@echo  "\trpm - Build the RPM (${PKG_FULLNAME}.rpm) package"
	@echo  "\trpmclean - Cleans the RPM build area"

# Package Builders

debian:
	GO_TAGS=
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
	gzip -n -9 -c distrib/manpages/man1/$(EXEC_AFFINE).1 > ${PKG_BUILD_DIR}/usr/share/man/man1/bellaso.1.gz
	gzip -n -9 -c distrib/manpages/man1/$(EXEC_AFFINE).1 > ${PKG_BUILD_DIR}/usr/share/man/man1/didimus.1.gz
	gzip -n -9 -c distrib/manpages/man1/$(EXEC_AFFINE).1 > ${PKG_BUILD_DIR}/usr/share/man/man1/fibonacci.1.gz
	gzip -n -9 -c distrib/manpages/man1/$(EXEC_AFFINE).1 > ${PKG_BUILD_DIR}/usr/share/man/man1/vigenere.1.gz
	strip --strip-unneeded ${BIN_OUT_1}
	cp ${BIN_OUT_1} ${PKG_BUILD_DIR}/usr/bin
	strip --strip-unneeded ${BIN_OUT_2}
	cp ${BIN_OUT_2} ${PKG_BUILD_DIR}/usr/bin
	strip --strip-unneeded ${BIN_OUT_3}
	cp ${BIN_OUT_3} ${PKG_BUILD_DIR}/usr/bin	
	strip --strip-unneeded ${BIN_OUT_4}
	cp ${BIN_OUT_4} ${PKG_BUILD_DIR}/usr/bin
	strip --strip-unneeded ${BIN_OUT_5}
	cp ${BIN_OUT_5} ${PKG_BUILD_DIR}/usr/bin
	cp distrib/DEBIAN/copyright ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
	cp docs/README.md ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
	cp LICENSE.md ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
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

rpm:
	GO_TAGS=
	mkdir -p ${PKG_BUILD_DIR}/rpmbuild/BUILD
	mkdir -p ${PKG_BUILD_DIR}/rpmbuild/RPMS/x86_64
	mkdir -p ${PKG_BUILD_DIR}/rpmbuild/SOURCES
	mkdir -p ${PKG_BUILD_DIR}/rpmbuild/SPECS
	mkdir -p ${PKG_BUILD_DIR}/rpmbuild/SRPMS
	#echo "%_topdir ${PKG_BUILD_DIR}/rpmbuild" > ~/.rpmmacros
	cp distrib/Fedora/${PKG_NAME}.spec ${PKG_BUILD_DIR}/rpmbuild/SPECS/
	sed -Ei "s/(^Version:[[:space:]]*).*/\1${PKG_FULL_VERSION}/" ${PKG_BUILD_DIR}/rpmbuild/SPECS/${PKG_NAME}.spec
	tar --exclude=".git*" --transform='s/^caesarx/caesarx-${PKG_FULL_VERSION}/' -cvzf ${PKG_BUILD_DIR}/rpmbuild/SOURCES/${PKG_NAME}-${PKG_FULL_VERSION}.tar.gz ../${PKG_NAME}/
	(cd ${PKG_BUILD_DIR}/rpmbuild/SPECS && rpmbuild --nodeps -bb ${PKG_NAME}.spec)

rpmclean:
	rm -fR ${PKG_BUILD_DIR}/rpmbuild
	@echo "Cleaned RPM build v$(PKG_FULL_VERSION)"
