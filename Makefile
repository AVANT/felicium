.PHONY: dev_tunnel run testing_site dev package clean upload staging prod backup restore enable_cross_compile ssh_prod ssh_staging

include $(shell readlink .env)

project_base_path := github.com/AVANT/felicium
go_build_flags := GOOS=linux GOARCH=amd64 CGO_ENABLED=0

dev:
	bash -c "( make dev_tunnel & ) ; make run"

testing_site:
	bash -c "( make testing_tunnel & ) ; make dev"

run:
	revel run $(project_base_path)/moonrakr

dev_tunnel:
	ssh -N -o StrictHostKeyChecking=no -i $(AVANT_KEY_PATH) -L 5984:localhost:5984 -L 9200:localhost:9200 $(AVANT_USER)@$(AVANT_STAGING_URL)

testing_tunnel:
	ssh -N -o StrictHostKeyChecking=no -R 5000:127.0.0.1:9000 -i $(AVANT_KEY_PATH) $(AVANT_USER)@$(AVANT_STAGING_URL)

package: clean
	$(go_build_flags) revel package $(project_base_path)/moonrakr
	mkdir -p dist/migrations/
	tar xf moonrakr.tar.gz -C  dist/
	$(go_build_flags) go build -o dist/moonrakr-tool $(project_base_path)/tools
	for i in moonrakr/app/db/migrations/*.go; do $(go_build_flags) go build -o dist/migrations/$$i $$i ; done
	tar czf moonrakr.latest.tar.gz dist/
	rm -rf dist moonrakr.tar.gz

clean:
	rm -rf dist moonrakr.latest.tar.gz

upload: package clean
	AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) aws s3 cp moonrakr.latest.tar.gz s3://fanny-pack/api-release/moonrakr.latest.tar.gz

staging:
	ssh -o StrictHostKeyChecking=no -i $(AVANT_KEY_PATH) $(AVANT_USER)@$(AVANT_STAGING_URL) sudo /usr/local/bin/update-vvvnt.sh api

ssh_staging:
	ssh -o StrictHostKeyChecking=no -i $(AVANT_KEY_PATH) $(AVANT_USER)@$(AVANT_STAGING_URL)

prod:
	ssh -o StrictHostKeyChecking=no -i $(AVANT_KEY_PATH) $(AVANT_USER)@$(AVANT_PRODUCTION_URL) sudo /usr/local/bin/update-vvvnt.sh api

ssh_prod:
	ssh -o StrictHostKeyChecking=no -i $(AVANT_KEY_PATH) $(AVANT_USER)@$(AVANT_PRODUCTION_URL)

restore:
	AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) aws s3 cp s3://fanny-pack/api-release/moonrakr.latest.tar.gz.frozen s3://fanny-pack/api-release/moonrakr.latest.tar.gz

backup:
	AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) AWS_DEFAULT_REGION=$(AWS_DEFAULT_REGION) aws s3 cp s3://fanny-pack/api-release/moonrakr.latest.tar.gz s3://fanny-pack/api-release/moonrakr.latest.tar.gz.frozen

enable_cross_compile:
	bash -c "cd $$GOROOT/src/ && $(go_build_flags) $$GOROOT/src/make.bash --no-clean"