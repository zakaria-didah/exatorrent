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
  import CollapsibleSection from './CollapsibleSection.svelte';
  import slocation from 'slocation';
  import { toast } from 'svelte-sonner';

  let ds = false;
  if ($dontstart === 'true') {
    ds = true;
  }

  let toggledontstart = (ds: boolean) => {
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

    if (socket?.readyState === WebSocket.OPEN) socket?.close();
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
  <CollapsibleSection title="Disk Usage" bind:open={diskstatsOpen} onToggle={diskusageaction} onRefresh={() => { Send({ command: 'diskusage' }); }}>
    <div class="space-y-2 text-sm text-slate-300">
      <div class="flex justify-between"><span class="text-slate-400">Total</span><span>{fileSize($diskstats?.total)}</span></div>
      <div class="flex justify-between"><span class="text-slate-400">Free</span><span>{fileSize($diskstats?.free)}</span></div>
      <div class="flex justify-between"><span class="text-slate-400">Used</span><span>{fileSize($diskstats?.used)} ({$diskstats?.usedPercent}%)</span></div>
    </div>
  </CollapsibleSection>

  {#if $isAdmin === true}
    <!-- Misc Settings Card (Admin) -->
    <CollapsibleSection title="Misc Settings" bind:open={miscOpen} onToggle={miscsettingsaction} onRefresh={() => { Send({ command: 'nooftrackersintrackerdb', aop: 1 }); }}>
      <div class="space-y-3">
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
    </CollapsibleSection>

    <!-- Engine Settings Card (Admin) -->
    <CollapsibleSection title="Engine Settings" bind:open={engsettingsOpen} onToggle={enginesettingsOpen} onRefresh={() => { Send({ command: 'getconfig', aop: 1 }); }}>
      <svelte:fragment slot="actions">
        <button type="button" class="px-3 py-1 rounded-lg bg-violet-600 hover:bg-violet-500 text-white text-xs font-medium transition-colors duration-150" on:click|stopPropagation={() => {
          if (engsettingsstring?.length === 0) {
            toast.error('Empty Config!');
          } else {
            let b64config = window.btoa(engsettingsstring);
            Send({ command: 'updateconfig', data1: b64config, aop: 1 });
          }
        }}>Update</button>
      </svelte:fragment>
      {#if $engconfig != null}
        <textarea
          class="w-full bg-white/5 border border-white/[0.06] text-slate-200 text-sm font-mono resize-y h-96 rounded-xl p-3 focus:outline-none focus:ring-2 focus:ring-violet-500/40"
          bind:value={engsettingsstring}></textarea>
      {/if}
    </CollapsibleSection>
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
