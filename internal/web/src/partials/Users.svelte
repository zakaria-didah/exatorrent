<script lang="ts">
  import { Send, userconnlist, userlist, isAdmin, signuprequests } from './core';
  import slocation from 'slocation';
  import CollapsibleSection from './CollapsibleSection.svelte';
  import { onMount } from 'svelte';
  import { toast } from 'svelte-sonner';
  import Useredit from './Useredit.svelte';

  const validUsernameRe = /^[a-zA-Z0-9_-]+$/;

  onMount(() => {
    if ($isAdmin === false) {
      slocation.goto('/');
    }
  });

  let signupRequestsOpen = false;
  let userconnlistOpen = false;
  let manageUsersOpen = false;
  let addUserOpen = false;

  let signupRequestsaction = () => {
    if (!signupRequestsOpen) {
      Send({ command: 'getsignuprequests', aop: 1 });
      signupRequestsOpen = true;
    } else {
      signupRequestsOpen = false;
    }
  };
  let newusername = '';
  let newpassword = '';
  let newusertype: 'user' | 'disabled' | 'admin' = 'user';
  let pwbox: HTMLInputElement;
  let pwvisible = false;

  let toggleinput = () => {
    pwvisible = !pwvisible;
    pwbox.type = pwvisible ? 'text' : 'password';
  };

  let addUser = () => {
    if (!newusername || !newpassword) {
      toast.error('Please fill in username and password');
      return;
    }
    if (newusername.length <= 5 || newpassword.length <= 5) {
      toast.error('Username and password must be longer than 5 characters');
      return;
    }
    if (!validUsernameRe.test(newusername)) {
      toast.error('Username must contain only letters, numbers, underscores, and hyphens');
      return;
    }
    Send({
      command: 'adduser',
      data1: newusername,
      data2: newpassword,
      data3: newusertype,
      aop: 1
    });
  };

  let typenotostring = (t: number): 'user' | 'disabled' | 'admin' => {
    switch (t) {
      case 0:
        return 'user';
      case 1:
        return 'admin';
      case -1:
        return 'disabled';
    }
  };

  let userconnlistaction = () => {
    if (userconnlistOpen === false) {
      Send({ command: 'listuserconns', aop: 1 });
      userconnlistOpen = true;
    } else {
      userconnlistOpen = false;
    }
  };

  let manageUsersaction = () => {
    if (manageUsersOpen === false) {
      Send({ command: 'getusers', aop: 1 });
      manageUsersOpen = true;
    } else {
      manageUsersOpen = false;
    }
  };

  let useraddaction = () => {
    addUserOpen = !addUserOpen;
  };
</script>

