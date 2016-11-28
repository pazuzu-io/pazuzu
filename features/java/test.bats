#!/usr/bin/env bats

@test "Check that Java is installed" {
    command java -version
}
