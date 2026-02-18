<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { downloadslist, Send, adminmode, isAdmin, terrormsg } from './core';
  import type { DlObject } from './core';
  import TorrentCard from './TorrentCard.svelte';

  onMount(() => {
    if ($adminmode === false) {
      Send({ command: 'listtorrents' });
      Send({ command: 'gettorrents' });
    } else {
      Send({ command: 'listalltorrents', aop: 1 });
      Send({ command: 'getalltorrents', aop: 1 });
    }
  });

  onDestroy(() => {
    console.log('on destroy');
    Send({ command: 'stopstream' });
  });

  let checam = (b: boolean) => {
    b ? (Send({ command: 'listalltorrents', aop: 1 }), Send({ command: 'getalltorrents', aop: 1 })) : Send({ command: 'gettorrents' });
  };

  $: checam($adminmode);

  let searchQuery = '';
  let stateFilter = 'all';
  let categoryFilter = 'all';
  let sortBy = 'name';

  let filteredList: DlObject[] = [];
  let availableCategories: string[] = [];

  $: {
    const allList = $downloadslist?.data || [];
    const cats = new Set<string>();
    allList.forEach((dl) => { if (dl.category) cats.add(dl.category); });
    availableCategories = [...cats].sort();
  }

  $: {
    let list = $downloadslist?.data || [];

    if (searchQuery.trim()) {
      const q = searchQuery.toLowerCase();
      list = list.filter(
        (dl) =>
          (dl.name || '').toLowerCase().includes(q) ||
          dl.infohash.toLowerCase().includes(q)
      );
    }

    if (stateFilter !== 'all') {
      if (stateFilter === 'seeding') {
        list = list.filter((dl) => dl.seeding === true);
      } else {
        list = list.filter((dl) => dl.state === stateFilter);
      }
    }

    if (categoryFilter !== 'all') {
      if (categoryFilter === 'uncategorized') {
        list = list.filter((dl) => !dl.category);
      } else {
        list = list.filter((dl) => dl.category === categoryFilter);
      }
    }

    list = [...list].sort((a, b) => {
      if (sortBy === 'name') {
        return (a.name || '').localeCompare(b.name || '');
      } else if (sortBy === 'progress') {
        const pa = a.length ? (a.bytescompleted || 0) / a.length : 0;
        const pb = b.length ? (b.bytescompleted || 0) / b.length : 0;
        return pb - pa;
      } else if (sortBy === 'size') {
        return (b.length || 0) - (a.length || 0);
      }
      return 0;
    });

    filteredList = list;
  }

  let totalCount: number;
  $: totalCount = ($downloadslist?.data || []).length;
</script>

<div class="mx-auto max-w-7xl px-3 sm:px-4 lg:px-6">
  {#if $terrormsg?.has === false}
    {#if $isAdmin === true && $adminmode === true}
      <div class="flex items-center justify-between bg-slate-900 rounded-lg mt-4 px-4 py-3 border border-slate-700/40">
        <span class="text-slate-300 text-sm font-medium">Admin Mode</span>
        <label class="relative inline-flex items-center cursor-pointer">
          <input type="checkbox" class="sr-only peer" bind:checked={$adminmode} />
          <div class="w-9 h-5 bg-slate-700 peer-focus:ring-2 peer-focus:ring-violet-500/40 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-violet-600"></div>
        </label>
      </div>
    {/if}

    <div class="mt-4 space-y-3">
      <div class="flex gap-2">
        <div class="flex-1 relative">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            type="text"
            placeholder="Search torrents..."
            bind:value={searchQuery}
            class="w-full bg-slate-900 border border-slate-700/40 rounded-lg pl-9 pr-3 py-2.5 text-sm text-slate-200 placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-violet-500/50 focus:border-violet-500/50 transition-all duration-150" />
        </div>
      </div>
      <div class="flex gap-2 items-center flex-wrap">
        <select
          bind:value={stateFilter}
          class="bg-slate-900 border border-slate-700/40 rounded-lg px-3 py-2 text-sm text-slate-300 focus:outline-none focus:ring-1 focus:ring-violet-500/50 transition-all duration-150 appearance-none cursor-pointer">
          <option value="all">All States</option>
          <option value="active">Active</option>
          <option value="inactive">Stopped</option>
          <option value="seeding">Seeding</option>
          <option value="loading">Loading</option>
          <option value="removed">Removed</option>
        </select>
        {#if availableCategories.length > 0}
          <select
            bind:value={categoryFilter}
            class="bg-slate-900 border border-slate-700/40 rounded-lg px-3 py-2 text-sm text-slate-300 focus:outline-none focus:ring-1 focus:ring-violet-500/50 transition-all duration-150 appearance-none cursor-pointer">
            <option value="all">All Categories</option>
            <option value="uncategorized">Uncategorized</option>
            {#each availableCategories as cat}
              <option value={cat}>{cat}</option>
            {/each}
          </select>
        {/if}
        <select
          bind:value={sortBy}
          class="bg-slate-900 border border-slate-700/40 rounded-lg px-3 py-2 text-sm text-slate-300 focus:outline-none focus:ring-1 focus:ring-violet-500/50 transition-all duration-150 appearance-none cursor-pointer">
          <option value="name">Sort: Name</option>
          <option value="progress">Sort: Progress</option>
          <option value="size">Sort: Size</option>
        </select>
        <span class="text-slate-500 text-xs ml-auto tabular-nums">
          {filteredList.length}{filteredList.length !== totalCount ? ` / ${totalCount}` : ''} torrent{totalCount !== 1 ? 's' : ''}
        </span>
      </div>
    </div>

    {#if filteredList.length === 0 && totalCount > 0}
      <p class="text-sm text-center text-slate-500 mt-12">No torrents match your filter</p>
    {/if}

    <div class="mt-4 grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-3">
      {#each filteredList as dl (dl.infohash)}
        <TorrentCard state={dl.state} name={dl?.name} infohash={dl.infohash} bytescompleted={dl?.bytescompleted} bytesmissing={dl?.bytesmissing} length={dl?.length} seeding={dl?.seeding} category={dl?.category} />
      {/each}
    </div>
  {:else}
    <p class="text-xl text-center text-slate-400 font-sans mt-8">{$terrormsg?.msg}</p>
  {/if}
</div>
