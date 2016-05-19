(ns pazuzu-cli.console-test
  (:require [midje.sweet :refer :all]
            [pazuzu-cli.console :refer :all]
            [clj-http.client :as client]
            [clojure.data.json :as json]
            [clojure.java.io :as io]))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; Helper fcns for testing console/fetch-and-compile-features
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(def test-url "http://localhost:8080/api/features")

; Parameters sent to client/get
(def params-map {:py-scala {"sorted" 1 "name" "Python,scala"}})

; Response map to be returned by the Pazuzu-registry
(def return-response
  {:py-scala [{:name "java" :docker_data "Java-Data" :test_instruction "jvm-test"}
              {:name "scala" :docker_data "scala" :test_instruction "scala-data"}
              {:name "Python" :docker_data "Python-Data" :test_instruction "py-test"}]})

; Generate the response for given key
(defn http-response
  [query-key]
  {:status 200
   :headers {}
   :body (json/write-str (query-key return-response))})

(def expected-docker-data
  (str "# Auto-generated DockerFile by Pazuzu2\n\n"
       "FROM ubuntu:latest\n\n# java\nJava-Data\n\n"
       "# scala\nscala\n\n# Python\nPython-Data\n\nCMD /bin/bash"))


;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; Tests for console/fetch-and-compile-features
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(fact "we have a proper docker file"
      (fetch-and-compile-features ["Python", "scala"]) => expected-docker-data
      (provided
        (client/get test-url {:query-params (:py-scala params-map )}) => (http-response :py-scala)))


;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; Helper fcns for testing console/save-docker-file.
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
(defn get-test-docker-files
  [dir-path]
  (filter #(re-find #"Dockerfile" %)
          (map #(str %) (file-seq (io/file dir-path)))))

(defn clean-up-docker-files
  [dir-path]
  (let
     [test-docker-files (get-test-docker-files dir-path)]
    (doall (map #(io/delete-file %) test-docker-files))))

(def docker-data "Hello This is docker Data")

(defn save-and-check-for-dockerfile
  [& [path]]
  (let [saved-file-name (save-docker-file docker-data path)]
    (.exists (io/as-file saved-file-name))))

(with-state-changes [(after :facts (clean-up-docker-files "."))]
                    (fact "Dockerfile is generated for given docker-data"
                          (save-and-check-for-dockerfile) => true
                          (save-and-check-for-dockerfile "./doc" => true)))
