include .env
export

runsonar:
	go test ./... -coverprofile=coverage.out
	docker compose -f .\compose.sonar.yaml up -d
	docker run --rm \
		-e SONAR_HOST_URL="http://host.docker.internal:9000" \
		-e SONAR_LOGIN=${SONAR_TOKEN} \
		-v ".:/usr/src" \
		sonarsource/sonar-scanner-cli \
		-Dsonar.projectBaseDir=/usr/src \
		-Dsonar.login=${SONAR_TOKEN}
