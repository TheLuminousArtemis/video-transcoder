<script lang="ts">
  import { link } from 'svelte-routing';
  
  // No need for explicit navigation with svelte-routing's Link component
  // We'll use the Link component directly in the template
</script>

<main class="home-container">
  <h1>Video Transcoder</h1>
  <p>Upload and transcode your videos with ease</p>
  
  <div class="action-buttons">
    <a href="/upload" use:link class="cta-button">
      Upload Video
    </a>
  </div>
</main>

<style>
  .home-container {
    max-width: 800px;
    margin: 2rem auto;
    padding: 2rem;
    text-align: center;
  }

  h1 {
    font-size: 2.5rem;
    color: #333;
    margin-bottom: 1rem;
  }

  p {
    font-size: 1.2rem;
    color: #666;
    margin-bottom: 2rem;
  }

  .cta-button {
    padding: 0.75rem 2rem;
    font-size: 1.1rem;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.3s;
  }

  .cta-button:hover {
    background-color: #45a049;
  }
</style>