<div class="mx-auto max-w-5xl px-2 sm:px-0">
  <!-- Signup Requests -->
  <CollapsibleSection title="Signup Requests" bind:open={signupRequestsOpen} onToggle={signupRequestsaction} onRefresh={() => { Send({ command: 'getsignuprequests', aop: 1 }); }}>
    <div class="space-y-2">
      {#if $signuprequests.length === 0}
          <p class="text-slate-500 text-sm text-center py-4">No pending signup requests</p>
        {:else}
          {#each $signuprequests as req (req.id)}
            <div class="text-slate-200 bg-white/5 px-4 py-3 rounded-xl border border-white/[0.06]">
              <div class="flex items-start justify-between gap-3">
                <div class="flex-1 min-w-0">
                  <p class="font-medium text-sm">{req.username}</p>
                  {#if req.message}
                    <p class="text-xs text-slate-400 mt-1 break-all">{req.message}</p>
                  {/if}
                  <p class="text-xs text-slate-500 mt-1">{new Date(req.created_at)?.toLocaleString()}</p>
                </div>
                <div class="flex gap-2 flex-shrink-0">
                  <button
                    type="button"
                    class="px-3 py-1.5 text-xs font-medium rounded-md bg-violet-600 hover:bg-violet-500 text-white transition-colors duration-150"
                    on:click={() => {
                      Send({ command: 'approvesignup', data1: String(req.id), aop: 1 });
                      setTimeout(() => { Send({ command: 'getsignuprequests', aop: 1 }); }, 500);
                    }}>
                    Approve
                  </button>
                  <button
                    type="button"
                    class="px-3 py-1.5 text-xs font-medium rounded-md bg-white/10 hover:bg-white/15 text-slate-300 transition-colors duration-150"
                    on:click={() => {
                      Send({ command: 'declinesignup', data1: String(req.id), aop: 1 });
                      setTimeout(() => { Send({ command: 'getsignuprequests', aop: 1 }); }, 500);
                    }}>
                    Decline
                  </button>
                </div>
              </div>
            </div>
          {/each}
        {/if}
    </div>
  </CollapsibleSection>

  <!-- User Connections -->
  <CollapsibleSection title="User Connections" bind:open={userconnlistOpen} onToggle={userconnlistaction} onRefresh={() => { Send({ command: 'listuserconns', aop: 1 }); }}>
    <div class="space-y-2">
      {#each $userconnlist as eachuser (eachuser?.username)}
          <div class="text-slate-200 bg-white/5 px-4 py-3 rounded-xl border border-white/[0.06]">
            <div class="flex items-center justify-between gap-2">
              <div class="flex-1 min-w-0">
                <p class="font-medium text-sm break-all">
                  {eachuser?.username}
                  {#if eachuser?.isadmin === true}<span class="text-xs text-violet-400 ml-1">(admin)</span>{/if}
                </p>
                <p class="text-xs text-slate-500 mt-0.5">Connected at {new Date(eachuser?.contime)?.toLocaleString()}</p>
              </div>
              <button
                type="button"
                aria-label="Kick user"
                class="p-1.5 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-150 flex-shrink-0"
                on:click={() => { Send({ command: 'kickuser', data1: eachuser?.username, aop: 1 }); }}>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>
        {/each}
    </div>
  </CollapsibleSection>

  <!-- Manage Users -->
  <CollapsibleSection title="Manage Users" bind:open={manageUsersOpen} onToggle={manageUsersaction} onRefresh={() => { Send({ command: 'getusers', aop: 1 }); }}>
    <div class="space-y-2">
      {#each $userlist as eachuser (eachuser?.Username)}
          <Useredit username={eachuser?.Username} selected={typenotostring(eachuser?.UserType)} timestring={new Date(eachuser?.CreatedAt)?.toLocaleString()} quotabytes={eachuser?.quotabytes ?? 0} />
        {/each}
    </div>
  </CollapsibleSection>

  <!-- Add User -->
  <CollapsibleSection title="Add User" bind:open={addUserOpen} onToggle={useraddaction}>
    <div class="space-y-3">
        <input id="username" name="email" type="text" bind:value={newusername} required
               class="bg-white/5 appearance-none rounded-xl w-full px-3.5 py-2.5 border border-white/[0.06] placeholder-slate-500 text-slate-200 text-sm focus:outline-none focus:ring-1 focus:ring-violet-500/40"
               placeholder="Username" />

        <div class="flex bg-white/5 rounded-xl border border-white/[0.06] overflow-hidden">
          <input
            id="password"
            name="password"
            type="password"
            autocomplete="current-password"
            bind:value={newpassword}
            bind:this={pwbox}
            required
            class="bg-transparent w-full flex-grow px-3.5 py-2.5 border-none placeholder-slate-500 text-slate-200 text-sm focus:outline-none"
            placeholder="Password" />
          <button type="button" class="flex items-center justify-center w-11 flex-shrink-0 transition-colors duration-150 hover:bg-white/10" on:click={toggleinput}>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4.5 w-4.5 text-slate-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              {#if pwvisible}
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
              {:else}
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              {/if}
            </svg>
          </button>
        </div>

        <select class="bg-white/5 rounded-xl w-full px-3.5 py-2.5 border border-white/[0.06] text-slate-200 text-sm focus:outline-none focus:ring-1 focus:ring-violet-500/40 appearance-none cursor-pointer" bind:value={newusertype}>
          <option value="user">User</option>
          <option value="disabled">Disabled</option>
          <option value="admin">Admin</option>
        </select>

        <button type="button" on:click={addUser}
                class="w-full py-2.5 px-4 text-sm font-medium rounded-xl text-white bg-violet-600 hover:bg-violet-500 active:bg-violet-700 transition-colors duration-150 focus:outline-none focus:ring-1 focus:ring-violet-500/50">
          Add User
        </button>
    </div>
  </CollapsibleSection>
</div>
