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

(defn save-dockerfile
  "Given Dockerfile data & path to save at (or defaults to Dockerfile<timestamp>),
  creates the Dockerfile on disk. Returns the saved path."
  [docker-file-data  & [ path]]
  (let [save-path (case path
                    nil (str "Dockerfile" (System/currentTimeMillis))
                    (str (clojure.string/join "" path) "/Dockerfile"))]
    (spit save-path docker-file-data)
    save-path))


(def feature-flags {"-f" :features "-p" :path "--dry-run" :dry-run})


(defn to-args-map
  "Convert the list of args into a Hash Map charting all the arguments & flags.
  Case empty?: Argument List is empty => Either no arguments provided or the arg-list is done processing.
  Case - : First arg starts with - => First arg is a flag and following it are arguments.
  Case default: First arg is a value to be set for the existing flag.
                  If a flag is not provided, they are values @ nil."
  [args]
  (loop [args-map {}
         flag nil
         args-list args]
    (cond
      (empty? args-list) args-map
      (= \- (->> args-list first first)) (recur (update
                                                  args-map
                                                  (get feature-flags (first args-list))
                                                  #(if (nil? %) [] %))
                                                (get feature-flags (first args-list))
                                                (rest args-list))

      :default (recur
                 (update args-map flag #(conj (if (nil? %) [] %) (first args-list)))
                 flag
                 (rest args-list)))))
