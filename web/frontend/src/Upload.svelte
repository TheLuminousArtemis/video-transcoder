<script lang="ts">
    import { onDestroy } from "svelte";
    import VideoPlayer from "./VideoPlayer.svelte";

    let file: File | null = null;
    let uploading = false;
    let processed = false;
    let videoId: string | null = null;
    let errorMsg = "";
    let ws: WebSocket | null = null;

    interface UploadResponse {
        videoId: string;
    }

    interface StatusMessage {
        id: string;
        processed: boolean;
        status: string;
    }

    const handleFileChange = (e: Event) => {
        const target = e.target as HTMLInputElement;
        if (target.files && target.files.length > 0) {
            file = target.files[0];
        }
    };

    const uploadVideo = async () => {
        if (!file) {
            errorMsg = "Please select a file before uploading.";
            return;
        }

        uploading = true;
        errorMsg = "";

        try {
            const formData = new FormData();
            formData.append("video", file);

            const res = await fetch("http://localhost:3030/api/v1/upload", {
                method: "POST",
                body: formData,
            });

            if (!res.ok) {
                throw new Error(`Upload failed: ${res.status}`);
            }

            const data: UploadResponse = await res.json();
            videoId = data.videoId;

            connectWebSocket(videoId);
        } catch (err: unknown) {
            errorMsg = err instanceof Error ? err.message : String(err);
            uploading = false;
        }
    };

    const connectWebSocket = (id: string) => {
        ws = new WebSocket(`ws://localhost:3030/api/v1/ws?id=${id}`);

        ws.onopen = () => {
            console.log("✅ WebSocket connected");
        };

        ws.onmessage = (event: MessageEvent) => {
            try {
                const msg: StatusMessage = JSON.parse(event.data);
                console.log("📨 Message from server:", msg);

                if (msg.processed) {
                    processed = true;
                    uploading = false;
                    ws?.close();
                }
            } catch (e) {
                console.error("Invalid WebSocket message:", e);
            }
        };

        ws.onerror = (err: Event) => {
            console.error("WebSocket error:", err);
            errorMsg = "WebSocket connection error.";
        };

        ws.onclose = () => {
            console.log("🔌 WebSocket closed");
        };
    };

    onDestroy(() => {
        ws?.close();
    });
</script>

<div class="upload-container">
    <input type="file" accept="video/*" on:change={handleFileChange} />
    <button on:click={uploadVideo} disabled={uploading || processed}>
        {uploading ? "Uploading..." : "Upload"}
    </button>

    {#if errorMsg}
        <p style="color:red;">{errorMsg}</p>
    {/if}

    {#if uploading && !processed}
        <div>
            <p>Processing video...</p>
            <div class="spinner"></div>
        </div>
    {/if}

    {#if processed && videoId}
        <VideoPlayer {videoId} />
    {/if}
</div>

<style>
    .upload-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1rem;
    }
    .spinner {
        border: 4px solid #f3f3f3;
        border-top: 4px solid #3b82f6;
        border-radius: 50%;
        width: 36px;
        height: 36px;
        animation: spin 1s linear infinite;
    }
    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>
