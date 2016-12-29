package main

import (
	"fmt"
	"github.com/zalando-incubator/pazuzu/storageconnector"
	"github.com/zalando-incubator/pazuzu/shared"
)

var (
	java = shared.NewFeature_str("java", "basic java feature", "Oracle", nil, "" +
		"RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys C2518248EEA14886 \n" +
		"RUN echo \"deb http://ppa.launchpad.net/webupd8team/java/ubuntu trusty main\" | tee /etc/apt/sources.list.d/webupd8team-java.list\n" +
		"RUN echo \"deb-src http://ppa.launchpad.net/webupd8team/java/ubuntu trusty main\" | tee -a /etc/apt/sources.list.d/webupd8team-java.list\n" +
		"RUN apt-get update\nRUN echo debconf shared/accepted-oracle-license-v1-1 select true | debconf-set-selections\n" +
		"RUN echo debconf shared/accepted-oracle-license-v1-1 seen true | debconf-set-selections\n" +
		"RUN apt-get update && apt-get install -y oracle-java8-installer\n" +
		"RUN update-java-alternatives -s java-8-oracle",
		"#!/usr/bin/env bats\n\n@test \"Check that Java is installed\" {\ncommand java -version\n}")

	leiningen = shared.NewFeature_str("leiningen", "demo feature with no dependency", "technomancy", nil, "" +
		"RUN wget https://raw.githubusercontent.com/technomancy/leiningen/stable/bin/lein -O /usr/bin/lein \\\n" +
		"&& chmod +x /usr/bin/lein",
		"#!/usr/bin/env bats\n\n@test \"Check that Leiningen is installed\" {\ncommand lein -v\n}")

	javalein = shared.NewFeature_str("java+leiningen", "a feature depending on java and leiningen", "some_author", []string{"java", "leiningen"}, "", "")

)

func main() {
	fmt.Printf("Running Pazuzu demo setup...\n")

	registry, err := storageconnector.NewRegistryStorage("localhost", 8080, nil)
	if err != nil {
		fmt.Errorf("Something bad happened: %v\n", err)
	}

	err = registry.AddFeature(java)
	err = registry.AddFeature(leiningen)
	err = registry.AddFeature(javalein)
	if err != nil {
		fmt.Errorf("Something bad happened: %v\n", err)
	}

	fmt.Printf("Done!\n")
}


