#!/usr/bin/env bash
set -euo pipefail

INPUT="$1"
UUID=$(basename "$INPUT" | cut -d. -f1)
WORKDIR="/outputs/$UUID"

RESOLUTIONS=("640x360" "854x480" "1280x720")
BITRATES=("500k" "1250k" "2500k")
LABELS=("360p" "480p" "720p")

mkdir -p "$WORKDIR"
MASTER="$WORKDIR/master.m3u8"

echo "Processing: $INPUT"
echo "Output: $WORKDIR"

# --- Detect Audio Tracks ---
AUDIO_TRACKS=$(ffprobe -v error -select_streams a -show_entries stream=index:stream_tags=language -of csv=p=0 "$INPUT")

if [[ -z "$AUDIO_TRACKS" ]]; then
  echo "No audio tracks detected, proceeding video-only!"
else
  echo "Detected audio tracks:"
  echo "$AUDIO_TRACKS"
fi

# --- Extract Each Audio Track ---
INDEX_COUNT=0
while IFS=, read -r INDEX LANG; do
  if [[ -z "$INDEX" ]]; then
    continue
  fi

  LANG=${LANG:-und}
  AUDIO_DIR="$WORKDIR/audio/$LANG"
  mkdir -p "$AUDIO_DIR"

  echo "Extracting audio track index $INDEX (lang=$LANG)..."
  ffmpeg -y -i "$INPUT" \
    -map 0:a:${INDEX}? -c:a aac -b:a 128k \
    -hls_time 4 -hls_playlist_type vod \
    -hls_segment_filename "$AUDIO_DIR/audio_%03d.ts" \
    "$AUDIO_DIR/audio.m3u8"

  # Add EXT-X-MEDIA entry (mark first audio as default)
  if [[ $INDEX_COUNT -eq 0 ]]; then
    echo "#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID=\"audio-$LANG\",NAME=\"$LANG\",DEFAULT=YES,AUTOSELECT=YES,URI=\"audio/$LANG/audio.m3u8\"" >> "$MASTER"
  else
    echo "#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID=\"audio-$LANG\",NAME=\"$LANG\",DEFAULT=NO,AUTOSELECT=YES,URI=\"audio/$LANG/audio.m3u8\"" >> "$MASTER"
  fi

  INDEX_COUNT=$((INDEX_COUNT+1))
done <<< "$AUDIO_TRACKS"

# --- Master Playlist Header ---
if [[ -s "$MASTER" ]]; then
  sed -i '1i#EXTM3U\n#EXT-X-VERSION:3' "$MASTER"
else
  echo -e "#EXTM3U\n#EXT-X-VERSION:3" > "$MASTER"
fi

# --- Generate Video Renditions ---
for i in "${!RESOLUTIONS[@]}"; do
  RES="${RESOLUTIONS[$i]}"
  BR="${BITRATES[$i]}"
  LABEL="${LABELS[$i]}"
  SEG_DIR="$WORKDIR/$LABEL"

  mkdir -p "$SEG_DIR"
  echo "Encoding video $LABEL ($RES @ $BR)..."

  ffmpeg -y -i "$INPUT" \
    -c:v libx264 -tag:v avc1 -preset medium -b:v "$BR" -s "$RES" \
    -pix_fmt yuv420p \
    -an \
    -hls_time 4 -hls_playlist_type vod \
    -hls_segment_filename "$SEG_DIR/segment_%03d.ts" \
    "$SEG_DIR/playlist.m3u8"

  BW=$(( ${BR%k} * 1000 ))
  for LANG in $(echo "$AUDIO_TRACKS" | cut -d, -f2 | sort -u); do
    LANG=${LANG:-und}
    echo "#EXT-X-STREAM-INF:BANDWIDTH=$BW,RESOLUTION=$RES,CODECS=\"avc1.4d401f,mp4a.40.2\",AUDIO=\"audio-$LANG\"" >> "$MASTER"
    echo "$LABEL/playlist.m3u8" >> "$MASTER"
  done
done

echo "HLS Transcoding Complete!"
echo "Master Playlist: $MASTER"
