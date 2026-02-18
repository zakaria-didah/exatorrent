<script lang="ts">
  import { toast } from 'svelte-sonner';
  import { Connect } from './core';
  import { onMount } from 'svelte';

  let mode: 'signin' | 'request' = 'signin';
  let pwvisible = false;
  let pwbox: HTMLInputElement;
  let reqpwbox: HTMLInputElement;
  let reqpwvisible = false;
  let submitting = false;

  function toggleinput() {
    pwvisible = !pwvisible;
    pwbox.type = pwvisible ? 'text' : 'password';
  }

  function togglereqpw() {
    reqpwvisible = !reqpwvisible;
    reqpwbox.type = reqpwvisible ? 'text' : 'password';
  }

  let exausername: string;
  let exapassword: string;

  let reqUsername = '';
  let reqPassword = '';
  let reqMessage = '';

  let signingIn = false;
  const validUsernameRe = /^[a-zA-Z0-9_-]+$/;

  onMount(() => {
    const un = localStorage.getItem('exausername');
    if (un) {
      Connect();
    }
  });

  async function signIn() {
    if (!exausername || !exapassword) {
      toast.error('Please fill in all fields');
      return;
    }
    if (exausername.length <= 5 || exapassword.length <= 5) {
      toast.error('Username and password must be longer than 5 characters');
      return;
    }
    if (!validUsernameRe.test(exausername)) {
      toast.error('Username must contain only letters, numbers, underscores, and hyphens');
      return;
    }
    signingIn = true;
    try {
      const res = await fetch('/api/auth', {
        method: 'POST',
        body: JSON.stringify({ data1: exausername, data2: exapassword })
      });
      if (res.status < 200 || res.status > 299) {
        toast.error('Invalid credentials');
        return;
      }
      const data = await res.json();
      localStorage.setItem('exausername', exausername);
      localStorage.setItem('exausertype', data?.usertype);
      Connect();
    } catch {
      toast.error('Network error');
    } finally {
      signingIn = false;
    }
  }

  async function submitRequest() {
    if (!reqUsername || !reqPassword) {
      toast.error('Please fill in username and password');
      return;
    }
    if (reqUsername.length <= 5 || reqPassword.length <= 5) {
      toast.error('Username and password must be longer than 5 characters');
      return;
    }
    if (!validUsernameRe.test(reqUsername)) {
      toast.error('Username must contain only letters, numbers, underscores, and hyphens');
      return;
    }
    submitting = true;
    try {
      const res = await fetch('/api/signup-request', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: reqUsername, password: reqPassword, message: reqMessage })
      });
      const data = await res.json();
      if (res.ok) {
        toast.success(data.message || 'Request submitted');
        reqUsername = '';
        reqPassword = '';
        reqMessage = '';
        mode = 'signin';
      } else {
        toast.error(data.error || 'Request failed');
      }
    } catch {
      toast.error('Network error');
    } finally {
      submitting = false;
    }
  }

  let entertosignin = (event: KeyboardEvent) => {
    if (event.code === 'Enter') signIn();
  };

  let entertorequest = (event: KeyboardEvent) => {
    if (event.code === 'Enter') submitRequest();
  };
</script>

