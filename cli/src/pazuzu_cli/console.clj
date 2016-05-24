(ns pazuzu-cli.console
  (:require
    [pazuzu-core.docker :as docker]
    [pazuzu-core.registry-client :as registry]))


(defn fetch-and-compile-features
  "CLI should be able to fetch the features & compile it into a single Dockerfile.
   Output is Dockerfile data."
  [features & any-url]
  (let [url (case any-url
              nil "http://localhost:8080"
              any-url)]
    ((comp docker/generate-dockerfile registry/get-sorted-features) url features)))

(defn save-docker-file
  "Given Dockerfile data & path to save at
  (or defaults to Dockerfile<timestamp>),
  creates the Dockerfile on disk.
  Returns the saved path."
  [docker-file-data  & [ path]]
  (let [save-path (case path
                    nil (str "Dockerfile" (System/currentTimeMillis))
                    (str (clojure.string/join "" path) "/Dockerfile"))]
    (spit save-path docker-file-data)
    save-path))


(def feature-flags {"-f" :feature "-p" :path})


(defn to-args-map
  "Convert the list of args into an "
  [args]
  (loop [args-map {}
         flag nil
         args-list args]
    (cond
      (empty? args-list) args-map
      (nil? flag) (recur
                    args-map
                    (get feature-flags (first args-list))
                    (rest args-list))

      ; If flag, update args-map to have k: flag & v: []
      (= \- (->> args-list first first)) (recur (update args-map (get feature-flags (first args-list)) (fn [_] []))
                                                (get feature-flags (first args-list))
                                                (rest args-list))


      :default (recur
                 (update args-map flag #(conj % (first args-list)))
                 flag
                 (rest args-list)))))
