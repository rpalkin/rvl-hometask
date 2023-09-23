include ./common.env
include ./gcp.env

IMAGE = $(REGISTRY)birthday-app
MIGRATION_IMAGE = $(REGISTRY)birthday-migration
DB_HOST = $(DB_SVC_NAME).$(DB_NAMESPACE).svc.cluster.local
DB_PORT = 5432

export

.PHONY: help
help: ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

app-test: ## Run application tests
	cd ./app && go test ./...

docker-up: ## Run local docker environment
	docker compose --project-directory ./app up --build -d

docker-down: ## Stop local docker environment
	docker compose --project-directory ./app down

build: ## Build docker images locally
	docker build -t $(IMAGE):$(IMAGE_TAG) ./app
	docker build -t $(MIGRATION_IMAGE):$(IMAGE_TAG) ./app/migrations

minikube-up: build ## Run application in minikube environment
	# Start minikube
	minikube status || minikube start --wait all

	# Load images
	minikube image load $(IMAGE):$(IMAGE_TAG)
	minikube image load $(MIGRATION_IMAGE):$(IMAGE_TAG)

	# Create environment (database)
	envsubst < ./local/env/postgres/postgres.env.template > ./local/env/postgres/postgres.env
	kustomize build ./local/env | kubectl apply -f -
	kubectl wait --for=condition=Ready pod -l app=postgres -n postgres --timeout=120s

	# Run application
	helm upgrade --install birthday-app ./app/chart \
		-n $(NAMESPACE) \
		--create-namespace \
		--set image.tag=$(IMAGE_TAG) \
		--set image.repository=$(IMAGE) \
		--set migration.image.repository=$(MIGRATION_IMAGE) \
		--set migration.image.tag=$(IMAGE_TAG) \
		--set appConfig.dbHost=$(DB_HOST) \
		--set appConfig.dbPort=$(DB_PORT) \
		--set appConfig.dbUser=$(DB_USERNAME) \
		--set appConfig.dbPassword=$(DB_PASSWORD) \
		--set appConfig.dbName=$(DB_NAME)

minikube-test: ## Run application functional tests in minikube environment
	bash ./local/minikube-test.sh

minikube-down: ## Stop minikube environment
	minikube stop
	minikube delete

gcp-up: ## Create GCP environment
	terraform -chdir=./gcp init
	terraform -chdir=./gcp apply \
		-var="gcloud_project_id=$(GOOGLE_PROJECT_ID)" \
		-var="gcloud_region=$(GOOGLE_REGION)" \
	 	-var="db_user=$(DB_USERNAME)" \
	 	-var="db_password=$(DB_PASSWORD)" \
	 	-var="db_name=$(DB_NAME)" \
	 	-var="dns_zone=$(DNS_ZONE)"
	$(MAKE) cluster-setup

cluster-setup: prepare-gcp-deploy ## Setup created cluster
	kustomize build ./gcp/in_cluster | envsubst | kubectl apply -f -

prepare-gcp-deploy: ## Prepare environment variables for GCP deploy
	$(eval export REGISTRY=$(shell terraform -chdir=./gcp output -raw registry_location))
	$(eval export DNS_SA=$(shell terraform -chdir=./gcp output -raw dns_sa))
	$(eval export CLOUDSQL_HOST=$(shell terraform -chdir=./gcp output -raw cloudsql_host))
	gcloud container clusters get-credentials birthday-cluster --location=$(GOOGLE_REGION)

gcp-down: ## Destroy GCP environment
	terraform -chdir=./gcp destroy -auto-approve \
		-var="gcloud_project_id=$(GOOGLE_PROJECT_ID)" \
		-var="gcloud_region=$(GOOGLE_REGION)" \
	 	-var="db_user=$(DB_USERNAME)" \
	 	-var="db_password=$(DB_PASSWORD)" \
	 	-var="db_name=$(DB_NAME)" \
		-var="dns_zone=$(DNS_ZONE)"

cloud-build: prepare-gcp-deploy ## Build docker images in GCP
	gcloud builds submit --tag $(IMAGE):$(IMAGE_TAG) ./app
	gcloud builds submit --tag $(MIGRATION_IMAGE):$(IMAGE_TAG) ./app/migrations

gcp-deploy: cloud-build ## Deploy and update application in GCP
	$(eval export)

	# Run/update application
	helm upgrade --install birthday-app ./app/chart \
		-n $(NAMESPACE) \
		--create-namespace \
		--set image.tag=$(IMAGE_TAG) \
		--set image.repository=$(IMAGE) \
		--set migration.image.repository=$(MIGRATION_IMAGE) \
		--set migration.image.tag=$(IMAGE_TAG) \
		--set appConfig.dbHost=$(DB_HOST) \
		--set appConfig.dbPort=$(DB_PORT) \
		--set appConfig.dbUser=$(DB_USERNAME) \
		--set appConfig.dbPassword=$(DB_PASSWORD) \
		--set appConfig.dbName=$(DB_NAME)