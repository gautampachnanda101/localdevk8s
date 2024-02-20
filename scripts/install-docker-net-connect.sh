#!/bin/sh
brew install chipmk/tap/docker-mac-net-connect || true
sudo brew services start chipmk/tap/docker-mac-net-connect || true
sudo brew services