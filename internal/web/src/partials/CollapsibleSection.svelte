<script lang="ts">
  export let title: string;
  export let open: boolean = false;
  export let onRefresh: (() => void) | undefined = undefined;
  export let onToggle: (() => void) | undefined = undefined;

  function toggle() {
    if (onToggle) {
      onToggle();
    } else {
      open = !open;
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      toggle();
    }
  }
</script>

<div class="bg-white/5 glass rounded-xl m-3 overflow-hidden border border-white/10 noHL">
  <div
    class="w-full flex items-center justify-between px-4 py-3 transition-colors duration-150 hover:bg-white/10 cursor-pointer"
    role="button"
    tabindex="0"
    aria-expanded={open}
    on:click={toggle}
    on:keydown={handleKeydown}>
    <h3 class="font-medium text-slate-200">{title}</h3>
    <div class="flex items-center gap-2">
      {#if open && onRefresh}
        <button
          type="button"
          aria-label="Refresh"
          class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-150"
          on:click|stopPropagation={onRefresh}>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      {/if}
      {#if open}
        <slot name="actions" />
      {/if}
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-slate-400 transition-transform duration-200 {open ? 'rotate-180' : ''}" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </div>
  </div>

  {#if open}
    <div class="px-4 pb-4">
      <slot />
    </div>
  {/if}
</div>
