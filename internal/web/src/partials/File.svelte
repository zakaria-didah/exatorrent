<script lang="ts">
  import { Send, fileviewpath, fileviewinfohash, fsfileinfo, fileType, socket, fileSize, filepagediscon } from './core';
  import slocation from 'slocation';
  import { onDestroy, onMount } from 'svelte';

  let stream = false;
  let ft = 'unknown';
  let copied = false;

  $: streamUrl = `${location.origin}/api/${stream ? 'stream' : 'torrent'}/${$fileviewinfohash}/${$fileviewpath}`;

  onMount(() => {
    if (socket == null || socket == undefined || socket?.readyState === WebSocket.CLOSED) {
      slocation.goto('/');
    }
    Send({
      command: 'getfsfileinfo',
      data1: $fileviewinfohash,
      data2: $fileviewpath
    });
    filepagediscon.set(true);
    ft = fileType($fileviewpath);
  });

  onDestroy(() => {
    filepagediscon.set(false);
    document.title = 'exatorrent';
  });

  function copyUrl() {
    navigator.clipboard.writeText(streamUrl).then(() => {
      copied = true;
      setTimeout(() => { copied = false; }, 2000);
    });
  }
</script>

<svelte:head>
  <title>{$fsfileinfo?.name}</title>
</svelte:head>

<div class="mx-auto max-w-2xl px-3">
  <!-- Header -->
  <div class="flex items-center justify-between h-14 mb-2">
    <button aria-label="Go back"
      class="flex items-center justify-center w-10 h-10 rounded-lg bg-white/5 border border-white/10 glass transition-colors duration-150 hover:bg-white/10 flex-shrink-0"
      on:click={() => { history.length > 2 ? history.back() : slocation.goto('/'); }}>
      <svg xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
      </svg>
    </button>
    <button class="flex items-center gap-2 px-4 py-2 rounded-lg noHL transition-colors duration-150 hover:bg-white/5"
      on:click={() => { slocation.goto('/'); }}>
      <span class="font-semibold text-lg text-slate-100 tracking-tight">exatorrent</span>
    </button>
    <div class="w-10"></div>
  </div>

  <!-- Player -->
  {#if ft === 'video' || ft === 'audio'}
    <div class="rounded-xl overflow-hidden bg-black/50 border border-white/10">
      <!-- svelte-ignore a11y-media-has-caption -->
      {#if ft === 'video'}
        <video controls class="w-full" src={streamUrl}></video>
      {:else}
        <audio controls class="w-full" src={streamUrl}></audio>
      {/if}
    </div>
  {/if}

  <!-- File info -->
  <div class="bg-white/5 glass border border-white/10 rounded-xl mt-3 px-4 py-3">
    <p class="text-slate-200 text-sm font-medium break-all">{$fsfileinfo?.name}</p>
    <p class="text-slate-500 text-xs mt-0.5 tabular-nums">{fileSize($fsfileinfo?.size)}</p>
  </div>

  <!-- Stream API toggle -->
  <div class="bg-white/5 glass border border-white/10 rounded-xl mt-2 px-4 py-2.5 flex items-center justify-between">
    <span class="text-slate-400 text-sm">Stream API</span>
    <label class="relative inline-flex items-center cursor-pointer">
      <input type="checkbox" class="sr-only peer" bind:checked={stream} />
      <div class="w-8 h-4 bg-slate-700 peer-focus:ring-2 peer-focus:ring-violet-500/40 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-3 after:w-3 after:transition-all peer-checked:bg-violet-600"></div>
    </label>
  </div>

  <!-- URL bar -->
  <div class="bg-white/5 glass border border-white/10 rounded-xl mt-2 flex items-center overflow-hidden">
    <input type="text" class="flex-1 bg-transparent px-4 py-2.5 text-xs text-slate-400 truncate border-none focus:outline-none" disabled value={streamUrl} />
    <button type="button" class="px-3 py-2.5 border-l border-white/[0.06] hover:bg-white/10 transition-colors duration-150 flex-shrink-0"
      on:click={copyUrl}>
      {#if copied}
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-violet-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" /></svg>
      {:else}
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" /></svg>
      {/if}
    </button>
  </div>

  <!-- Action buttons -->
  <div class="mt-3 space-y-2">
    {#if ft === 'video' || ft === 'audio'}
      <a href="vlc://{location.origin}/api/{stream ? 'stream' : 'torrent'}/{$fileviewinfohash}/{$fileviewpath}" target="_blank" rel="noopener noreferrer" class="block">
        <button type="button" class="w-full py-2.5 px-4 text-sm font-medium rounded-xl text-white bg-violet-600 hover:bg-violet-500 active:bg-violet-700 transition-colors duration-150 focus:outline-none flex items-center justify-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" /></svg>
          Open in VLC
        </button>
      </a>

      <a href="intent://{location.host}/api/{stream ? 'stream' : 'torrent'}/{$fileviewinfohash}/{$fileviewpath}#Intent;type=video/any;package=is.xyz.mpv;scheme={location.protocol.slice(0, -1)};end;" target="_blank" rel="noopener noreferrer" class="block md:hidden">
        <button type="button" class="w-full py-2.5 px-4 text-sm font-medium rounded-xl text-white bg-white/5 glass border border-white/10 hover:bg-white/10 transition-colors duration-150 focus:outline-none flex items-center justify-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" /></svg>
          Open in MPV
        </button>
      </a>
    {/if}

    <a href={streamUrl} target="_blank" rel="noopener noreferrer" download class="block">
      <button type="button" class="w-full py-2.5 px-4 text-sm font-medium rounded-xl text-white bg-white/5 glass border border-white/10 hover:bg-white/10 transition-colors duration-150 focus:outline-none flex items-center justify-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
        Download
      </button>
    </a>
  </div>
</div>
