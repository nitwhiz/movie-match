#!/bin/sh

echo 'y' | pnpm install --frozen-lockfile && exec pnpm dev
