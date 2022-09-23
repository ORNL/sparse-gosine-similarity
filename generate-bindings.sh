#!/bin/bash

VERSION="1.18.4"
OS="$(uname -s)"
ARCH="$(uname -m)"

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
BINDINGS_DIR="$SCRIPT_DIR/gobindings"
BUILD_DIR="$SCRIPT_DIR/build"

mkdir -p $BUILD_DIR && cd $BUILD_DIR

echo "BUILD_DIR: $BUILD_DIR"
echo "BINDINGS_DIR: $BINDINGS_DIR"
echo "SCRIPT_DIR: $SCRIPT_DIR"

if [ -d "$BINDINGS_DIR" ] ; then
    echo "Go bindings directory ($BINDINGS_DIR) already present. Please remove before running this script"
    exit 1
fi

case $OS in
    "Linux")
        case $ARCH in
        "x86_64")
            ARCH=amd64
            ;;
        "aarch64")
            ARCH=arm64
            ;;
        "armv6" | "armv7l")
            ARCH=armv6l
            ;;
        "armv8")
            ARCH=arm64
            ;;
        .*386.*)
            ARCH=386
            ;;
        esac
        PLATFORM="linux-$ARCH"
    ;;
    "Darwin")
          case $ARCH in
          "x86_64")
              ARCH=amd64
              ;;
          "arm64")
              ARCH=arm64
              ;;
          esac
        PLATFORM="darwin-$ARCH"
    ;;
esac

if [ -z "$PLATFORM" ]; then
    echo "Your operating system is not supported by the script."
    exit 1
fi

PACKAGE_NAME="go$VERSION.$PLATFORM.tar.gz"

if [ ! -d "$BUILD_DIR/go" ] ; then
    echo "Downloading $PACKAGE_NAME ..."
    if hash wget 2>/dev/null; then
        wget -q https://storage.googleapis.com/golang/$PACKAGE_NAME -O "$BUILD_DIR/go.tar.gz"
    else
        curl -o "$BUILD_DIR/go.tar.gz" https://storage.googleapis.com/golang/$PACKAGE_NAME
    fi

    if [ $? -ne 0 ]; then
        echo "Download failed! Exiting."
        exit 1
    fi

    echo "Extracting go package..."
    tar -xf $BUILD_DIR/go.tar.gz
else
    echo "Using existing go installation at $BUILD_DIR/go"
fi

echo "Ensuring environment is set up correctly..."
export GOROOT=$BUILD_DIR/go
mkdir -p $BUILD_DIR/gopath
export GOPATH=$BUILD_DIR/gopath
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH

if [ ! -f "$BUILD_DIR/go.mod" ] ; then
    echo "Creating go.mod..."
    go mod init bindings
    echo -e "\nreplace github.com/ORNL/sparse-gosine-similarity => ../\n" >> go.mod  
else
    echo "Using existing go.mod at $BUILD_DIR/go.mod"
fi

if [ ! -d "$BUILD_DIR/venv" ] ; then
    echo "Creating virtual environment..."
    python3 -m venv $BUILD_DIR/venv
else
    echo "Using existing python venv at $BUILD_DIR/venv"
fi

source $BUILD_DIR/venv/bin/activate

pybindcount=$(python -m pip list | grep -i pybindgen | wc -l)
if [ "$pybindcount" -eq 0 ] ; then
    echo "Installing pip packages..."
    python -m pip install --upgrade pip
    python -m pip install pybindgen
    python -m pip install --upgrade setuptools wheel
else
    echo "Using pip and pybindgen installations present in $BUILD_DIR/venv"
fi

# Check for required binaries.
which gopy > /dev/null
if [ $? -eq 1 ]; then
    echo "Installing go packages and upgrading setuptools..."
    go get golang.org/x/tools/cmd/goimports
    go install golang.org/x/tools/cmd/goimports
    go get github.com/go-python/gopy
    go install github.com/go-python/gopy
else
    echo "Using go package installations present in $BUILD_DIR/gopath"
fi

export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:.

echo "Creating bindings..."
go get github.com/ORNL/sparse-gosine-similarity
gopy build -output=$BINDINGS_DIR -vm=python3 github.com/ORNL/sparse-gosine-similarity