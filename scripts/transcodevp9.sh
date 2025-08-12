#!/bin/bash

INPUT="$1"
BASENAME=$(basename "$INPUT" .mp4)

OUTDIR="vp9_renditions"
mkdir -p "$OUTDIR"

declare -A BITRATES=(
  ["360p"]="800k"
  ["480p"]="1400k"
  ["720p"]="2800k"
)

declare -A SCALES=(
  ["360p"]="640:360"
  ["480p"]="854:480"
  ["720p"]="1280:720"
)

for RES in 360p 480p 720p; do
    SCALE=${SCALES[$RES]}
    BITRATE=${BITRATES[$RES]}
    OUTFILE="${OUTDIR}/${BASENAME}_${RES}.webm"

    echo "=== Encoding $RES (${BITRATE}) ==="

    ffmpeg -y -i "$INPUT" \
        -c:v libvpx-vp9 \
        -b:v $BITRATE \
        -vf "scale=${SCALE}" \
        -row-mt 1 -tile-columns 2 -frame-parallel 1 \
        -threads $(nproc) \
        -speed 4 \
        -an -pass 1 -f null /dev/null

    ffmpeg -i "$INPUT" \
        -c:v libvpx-vp9 \
        -b:v $BITRATE \
        -vf "scale=${SCALE}" \
        -row-mt 1 -tile-columns 2 -frame-parallel 1 \
        -threads $(nproc) \
        -speed 2 \
        -c:a libopus -b:a 96k \
        -pass 2 "$OUTFILE"

    echo "Finished $RES"
done

rm -f ffmpeg2pass-0.log ffmpeg2pass-0.log.mbtree
