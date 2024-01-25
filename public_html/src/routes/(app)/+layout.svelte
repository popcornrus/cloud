<script>
    import {onMount} from "svelte";
    import {User} from "$lib/classes/User.js";
    import Image from "$lib/components/Image.svelte";
    import {authToken} from "$lib/stores/token.js";
    import { afterNavigate } from "$app/navigation";

    let userData = null;

    afterNavigate(() => {
        HSStaticMethods.autoInit();
    });

    onMount(() => {

        if ($authToken === null) {
            window.location.href = '/auth/sign-in';
        }

        userData = JSON.parse(localStorage.getItem('user'));

        if (userData === null) {
            const user = new User($authToken);

            user.me().then((response) => {
                localStorage.setItem('user', JSON.stringify(response.data.data));
                userData = response.data.data;
            });
        }
    })
</script>

<div class="fixed h-full hs-overlay hs-overlay-open:translate-x-0 -translate-x-full transition-all duration-300 transform hidden z-[60] w-64 bg-white border-e border-gray-200 pt-7 pb-10 overflow-y-auto lg:block lg:translate-x-0 lg:end-auto lg:bottom-0 [&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-track]:bg-gray-100 [&::-webkit-scrollbar-thumb]:bg-gray-300 dark:[&::-webkit-scrollbar-track]:bg-slate-700 dark:[&::-webkit-scrollbar-thumb]:bg-slate-500 dark:bg-gray-800 dark:border-gray-700">
    <div class="px-6">
        <a class="flex-none text-xl font-semibold dark:text-white" href="/" aria-label="Brand">Brand</a>
    </div>
    <nav class="hs-accordion-group p-6 w-full flex flex-col flex-wrap" data-hs-accordion-always-open>
        <ul class="space-y-1.5">
            <li>
                <a class="flex items-center gap-x-3.5 py-2 px-2.5 bg-gray-100 text-sm text-slate-700 rounded-lg hover:bg-gray-100 dark:bg-gray-900 dark:text-white dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600"
                   href="/">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                         stroke="currentColor" class="w-6 h-6">
                        <path stroke-linecap="round" stroke-linejoin="round"
                              d="m2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25"/>
                    </svg>
                    Dashboard
                </a>
            </li>

            <li class="hs-accordion" id="projects-accordion">
                <button type="button"
                        class="hs-accordion-toggle hs-accordion-active:text-blue-600 hs-accordion-active:hover:bg-transparent w-full text-start flex items-center gap-x-3.5 py-2 px-2.5 text-sm text-slate-700 rounded-lg hover:bg-gray-100 dark:bg-gray-800 dark:hover:bg-gray-900 dark:text-slate-400 dark:hover:text-slate-300 dark:hs-accordion-active:text-white dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                         stroke="currentColor" class="w-6 h-6">
                        <path stroke-linecap="round" stroke-linejoin="round"
                              d="M2.25 15a4.5 4.5 0 0 0 4.5 4.5H18a3.75 3.75 0 0 0 1.332-7.257 3 3 0 0 0-3.758-3.848 5.25 5.25 0 0 0-10.233 2.33A4.502 4.502 0 0 0 2.25 15Z"/>
                    </svg>
                    Drive

                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                         stroke="currentColor"
                         class="w-4 h-4  hs-accordion-active:block ms-auto hidden text-gray-600 group-hover:text-gray-500 dark:text-gray-400">
                        <path stroke-linecap="round" stroke-linejoin="round" d="m4.5 15.75 7.5-7.5 7.5 7.5"/>
                    </svg>

                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                         stroke="currentColor"
                         class="w-4 h-4 hs-accordion-active:hidden ms-auto block text-gray-600 group-hover:text-gray-500 dark:text-gray-400">
                        <path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5"/>
                    </svg>
                </button>

                <div id="projects-accordion"
                     class="hs-accordion-content w-full overflow-hidden transition-[height] duration-300 hidden">
                    <ul class="pt-2 ps-2">
                        <li>
                            <a class="flex items-center gap-x-3.5 py-2 px-2.5 text-sm text-slate-700 rounded-lg hover:bg-gray-100 dark:bg-gray-800 dark:text-slate-400 dark:hover:text-slate-300 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600"
                               href="/explorer">
                                Explore
                            </a>
                        </li>
                    </ul>
                </div>
            </li>
        </ul>
    </nav>
    <div class="absolute bottom-0 left-0 p-4 w-full">
        {#if userData !== null}
            <div class="flex justify-between">
                <div class="relative inline-block">
                    <Image
                            src="https://images.unsplash.com/photo-1568602471122-7832951cc4c5?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=facearea&facepad=2&w=300&h=300&q=80"
                            imgClass="inline-block rounded-full"
                            size="h-[3.875rem] w-[3.875rem]"
                            alt="Profile image" />
                </div>

                <div class="flex flex-col justify-center ms-3">
                    <span class="text-sm font-semibold text-gray-700 dark:text-gray-300">Hi, { userData.username }</span>
                    <span class="text-xs text-gray-500 dark:text-gray-400">Welcome back!</span>
                </div>

                <button type="button" title="Log out"
                        class="flex items-center justify-center gap-2.5 ms-auto text-gray-500 hover:text-gray-600 dark:text-gray-400 dark:hover:text-gray-300 rotate-180">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                         stroke="currentColor" class="w-6 h-6">
                        <path stroke-linecap="round" stroke-linejoin="round"
                              d="M8.25 9V5.25A2.25 2.25 0 0 1 10.5 3h6a2.25 2.25 0 0 1 2.25 2.25v13.5A2.25 2.25 0 0 1 16.5 21h-6a2.25 2.25 0 0 1-2.25-2.25V15m-3 0-3-3m0 0 3-3m-3 3H15"/>
                    </svg>
                </button>
            </div>
        {:else}
            <div class="flex items-center justify-center">
                <div class="animate-spin inline-block w-6 h-6 border-[3px] border-current border-t-transparent text-blue-600 rounded-full dark:text-blue-500" role="status" aria-label="loading">
                    <span class="sr-only">Loading...</span>
                </div>
            </div>
        {/if}
    </div>
</div>

<main class="h-[100vh] grid lg:grid-cols-[16rem,auto]">
    <div class="lg:block hidden"></div>
    <div class="p-4">
        <slot/>
    </div>
</main>