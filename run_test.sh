#!/bin/bash

go build && ./logdel -conf-dir testdata/rules -no-delete -today 2019-10-13
