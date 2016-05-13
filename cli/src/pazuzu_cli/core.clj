(ns pazuzu-cli.core
  (:gen-class)
  (:require [pazuzu-core.docker :as docker])
  (:require [pazuzu-core.registry-client :as registry]))


(defn fetch-and-compile-features
  "CLI should be able to fetch the features & compile it into a single Dockerfile.
   Output is Dockerfile data."
  [features & any-url]
  (let [url (case any-url
                  nil "http://localhost"
                  any-url)]
    ((comp docker/generate-dockerfile registry/get-sorted-features) url features)))


(defn save-docker-file
  [docker-file-data  & path]
  (let [save-path (case path
                    nil (str "Dockerfile-" (System/currentTimeMillis))
                    path)]
    (spit save-path docker-file-data)
    save-path))


(defn -main
  "I don't do a whole lot ... yet."
  [& args]
  (println "Hello, World!"))
