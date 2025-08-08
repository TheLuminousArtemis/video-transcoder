<script>
  import { Router, Route, Link } from 'svelte-routing';
  import Home from './routes/Home.svelte';
  import Upload from './routes/Upload.svelte';

  // Global styles
  import './app.css';
  
  // Base URL for routing
  export let url = '';
</script>

<Router {url}>
  <div class="app">
    <nav class="navbar">
      <Link to="/" class="nav-logo">Video Transcoder</Link>
      <div class="nav-links">
        <Link to="/" class="nav-link">Home</Link>
        <Link to="/upload" class="nav-link">Upload</Link>
      </div>
    </nav>

    <main class="main-content">
      <Route path="/" let:params>
        <Home />
      </Route>
      <Route path="/upload" let:params>
        <Upload />
      </Route>
    </main>

    <footer class="footer">
      <p> 2025 Video Transcoder</p>
    </footer>
  </div>
</Router>
