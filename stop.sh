#!/bin/sh

pid=$(cat chatbot.pid)
kill -2 $pid
