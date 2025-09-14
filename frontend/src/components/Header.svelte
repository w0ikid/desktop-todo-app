<script lang="ts">
  export let isDarkMode: boolean;
  export let toggleTheme: () => void;
</script>

<header class="app-header">
  <div class="header-content">
    <div class="header-top">
      <h1 class="app-title">
        <span class="title-icon">✨</span>
        Todo Dashboard
      </h1>
      <button 
        class="theme-toggle" 
        on:click={toggleTheme}
        aria-label="Toggle theme"
        title={isDarkMode ? 'Switch to light mode' : 'Switch to dark mode'}
      >
        <span class="theme-icon">
          {#if isDarkMode}
            ☀️
          {:else}
            🌙
          {/if}
        </span>
      </button>
    </div>
    <p class="app-subtitle">Manage your tasks with style and efficiency</p>
  </div>
</header>