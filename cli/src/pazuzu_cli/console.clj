(ns pazuzu-cli.console
  (:require [pazuzu-core.docker :as docker]
                                [pazuzu-core.registry-client :as registry]
                                [midje.sweet :refer :all]
                                [clojure.data.json :as json]))

(defn fetch-and-compile-features
  "CLI should be able to fetch the features & compile it into a single Dockerfile.
   Output is Dockerfile data."
  [features & any-url]
  (let [url (case any-url
              nil "http://localhost:8080"
              any-url)]
    ((comp docker/generate-dockerfile registry/get-sorted-features) url features)))

;TODO: Figure out how to temporarily save and check for saved docker file.
(defn save-docker-file
  [docker-file-data  & path]
  (let [save-path (case path
                    nil (str "Dockerfile-" (System/currentTimeMillis))
                    path)]
    (spit save-path docker-file-data)
    save-path))