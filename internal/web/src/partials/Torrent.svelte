<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Send, adminmode, torrentstats, isAdmin, usersfortorrent, torctime, torrentinfo, fileSize, fsdirinfo, torrentfiles, fileviewpath, fileviewinfohash, istrntlocked } from './core';
  import slocation from 'slocation';
  import TorrentCard from './TorrentCard.svelte';
  import CollapsibleSection from './CollapsibleSection.svelte';
  import type { DlObject } from './core';
  import ProgStat from './ProgStat.svelte';
  import { toast } from 'svelte-sonner';

  onMount(() => {
    torrentinfo.set({} as DlObject);
    Send({ command: 'listtorrentinfo', data1: infohash });
    Send({ command: 'gettorrentinfo', data1: infohash });
    Send({ command: 'istorrentlocked', data1: infohash });
  });

  onDestroy(() => {
    Send({ command: 'stopstream' });
    document.title = 'exatorrent';
  });

  let infohash = $slocation.pathname?.split('/').reverse()[0];
  let sequentialMode = false;
  let categoryInput = '';
  let categoryOptions = ['Movies', 'TV Shows', 'Music', 'Games', 'Software', 'Other'];

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

  let fileProgressOpen = false;
  let browseFilesOpen = false;
  let detailsOpen = false;
  let trntUsersOpen = false;
  let detailsLoaded = false;
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

  let fileProgressaction = () => {
    if (fileProgressOpen === false) {
      Send({ command: 'gettorrentfiles', data1: infohash });
      fileProgressOpen = true;
    } else {
      fileProgressOpen = false;
    }
  };

  let browseFilesaction = () => {
    if (browseFilesOpen === false) {
      Send({ command: 'getfsdirinfo', data1: infohash });
      browseFilesOpen = true;
    } else {
      browseFilesOpen = false;
    }
  };

  let detailsaction = () => {
    if (detailsOpen === false) {
      if (!detailsLoaded) {
        Send({ command: 'gettorrentstats', data1: infohash });
        Send({ command: 'gettorrentinfostat', data1: infohash });
        detailsLoaded = true;
      }
      detailsOpen = true;
    } else {
      detailsOpen = false;
    }
  };

  let uwottaction = () => {
    if (trntUsersOpen === false) {
      Send({ command: 'listusersfortorrent', data1: infohash, aop: 1 });
      trntUsersOpen = true;
    } else {
      trntUsersOpen = false;
    }
  };
</script>

<svelte:head>
  <title>{$torrentinfo?.name} ({$torrentinfo?.infohash})</title>
</svelte:head>

