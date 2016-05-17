(ns pazuzu-cli.console-test
  (:require [midje.sweet :refer :all]
            [pazuzu-cli.console :refer :all]
            [clj-http.client :as client]
            [clojure.data.json :as json]))

(def test-url "http://localhost:8080/api/features")

; Parameters sent to client/get
(def params-map {:py-scala {"sorted" 1 "name" "Python,scala"}
                 :empty ""})

; Response map to be returned by the Pazuzu-registry
(def return-response
  {:py-scala [{:name "java" :docker_data "Java-Data" :test_instruction "jvm-test"}
              {:name "scala" :docker_data "scala" :test_instruction "scala-data"}
              {:name "Python" :docker_data "Python-Data" :test_instruction "py-test"}]
  :empty    ""})

; Generate the response for given key
(defn http-response
  [query-key]
  {:status 200
   :headers {}
   :body (json/write-str (query-key return-response))})

; Make this work! Macros?
;(defn test-get [key]
;  #(client/get test-url {:query-params (key params-map)}) => (http-response key))

(fact "we have a proper docker file"
      (fetch-and-compile-features ["Python", "scala"]) => (str "# Auto-generated DockerFile by Pazuzu2\n\n"
                                                               "FROM ubuntu:latest\n\n# java\nJava-Data\n\n"
                                                               "# scala\nscala\n\n# Python\nPython-Data\n\nCMD /bin/bash")
      ;(fetch-and-compile-features []) => ""
      (provided
        (client/get test-url {:query-params (:py-scala params-map)}) => (http-response :py-scala)))
      ;(client/get test-url {:query-params (:empty params-map)}) => (http-response :empty)





