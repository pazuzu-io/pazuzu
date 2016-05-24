(ns pazuzu-cli.core-test
  (:require [midje.sweet :refer :all]
            [pazuzu-cli.core :refer :all]
            [pazuzu-cli.console :as console]))


(fact "main fcn triggers only for build"
      (execute-command ["build" "-f" "a" "b" "-p" "123/f"]) => {:feature `("b" "a") :path ["123/f"]}
      (execute-command ["wrong-command"]) => "wrong-command not found")