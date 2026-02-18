<script lang="ts">
  import {
    dontstart,
    Send,
    SignOut,
    socket,
    adminmode,
    isAdmin,
    engconfig,
    fileSize,
    diskstats,
    nooftrackersintrackerdb
  } from './core';
  import slocation from 'slocation';
  import { toast } from 'svelte-sonner';

  let ds = false;
  if ($dontstart === 'true') {
    ds = true;
  }

  let toggledontstart = (ds: boolean) => {
    console.log('dontstart changed to', ds);
    ds ? (localStorage.setItem('dontstart', 'true'), ($dontstart = 'true')) : (localStorage.setItem('dontstart', 'false'), ($dontstart = 'false'));
  };

  $: toggledontstart(ds);

  let editmode: boolean = false;
  let newpassword: string = '';

  let diskstatsOpen = false;
  let miscOpen = false;
  let miscfirsttime = true;

  let trdelno: string;
  let srno: string;

  let changepw = () => {
    if (newpassword.length < 5) {
      toast.error('Size of New Password must be more than 5');
      return;
    }
    Send({
      command: 'updatepw',
      data1: newpassword
    });

    socket?.readyState === WebSocket.OPEN ? socket?.close() : console.log('socket already closed');
    SignOut();
  };

  let engsettingsOpen = false;
  let engsettingsstring = '';

  let whenengconfigChange = (ec: Object) => {
    engsettingsstring = JSON.stringify(ec, null, 2);
  };

  $: whenengconfigChange($engconfig);

  let diskusageaction = () => {
    if (diskstatsOpen === false) {
      Send({ command: 'diskusage' });
      diskstatsOpen = true;
    } else {
      diskstatsOpen = false;
    }
  };

  let miscsettingsaction = () => {
    miscOpen = !miscOpen;
    if (miscfirsttime === true) {
      Send({ command: 'nooftrackersintrackerdb', aop: 1 });
    }
  };

  let enginesettingsOpen = () => {
    if (engsettingsOpen === false) {
      Send({ command: 'getconfig', aop: 1 });
      engsettingsOpen = true;
    } else {
      engsettingsOpen = false;
    }
  };
</script>

