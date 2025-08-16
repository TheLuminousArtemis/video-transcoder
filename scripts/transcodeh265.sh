#!/usr/bin/env bash
set -euo pipefail

INPUT="$1"
UUID=$(basename "$INPUT" | cut -d. -f1)
WORKDIR="/outputs/$UUID"
mkdir -p "$WORKDIR"

RESOLUTIONS=("640x360" "854x480" "1280x720")
BITRATES=("500k" "1250k" "2500k")
LABELS=("360p" "480p" "720p")

MASTER="$WORKDIR/master.m3u8"
echo "#EXTM3U" > "$MASTER"
echo "#EXT-X-VERSION:7" >> "$MASTER"

# --- Detect audio tracks safely ---
AUDIO_TRACKS=$(ffprobe -v error -select_streams a \
  -show_entries stream=index:stream_tags=language \
  -of csv=p=0 "$INPUT" || true)

INDEX_COUNT=0
while IFS=, read -r INDEX LANG; do
  [ -z "$INDEX" ] && continue
  LANG=${LANG:-und}
  AUDIO_DIR="$WORKDIR/audio/$LANG"
  mkdir -p "$AUDIO_DIR"

  echo "Extracting audio track index $INDEX (lang=$LANG)..."

  ffmpeg -y -i "$INPUT" \
    -map 0:a:${INDEX}? -vn \
    -c:a aac -profile:a aac_low -b:a 128k \
    -f hls -hls_time 4 -hls_playlist_type vod \
    -hls_segment_type fmp4 -hls_flags independent_segments \
    -hls_fmp4_init_filename "init.mp4" \
    -hls_segment_filename "$AUDIO_DIR/audio_%03d.m4s" \
    "$AUDIO_DIR/audio.m3u8"

  if [[ $INDEX_COUNT -eq 0 ]]; then
    echo "#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID=\"audio-$LANG\",NAME=\"$LANG\",DEFAULT=YES,AUTOSELECT=YES,URI=\"audio/$LANG/audio.m3u8\"" >> "$MASTER"
  else
    echo "#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID=\"audio-$LANG\",NAME=\"$LANG\",DEFAULT=NO,AUTOSELECT=YES,URI=\"audio/$LANG/audio.m3u8\"" >> "$MASTER"
  fi
  INDEX_COUNT=$((INDEX_COUNT+1))
done <<< "$AUDIO_TRACKS"

# --- Generate video renditions ---
for i in "${!RESOLUTIONS[@]}"; do
  RES="${RESOLUTIONS[$i]}"
  BR="${BITRATES[$i]}"
  LABEL="${LABELS[$i]}"
  SEG_DIR="$WORKDIR/$LABEL"
  mkdir -p "$SEG_DIR"

  echo "Encoding video $LABEL ($RES @ $BR)..."

  ffmpeg -y -i "$INPUT" \
    -c:v libx265 -tag:v hvc1 -preset medium -b:v "$BR" -s "$RES" \
    -pix_fmt yuv420p -an \
    -f hls -hls_time 4 -hls_playlist_type vod \
    -hls_segment_type fmp4 -hls_flags independent_segments \
    -hls_fmp4_init_filename "init.mp4" \
    -hls_segment_filename "$SEG_DIR/segment_%03d.m4s" \
    "$SEG_DIR/playlist.m3u8"

  BW=$(( ${BR%k} * 1000 ))
  if [[ -n "$AUDIO_TRACKS" ]]; then
    for LANG in $(echo "$AUDIO_TRACKS" | cut -d, -f2 | sort -u); do
      LANG=${LANG:-und}
      echo "#EXT-X-STREAM-INF:BANDWIDTH=$BW,RESOLUTION=$RES,CODECS=\"hvc1.1.6.L123,mp4a.40.2\",AUDIO=\"audio-$LANG\"" >> "$MASTER"
      echo "$LABEL/playlist.m3u8" >> "$MASTER"
    done
  else
    # fallback: no audio
    echo "#EXT-X-STREAM-INF:BANDWIDTH=$BW,RESOLUTION=$RES,CODECS=\"hvc1.1.6.L123\"" >> "$MASTER"
    echo "$LABEL/playlist.m3u8" >> "$MASTER"
  fi
done

echo "HLS Transcoding Complete!"
echo "Master Playlist: $MASTER"
