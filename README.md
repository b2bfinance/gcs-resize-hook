# Google Cloud Storage resize hook.

## Deployment

Below are variables which you can override when using the `Make` cli, see below for more information.

- **BUCKET** - The name of the bucket to resize new finalised images in.
- **PREFIX** - A path prefix that images must be in before we consider resizing.
- **TRIGGER_NAME** - The name of the trigger, this defaults to `resizehook`
- **REGION** - The region you want to publish the function in, this defaults to `europe-west2`
- **WIDTH** - The desired width, set this to 0 if you want to stay in ratio with the defined height.
- **HEIGHT** - The desired height, set this to 0 if you want to stay in ration with the defined width.

### Deploy with Make

```bash
make publish VAR=VALUE
```

### Configuration Management

If you would like to store the setup process in version control then we have a stub shell script below.

```bash
#!/usr/bin/env sh

set -e

BUCKET=my-bucket-name
PREFIX=/my/prefix
TRIGGER_NAME=resizehook
REGION=europe-west2
WIDTH=500
HEIGHT=0

echo " [x] Cloning repository: https://github.com/legalweb/gcs-resize-hook.git"
git clone https://github.com/legalweb/gcs-resize-hook.git

pushd gcs-resize-hook >/dev/null

echo " [x] Deploying function to GCP"

gcloud functions deploy "${TRIGGER_NAME}" --set-env-vars "CFG_PREFIX=${PREFIX}" --set-env-vars "CFG_WIDTH=${WIDTH}" --set-env-vars "CFG_HEIGHT=${HEIGHT}" --runtime go111 --entry-point Resize --trigger-resource ${BUCKET} --trigger-event google.storage.object.finalize --memory 128MB --retry --region ${REGION}

popd >/dev/null
```

## License

This trigger is open-source software licensed under the [MIT license](https://opensource.org/licenses/MIT).
