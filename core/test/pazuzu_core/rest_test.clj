(ns pazuzu-core.rest-test
  (:require [clojure.test :refer :all]
            [clj-http.client :as client]
            [pazuzu-core.rest :as rest]
            [midje.sweet :refer :all]))

(def correct-sorted-features
  [{"name" "cool-feature-0", "docker_data" "RUN echo 1", "test_instruction" "RUN echo  \"test echo\""}
   {"name" "cool-feature-1", "docker_data" "RUN echo 1", "test_instruction" "RUN echo  \"test echo\""}
   {"name" "cool-feature-2", "docker_data" "RUN echo 2", "test_instruction" "RUN echo  \"test echo\""}])

(def sorted-features-response-body
  "[{\"name\":\"cool-feature-0\",\"docker_data\":\"RUN echo 1\",\"test_instruction\":\"RUN echo  \\\"test echo\\\"\"},{\"name\":\"cool-feature-1\",\"docker_data\":\"RUN echo 1\",\"test_instruction\":\"RUN echo  \\\"test echo\\\"\"},{\"name\":\"cool-feature-2\",\"docker_data\":\"RUN echo 2\",\"test_instruction\":\"RUN echo  \\\"test echo\\\"\"}]")


(fact "get-sorted-features method should return list of sorted features"
      (rest/get-sorted-features "http://localhost:8080" ["cool-feature-2" "cool-feature-1"])
      => correct-sorted-features
      (provided
        (client/get "http://localhost:8080/api/features"
                    {:query-params {"sorted" 1 "name" "cool-feature-2,cool-feature-1"}})
        => {:status 200 :headers {} :body sorted-features-response-body}))