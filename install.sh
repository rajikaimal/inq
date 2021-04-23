#!/bin/sh

echo "Installing inq ..."
go build inq.go
echo "Copying binary ..."
cp inq /usr/local/bin
echo "Intalled inq, run inq -h to view commands"
