<script lang="ts">
  import { fileSize, Send, adminmode, slocation } from './core';

  export let state: string = '';
  export let name: string | undefined = '';
  export let infohash: string = '';
  export let bytescompleted: number | undefined = 0;
  export let bytesmissing: number | undefined = 0;
  export let length: number | undefined = 0;
  export let seeding: boolean | undefined = false;
  export let isTorrentPage: boolean = false;
  export let locked: boolean = false;
  export let category: string | undefined = '';

  let progpercentage = 0;
  let posterUrl = '';
  let fetchedFor = '';

  const posterCache: Map<string, string> = (globalThis as any).__posterCache ??= new Map();

  let refresh = () => {
    if (isTorrentPage === false) {
      Send({ command: 'listtorrents' });
    } else {
      Send({ command: 'listtorrentinfo', data1: infohash });
    }
  };

  $: progpercentage = (bytescompleted / length) * 100;

  $: stateLabel =
    state === 'active'
      ? seeding ? 'Seeding' : 'Downloading'
      : state === 'inactive'
        ? 'Stopped'
        : state === 'loading'
          ? 'Loading'
          : state === 'removed'
            ? 'Removed'
            : state;

  $: isActive = state === 'active' && !seeding;
  $: isSeeding = state === 'active' && seeding;

  function hashColor(s: string): string {
    let h = 0;
    for (let i = 0; i < s.length; i++) h = ((h << 5) - h + s.charCodeAt(i)) | 0;
    const hue = ((h % 360) + 360) % 360;
    return `hsl(${hue}, 30%, 15%)`;
  }

  async function fetchPoster(n: string, ih: string) {
    if (!n || !ih || fetchedFor === ih) return;
    fetchedFor = ih;

    if (posterCache.has(ih)) {
      posterUrl = posterCache.get(ih)!;
      return;
    }

    posterUrl = '';
    try {
      const resp = await fetch(`/api/poster?name=${encodeURIComponent(n)}`);
      if (resp.ok) {
        const data = await resp.json();
        let url = data.poster_url || '';
        if (!url) {
          const thumbResp = await fetch(`/api/thumbnail/${ih}`, { method: 'HEAD' });
          if (thumbResp.ok) {
            url = `/api/thumbnail/${ih}`;
          }
        }
        posterCache.set(ih, url);
        if (fetchedFor === ih) posterUrl = url;
      }
    } catch {}
  }

  $: if (infohash !== fetchedFor) posterUrl = '';
  $: fetchPoster(name, infohash);
</script>

