(ns pazuzu-cli.core
  (:gen-class)
  (:require [clojure.data.json :as json]))



(defn -main
  "I don't do a whole lot ... yet."
  [& args]
  ;(println "Hello, World!")
  ((comp println json/write-str) {:a 1 :b 2}))
