<script lang="ts">
  import { onMount } from 'svelte';
  import { socket, Connect, isDisConnected } from './core';
  import slocation from 'slocation';

  function closeSocket() {
    if (socket?.readyState === WebSocket.OPEN) {
      socket.close();
    }
  }

  onMount(() => {
    if (socket == null || socket == undefined || socket?.readyState === WebSocket.CLOSED) {
      console.log('Attempting to Connect');
      slocation.goto('/');
      Connect();
    }
  });
</script>

<nav class="max-w-7xl mx-auto px-3 sm:px-4 lg:px-6">
  <div class="flex items-center justify-between h-14 gap-2">
    {#if $slocation.pathname !== '/'}
      <button
        aria-label="Go back"
        class="flex items-center justify-center w-10 h-10 rounded-lg bg-white/5 border border-white/10 glass transition-colors duration-150 hover:bg-white/10 active:bg-white/15 flex-shrink-0"
        on:click={() => {
          history.length > 2 ? history.back() : slocation.goto('/');
        }}
        title="Go back">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
      </button>
    {:else}
      <button
        aria-label="Notifications"
        class="flex items-center justify-center w-10 h-10 rounded-lg bg-white/5 border border-white/10 glass transition-colors duration-150 hover:bg-white/10 active:bg-white/15 flex-shrink-0"
        on:click={() => {
          slocation.goto('/notifications');
        }}
        title="Notifications">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
      </button>
    {/if}

    <button
      class="flex items-center gap-2 px-4 py-2 rounded-lg noHL transition-colors duration-150 hover:bg-white/5"
      title="exatorrent"
      on:click={() => {
        if ($slocation.pathname === '/') {
          slocation.goto('/about');
        } else {
          slocation.goto('/');
        }
      }}>
      <span class="font-semibold text-lg text-slate-100 tracking-tight">exatorrent</span>
      <span class="w-2 h-2 rounded-full {$isDisConnected ? 'bg-slate-500' : 'bg-violet-500'} transition-colors duration-300"></span>
    </button>

    <button
      aria-label="Disconnect"
      class="flex items-center justify-center w-10 h-10 rounded-lg bg-white/5 border border-white/10 glass transition-colors duration-150 hover:bg-white/10 active:bg-white/15 flex-shrink-0"
      on:click={closeSocket}
      title="Disconnect">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.636 18.364a9 9 0 010-12.728m12.728 0a9 9 0 010 12.728m-9.9-2.829a5 5 0 010-7.07m7.072 0a5 5 0 010 7.07M13 12a1 1 0 11-2 0 1 1 0 012 0z" />
      </svg>
    </button>
  </div>
</nav>
