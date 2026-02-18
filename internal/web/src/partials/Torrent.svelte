<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Send, adminmode, torrentstats, isAdmin, usersfortorrent, torctime, torrentinfo, fileSize, fileType, torrentfiles, fileviewpath, fileviewinfohash, istrntlocked } from './core';
  import slocation from 'slocation';
  import TorrentCard from './TorrentCard.svelte';
  import CollapsibleSection from './CollapsibleSection.svelte';
  import type { DlObject } from './core';
  import ProgStat from './ProgStat.svelte';
  import { toast } from 'svelte-sonner';

  let infohash = $slocation.pathname?.split('/').reverse()[0];
  let sequentialMode = false;
  let categoryInput = '';
  let categoryOptions = ['Movies', 'TV Shows', 'Music', 'Games', 'Software', 'Other'];

  onMount(() => {
    torrentinfo.set({} as DlObject);
    Send({ command: 'listtorrentinfo', data1: infohash });
    Send({ command: 'gettorrentinfo', data1: infohash });
    Send({ command: 'istorrentlocked', data1: infohash });
    Send({ command: 'gettorrentstats', data1: infohash });
    Send({ command: 'gettorrentinfostat', data1: infohash });
  });

  onDestroy(() => {
    Send({ command: 'stopstream' });
    document.title = 'exatorrent';
  });

  $: if ($torrentinfo?.category) {
    categoryInput = $torrentinfo.category;
  }

  let toggleSequential = () => {
    sequentialMode = !sequentialMode;
    Send({ command: 'setsequential', data1: infohash, data2: sequentialMode ? 'true' : 'false' });
  };

  let setCategory = (cat: string) => {
    categoryInput = cat;
    Send({ command: 'setcategory', data1: infohash, data2: cat });
  };

  let filesOpen = false;
  let detailsOpen = false;
  let trntUsersOpen = false;
  let trackerfilestring = '';
  let trckrfileinput: HTMLInputElement;

  let readtracker = (e: Event) => {
    trackerfilestring = '';
    let f = (e.target as HTMLInputElement).files[0];
    if (f.size > 20971520) {
      toast.error('Error: Maximum Tracker File Size is 20MB');
      return;
    }
    let reader = new FileReader();
    reader.onload = (e) => {
      trackerfilestring = btoa(e.target.result as string);
      Send({
        command: 'addtrackerstotorrent',
        data1: infohash,
        data2: trackerfilestring,
        ...($adminmode === true && { aop: 1 })
      });
    };
    reader.readAsBinaryString(f);
    (e.target as HTMLInputElement).value = null;
  };

  let filesaction = () => {
    if (!filesOpen) {
      Send({ command: 'gettorrentfiles', data1: infohash });
      filesOpen = true;
    } else {
      filesOpen = false;
    }
  };

  let detailsaction = () => {
    if (!detailsOpen) {
      Send({ command: 'gettorrentstats', data1: infohash });
      Send({ command: 'gettorrentinfostat', data1: infohash });
      detailsOpen = true;
    } else {
      detailsOpen = false;
    }
  };

  let uwottaction = () => {
    if (!trntUsersOpen) {
      Send({ command: 'listusersfortorrent', data1: infohash, aop: 1 });
      trntUsersOpen = true;
    } else {
      trntUsersOpen = false;
    }
  };

  $: seedRatio = ($torrentstats?.BytesWrittenData && $torrentstats?.BytesReadData)
    ? ($torrentstats.BytesWrittenData / $torrentstats.BytesReadData).toLocaleString('en-US', { maximumFractionDigits: 2, minimumFractionDigits: 2 })
    : '—';
</script>

<svelte:head>
  <title>{$torrentinfo?.name} ({$torrentinfo?.infohash})</title>
</svelte:head>

