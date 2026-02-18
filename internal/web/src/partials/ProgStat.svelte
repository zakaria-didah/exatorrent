<script lang="ts">
  import { fileSize } from './core';

  export let bytescompleted: number;
  export let length: number;

  let progress: number;

  $: progress = length > 0 ? (bytescompleted / length) * 100 : 0;
</script>

<div class="space-y-1 mt-1">
  <div class="bg-white/10 rounded-full overflow-hidden h-1">
    <div class="bg-violet-500 h-full rounded-full transition-all duration-500 ease-out" style="width:{progress}%"></div>
  </div>
  <div class="text-slate-500 flex justify-between text-xs tabular-nums">
    <span>{fileSize(bytescompleted)} / {fileSize(length)}</span>
    <span class="font-medium {progress >= 100 ? 'text-violet-400' : 'text-slate-400'}">
      {progress.toLocaleString('en-US', { maximumFractionDigits: 1, minimumFractionDigits: 1 })}%
    </span>
  </div>
</div>
