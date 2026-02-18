<script lang="ts">
  import { dontstart, Send, isAdmin } from './core';
  import slocation from 'slocation';
  import { toast } from 'svelte-sonner';
  let ismetainfo = true;
  let torrentinput = '';
  let trntfilestring = '';
  let trntfileinput: HTMLInputElement;

  let addfunc = () => {
    if (ismetainfo === true) {
      if (torrentinput === '') {
        toast.error('Empty Input');
        return;
      }
      if (torrentinput.startsWith('magnet:')) {
        console.log('Adding Magnet', torrentinput, $dontstart);
        Send({
          command: 'addmagnet',
          data1: torrentinput,
          data2: $dontstart === 'true' ? 'true' : 'false'
        });
        torrentinput = '';
      } else {
        console.log('Adding Infohash', torrentinput, $dontstart);
        Send({
          command: 'addinfohash',
          data1: torrentinput,
          data2: $dontstart === 'true' ? 'true' : 'false'
        });
        torrentinput = '';
      }
    } else {
      if (trntfilestring === '') {
        toast.error('File Invalid');
        return;
      }
      console.log('Adding Torrent', trntfilestring, $dontstart);
      Send({
        command: 'addtorrent',
        data1: trntfilestring,
        data2: $dontstart === 'true' ? 'true' : 'false'
      });
      trntfilestring = '';
    }
  };

  let entertoadd = (event: KeyboardEvent) => {
    if (event.code === 'Enter') {
      addfunc();
    }
  };
  function toggleismetainfo() {
    ismetainfo = !ismetainfo;
  }

  function readtrnt(e: Event) {
    let f = (e.target as HTMLInputElement).files[0];
    if (f.size > 20971520) {
      toast.error('Error: Maximum Torrent File Size is 20MB');
      return;
    }
    let reader = new FileReader();
    reader.onload = (e) => {
      trntfilestring = btoa(e.target.result as string);
    };
    reader.readAsBinaryString(f);
    (e.target as HTMLInputElement).value = null;
  }
</script>

<div class="mt-8 flex items-center justify-center px-4">
  <div class="max-w-lg w-full">
    <div>
      <h2 class="text-center text-2xl sm:text-3xl font-bold text-slate-200">
        {#if ismetainfo}Enter Magnet or Infohash{:else}Select Torrent File{/if}
      </h2>
    </div>

    <div class="mt-6 space-y-3">
      <div class="flex bg-white/5 rounded-lg border border-white/10 glass overflow-hidden transition-all duration-150 focus-within:ring-1 focus-within:ring-violet-500/50 focus-within:border-violet-500/50">
        {#if ismetainfo}
          <input
            id="torrentinput"
            type="text"
            required
            class="bg-transparent w-full flex-grow px-4 py-3 border-none placeholder-slate-500 text-slate-200 focus:outline-none text-sm"
            placeholder="Magnet link or Infohash..."
            bind:value={torrentinput}
            on:keydown={entertoadd} />
        {:else}
          <label class="w-full flex-grow px-4 py-3 cursor-pointer">
            <div class="text-slate-300 flex items-center gap-2 text-sm">
              <svg xmlns="http://www.w3.org/2000/svg" class="text-slate-500 h-5 w-5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
              </svg>
              <span>{trntfilestring ? 'File selected' : 'Select a .torrent file'}</span>
            </div>
            <input accept=".torrent,application/x-bittorrent" bind:this={trntfileinput} on:change={(e) => readtrnt(e)} id="torrentfile" name="torrentfile" type="file" class="hidden" />
          </label>
        {/if}
        <button
          type="button"
          class="flex items-center justify-center w-11 flex-shrink-0 transition-colors duration-150 hover:bg-white/10"
          on:click={toggleismetainfo}
          title={ismetainfo ? 'Switch to file upload' : 'Switch to magnet/infohash'}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            {#if ismetainfo}
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
            {:else}
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
            {/if}
          </svg>
        </button>
      </div>

      <button
        type="button"
        class="w-full flex justify-center py-2.5 px-4 text-sm font-medium rounded-lg text-white bg-violet-600 hover:bg-violet-500 active:bg-violet-700 transition-colors duration-150 focus:outline-none focus:ring-1 focus:ring-violet-500/50"
        on:click={addfunc}>
        Add Torrent
      </button>
    </div>
  </div>
</div>

<div class="mx-auto max-w-lg px-4 mt-8">
  <div class="grid grid-cols-2 gap-3">
    <button
      class="bg-white/5 text-slate-200 px-4 py-4 rounded-lg border border-white/10 glass transition-all duration-150 hover:bg-white/10 hover:border-white/15 active:bg-white/15 flex items-center gap-3"
      on:click={() => {
        slocation.goto('/torrents');
      }}>
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
      </svg>
      <span class="text-sm font-medium">Torrents</span>
    </button>
    <button
      class="bg-white/5 text-slate-200 px-4 py-4 rounded-lg border border-white/10 glass transition-all duration-150 hover:bg-white/10 hover:border-white/15 active:bg-white/15 flex items-center gap-3"
      on:click={() => {
        slocation.goto('/settings');
      }}>
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
      </svg>
      <span class="text-sm font-medium">Settings</span>
    </button>
  </div>
  {#if $isAdmin}
    <div class="grid grid-cols-2 gap-3 mt-3">
      <button
        class="bg-white/5 text-slate-200 px-4 py-4 rounded-lg border border-white/10 glass transition-all duration-150 hover:bg-white/10 hover:border-white/15 active:bg-white/15 flex items-center gap-3"
        on:click={() => {
          slocation.goto('/users');
        }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
        </svg>
        <span class="text-sm font-medium">Users</span>
      </button>
      <button
        class="bg-white/5 text-slate-200 px-4 py-4 rounded-lg border border-white/10 glass transition-all duration-150 hover:bg-white/10 hover:border-white/15 active:bg-white/15 flex items-center gap-3"
        on:click={() => {
          slocation.goto('/stats');
        }}>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z" />
        </svg>
        <span class="text-sm font-medium">Stats</span>
      </button>
    </div>
  {/if}
</div>
