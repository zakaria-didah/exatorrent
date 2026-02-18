<script lang="ts">
  import slocation from 'slocation';
  import { resplist } from './core';
</script>

<div class="mx-auto max-w-5xl px-3 ">
  {#if $resplist?.has === true}
    <button
      type="button"
      class="w-full my-2 flex justify-center py-2 px-4 border border-slate-700/40 text-sm font-medium rounded-lg text-slate-200 bg-slate-900 hover:bg-slate-800 focus:outline-none transition-colors duration-150"
      on:click={() => {
        resplist.set({ has: false, data: [] });
      }}>
      Clear All
    </button>
    {#each $resplist?.data as resp}
      <button
        type="button"
        class="bg-slate-900 text-slate-200 px-3 py-5 rounded-lg m-3 noHL w-full text-left"
        on:click={() => {
          if (typeof resp?.infohash === 'string' && resp?.infohash.length > 0) {
            slocation.goto(`/torrent/${resp?.infohash}`);
          }
        }}>
        <div class="break-all mx-1 mb-1 font-bold">
          {resp?.message}
        </div>
        <div class="text-slate-300 flex justify-between text-sm font-medium tabular-nums">
          <div class="break-all mx-1">
            {resp?.type}
            {resp?.state}
            {resp?.infohash ? `( ${resp?.infohash} )` : ''}
          </div>
        </div>
      </button>
    {/each}
  {:else}
    <p class="text-xl text-center text-slate-400 font-sans">No Notifications</p>
  {/if}
</div>
