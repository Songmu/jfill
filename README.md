jfill
=======

[![Build Status](https://travis-ci.org/Songmu/jfill.png?branch=master)][travis]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/jfill?status.svg)][godoc]

[travis]: https://travis-ci.org/Songmu/jfill
[license]: https://github.com/Songmu/jfill/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/Songmu/jfill

## Description

fill the command line from json via STDIN and exec the command

## Synopsis

```console
% echo '{"age":10}' | jfill echo {{age}}
```

## Author

[Songmu](https://github.com/Songmu)
