<script lang="ts">
  import { onMount } from 'svelte';
  import { Send, slocation } from './core';

  onMount(() => {
    oldselected = selected;
    utchangeallowed = true;
  });
  import { toast } from 'svelte-sonner';

  export let username = '';
  export let timestring = '';
  export let selected: 'user' | 'disabled' | 'admin' = 'disabled';
  export let quotabytes: number = 0;

  let oldselected: 'user' | 'disabled' | 'admin' = 'disabled';
  let usereditmode = false;
  let newpw = '';
  let utchangeallowed = false;
  let quotaGB = quotabytes > 0 ? +(quotabytes / (1024 * 1024 * 1024)).toFixed(1) : 0;

  let updateuser = () => {
    if (newpw?.length > 5) {
      Send({
        command: 'updatepw',
        data1: username,
        data2: newpw,
        aop: 1
      });
    }
    if (utchangeallowed === true) {
      if (selected !== oldselected) {
        Send({
          command: 'changeusertype',
          data1: username,
          data2: selected,
          aop: 1
        });
      } else {
        toast.error('Same Usertype Selected');
      }
    }
    const newQuotaBytes = quotaGB * 1024 * 1024 * 1024;
    if (newQuotaBytes !== quotabytes) {
      Send({
        command: 'setquota',
        data1: username,
        data2: String(quotaGB),
        aop: 1
      });
    }
  };
</script>

<div class="text-slate-200 bg-white/5 px-4 py-3 rounded-xl border border-white/[0.06]">
  <div class="flex flex-col justify-between flex-wrap py-1">
    <div class="flex flex-row justify-between items-center">
      <div class="font-medium text-sm break-all text-left">{username}</div>
      <div class="flex gap-1">
        <button
          type="button"
          aria-label="Edit user"
          class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-150 noHL"
          on:click={() => { usereditmode = !usereditmode; }}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
          </svg>
        </button>

        <button
          type="button"
          aria-label="Remove user"
          class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-150 noHL"
          on:click={() => { Send({ command: 'removeuser', data1: username, aop: 1 }); }}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>

        <button
          type="button"
          aria-label="View user torrents"
          class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-150 noHL"
          on:click={() => { slocation.goto(`/user/${username}`); }}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
          </svg>
        </button>
      </div>
    </div>
    <p class="text-xs text-slate-500 mt-0.5">Created at {timestring} · Quota: {quotabytes > 0 ? (quotabytes / (1024 * 1024 * 1024)).toFixed(1) + ' GB' : 'Unlimited'}</p>
  </div>
  {#if usereditmode === true}
    <div class="flex items-center gap-2 mt-2 pt-2 border-t border-white/[0.06]">
      <input id="changepw" type="text" bind:value={newpw} required
             class="bg-white/5 appearance-none rounded-lg flex-grow px-3 py-2 border border-white/[0.06] placeholder-slate-500 text-slate-200 text-sm focus:outline-none focus:ring-1 focus:ring-violet-500/40"
             placeholder="New Password" />
      <button
        type="button"
        aria-label="Revoke token"
        class="p-2 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-150 noHL flex-shrink-0"
        on:click={() => { Send({ command: 'revoketoken', data1: username, aop: 1 }); }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
      </button>
      <div class="flex items-center gap-1 flex-shrink-0">
        <input
          type="number"
          min="0"
          step="0.5"
          bind:value={quotaGB}
          class="bg-white/5 appearance-none rounded-lg w-16 px-2 py-2 border border-white/[0.06] text-slate-200 text-sm text-center focus:outline-none focus:ring-1 focus:ring-violet-500/40"
          title="Storage quota in GB (0 = unlimited)" />
        <span class="text-slate-500 text-xs">GB</span>
      </div>
      <select class="bg-white/5 text-slate-200 text-sm rounded-lg px-2 py-2 border border-white/[0.06] focus:outline-none appearance-none cursor-pointer" bind:value={selected}>
        <option value="user">User</option>
        <option value="disabled">Disabled</option>
        <option value="admin">Admin</option>
      </select>
      <button class="bg-violet-600 hover:bg-violet-500 text-white text-sm px-3 py-2 rounded-lg transition-colors duration-150 flex-shrink-0" on:click={updateuser}>Update</button>
    </div>
  {/if}
</div>
