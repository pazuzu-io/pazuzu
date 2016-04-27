(ns pazuzu-core.docker-test
  (:require [pazuzu-core.docker :as docker]
            [midje.sweet :refer :all]))

(def features
  [{"name" "cool-feature-0", "docker_data" "RUN echo 1", "test_instruction" "RUN echo  \"test echo\""}
   {"name" "cool-feature-2", "docker_data" "RUN echo 2", "test_instruction" "RUN echo  \"test echo\""}
   {"name" "cool-feature-1", "docker_data" "RUN echo 1", "test_instruction" "RUN echo  \"test echo\""}])

(defn long-str [& strings] (clojure.string/join "\n\n" strings))

(def dockerfile
  (long-str "# Auto-generated DockerFile by Pazuzu2"
            "FROM ubuntu:latest"
            "# cool-feature-0\nRUN echo 1"
            "# cool-feature-2\nRUN echo 2"
            "# cool-feature-1\nRUN echo 1"
            "CMD /bin/bash"))

(fact "generate-dockerfile should return the string with dockerfile based on the list of features"
      (docker/generate-dockerfile features) => dockerfile)

