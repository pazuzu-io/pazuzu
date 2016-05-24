(ns pazuzu-cli.core-test
  (:require [midje.sweet :refer :all]
            [pazuzu-cli.core :refer :all]
            [pazuzu-cli.console :as console]))


(fact "main fcn triggers only for build"
      (execute-command ["wrong-command"]) => "wrong-command not found")


(fact "build-dockerfile sends a correct arguments."
      (execute-command ["build" "-f" "Python" "Scala"]) => true
      (provided
        (create-dockerfile ["Python" "Scala"]) => true
        (create-dockerfile anything) => false :times 0))


(fact "build-dockerfile doesn't build on --dry-run flag."
      (execute-command ["build" "-f" "Python" "Scala" "--dry-run"]) => true
      (provided
        (build-docker-image anything) => true))
