(ns pazuzu-core.rest-test
  (:require [clj-http.client :as client]
            [clojure.test :refer :all]
            [pazuzu-core.rest :as rest])
  (:use clj-http.fake))

(def correct-sorted-features
  [{"name" "cool-feature-0", "docker_data" "RUN echo 1", "test_instruction" "RUN echo  \"test echo\""}
   {"name" "cool-feature-2", "docker_data" "RUN echo 2", "test_instruction" "RUN echo  \"test echo\""}
   {"name" "cool-feature-1", "docker_data" "RUN echo 1", "test_instruction" "RUN echo  \"test echo\""}])

(def sorted-features-response-body
  "[{\"name\":\"cool-feature-0\",\"docker_data\":\"RUN echo 1\",\"test_instruction\":\"RUN echo  \\\"test echo\\\"\"},{\"name\":\"cool-feature-1\",\"docker_data\":\"RUN echo 1\",\"test_instruction\":\"RUN echo  \\\"test echo\\\"\"},{\"name\":\"cool-feature-2\",\"docker_data\":\"RUN echo 2\",\"test_instruction\":\"RUN echo  \\\"test echo\\\"\"}]")

(with-fake-routes {
   "http://localhost:8080/api/features?sorted=1&name=cool-feature-2,cool-feature-1"
   (fn [request] {:status 200 :headers {} :body sorted-features-response-body})
  }
  (deftest get-sorted-features
    (is (=
          (rest/get-sorted-features "http://localhost:8080" ["cool-feature-2" "cool-feature-1"])
          correct-sorted-features))))
