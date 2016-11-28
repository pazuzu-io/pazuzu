#!/usr/bin/env bats

@test "Check that Leiningen is installed" {
    command lein -v
}
