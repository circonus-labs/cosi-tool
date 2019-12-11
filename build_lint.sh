#!/usr/bin/env bash

echo
echo "Running GO Pull Request Lint (generic, bulk)"
goprcheck || exit 1
echo "Running GO Pull Request Lint (Linux specific)"
GOOS=linux goprcheck || exit 1
echo "Running GO Pull Request Lint (Freebsd specific)"
GOOS=freebsd goprcheck || exit 1
echo "Running GO Pull Request Lint (Windows specific)"
GOOS=windows goprcheck || exit 1
echo "Running GO Pull Request Lint (Solaris specific)"
GOOS=solaris goprcheck || exit 1

