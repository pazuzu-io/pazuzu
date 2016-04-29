(defproject pazuzu-core "0.1.0-SNAPSHOT"
  :description "Pazuzu, the Docker Maker"
  :url "http://example.com/FIXME"
  :license {:name "MIT License"
            :url "http://www.opensource.org/licenses/mit-license.php"}
  :dependencies [[org.clojure/clojure "1.8.0"]
                 [clj-http "2.1.0"]
                 [org.clojure/data.json "0.2.6"]
                 [com.github.docker-java/docker-java "2.1.2"]
                 [org.apache.commons/commons-compress "1.8"]]
  :profiles {:dev {:plugins      [[lein-midje "3.2"]]
                   :dependencies [[midje "1.8.3"]]}})