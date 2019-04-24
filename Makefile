BUCKET ?= my-bucket-name
PREFIX ?= /
TRIGGER_NAME ?= resizehook
REGION ?= europe-west2
WIDTH ?= 500
HEIGHT ?= 0

.PHONY: publish
publish:
	gcloud functions deploy $(TRIGGER_NAME) --set-env-vars "CFG_PREFIX=$(PREFIX)" --set-env-vars "CFG_WIDTH=$(WIDTH)" --set-env-vars "CFG_HEIGHT=$(HEIGHT)" --runtime go111 --entry-point Resize --trigger-resource $(BUCKET) --trigger-event google.storage.object.finalize --memory 128MB --retry --region $(REGION)
