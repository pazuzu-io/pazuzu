(ns pazuzu-core.docker)

(defn generate-dockerfile
  "generates dockerfile from list of features"
  [features]
  (let [comment-str "# Auto-generated DockerFile by Pazuzu2"
        from-str "FROM ubuntu:latest"
        cmd-str "CMD /bin/bash"
        feature->string (fn [feature]
                          (str "# " (get feature "name") "\n"
                               (get feature "docker_data")))]
    (clojure.string/join "\n\n" (concat
                                  [comment-str from-str]
                                  (map feature->string features)
                                  [cmd-str]))))


