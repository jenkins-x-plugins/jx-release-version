#!/bin/sh -le

version=$(jx-release-version)
echo "version=$version" >> $GITHUB_OUTPUT

previous_version=$(jx-release-version --print-previous-version)
echo "previous-version=$previous_version" >> $GITHUB_OUTPUT
