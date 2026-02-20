<script lang="ts">
  import { onMount } from 'svelte';
  import { torrentsforuser, Send, adminmode, isAdmin, slocation } from './core';

  let username = $slocation.pathname?.split('/').reverse()[0];

  onMount(() => {
    if ($isAdmin === false) {
      slocation.goto('/');
    }
    torrentsforuser.set([]);
    adminmode.set(true);
    Send({ command: 'listtorrentsforuser', data1: username, aop: 1 });
  });
</script>

<div class="mx-auto max-w-3xl ">
  {#if Array.isArray($torrentsforuser) && $torrentsforuser?.length}
    {#each $torrentsforuser as trnt (trnt)}
      <button
        type="button"
        class="bg-white/5 glass border border-white/10 text-slate-200 px-3 py-5 rounded-xl m-3 noHL w-full text-left hover:bg-white/[0.08]"
        on:click={() => {
          if (typeof trnt === 'string' && trnt?.length > 0) {
            slocation.goto(`/torrent/${trnt}`);
          }
        }}>
        <div class="break-all mx-1 mb-1 font-bold">
          {trnt}
        </div>
      </button>
    {/each}
  {:else}
    <p class="text-xl text-center text-slate-400 font-sans">User owns no Torrents</p>
  {/if}
</div>
