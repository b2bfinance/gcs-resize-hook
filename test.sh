#!/usr/bin/env bash

# Example:
# ```bash
# wget -O ./image.jpg https://cdn.pixabay.com/photo/2018/07/16/16/08/island-3542290_960_720.jpg;
# ./test.sh "capis" ./image.jpg "/images/providers" 240 240;
# ```

set -e

command -v imageinfo >/dev/null 2>&1 || { echo >&2 "Aborting, please make sure imageinfo is in the PATH."; exit 1; }
command -v gsutil >/dev/null 2>&1 || { echo >&2 "Aborting, please make sure Google Cloud SDK is installed and is in the PATH."; exit 1; }

BUCKET=$1
IMG_FILE=$2
DESTINATION=$3
EXPECTED_WIDTH=$4
EXPECTED_HEIGHT=$5
IMG_FILE_BASE=$(basename "${IMG_FILE}")
TEST_FILE_LOCATION="/tmp/gcs-resize-hook-${IMG_FILE_BASE}"

test -f "${IMG_FILE}" || { echo >&2 "Aborting, the source image file does not exist."; exit 1; }

ORIGINAL_WIDTH=$(imageinfo --width "${IMG_FILE}")
ORIGINAL_HEIGHT=$(imageinfo --height "${IMG_FILE}")

echo "The source image: width ${ORIGINAL_WIDTH}px height ${ORIGINAL_HEIGHT}px"

gsutil cp "${IMG_FILE}" "gs://${BUCKET}${DESTINATION}/${IMG_FILE_BASE}"

echo "Image uploaded, waiting for 10 seconds to start checking for an updated image."
count=0  
while [ $count -le 10 ]  
do  
    echo "Result test attempt: $count"
    gsutil cp "gs://${BUCKET}${DESTINATION}/${IMG_FILE_BASE}" "${TEST_FILE_LOCATION}"

    ACTUAL_WIDTH=$(imageinfo --width "${TEST_FILE_LOCATION}")
    ACTUAL_HEIGHT=$(imageinfo --height "${TEST_FILE_LOCATION}")

    echo "The new image: width ${ACTUAL_WIDTH}px height ${ACTUAL_HEIGHT}px"

    count=$(( $count + 1 ))
done

gsutil rm "gs://${BUCKET}${DESTINATION}/${IMG_FILE_BASE}"