<div class="mx-auto max-w-5xl px-2 sm:px-0">
  {#key $torrentinfo?.infohash}
    <TorrentCard state={$torrentinfo?.state} name={$torrentinfo?.name} infohash={$torrentinfo?.infohash} bytescompleted={$torrentinfo?.bytescompleted} bytesmissing={$torrentinfo?.bytesmissing} length={$torrentinfo?.length} seeding={$torrentinfo?.seeding} locked={$istrntlocked} isTorrentPage={true} category={$torrentinfo?.category} />
  {/key}

  <!-- Inline Stats -->
  {#if $torrentinfo?.state === 'active' || $torrentinfo?.state === 'inactive' || $torrentinfo?.state === 'loading'}
    <div class="grid grid-cols-3 gap-2 mx-3 mt-2">
      <div class="bg-white/5 glass border border-white/10 rounded-lg px-3 py-2 text-center">
        <p class="text-xs text-slate-500">Peers</p>
        <p class="text-sm text-slate-200 font-medium tabular-nums">{$torrentstats?.ActivePeers ?? '—'}<span class="text-slate-500 font-normal text-xs">/{$torrentstats?.TotalPeers ?? '—'}</span></p>
      </div>
      <div class="bg-white/5 glass border border-white/10 rounded-lg px-3 py-2 text-center">
        <p class="text-xs text-slate-500">Seeders</p>
        <p class="text-sm text-slate-200 font-medium tabular-nums">{$torrentstats?.ConnectedSeeders ?? '—'}</p>
      </div>
      <div class="bg-white/5 glass border border-white/10 rounded-lg px-3 py-2 text-center">
        <p class="text-xs text-slate-500">Ratio</p>
        <p class="text-sm text-slate-200 font-medium tabular-nums">{seedRatio}</p>
      </div>
    </div>
  {/if}

  <!-- Actions + Category + Sequential — all in one compact card -->
  <div class="bg-white/5 border border-white/10 glass rounded-xl m-3 overflow-hidden">
    <!-- Quick Actions Row -->
    <div class="flex items-center justify-around px-2 py-2 border-b border-white/[0.06]">
      {#if $adminmode === false}
        <button type="button" class="flex flex-col items-center gap-0.5 px-2 py-1 rounded-lg transition-colors duration-150 hover:bg-white/10 min-w-[3rem]"
          on:click={() => { Send({ command: 'abandontorrent', data1: infohash }); slocation.goto('/'); }}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7a4 4 0 11-8 0 4 4 0 018 0zM9 14a6 6 0 00-6 6v1h12v-1a6 6 0 00-6-6zM21 12h-6" /></svg>
          <span class="text-slate-500 text-[10px]">Abandon</span>
        </button>
      {/if}
      <a href="/api/torrent/{infohash}/?dl=tar" target="_blank" rel="noopener noreferrer" download
        class="flex flex-col items-center gap-0.5 px-2 py-1 rounded-lg transition-colors duration-150 hover:bg-white/10 min-w-[3rem]">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" /></svg>
        <span class="text-slate-500 text-[10px]">Tar</span>
      </a>
      <a href="/api/torrent/{infohash}/?dl=zip" target="_blank" rel="noopener noreferrer" download
        class="flex flex-col items-center gap-0.5 px-2 py-1 rounded-lg transition-colors duration-150 hover:bg-white/10 min-w-[3rem]">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" viewBox="0 0 20 20" fill="currentColor"><path d="M4 3a2 2 0 100 4h12a2 2 0 100-4H4z" /><path fill-rule="evenodd" d="M3 8h14v7a2 2 0 01-2 2H5a2 2 0 01-2-2V8zm5 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z" clip-rule="evenodd" /></svg>
        <span class="text-slate-500 text-[10px]">Zip</span>
      </a>
      <button type="button" class="flex flex-col items-center gap-0.5 px-2 py-1 rounded-lg transition-colors duration-150 hover:bg-white/10 min-w-[3rem]"
        on:click={() => { Send({ command: 'gettorrentmetainfo', data1: infohash }); }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
        <span class="text-slate-500 text-[10px]">.torrent</span>
      </button>
      <label class="flex flex-col items-center gap-0.5 px-2 py-1 rounded-lg transition-colors duration-150 hover:bg-white/10 min-w-[3rem] cursor-pointer">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" /></svg>
        <span class="text-slate-500 text-[10px]">Trackers</span>
        <input accept=".txt" bind:this={trckrfileinput} on:change={(e) => readtracker(e)} type="file" class="hidden" />
      </label>
    </div>

    <!-- Category Row -->
    <div class="flex items-center gap-2 px-4 py-2.5 border-b border-white/[0.06]">
      <span class="text-xs text-slate-500 flex-shrink-0">Category</span>
      <div class="flex items-center gap-1 flex-wrap flex-1 justify-end">
        {#each categoryOptions as cat}
          <button type="button"
            class="px-2 py-0.5 text-[11px] rounded transition-colors duration-150 border {categoryInput === cat ? 'bg-violet-600/20 border-violet-500/40 text-violet-300' : 'bg-white/5 border-white/[0.06] text-slate-500 hover:text-slate-300'}"
            on:click={() => setCategory(categoryInput === cat ? '' : cat)}>
            {cat}
          </button>
        {/each}
      </div>
    </div>

    <!-- Sequential Toggle Row -->
    {#if $torrentinfo?.state === 'active' || $torrentinfo?.state === 'inactive'}
      <div class="flex items-center justify-between px-4 py-2.5">
        <span class="text-xs text-slate-500">Sequential Download</span>
        <label class="relative inline-flex items-center cursor-pointer">
          <input type="checkbox" class="sr-only peer" bind:checked={sequentialMode} on:change={toggleSequential} />
          <div class="w-8 h-4 bg-slate-700 peer-focus:ring-2 peer-focus:ring-violet-500/40 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-3 after:w-3 after:transition-all peer-checked:bg-violet-600"></div>
        </label>
      </div>
    {/if}
  </div>

  <!-- Files (merged File Progress + Browse Files) -->
  {#if $torrentinfo?.state === 'active' || $torrentinfo?.state === 'inactive' || $torrentinfo?.state === 'removed'}
    <CollapsibleSection title="Files" bind:open={filesOpen} onToggle={filesaction} onRefresh={() => { Send({ command: 'gettorrentfiles', data1: infohash }); }}>
      <div class="space-y-1.5">
        {#each $torrentfiles as file (file.path)}
          {@const ft = fileType(file?.path || file?.displaypath || '')}
          <div class="text-slate-200 bg-white/5 rounded-lg border border-white/[0.06] px-3 py-2.5">
            <div class="flex items-center gap-2">
              <!-- File type icon -->
              <div class="flex-shrink-0 w-7 h-7 rounded-md bg-white/5 flex items-center justify-center">
                {#if ft === 'video'}
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-violet-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                {:else if ft === 'audio'}
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-violet-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" /></svg>
                {:else}
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" /></svg>
                {/if}
              </div>

              <!-- File name + size -->
              <div class="flex-1 min-w-0">
                <p class="text-sm truncate">{file?.displaypath}</p>
                <p class="text-xs text-slate-500 tabular-nums">{fileSize(file?.length)}</p>
              </div>

              <!-- Action buttons -->
              <div class="flex items-center gap-1 flex-shrink-0">
                {#if ft === 'video' || ft === 'audio'}
                  <button type="button" aria-label="Play" class="p-1.5 rounded-md hover:bg-white/10 transition-colors duration-150"
                    on:click={() => { fileviewpath.set(file?.path); fileviewinfohash.set(infohash); slocation.goto('/file'); }}>
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-violet-400" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z" /></svg>
                  </button>
                  <a href="vlc://{location.origin}/api/torrent/{infohash}/{file?.path}" target="_blank" rel="noopener noreferrer"
                    class="p-1.5 rounded-md hover:bg-white/10 transition-colors duration-150" aria-label="Open in VLC">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" /></svg>
                  </a>
                {:else}
                  <a href="/api/torrent/{infohash}/{file?.path}" target="_blank" rel="noopener noreferrer" download
                    class="p-1.5 rounded-md hover:bg-white/10 transition-colors duration-150" aria-label="Download">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
                  </a>
                {/if}

                {#if file?.priority === 1}
                  <button type="button" aria-label="Pause file" class="p-1.5 rounded-md hover:bg-white/10 transition-colors duration-150"
                    on:click={() => { Send({ command: 'stopfile', data1: infohash, data2: file?.path, ...($adminmode === true && { aop: 1 }) }); setTimeout(() => { Send({ command: 'gettorrentfiles', data1: infohash }); }, 1000); }}>
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6" /></svg>
                  </button>
                {:else if file?.priority === 0}
                  <button type="button" aria-label="Resume file" class="p-1.5 rounded-md hover:bg-white/10 transition-colors duration-150"
                    on:click={() => { Send({ command: 'startfile', data1: infohash, data2: file?.path, ...($adminmode === true && { aop: 1 }) }); setTimeout(() => { Send({ command: 'gettorrentfiles', data1: infohash }); }, 1000); }}>
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-violet-400" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z" /></svg>
                  </button>
                {/if}

                {#if $torrentinfo?.state === 'removed'}
                  <button type="button" aria-label="Delete file" class="p-1.5 rounded-md hover:bg-white/10 transition-colors duration-150"
                    on:click={() => { Send({ command: 'deletefilepath', data1: infohash, data2: file?.path, ...($adminmode === true && { aop: 1 }) }); }}>
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                  </button>
                {/if}
              </div>
            </div>

            {#if $torrentinfo?.state !== 'removed'}
              <ProgStat bytescompleted={file?.bytescompleted} length={file?.length} />
            {/if}
          </div>
        {/each}
      </div>
    </CollapsibleSection>
  {/if}

  <!-- Details (deep stats) -->
  {#if $torrentinfo?.state === 'active' || $torrentinfo?.state === 'inactive' || $torrentinfo?.state === 'loading'}
    <CollapsibleSection title="Details" bind:open={detailsOpen} onToggle={detailsaction} onRefresh={() => { Send({ command: 'gettorrentstats', data1: infohash }); Send({ command: 'gettorrentinfostat', data1: infohash }); }}>
      <div class="space-y-3">
        <div class="flex flex-wrap gap-x-6 gap-y-1 text-xs text-slate-500">
          <span>Added {$torctime.addedat}</span>
          {#if $torrentinfo?.state === 'active'}
            <span>Started {$torctime.startedat}</span>
          {/if}
        </div>

        <div class="grid grid-cols-2 gap-2">
          <div class="bg-white/5 rounded-lg border border-white/[0.06] px-3 py-2">
            <p class="text-xs text-slate-500 mb-0.5">Written</p>
            <p class="text-sm text-slate-200 font-medium tabular-nums">{fileSize($torrentstats?.BytesWrittenData)}</p>
          </div>
          <div class="bg-white/5 rounded-lg border border-white/[0.06] px-3 py-2">
            <p class="text-xs text-slate-500 mb-0.5">Read</p>
            <p class="text-sm text-slate-200 font-medium tabular-nums">{fileSize($torrentstats?.BytesReadUsefulData)}</p>
          </div>
          <div class="bg-white/5 rounded-lg border border-white/[0.06] px-3 py-2">
            <p class="text-xs text-slate-500 mb-0.5">Pieces</p>
            <p class="text-sm text-slate-200 font-medium tabular-nums">{$torrentstats?.PiecesDirtiedGood ?? '—'} <span class="text-slate-500 font-normal">good</span> / {$torrentstats?.PiecesDirtiedBad ?? '0'} <span class="text-slate-500 font-normal">bad</span></p>
          </div>
          <div class="bg-white/5 rounded-lg border border-white/[0.06] px-3 py-2">
            <p class="text-xs text-slate-500 mb-0.5">Half-Open</p>
            <p class="text-sm text-slate-200 font-medium tabular-nums">{$torrentstats?.HalfOpenPeers ?? '—'}</p>
          </div>
        </div>

        {#if $isAdmin === true && $adminmode === true}
          <div class="grid grid-cols-2 gap-2">
            <button class="bg-white/5 rounded-lg border border-white/[0.06] py-2 px-3 text-xs text-slate-400 hover:bg-white/10 transition-colors duration-150"
              on:click={() => { Send({ command: 'changedataload', data1: infohash, data2: 'upload', data3: 'allow', aop: 1 }); }}>Allow Upload</button>
            <button class="bg-white/5 rounded-lg border border-white/[0.06] py-2 px-3 text-xs text-slate-400 hover:bg-white/10 transition-colors duration-150"
              on:click={() => { Send({ command: 'changedataload', data1: infohash, data2: 'upload', data3: 'disallow', aop: 1 }); }}>Disallow Upload</button>
            <button class="bg-white/5 rounded-lg border border-white/[0.06] py-2 px-3 text-xs text-slate-400 hover:bg-white/10 transition-colors duration-150"
              on:click={() => { Send({ command: 'changedataload', data1: infohash, data2: 'download', data3: 'allow', aop: 1 }); }}>Allow Download</button>
            <button class="bg-white/5 rounded-lg border border-white/[0.06] py-2 px-3 text-xs text-slate-400 hover:bg-white/10 transition-colors duration-150"
              on:click={() => { Send({ command: 'changedataload', data1: infohash, data2: 'download', data3: 'disallow', aop: 1 }); }}>Disallow Download</button>
          </div>
        {/if}
      </div>
    </CollapsibleSection>
  {/if}

  <!-- Owners (admin-only) -->
  {#if $isAdmin === true}
    <CollapsibleSection title="Owners" bind:open={trntUsersOpen} onToggle={uwottaction} onRefresh={() => { Send({ command: 'listusersfortorrent', data1: infohash, aop: 1 }); }}>
      <div class="space-y-2">
        {#each $usersfortorrent as eachuser (eachuser)}
          <div class="text-slate-200 bg-white/5 rounded-xl border border-white/[0.06] px-3 py-3">
            <div class="flex items-center justify-between py-1">
              <p class="font-medium break-all flex-1 min-w-0">{eachuser}</p>
              <button type="button" aria-label="Remove user" class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-150 flex-shrink-0 ml-2"
                on:click={() => { Send({ command: 'abandontorrent', data1: infohash, data2: eachuser, aop: 1 }); }}>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
          </div>
        {/each}
      </div>
    </CollapsibleSection>
  {/if}
</div>
