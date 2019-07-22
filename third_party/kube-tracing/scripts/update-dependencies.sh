#!/bin/bash

set -e

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir=$script_dir/..

cd $root_dir
copy_dependencies() {
    mkdir -p dependencies
    mv vendor dependencies/
}

copy_dependencies