<div class="group bg-white/5 rounded-lg overflow-hidden transition-all duration-200 hover:bg-white/[0.08] border border-white/10 hover:border-white/15 glass">
  <div class="flex">
    {#if isTorrentPage === false}
      <button
        type="button"
        class="relative flex-shrink-0 w-24 sm:w-32 bg-transparent border-none p-0 cursor-pointer"
        on:click={() => slocation.goto(`/torrent/${infohash}`)}>
        {#if posterUrl}
          <img src={posterUrl} alt="" class="w-full h-full object-cover min-h-[8rem] sm:min-h-[10rem]" />
        {:else}
          <div class="w-full h-full min-h-[8rem] sm:min-h-[10rem] flex items-center justify-center" style="background:{hashColor(infohash)}">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-slate-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a1 1 0 001-1V5a1 1 0 00-1-1H4a1 1 0 00-1 1v14a1 1 0 001 1z" />
            </svg>
          </div>
        {/if}
      </button>
    {:else}
      <div class="relative flex-shrink-0 w-24 sm:w-32">
        {#if posterUrl}
          <img src={posterUrl} alt="" class="w-full h-full object-cover min-h-[8rem] sm:min-h-[10rem]" />
        {:else}
          <div class="w-full h-full min-h-[8rem] sm:min-h-[10rem] flex items-center justify-center" style="background:{hashColor(infohash)}">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-slate-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a1 1 0 001-1V5a1 1 0 00-1-1H4a1 1 0 00-1 1v14a1 1 0 001 1z" />
            </svg>
          </div>
        {/if}
      </div>
    {/if}

    <div class="flex-1 min-w-0 p-3 sm:p-4 flex flex-col justify-between">
      <div>
        <div class="flex items-center justify-between gap-2 mb-1">
          {#if isTorrentPage === false && state !== 'removed'}
            <button
              type="button"
              class="text-slate-100 text-sm sm:text-base font-medium truncate text-left bg-transparent border-none p-0 cursor-pointer hover:text-violet-400 transition-colors duration-150"
              on:click={() => slocation.goto(`/torrent/${infohash}`)}>
              {name || 'Loading...'}
            </button>
          {:else}
            <h2 class="text-slate-100 text-sm sm:text-base font-medium truncate">{name || 'Loading...'}</h2>
          {/if}
          <div class="flex items-center gap-1.5 flex-shrink-0">
            {#if category}
              <span class="text-xs px-1.5 py-0.5 rounded bg-violet-900/40 text-violet-300 border border-violet-700/30">{category}</span>
            {/if}
            <span class="inline-flex items-center gap-1.5 text-xs px-2 py-0.5 rounded bg-white/10 text-slate-400">
              <span class="w-1.5 h-1.5 rounded-full {isSeeding ? 'bg-violet-400' : isActive ? 'bg-violet-500' : 'bg-slate-500'}"></span>
              {stateLabel}
            </span>
          </div>
        </div>

        <p class="text-slate-500 text-xs font-mono truncate mb-2">{infohash}</p>

        {#if state === 'active' || state === 'inactive'}
          <div class="space-y-1">
            <div class="bg-white/10 rounded-full overflow-hidden h-1">
              <div
                class="h-full rounded-full transition-all duration-500 ease-out bg-violet-500"
                style="width:{progpercentage ? progpercentage : 0}%"></div>
            </div>
            <div class="flex justify-between text-xs text-slate-500 tabular-nums">
              <span>{fileSize(bytescompleted)} / {fileSize(length)}</span>
              <span class="font-medium {progpercentage >= 100 ? 'text-violet-400' : 'text-slate-300'}">{progpercentage?.toLocaleString('en-US', { maximumFractionDigits: 1, minimumFractionDigits: 1 })}%</span>
            </div>
          </div>
        {/if}
      </div>

      <div class="flex items-center gap-1 mt-2 -mb-1">
        {#if state === 'active' || state === 'loading' || state === 'inactive'}
          <button type="button" aria-label="Remove" class="p-1.5 rounded transition-colors duration-150 text-slate-500 hover:text-violet-400 hover:bg-white/10"
            on:click={() => { Send({ command: 'removetorrent', data1: infohash, ...($adminmode === true && { aop: 1 }) }); }}>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        {:else if state === 'removed'}
          <button type="button" aria-label="Delete" class="p-1.5 rounded transition-colors duration-150 text-slate-500 hover:text-violet-400 hover:bg-white/10"
            on:click={() => { Send({ command: 'deletetorrent', data1: infohash, ...($adminmode === true && { aop: 1 }) }); if (isTorrentPage) { slocation.goto('/'); } else { refresh(); Send({ command: 'gettorrents' }); } }}>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
          </button>
        {/if}

        {#if state === 'loading'}
          <div class="p-1.5">
            <svg class="animate-spin h-4 w-4 text-violet-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
        {:else if state === 'active'}
          <button type="button" aria-label="Stop" class="p-1.5 rounded transition-colors duration-150 text-slate-500 hover:text-violet-400 hover:bg-white/10"
            on:click={() => { Send({ command: 'stoptorrent', data1: infohash, ...($adminmode === true && { aop: 1 }) }); refresh(); }}>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><rect x="6" y="6" width="12" height="12" rx="1" stroke-width="2" /></svg>
          </button>
        {:else if state === 'removed'}
          <button type="button" aria-label="Re-add" class="p-1.5 rounded transition-colors duration-150 text-slate-500 hover:text-violet-400 hover:bg-white/10"
            on:click={() => { Send({ command: 'addinfohash', data1: infohash }); refresh(); }}>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg>
          </button>
        {:else if state === 'inactive'}
          <button type="button" aria-label="Start" class="p-1.5 rounded transition-colors duration-150 text-slate-500 hover:text-violet-400 hover:bg-white/10"
            on:click={() => { Send({ command: 'starttorrent', data1: infohash, ...($adminmode === true && { aop: 1 }) }); refresh(); }}>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z" /></svg>
          </button>
        {/if}

        {#if isTorrentPage === false}
          <button type="button" aria-label="View" class="p-1.5 rounded transition-colors duration-150 text-slate-500 hover:text-violet-400 hover:bg-white/10"
            on:click={() => slocation.goto(`/torrent/${infohash}`)}>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
          </button>
        {:else}
          <button type="button" aria-label={locked ? 'Locked' : 'Unlocked'} class="p-1.5 rounded transition-colors duration-150 text-slate-500 hover:text-violet-400 hover:bg-white/10"
            on:click={() => { Send({ command: 'toggletorrentlock', data1: infohash }); setTimeout(() => { Send({ command: 'istorrentlocked', data1: infohash }); }, 1000); }}>
            {#if locked}
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-violet-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" /></svg>
            {:else}
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 11V7a4 4 0 118 0m-4 8v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z" /></svg>
            {/if}
          </button>
        {/if}

        {#if state === 'active' || state === 'inactive'}
          <span class="ml-auto text-xs text-slate-600 tabular-nums hidden sm:inline">{fileSize(bytesmissing == null ? 0 : bytesmissing)} left</span>
        {/if}
      </div>
    </div>
  </div>
</div>
