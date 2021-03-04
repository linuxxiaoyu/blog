#!/bin/bash
ps | grep -E 'account|comment|article' | awk -F ' ' '{print $1}' | xargs kill -9
cd ..
