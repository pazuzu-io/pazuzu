(ns pazuzu-core.registry-client
  (:require
    [clj-http.client :as client]
    [clojure.data.json :as json]
    ))

(def features-resource "/api/features")

(defn get-features
  "return an array features based on server-url"
  [server-url query-params]
  (let [endpoint (str server-url features-resource)
        {body :body} (client/get endpoint {:query-params query-params})]
    (json/read-str body)))

(defn get-sorted-features
  "returns an array of topollogicaly sorted features by the list of names"
  [server-url names]
  (let [names-str (clojure.string/join "," names)
        query-params {"sorted" 1 "name" names-str}]
    (get-features server-url query-params)))
