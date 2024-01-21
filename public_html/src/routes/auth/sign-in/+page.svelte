<script>
    import {onMount} from "svelte";
    import Toast from "$lib/components/Toast.svelte";
    import {User} from "$lib/classes/User.js";

    let formInProcess = false;

    let toast = {
        status: null,
        message: null,
    };

    onMount(() => {
        const user = new User();

        const form = document.querySelector('form');
        const email = document.querySelector('#email');
        const password = document.querySelector('#password');

        form.addEventListener('submit', async (e) => {
            e.preventDefault();

            email.disabled = true;
            password.disabled = true;
            formInProcess = true;

            toast = await user.login(email.value, password.value)
            setTimeout(() => {
                if (toast.status === 'success') {
                    window.location.href = '/';
                }
            }, 2500)
        })
    })
</script>

<form action="/auth?/login" method="post" class="shadow rounded-3xl min-w-[25%] max-w-[30%] p-10">
    <h3 class="text-2xl dark:text-white font-bold">Hi, Welcome Back!</h3>
    <div class="my-8 grid gap-y-2">
        <div class="mb-2">
            <label for="email" class="block mb-1 text-gray-400">Email</label>
            <fieldset class="relative">
                <input type="email" id="email" value="root@domain.com"
                       class="peer py-3 px-4 ps-11 block w-full bg-gray-100 border-transparent rounded-lg text-sm focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none dark:bg-gray-700 dark:border-transparent dark:text-gray-400 dark:focus:ring-gray-600"
                       placeholder="Enter email">
                <div class="absolute inset-y-0 start-0 flex items-center pointer-events-none ps-4 peer-disabled:opacity-50 peer-disabled:pointer-events-none">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                         stroke="currentColor" class="flex-shrink-0 w-4 h-4 text-gray-500">
                        <path stroke-linecap="round" stroke-linejoin="round"
                              d="M16.5 12a4.5 4.5 0 1 1-9 0 4.5 4.5 0 0 1 9 0Zm0 0c0 1.657 1.007 3 2.25 3S21 13.657 21 12a9 9 0 1 0-2.636 6.364M16.5 12V8.25"/>
                    </svg>
                </div>
            </fieldset>
        </div>
        <div class="">
            <label for="password" class="block mb-1 text-gray-400">Password</label>
            <fieldset class="relative">
                <input type="password" value="pa$$w0rd!" id="password"
                       class="peer py-3 px-4 ps-11 block w-full bg-gray-100 border-transparent rounded-lg text-sm focus:border-blue-500 focus:ring-blue-500 disabled:opacity-50 disabled:pointer-events-none dark:bg-gray-700 dark:border-transparent dark:text-gray-400 dark:focus:ring-gray-600"
                       placeholder="Enter password">
                <div class="absolute inset-y-0 start-0 flex items-center pointer-events-none ps-4 peer-disabled:opacity-50 peer-disabled:pointer-events-none">
                    <svg class="flex-shrink-0 w-4 h-4 text-gray-500" xmlns="http://www.w3.org/2000/svg" width="24"
                         height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                         stroke-linecap="round" stroke-linejoin="round">
                        <path d="M2 18v3c0 .6.4 1 1 1h4v-3h3v-3h2l1.4-1.4a6.5 6.5 0 1 0-4-4Z"/>
                        <circle cx="16.5" cy="7.5" r=".5"/>
                    </svg>
                </div>
            </fieldset>
        </div>
    </div>

    <div class="flex justify-center flex-wrap gap-y-2">
        {#if !formInProcess}
            <div class="w-full flex justify-center mb-3">
                <button type="submit" id="submit"
                        class="w-full flex py-3 justify-center items-center gap-x-2 font-semibold rounded-lg border border-transparent bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50 disabled:pointer-events-none dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-gray-600">
                    Sign In
                </button>
            </div>
            <div class="w-full text-center flex justify-center mb-3">
                <a href="/auth/reset" class="text-blue-600 underline font-bold">Forgot password?</a>
            </div>
            <div class="w-full text-center">
                <hr class="border-1 mb-3">
                <p class="text-gray-400">Don't have an account? <a href="/auth/sign-up"
                                                                   class="text-blue-600 underline font-bold">Sign Up</a>
                </p>
            </div>
        {:else}
            <div class="animate-spin inline-block w-6 h-6 border-[3px] border-current border-t-transparent text-blue-600 rounded-full dark:text-blue-500"
                 role="status"
                 aria-label="loading">
                <span class="sr-only">Loading...</span>
            </div>
        {/if}
    </div>
</form>

<Toast
        bind:toast={toast}
        visible="{toast.status !== null}"
/>