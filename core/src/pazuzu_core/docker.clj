(ns pazuzu-core.docker
  (:require [clojure.java.io :as io])
  (:import (com.github.dockerjava.core DockerClientConfig DockerClientBuilder)
           (com.github.dockerjava.api DockerClient)
           (com.github.dockerjava.core.command BuildImageResultCallback)
           (java.io File ByteArrayOutputStream)
           (org.apache.commons.compress.archivers.tar TarArchiveEntry TarArchiveOutputStream)))

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


(defn docker-client [uri]
  (let [config (.. DockerClientConfig createDefaultConfigBuilder
                   (withUri uri)
                   build)]
    (.. DockerClientBuilder (getInstance config) build)))

(defn docker-info [uri]
  (.. (docker-client uri) infoCmd exec))

(defn put-archive-entry
  "Adds an entry to a tar archive."
  [^TarArchiveOutputStream tar-output-stream filename content]
  (let [entry (doto (TarArchiveEntry. filename)
                (.setSize (count content)))]
    (doto tar-output-stream
      (.putArchiveEntry entry)
      (.write (.getBytes content))
      (.closeArchiveEntry))))

(defn build-docker-tar
  "Returns a docker context tar data in-memory,
  containing only a Dockerfile with the provided contents."
  [dockerfile-str]
  (let [buffer (ByteArrayOutputStream.)]
    (with-open [os (TarArchiveOutputStream. buffer)]
      (put-archive-entry os "Dockerfile" dockerfile-str))
    (.toByteArray buffer)))

(defn docker-build [docker-uri dockerfile-str]
  (let [docker-context-istream (io/input-stream (build-docker-tar dockerfile-str))]
    (.. (docker-client docker-uri)
        (buildImageCmd docker-context-istream)
        (exec (BuildImageResultCallback.))
        awaitImageId)))
