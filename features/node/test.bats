#!/usr/bin/env bats

@test "Check that Node.js is installed" {
    command node --version
}
