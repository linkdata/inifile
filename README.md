[![build](https://github.com/linkdata/inifile/actions/workflows/go.yml/badge.svg)](https://github.com/linkdata/inifile/actions/workflows/go.yml)
[![coverage](https://coveralls.io/repos/github/linkdata/inifile/badge.svg?branch=main)](https://coveralls.io/github/linkdata/inifile?branch=main)
[![goreport](https://goreportcard.com/badge/github.com/linkdata/inifile)](https://goreportcard.com/report/github.com/linkdata/inifile)
[![Docs](https://godoc.org/github.com/linkdata/inifile?status.svg)](https://godoc.org/github.com/linkdata/inifile)

# inifile

Simple INI file reader.

Section and key names are case-insensitive and ignore leading and trailing whitespace.

Supports line comments, trailing comments and quoted values.

Keys appearing more than once will either replace previous values or append to them with a separator.