<div class="mx-auto max-w-5xl px-2 sm:px-0">
  <TorrentCard state={$torrentinfo?.state} name={$torrentinfo?.name} infohash={$torrentinfo?.infohash} bytescompleted={$torrentinfo?.bytescompleted} bytesmissing={$torrentinfo?.bytesmissing} length={$torrentinfo?.length} seeding={$torrentinfo?.seeding} locked={$istrntlocked} isTorrentPage={true} category={$torrentinfo?.category} />

  <!-- Quick Actions -->
  <div class="bg-white/5 border border-white/10 glass flex items-center justify-around text-white rounded-xl m-3 py-2 px-2">
    {#if $adminmode === false}
      <button
        type="button"
        class="flex flex-col items-center gap-0.5 px-3 py-1.5 rounded-lg transition-colors duration-150 hover:bg-white/10 active:bg-white/15 min-w-[3.5rem]"
        on:click={() => {
          Send({ command: 'abandontorrent', data1: infohash });
          slocation.goto('/');
        }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7a4 4 0 11-8 0 4 4 0 018 0zM9 14a6 6 0 00-6 6v1h12v-1a6 6 0 00-6-6zM21 12h-6" />
        </svg>
        <span class="text-slate-400 text-xs">Abandon</span>
      </button>
    {/if}

    <a href="/api/torrent/{infohash}/?dl=tar" target="_blank" rel="noopener noreferrer" download
       class="flex flex-col items-center gap-0.5 px-3 py-1.5 rounded-lg transition-colors duration-150 hover:bg-white/10 active:bg-white/15 min-w-[3.5rem]">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
      </svg>
      <span class="text-slate-400 text-xs">Tar</span>
    </a>

    <a href="/api/torrent/{infohash}/?dl=zip" target="_blank" rel="noopener noreferrer" download
       class="flex flex-col items-center gap-0.5 px-3 py-1.5 rounded-lg transition-colors duration-150 hover:bg-white/10 active:bg-white/15 min-w-[3.5rem]">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400" viewBox="0 0 20 20" fill="currentColor">
        <path d="M4 3a2 2 0 100 4h12a2 2 0 100-4H4z" />
        <path fill-rule="evenodd" d="M3 8h14v7a2 2 0 01-2 2H5a2 2 0 01-2-2V8zm5 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z" clip-rule="evenodd" />
      </svg>
      <span class="text-slate-400 text-xs">Zip</span>
    </a>
  </div>

  <!-- Category & Sequential -->
  <div class="bg-white/5 border border-white/10 glass rounded-xl m-3 px-4 py-3 space-y-3">
    <div class="flex items-center justify-between gap-3">
      <span class="text-sm text-slate-400 flex-shrink-0">Category</span>
      <div class="flex items-center gap-1.5 flex-wrap justify-end">
        {#each categoryOptions as cat}
          <button
            type="button"
            class="px-2.5 py-1 text-xs rounded-md transition-colors duration-150 border {categoryInput === cat ? 'bg-violet-600/20 border-violet-500/40 text-violet-300' : 'bg-white/5 border-white/[0.06] text-slate-400 hover:border-white/15 hover:text-slate-300'}"
            on:click={() => setCategory(categoryInput === cat ? '' : cat)}>
            {cat}
          </button>
        {/each}
      </div>
    </div>

    {#if $torrentinfo?.state === 'active' || $torrentinfo?.state === 'inactive'}
      <div class="flex items-center justify-between border-t border-white/[0.06] pt-3">
        <div>
          <span class="text-sm text-slate-400">Sequential Download</span>
          <p class="text-xs text-slate-600 mt-0.5">Downloads pieces in order for streaming</p>
        </div>
        <label class="relative inline-flex items-center cursor-pointer">
          <input type="checkbox" class="sr-only peer" bind:checked={sequentialMode} on:change={toggleSequential} />
          <div class="w-9 h-5 bg-slate-700 peer-focus:ring-2 peer-focus:ring-violet-500/40 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-violet-600"></div>
        </label>
      </div>
    {/if}
  </div>

  <!-- File Progress -->
  {#if $torrentinfo?.state === 'active' || $torrentinfo?.state === 'inactive'}
    <CollapsibleSection title="File Progress" bind:open={fileProgressOpen} onToggle={fileProgressaction} onRefresh={() => { Send({ command: 'gettorrentfiles', data1: infohash }); }}>
      <div class="space-y-2">
        {#each $torrentfiles as file (file.path)}
          <div class="text-slate-200 bg-white/5 px-3 py-3 rounded-xl border border-white/[0.06]">
            <div class="flex items-center justify-between gap-2 py-1">
              <button
                class="flex-1 min-w-0 text-left"
                on:click={() => {
                  fileviewpath.set(file?.path);
                  fileviewinfohash.set(infohash);
                  slocation.goto('/file');
                }}>
                <p class="font-medium text-sm truncate hover:text-violet-400 transition-colors duration-150">{file?.displaypath}</p>
              </button>

              {#if file?.priority === 1}
                <button
                  aria-label="Stop file"
                  class="p-2 rounded-lg bg-white/5 hover:bg-white/10 flex-shrink-0 transition-colors duration-150"
                  on:click={() => {
                    Send({ command: 'stopfile', data1: infohash, data2: file?.path, ...($adminmode === true && { aop: 1 }) });
                    setTimeout(() => { Send({ command: 'gettorrentfiles', data1: infohash }); }, 1000);
                  }}>
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </button>
              {:else if file?.priority === 0}
                <button
                  aria-label="Start file"
                  class="p-2 rounded-lg bg-white/5 hover:bg-white/10 flex-shrink-0 transition-colors duration-150"
                  on:click={() => {
                    Send({ command: 'startfile', data1: infohash, data2: file?.path, ...($adminmode === true && { aop: 1 }) });
                    setTimeout(() => { Send({ command: 'gettorrentfiles', data1: infohash }); }, 1000);
                  }}>
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-violet-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </button>
              {/if}
            </div>
            <ProgStat bytescompleted={file?.bytescompleted} length={file?.length} offset={file?.offset} />
          </div>
        {/each}
      </div>
    </CollapsibleSection>
  {/if}

  <!-- Browse Files -->
  <CollapsibleSection title="Browse Files" bind:open={browseFilesOpen} onToggle={browseFilesaction} onRefresh={() => { Send({ command: 'getfsdirinfo', data1: infohash }); }}>
      <div class="space-y-2">
        {#each $fsdirinfo as file (file.path)}
          <div class="text-slate-200 bg-white/5 rounded-xl border border-white/[0.06] px-3 py-3">
            <div class="flex flex-col justify-between flex-wrap py-1">
              <button
                type="button"
                class="w-full bg-transparent border-none p-0 text-left text-inherit cursor-pointer"
                on:click={() => {
                  if (file?.isdir === true) {
                    Send({ command: 'getfsdirinfo', data1: infohash, data2: file?.path });
                  } else if (file?.isdir === false) {
                    fileviewpath.set(file?.path);
                    fileviewinfohash.set(infohash);
                    slocation.goto('/file');
                  }
                }}>
                <p class="font-medium break-all text-left">{file?.name}</p>
              </button>

              <div class="grid grid-flow-col">
                <p class="py-2 text-sm text-slate-300 break-all">
                  {#if file?.isdir === false}{fileSize(file?.size)}{:else}Directory{/if}
                </p>

                <div class="flex flex-row justify-end gap-1 flex-wrap py-1">
                  {#if file?.isdir === true}
                    <a href="/api/torrent/{infohash}/{file?.path}/?dl=zip" target="_blank" rel="noopener noreferrer" aria-label="Download as zip" class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-150 noHL" download>
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" viewBox="0 0 20 20" fill="currentColor">
                        <path d="M4 3a2 2 0 100 4h12a2 2 0 100-4H4z" />
                        <path fill-rule="evenodd" d="M3 8h14v7a2 2 0 01-2 2H5a2 2 0 01-2-2V8zm5 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z" clip-rule="evenodd" />
                      </svg>
                    </a>

                    <a href="/api/torrent/{infohash}/{file?.path}/?dl=tar" target="_blank" rel="noopener noreferrer" aria-label="Download as tar" class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-150 noHL" download>
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
                      </svg>
                    </a>
                  {/if}

                  {#if $torrentinfo?.state === 'removed'}
                    <button
                      type="button"
                      aria-label="Delete file"
                      class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-150"
                      on:click={() => {
                        Send({
                          command: 'deletefilepath',
                          data1: infohash,
                          data2: file?.path,
                          ...($adminmode === true && { aop: 1 })
                        });
                      }}>
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                    </button>
                  {/if}
                </div>
              </div>
            </div>
          </div>
        {/each}
      </div>
  </CollapsibleSection>

  <!-- Details (merged Stats + Misc) -->
  {#if $torrentinfo?.state === 'active' || $torrentinfo?.state === 'inactive' || $torrentinfo?.state === 'loading'}
    <CollapsibleSection title="Details" bind:open={detailsOpen} onToggle={detailsaction} onRefresh={() => { Send({ command: 'gettorrentstats', data1: infohash }); Send({ command: 'gettorrentinfostat', data1: infohash }); }}>
        <div class="space-y-3">
          <!-- Timestamps -->
          <div class="flex flex-wrap gap-x-6 gap-y-1 text-xs text-slate-500">
            <span>Added {$torctime.addedat}</span>
            {#if $torrentinfo?.state === 'active'}
              <span>Started {$torctime.startedat}</span>
            {/if}
          </div>

          <!-- Stats Grid -->
          <div class="grid grid-cols-2 gap-2">
            <div class="bg-white/5 rounded-lg border border-white/[0.06] px-3 py-2.5">
              <p class="text-xs text-slate-500 mb-0.5">Active Peers</p>
              <p class="text-sm text-slate-200 font-medium tabular-nums">{$torrentstats?.ActivePeers ?? '—'} <span class="text-slate-500 font-normal">/ {$torrentstats?.TotalPeers ?? '—'} total</span></p>
            </div>
            <div class="bg-white/5 rounded-lg border border-white/[0.06] px-3 py-2.5">
              <p class="text-xs text-slate-500 mb-0.5">Connected Seeders</p>
              <p class="text-sm text-slate-200 font-medium tabular-nums">{$torrentstats?.ConnectedSeeders ?? '—'}</p>
            </div>
            <div class="bg-white/5 rounded-lg border border-white/[0.06] px-3 py-2.5">
              <p class="text-xs text-slate-500 mb-0.5">Written</p>
              <p class="text-sm text-slate-200 font-medium tabular-nums">{fileSize($torrentstats?.BytesWrittenData)}</p>
            </div>
            <div class="bg-white/5 rounded-lg border border-white/[0.06] px-3 py-2.5">
              <p class="text-xs text-slate-500 mb-0.5">Read</p>
              <p class="text-sm text-slate-200 font-medium tabular-nums">{fileSize($torrentstats?.BytesReadUsefulData)}</p>
            </div>
            <div class="bg-white/5 rounded-lg border border-white/[0.06] px-3 py-2.5">
              <p class="text-xs text-slate-500 mb-0.5">Seed Ratio</p>
              <p class="text-sm text-slate-200 font-medium tabular-nums">{($torrentstats?.BytesWrittenData && $torrentstats?.BytesReadData) ? ($torrentstats.BytesWrittenData / $torrentstats.BytesReadData).toLocaleString('en-US', { maximumFractionDigits: 3, minimumFractionDigits: 3 }) : '—'}</p>
            </div>
            <div class="bg-white/5 rounded-lg border border-white/[0.06] px-3 py-2.5">
              <p class="text-xs text-slate-500 mb-0.5">Pieces</p>
              <p class="text-sm text-slate-200 font-medium tabular-nums">{$torrentstats?.PiecesDirtiedGood ?? '—'} <span class="text-slate-500 font-normal">good</span> / {$torrentstats?.PiecesDirtiedBad ?? '0'} <span class="text-slate-500 font-normal">bad</span></p>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex flex-wrap gap-2">
            <button
              class="flex-1 min-w-[140px] bg-white/5 rounded-lg border border-white/[0.06] py-2.5 px-4 text-sm text-slate-300 hover:bg-white/10 transition-colors duration-150"
              on:click={() => { Send({ command: 'gettorrentmetainfo', data1: infohash }); }}>
              Download .torrent
            </button>
            <label class="flex-1 min-w-[140px] appearance-none block cursor-pointer">
              <div class="flex items-center justify-center gap-2 bg-white/5 rounded-lg border border-white/[0.06] py-2.5 px-4 text-sm text-slate-300 hover:bg-white/10 transition-colors duration-150">
                <svg xmlns="http://www.w3.org/2000/svg" class="text-slate-400 h-4 w-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
                </svg>
                Add Trackers
              </div>
              <input accept=".txt" bind:this={trckrfileinput} on:change={(e) => readtracker(e)} id="torrentfile" name="torrentfile" type="file" class="hidden" />
            </label>
          </div>

          <!-- Admin data toggles -->
          {#if $isAdmin === true && $adminmode === true}
            <div class="grid grid-cols-2 gap-2">
              <button
                class="bg-white/5 rounded-lg border border-white/[0.06] py-2 px-3 text-xs text-slate-400 hover:bg-white/10 transition-colors duration-150"
                on:click={() => { Send({ command: 'changedataload', data1: infohash, data2: 'upload', data3: 'allow', aop: 1 }); }}>Allow Upload</button>
              <button
                class="bg-white/5 rounded-lg border border-white/[0.06] py-2 px-3 text-xs text-slate-400 hover:bg-white/10 transition-colors duration-150"
                on:click={() => { Send({ command: 'changedataload', data1: infohash, data2: 'upload', data3: 'disallow', aop: 1 }); }}>Disallow Upload</button>
              <button
                class="bg-white/5 rounded-lg border border-white/[0.06] py-2 px-3 text-xs text-slate-400 hover:bg-white/10 transition-colors duration-150"
                on:click={() => { Send({ command: 'changedataload', data1: infohash, data2: 'download', data3: 'allow', aop: 1 }); }}>Allow Download</button>
              <button
                class="bg-white/5 rounded-lg border border-white/[0.06] py-2 px-3 text-xs text-slate-400 hover:bg-white/10 transition-colors duration-150"
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
            <div class="flex items-center justify-between flex-wrap py-1">
              <div class="w-0 flex-1 flex">
                <p class="font-medium break-all">{eachuser}</p>
              </div>

              <button
                type="button"
                aria-label="Remove user"
                class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-150 flex-shrink-0"
                on:click={() => {
                  Send({
                    command: 'abandontorrent',
                    data1: infohash,
                    data2: eachuser,
                    aop: 1
                  });
                }}>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>
        {/each}
      </div>
    </CollapsibleSection>
  {/if}
</div>
