BUCKET ?= my-bucket-name
PREFIX ?= /
TRIGGER_NAME ?= resizehook
REGION ?= europe-west2

.PHONY: publish
publish:
	gcloud functions deploy $(TRIGGER_NAME) --set-env-vars "CFG_NAME_PREFIX=$(PREFIX)" --runtime go111 --entry-point Resize --trigger-resource $(BUCKET) --trigger-event google.storage.object.finalize --memory 128MB --retry --region $(REGION)
