<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import videojs from "video.js";
    import type Player from "video.js/dist/types/player";

    import "video.js/dist/video-js.css";
    import "jb-videojs-hls-quality-selector/dist/videojs-hls-quality-selector.css";

    import "videojs-contrib-quality-levels";
    import "jb-videojs-hls-quality-selector";

    export let videoId: string;

    let player: Player | null = null;
    let videoElement: HTMLVideoElement;
    const backendUrl = "http://localhost:3030";

    onMount(() => {
        if (!videoElement) return;

        player = videojs(videoElement, {
            controls: true,
            autoplay: false,
            preload: "auto",
            fluid: true,
            aspectRatio: "16:9",
            controlBar: {
                children: [
                    "playToggle",
                    "progressControl",
                    "volumePanel",
                    "subsCapsButton",
                    "audioTrackButton",
                    "qualitySelector",
                    "fullscreenToggle",
                ],
            },
        }) as Player;

        player.src({
            src: `${backendUrl}/videos/${videoId}/master.m3u8`,
            type: "application/x-mpegURL",
        });

        player.ready(() => {
            (player as any).hlsQualitySelector({
                displayCurrentQuality: true,
            });
        });

        player.on("loadedmetadata", () => {
            const ql = (player as any).qualityLevels?.();
            console.log("Detected quality levels:", ql?.length);
            if (ql) {
                for (let i = 0; i < ql.length; i++) {
                    console.log(
                        `Level ${i}: ${ql[i].height}p @ ${ql[i].bitrate}`,
                    );
                }
            }

            const tracks = (player as any).audioTracks?.();
            if (tracks?.length) {
                console.log("Available audio tracks:");
                for (let i = 0; i < tracks.length; i++) {
                    console.log(
                        `${i}: ${tracks[i]?.label || tracks[i]?.id || "Unnamed"}`,
                    );
                }
            }
        });
    });

    onDestroy(() => {
        if (player) {
            player.dispose();
            player = null;
        }
    });
</script>

<div class="player-shell">
    <!-- svelte-ignore a11y-media-has-caption -->
    <video
        bind:this={videoElement}
        class="video-js vjs-default-skin"
        playsinline
    ></video>
</div>

<style>
    .player-shell {
        width: min(95vw, 960px);
        margin: 0 auto;
    }
    .video-js {
        display: block;
        width: 100%;
        height: auto;
        margin: 0 auto;
    }
</style>
