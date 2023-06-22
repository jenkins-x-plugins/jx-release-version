#!/bin/sh -le

version=$(jx-release-version)
echo "version=$version" >> $GITHUB_OUTPUT
