import { writable } from 'svelte/store';
import {browser} from "$app/environment";

export const authToken = writable(browser ? localStorage.getItem("token") : null);