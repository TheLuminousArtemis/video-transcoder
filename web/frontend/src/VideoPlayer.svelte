<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import videojs from "video.js";
    import type Player from "video.js/dist/types/player";
    import "video.js/dist/video-js.css";

    export let videoId: string;

    let player: Player | null = null;
    let videoElement: HTMLVideoElement;
    let currentRes: "auto" | 360 | 480 | 720 = "auto";

    const backendUrl = "http://localhost:3030";

    // Switch resolution by loading the variant playlist, preserving time/state
    function switchResolution(res: "auto" | 360 | 480 | 720) {
        if (!player) return;
        const base = `${backendUrl}/videos/${videoId}`;
        const newSrc =
            res === "auto"
                ? `${base}/master.m3u8`
                : `${base}/${res}p/playlist.m3u8`;

        const currentTime = player.currentTime?.() ?? 0;
        const wasPaused = player.paused?.();

        player.src({ src: newSrc, type: "application/x-mpegURL" });
        player.one?.("loadedmetadata", () => {
            try {
                player!.currentTime?.(currentTime);
            } catch {}
            if (wasPaused === false) {
                player!.play?.();
            }
        });
        currentRes = res;
    }

    onMount(() => {
        if (!videoElement) return;

        player = videojs(videoElement, {
            controls: true,
            autoplay: false,
            preload: "auto",
            fluid: false,
            width: 960,
            height: 540,
            controlBar: {
                children: [
                    "playToggle",
                    "progressControl",
                    "volumePanel",
                    "subsCapsButton",
                    "audioTrackButton",
                    "fullscreenToggle",
                ],
            },
        }) as Player;

        // Default to auto quality on load
        switchResolution("auto");

        player.on("loadedmetadata", () => {
            try {
                const tracks = (player as any).audioTracks?.();
                const len: number =
                    typeof tracks?.length === "function"
                        ? tracks.length()
                        : (tracks?.length ?? 0);
                if (len > 0) {
                    console.log("Available audio tracks:");
                    for (let i = 0; i < len; i++) {
                        const t =
                            typeof tracks.item === "function"
                                ? tracks.item(i)
                                : tracks[i];
                        console.log(`${i}: ${t?.label || t?.id || "Unnamed"}`);
                    }
                }
            } catch (e) {
                console.warn("Audio track enumeration failed", e);
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
    <video bind:this={videoElement} class="video-js vjs-default-skin" playsinline
    ></video>

    <div class="resolution-controls">
        <button
            class:active={currentRes === "auto"}
            on:click={() => switchResolution("auto")}
            disabled={currentRes === "auto"}
        >
            Auto
        </button>
        <button
            class:active={currentRes === 360}
            on:click={() => switchResolution(360)}
            disabled={currentRes === 360}
        >
            360p
        </button>
        <button
            class:active={currentRes === 480}
            on:click={() => switchResolution(480)}
            disabled={currentRes === 480}
        >
            480p
        </button>
        <button
            class:active={currentRes === 720}
            on:click={() => switchResolution(720)}
            disabled={currentRes === 720}
        >
            720p
        </button>
    </div>
</div>

<style>
    .player-shell {
        max-width: 960px;
        margin: 0 auto;
    }
    .video-js {
        display: block;
        margin: 0 auto;
        max-width: 960px;
        max-height: 540px;
    }
    .resolution-controls {
        display: flex;
        gap: 0.5rem;
        justify-content: center;
        align-items: center;
        margin-top: 0.75rem;
    }
    .resolution-controls button {
        padding: 0.4rem 0.8rem;
        border: 1px solid #d1d5db;
        background: #fff;
        color: #111;
        border-radius: 6px;
        cursor: pointer;
    }
    .resolution-controls button:hover {
        background: #f3f4f6;
        color: #111;
    }
    .resolution-controls button.active {
        background: #3b82f6;
        color: white;
        border-color: #3b82f6;
    }
</style>