<div class="mx-auto max-w-2xl px-2 sm:px-0">
  <!-- User Settings Card -->
  <div class="bg-white/5 glass rounded-xl m-3 overflow-hidden border border-white/10">
    <div class="flex items-center justify-between px-4 py-3 border-b border-white/[0.06]">
      <div class="flex items-center gap-2">
        <h3 class="font-medium text-slate-200">User Settings</h3>
        {#if localStorage.getItem('exausertype') === 'admin'}
          <span class="text-xs font-medium px-2 py-0.5 rounded-full bg-violet-600/20 text-violet-400">admin</span>
        {/if}
      </div>
    </div>

    <div class="p-4 space-y-3">
      <div class="flex bg-white/5 rounded-xl border border-white/[0.06] overflow-hidden">
        {#if editmode === false}
          <input id="username" name="username" type="text"
                 class="bg-transparent w-full flex-grow px-3.5 py-2.5 text-sm border-none placeholder-slate-400 text-slate-100 focus:outline-none"
                 placeholder={localStorage.getItem('exausername')} disabled />
          <button
            type="button"
            aria-label="Edit password"
            class="flex items-center justify-center w-11 flex-shrink-0 transition-colors duration-150 hover:bg-white/10"
            on:click={() => { editmode = true; }}>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
            </svg>
          </button>
        {:else}
          <input id="password" name="password" type="text"
                 class="bg-transparent w-full flex-grow px-3.5 py-2.5 text-sm border-none placeholder-slate-500 text-slate-200 focus:outline-none"
                 placeholder="Enter New Password" bind:value={newpassword} />
          <button
            type="button"
            aria-label="Cancel"
            class="flex items-center justify-center w-10 flex-shrink-0 transition-colors duration-150 hover:bg-white/10"
            on:click={() => { editmode = false; }}>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
          <button type="button" aria-label="Save password" class="flex items-center justify-center w-10 flex-shrink-0 transition-colors duration-150 hover:bg-white/10" on:click={changepw}>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5 text-violet-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </button>
        {/if}
      </div>

      <div class="flex items-center justify-between bg-white/5 rounded-xl border border-white/[0.06] px-4 py-3">
        <span class="text-sm text-slate-300">Don't Start Torrents on Add</span>
        <label class="relative inline-flex items-center cursor-pointer">
          <input type="checkbox" class="sr-only peer" bind:checked={ds} />
          <div class="w-9 h-5 bg-slate-600 peer-focus:ring-2 peer-focus:ring-violet-500/40 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-violet-600"></div>
        </label>
      </div>

      {#if $isAdmin === true}
        <div class="flex items-center justify-between bg-white/5 rounded-xl border border-white/[0.06] px-4 py-3">
          <span class="text-sm text-slate-300">Admin Mode</span>
          <label class="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" class="sr-only peer" bind:checked={$adminmode} />
            <div class="w-9 h-5 bg-slate-600 peer-focus:ring-2 peer-focus:ring-violet-500/40 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-violet-600"></div>
          </label>
        </div>
      {/if}
    </div>
  </div>

  <!-- Disk Usage Card -->
  <div class="bg-white/5 glass rounded-xl m-3 overflow-hidden border border-white/10">
    <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
    <div class="w-full flex items-center justify-between px-4 py-3 transition-colors duration-150 hover:bg-white/10 cursor-pointer" role="button" tabindex="0" on:click={diskusageaction} on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); diskusageaction(); } }}>
      <h3 class="font-medium text-slate-200">Disk Usage</h3>
      <div class="flex items-center gap-2">
        {#if diskstatsOpen}
          <button
            type="button"
            aria-label="Refresh"
            class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-150"
            on:click|stopPropagation={() => { Send({ command: 'diskusage' }); }}>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
        {/if}
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400 transition-transform duration-200 {diskstatsOpen ? 'rotate-180' : ''}" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </div>
    </div>

    {#if diskstatsOpen}
      <div class="px-4 pb-4 space-y-2 text-sm text-slate-300">
        <div class="flex justify-between"><span class="text-slate-400">Total</span><span>{fileSize($diskstats?.total)}</span></div>
        <div class="flex justify-between"><span class="text-slate-400">Free</span><span>{fileSize($diskstats?.free)}</span></div>
        <div class="flex justify-between"><span class="text-slate-400">Used</span><span>{fileSize($diskstats?.used)} ({$diskstats?.usedPercent}%)</span></div>
      </div>
    {/if}
  </div>

  {#if $isAdmin === true}
    <!-- Misc Settings Card (Admin) -->
    <div class="bg-white/5 glass rounded-xl m-3 overflow-hidden border border-white/10">
      <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
      <div class="w-full flex items-center justify-between px-4 py-3 transition-colors duration-150 hover:bg-white/10 cursor-pointer" role="button" tabindex="0" on:click={miscsettingsaction} on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); miscsettingsaction(); } }}>
        <h3 class="font-medium text-slate-200">Misc Settings</h3>
        <div class="flex items-center gap-2">
          {#if miscOpen}
            <button
              type="button"
              aria-label="Refresh"
              class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-150"
              on:click|stopPropagation={() => { Send({ command: 'nooftrackersintrackerdb', aop: 1 }); }}>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </button>
          {/if}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400 transition-transform duration-200 {miscOpen ? 'rotate-180' : ''}" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </div>
      </div>

      {#if miscOpen}
        <div class="px-4 pb-4 space-y-3">
          <div class="text-sm text-slate-300">
            Total Trackers in TrackerDB: <span class="font-semibold text-slate-100">{$nooftrackersintrackerdb}</span>
          </div>
          <div class="grid gap-2">
            <button
              class="w-full bg-white/5 hover:bg-white/10 text-slate-200 py-2.5 text-sm rounded-xl border border-white/[0.06] transition-colors duration-150"
              on:click={() => { Send({ command: 'deletetrackersintrackerdb', data1: 'all', aop: 1 }); }}>Delete All Trackers in TrackerDB</button>
            <button
              class="w-full bg-white/5 hover:bg-white/10 text-slate-200 py-2.5 text-sm rounded-xl border border-white/[0.06] transition-colors duration-150"
              on:click={() => { Send({ command: 'trackerdbrefresh', aop: 1 }); }}>Refresh TrackerDB</button>
            <button
              class="w-full bg-white/5 hover:bg-white/10 text-slate-200 py-2.5 text-sm rounded-xl border border-white/[0.06] transition-colors duration-150"
              on:click={() => { Send({ command: 'stoponseedratio', aop: 1 }); }}>Seed Ratio Check</button>
          </div>
          <div class="space-y-2">
            <input type="text" bind:value={trdelno} required
                   class="w-full bg-white/5 rounded-xl px-3.5 py-2.5 border border-violet-700/30 placeholder-slate-500 text-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500/40"
                   placeholder="Delete N trackers from TrackerDB" />
            <button
              type="button"
              class="w-full bg-white/5 hover:bg-white/10 text-slate-200 py-2.5 text-sm rounded-xl border border-white/[0.06] transition-colors duration-150"
              on:click={() => { Send({ command: 'deletetrackersintrackerdb', data1: trdelno, aop: 1 }); }}>Delete Trackers</button>
          </div>
          <div class="space-y-2">
            <input type="text" bind:value={srno} required
                   class="w-full bg-white/5 rounded-xl px-3.5 py-2.5 border border-violet-700/30 placeholder-slate-500 text-slate-200 text-sm focus:outline-none focus:ring-2 focus:ring-violet-500/40"
                   placeholder="Stop on reaching this seed ratio" />
            <button
              type="button"
              class="w-full bg-white/5 hover:bg-white/10 text-slate-200 py-2.5 text-sm rounded-xl border border-white/[0.06] transition-colors duration-150"
              on:click={() => { Send({ command: 'stoponseedratio', data1: srno, aop: 1 }); }}>Stop Torrents</button>
          </div>
        </div>
      {/if}
    </div>

    <!-- Engine Settings Card (Admin) -->
    <div class="bg-white/5 glass rounded-xl m-3 overflow-hidden border border-white/10">
      <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
      <div class="w-full flex items-center justify-between px-4 py-3 transition-colors duration-150 hover:bg-white/10 cursor-pointer" role="button" tabindex="0" on:click={enginesettingsOpen} on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); enginesettingsOpen(); } }}>
        <h3 class="font-medium text-slate-200">Engine Settings</h3>
        <div class="flex items-center gap-2">
          {#if engsettingsOpen}
            <button
              type="button"
              class="px-3 py-1 rounded-lg bg-violet-600 hover:bg-violet-500 text-white text-xs font-medium transition-colors duration-150"
              on:click|stopPropagation={() => {
                if (engsettingsstring?.length === 0) {
                  toast.error('Empty Config!');
                } else {
                  let b64config = window.btoa(engsettingsstring);
                  Send({ command: 'updateconfig', data1: b64config, aop: 1 });
                }
              }}>Update</button>
            <button
              type="button"
              aria-label="Refresh"
              class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-150"
              on:click|stopPropagation={() => { Send({ command: 'getconfig', aop: 1 }); }}>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </button>
          {/if}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400 transition-transform duration-200 {engsettingsOpen ? 'rotate-180' : ''}" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </div>
      </div>

      {#if engsettingsOpen && $engconfig != null}
        <div class="px-4 pb-4">
          <textarea
            class="w-full bg-white/5 border border-white/[0.06] text-slate-200 text-sm font-mono resize-y h-96 rounded-xl p-3 focus:outline-none focus:ring-2 focus:ring-violet-500/40"
            bind:value={engsettingsstring}></textarea>
        </div>
      {/if}
    </div>
  {/if}

  <div class="mx-3 mt-2 mb-6">
    <button
      type="button"
      class="w-full py-3 text-sm font-medium rounded-xl text-slate-400 bg-white/5 border border-white/10 glass hover:bg-white/10 hover:text-slate-300 transition-all duration-150 focus:outline-none"
      on:click={() => { slocation.goto('/about'); }}>
      About exatorrent
    </button>
  </div>
</div>
