#!/bin/sh -le

version=$(jx-release-version)
echo "::set-output name=version::$version"
