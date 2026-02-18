<script lang="ts">
  import { onMount } from 'svelte';
  import { versionstr, versionchecked, Send } from './core';

  onMount(() => {
    if ($versionchecked === false) {
      Send({ command: 'version' });
    }
  });
</script>

<div class="mx-auto max-w-5xl text-slate-200 mt-3 px-3">
  <p class="text-center text-lg mb-1 text-violet-300">exatorrent is torrent client</p>
  <p class="text-center mb-2 font-mono text-slate-400">License: GPLv3</p>
  <p class="text-center mb-5 font-mono text-slate-400">Version: {$versionstr}</p>
  <div class="flex justify-around">
    <a href="https://github.com/zakaria-didah/exatorrent" target="_blank" rel="noopener noreferrer">
      <button type="button" class="my-2 mx-2 flex justify-center py-3 px-4 border border-white/10 text-sm font-medium rounded-lg text-slate-200 bg-white/5 glass hover:bg-white/10 noHL font-mono focus:outline-none transition-colors duration-150">Source</button>
    </a>
  </div>
</div>
