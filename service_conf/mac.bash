#!/bin/bash

exec 1> >(logger -s -t git-notes) 2>&1

eval "$@"