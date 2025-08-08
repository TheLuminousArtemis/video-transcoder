<script lang="ts">
  import { link } from 'svelte-routing';
  import { onMount } from 'svelte';
  
  let file: File | null = null;
  let isUploading = false;
  let uploadStatus: string = '';
  let uploadProgress: number = 0;
  
  // TypeScript interface for the event target with files
  interface HTMLInputEvent extends Event {
    target: HTMLInputElement & {
      files: FileList | null;
    };
  }

  const handleFileChange = (event: Event) => {
    const target = event.target as HTMLInputElement;
    if (target.files && target.files.length > 0) {
      file = target.files[0];
    }
  };
  
  // Handle form submission
  const handleSubmit = (e: Event) => {
    e.preventDefault();
    if (file) {
      uploadFile();
    }
  };

  const uploadFile = async () => {
    if (!file) return;

    const formData = new FormData();
    formData.append('video', file);

    isUploading = true;
    uploadStatus = 'Uploading...';
    uploadProgress = 0;

    try {
      const xhr = new XMLHttpRequest();
      
      xhr.upload.onprogress = (event) => {
        if (event.lengthComputable) {
          uploadProgress = Math.round((event.loaded / event.total) * 100);
        }
      };

      xhr.onload = () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          uploadStatus = 'Upload successful!';
          setTimeout(() => uploadStatus = '', 3000);
        } else {
          uploadStatus = `Error: ${xhr.statusText}`;
        }
        isUploading = false;
      };

      xhr.onerror = () => {
        uploadStatus = 'Upload failed. Please try again.';
        isUploading = false;
      };

      xhr.open('POST', '/api/v1/upload', true);
      xhr.send(formData);
    } catch (error) {
      console.error('Upload error:', error);
      uploadStatus = 'An error occurred during upload';
      isUploading = false;
    }
  };
  
  // No default export needed for Svelte components
</script>

<div class="upload-container">
  <h2>Upload Video</h2>
  
  <form on:submit|preventDefault={handleSubmit} class="file-upload">
    <input 
      type="file" 
      id="video-upload" 
      accept="video/*" 
      on:change={handleFileChange}
      disabled={isUploading}
    />
    <label for="video-upload" class="upload-button">
      {#if file}
        {file.name} ({(file.size / (1024 * 1024)).toFixed(2)} MB)
      {:else}
        Choose a video file
      {/if}
    </label>

    {#if file}
      <div class="upload-actions">
        <button 
          type="submit"
          class="upload-button"
          disabled={isUploading}
        >
          {isUploading ? 'Uploading...' : 'Upload Video'}
        </button>
      </div>
    {/if}

    {#if uploadStatus}
      <div class="status">
        <p>{uploadStatus}</p>
        {#if uploadProgress > 0}
          <progress value={uploadProgress} max="100">{uploadProgress}%</progress>
        {/if}
      </div>
    {/if}
    
    <div class="back-link">
      <a href="/" use:link>← Back to Home</a>
    </div>
  </form>
</div>

<style>
  .upload-container {
    max-width: 600px;
    margin: 2rem auto;
    padding: 2rem;
    background: #f5f5f5;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .file-upload {
    margin: 1.5rem 0;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  input[type="file"] {
    display: none;
  }

  .upload-button {
    display: inline-block;
    padding: 0.75rem 1.5rem;
    background: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
    transition: background-color 0.3s;
    text-align: center;
  }

  .upload-button:hover:not(:disabled) {
    background: #45a049;
  }

  .upload-button:disabled {
    background: #cccccc;
    cursor: not-allowed;
  }

  .status {
    margin-top: 1.5rem;
    padding: 1rem;
    background: #e8f5e9;
    border-radius: 4px;
  }

  progress {
    width: 100%;
    margin-top: 0.5rem;
    height: 8px;
    border-radius: 4px;
  }

  progress::-webkit-progress-bar {
    background-color: #f0f0f0;
    border-radius: 4px;
  }

  progress::-webkit-progress-value {
    background-color: #4CAF50;
    border-radius: 4px;
  }
  
  .back-link {
    margin-top: 1.5rem;
    text-align: center;
  }
  
  .back-link a {
    color: #4CAF50;
    text-decoration: none;
    font-weight: 500;
  }
  
  .back-link a:hover {
    text-decoration: underline;
  }
</style>
