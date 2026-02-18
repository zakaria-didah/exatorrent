<script lang="ts">
  import slocation from 'slocation';
  import { resplist } from './core';
</script>

<div class="mx-auto max-w-5xl px-3">
  {#if $resplist?.has === true}
    <button
      type="button"
      class="w-full my-2 flex justify-center py-2.5 px-4 border border-white/10 text-sm font-medium rounded-lg text-slate-200 bg-white/5 glass hover:bg-white/10 focus:outline-none transition-colors duration-150"
      on:click={() => { resplist.set({ has: false, data: [] }); }}>
      Clear All
    </button>
    {#each $resplist?.data as resp}
      <button
        type="button"
        class="bg-white/5 glass text-slate-200 px-4 py-4 rounded-lg border border-white/10 my-2 noHL w-full text-left hover:bg-white/[0.08] transition-colors duration-150"
        on:click={() => {
          if (typeof resp?.infohash === 'string' && resp?.infohash.length > 0) {
            slocation.goto(`/torrent/${resp?.infohash}`);
          }
        }}>
        <div class="break-all mb-1 font-bold text-sm">{resp?.message}</div>
        <div class="text-slate-400 flex justify-between text-xs font-medium tabular-nums">
          <div class="break-all">
            {resp?.type}
            {resp?.state}
            {resp?.infohash ? `( ${resp?.infohash} )` : ''}
          </div>
        </div>
      </button>
    {/each}
  {:else}
    <p class="text-xl text-center text-slate-400 font-sans mt-8">No Notifications</p>
  {/if}
</div>
