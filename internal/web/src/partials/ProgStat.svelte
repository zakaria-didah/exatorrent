<script lang="ts">
  import { fileSize } from './core';

  export let bytescompleted: number;
  export let length: number;
  export let offset: number;

  let progress: number;

  $: progress = (bytescompleted / length) * 100;
</script>

<div class="space-y-1.5 mt-1">
  <div class="bg-slate-800 rounded-full overflow-hidden h-1.5">
    <div class="bg-violet-500 h-full rounded-full transition-all duration-500 ease-out" style="width:{progress ? progress : 0}%"></div>
  </div>
  <div class="text-slate-400 flex justify-between text-xs font-medium">
    <div class="break-all">
      {fileSize(bytescompleted)} / {fileSize(length)} <span class="text-slate-500">(Off. {offset})</span>
    </div>
    <div class="font-semibold {progress >= 100 ? 'text-violet-400' : 'text-slate-300'}">
      {progress?.toLocaleString('en-US', {
        maximumFractionDigits: 1,
        minimumFractionDigits: 1
      })}%
    </div>
  </div>
</div>