<div class="min-h-screen flex items-center justify-center px-4">
    <div class="max-w-sm w-full">
      <div class="text-center mb-10">
        <a href="https://github.com/zakaria-didah/exatorrent" target="_blank" rel="noopener noreferrer" class="inline-block group">
          <h1 class="text-5xl sm:text-6xl font-bold tracking-tight text-white" style="text-shadow: 0 0 40px rgba(139, 92, 246, 0.3), 0 0 80px rgba(139, 92, 246, 0.1);">
            exa<span class="text-violet-400">torrent</span>
          </h1>
        </a>
        <p class="mt-3 text-sm text-slate-500 tracking-wide uppercase font-medium">Stream torrent</p>
      </div>

      <div class="glass bg-slate-900/30 rounded-2xl border border-slate-700/30 p-7 shadow-2xl shadow-violet-950/20">
        {#if mode === 'signin'}
          <div>
            <h2 class="text-lg font-semibold text-slate-200 mb-5">Sign in</h2>

            <div class="space-y-4">
              <div>
                <label for="username" class="block text-xs font-medium text-slate-400 mb-1.5 uppercase tracking-wider">Username</label>
                <input
                  id="username"
                  name="email"
                  type="text"
                  bind:value={exausername}
                  required
                  class="bg-slate-950/50 w-full px-4 py-3 rounded-xl border border-slate-700/40 placeholder-slate-600 text-slate-200 text-sm focus:outline-none focus:ring-1 focus:ring-violet-500/40 focus:border-violet-500/40 transition-all duration-200"
                  placeholder="Enter username" />
              </div>

              <div>
                <label for="password" class="block text-xs font-medium text-slate-400 mb-1.5 uppercase tracking-wider">Password</label>
                <div class="flex bg-slate-950/50 rounded-xl border border-slate-700/40 overflow-hidden transition-all duration-200 focus-within:ring-1 focus-within:ring-violet-500/40 focus-within:border-violet-500/40">
                  <input
                    id="password"
                    name="password"
                    type="password"
                    bind:value={exapassword}
                    bind:this={pwbox}
                    on:keydown={entertosignin}
                    required
                    class="bg-transparent w-full flex-grow px-4 py-3 border-none placeholder-slate-600 text-slate-200 text-sm focus:outline-none"
                    placeholder="Enter password" />
                  <button type="button" class="flex items-center justify-center w-12 flex-shrink-0 transition-colors duration-150 hover:bg-slate-800/50" on:click={toggleinput}>
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      {#if pwvisible}
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                      {:else}
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                      {/if}
                    </svg>
                  </button>
                </div>
              </div>

              <button
                type="button"
                disabled={signingIn}
                on:click={signIn}
                class="w-full py-3 px-4 text-sm font-semibold rounded-xl text-white bg-violet-600 hover:bg-violet-500 active:bg-violet-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-violet-500/40 focus:ring-offset-2 focus:ring-offset-slate-900 disabled:opacity-50 disabled:cursor-not-allowed mt-1">
                {signingIn ? 'Signing in...' : 'Sign in'}
              </button>
            </div>

            <div class="mt-6 pt-5 border-t border-slate-700/30 text-center">
              <p class="text-sm text-slate-500">Don't have an account?</p>
              <button
                type="button"
                on:click={() => { mode = 'request'; }}
                class="mt-2 text-sm font-medium text-violet-400 hover:text-violet-300 transition-colors duration-150 bg-transparent border-none cursor-pointer">
                Request Access
              </button>
            </div>
          </div>
        {:else}
          <div>
            <div class="flex items-center gap-3 mb-5">
              <button
                type="button"
                aria-label="Back to sign in"
                on:click={() => { mode = 'signin'; }}
                class="p-1.5 rounded-lg hover:bg-slate-800/50 transition-colors duration-150 bg-transparent border-none cursor-pointer text-slate-400 hover:text-slate-200">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                </svg>
              </button>
              <h2 class="text-lg font-semibold text-slate-200">Request Access</h2>
            </div>

            <div class="space-y-4">
              <div>
                <label for="req-username" class="block text-xs font-medium text-slate-400 mb-1.5 uppercase tracking-wider">Username</label>
                <input
                  id="req-username"
                  type="text"
                  bind:value={reqUsername}
                  class="bg-slate-950/50 w-full px-4 py-3 rounded-xl border border-slate-700/40 placeholder-slate-600 text-slate-200 text-sm focus:outline-none focus:ring-1 focus:ring-violet-500/40 focus:border-violet-500/40 transition-all duration-200"
                  placeholder="Choose a username" />
              </div>

              <div>
                <label for="req-password" class="block text-xs font-medium text-slate-400 mb-1.5 uppercase tracking-wider">Password</label>
                <div class="flex bg-slate-950/50 rounded-xl border border-slate-700/40 overflow-hidden transition-all duration-200 focus-within:ring-1 focus-within:ring-violet-500/40 focus-within:border-violet-500/40">
                  <input
                    id="req-password"
                    type="password"
                    bind:value={reqPassword}
                    bind:this={reqpwbox}
                    on:keydown={entertorequest}
                    class="bg-transparent w-full flex-grow px-4 py-3 border-none placeholder-slate-600 text-slate-200 text-sm focus:outline-none"
                    placeholder="Choose a password" />
                  <button type="button" class="flex items-center justify-center w-12 flex-shrink-0 transition-colors duration-150 hover:bg-slate-800/50" on:click={togglereqpw}>
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      {#if reqpwvisible}
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                      {:else}
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                      {/if}
                    </svg>
                  </button>
                </div>
              </div>

              <div>
                <label for="req-message" class="block text-xs font-medium text-slate-400 mb-1.5 uppercase tracking-wider">Message <span class="text-slate-600 normal-case">(optional)</span></label>
                <textarea
                  id="req-message"
                  bind:value={reqMessage}
                  rows="3"
                  maxlength="500"
                  class="bg-slate-950/50 w-full px-4 py-3 rounded-xl border border-slate-700/40 placeholder-slate-600 text-slate-200 text-sm focus:outline-none focus:ring-1 focus:ring-violet-500/40 focus:border-violet-500/40 transition-all duration-200 resize-none"
                  placeholder="Why would you like access?"></textarea>
              </div>

              <button
                type="button"
                disabled={submitting}
                on:click={submitRequest}
                class="w-full py-3 px-4 text-sm font-semibold rounded-xl text-white bg-violet-600 hover:bg-violet-500 active:bg-violet-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-violet-500/40 focus:ring-offset-2 focus:ring-offset-slate-900 disabled:opacity-50 disabled:cursor-not-allowed mt-1">
                {submitting ? 'Submitting...' : 'Submit Request'}
              </button>
            </div>

            <div class="mt-6 pt-5 border-t border-slate-700/30 text-center">
              <p class="text-sm text-slate-500">Already have an account?</p>
              <button
                type="button"
                on:click={() => { mode = 'signin'; }}
                class="mt-2 text-sm font-medium text-violet-400 hover:text-violet-300 transition-colors duration-150 bg-transparent border-none cursor-pointer">
                Sign in
              </button>
            </div>
          </div>
        {/if}
      </div>

      <p class="text-center text-slate-700 text-xs mt-8">Powered by exatorrent</p>
    </div>
</div>
