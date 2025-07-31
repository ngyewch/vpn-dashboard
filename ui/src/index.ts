import {mount} from 'svelte';
import App from './App.svelte';

import '@picocss/pico/css/pico.css';
import './global.css';

const app = mount(App, {
	target: document.getElementById('app')!,
});

export default app;
