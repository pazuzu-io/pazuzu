(ns pazuzu-cli.core
  (:gen-class)
  (:require [pazuzu-cli.console :as console]))


(defn create-dockerfile
  [& [features]]
  (println "Dockerfile created @ "
           (->> features console/fetch-and-compile-features console/save-dockerfile)))


; TODO: STUB!! Rewrite this function once we have function that will actually build a dockerfile.
(def build-docker-image create-dockerfile)


(def does-not-contain? (complement contains?))


(defn build-dockerfile
  [args]
  (let [args-map (console/to-args-map args)
        features (:features args-map)
        to-build? (does-not-contain? args-map :dry-run)]
    (if to-build?
      (create-dockerfile features)
      (build-docker-image features))))


(def command-map {"build" build-dockerfile})


(defn execute-command
  [args]
  (let [command (get command-map
                     (first args)
                     (fn [_] (str (first args) " not found")))]
    (->> args rest command)))


(defn -main
  "Start of a beautiful CLI."
  [& args]
  (execute-command args))
