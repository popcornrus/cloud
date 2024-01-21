import { writable } from 'svelte/store';
import {browser} from "$app/environment";

export const user = writable(browser ? JSON.parse(localStorage.getItem("user")) : null);