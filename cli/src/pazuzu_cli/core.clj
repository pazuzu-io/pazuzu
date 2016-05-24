(ns pazuzu-cli.core
  (:gen-class)
  (:require [pazuzu-cli.console :as console]))

(defn create-docker-features
  [& [features]]
  (println "Dockerfile created @ "
           ((comp
              console/save-docker-file
              console/fetch-and-compile-features) features)))
;
;(defn build-docker-image
;  [& features]
;  (sh "docker build " (create-docker-features features)))

(def command-map {"build" console/to-args-map})

(defn execute-command
  [args]
  (let [command (get command-map (first args) (fn [_] (str (first args) " not found")))]
    (command (rest args))))


(defn -main
  "Start of a beautiful CLI."
  [& args]
  (execute-command args))