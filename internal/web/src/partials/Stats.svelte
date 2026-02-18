<script lang="ts">
  import { torcstatus, machinfo, machstats, Send, fileSize, hasMachinfo, isAdmin } from './core';
  import { onMount } from 'svelte';
  import slocation from 'slocation';
  import CollapsibleSection from './CollapsibleSection.svelte';

  let deviceinfoOpen = false;
  let devicestatsOpen = false;
  let torcstatsOpen = false;

  onMount(() => {
    if ($isAdmin === false) {
      slocation.goto('/');
    }
    Send({ command: 'machstats', aop: 1 });
    Send({ command: 'torcstatus', aop: 1 });
  });

  let devinfoaction = () => {
    deviceinfoOpen = !deviceinfoOpen;
    if ($hasMachinfo === false) {
      Send({ command: 'machinfo', aop: 1 });
    }
  };

  let devicestatsaction = () => {
    if (devicestatsOpen === false) {
      Send({ command: 'machstats', aop: 1 });
      devicestatsOpen = true;
    } else {
      devicestatsOpen = false;
    }
  };

  let torcstatsaction = () => {
    if (torcstatsOpen === false) {
      Send({ command: 'torcstatus', aop: 1 });
      torcstatsOpen = true;
    } else {
      torcstatsOpen = false;
    }
  };
</script>

<div class="mx-auto max-w-5xl">
  <CollapsibleSection title="Device Info" bind:open={deviceinfoOpen} onToggle={devinfoaction}>
    <div class="flex flex-col">
      {#if $hasMachinfo === true}
        <div class="m-1 p-1 break-all text-slate-300">Arch: {$machinfo?.arch}</div>
        <div class="m-1 p-1 break-all text-slate-300">CPU Model: {$machinfo?.cpumodel}</div>
        <div class="m-1 p-1 break-all text-slate-300">Go Version: {$machinfo?.goversion}</div>
        <div class="m-1 p-1 break-all text-slate-300">Hostname: {$machinfo?.hostname}</div>
        <div class="m-1 p-1 break-all text-slate-300">CPU No: {$machinfo?.numbercpu}</div>
        <div class="m-1 p-1 break-all text-slate-300">OS: {$machinfo?.os}</div>
        <div class="m-1 p-1 break-all text-slate-300">Platform: {$machinfo?.platform}</div>
        <div class="m-1 p-1 break-all text-slate-300">Started at {new Date($machinfo?.startedat)?.toLocaleString()}</div>
        <div class="m-1 p-1 break-all text-slate-300">Total Mem: {fileSize($machinfo?.totalmem)}</div>
      {/if}
    </div>
  </CollapsibleSection>
</div>

<div class="mx-auto max-w-5xl">
  <CollapsibleSection title="Device Stats" bind:open={devicestatsOpen} onToggle={devicestatsaction} onRefresh={() => Send({ command: 'machstats', aop: 1 })}>
    <div class="flex flex-col">
      <div class="m-1 p-1 break-all text-slate-300">CPU Cycles: {$machstats?.cpucycles}</div>
      <div class="m-1 p-1 break-all text-slate-300">Disk Free: {fileSize($machstats?.diskfree)}</div>
      <div class="m-1 p-1 break-all text-slate-300">Disk Free(Bytes): {$machstats?.diskfree}</div>
      <div class="m-1 p-1 break-all text-slate-300">Disk Percent: {$machstats?.diskpercent} %</div>
      <div class="m-1 p-1 break-all text-slate-300">Memory Percent: {$machstats?.mempercent} %</div>
      <div class="m-1 p-1 break-all text-slate-300">Go Mem: {fileSize($machstats?.gomem)}</div>
      <div class="m-1 p-1 break-all text-slate-300">Go Mem(sys): {fileSize($machstats?.gomemsys)}</div>
      <div class="m-1 p-1 break-all text-slate-300">Goroutines: {$machstats?.goroutines}</div>
    </div>
  </CollapsibleSection>
</div>

<div class="mx-auto max-w-5xl">
  <CollapsibleSection title="Torrent Client Status" bind:open={torcstatsOpen} onToggle={torcstatsaction} onRefresh={() => Send({ command: 'torcstatus', aop: 1 })}>
    <div class="flex flex-col overflow-x-auto">
      <pre class="whitespace-pre mb-2 p-2 text-slate-300">{$torcstatus} </pre>
    </div>
  </CollapsibleSection>
</div>
