#!/bin/sh

./chatbot > stdout.log 2>&1 &
echo $! > chatbot.pid
