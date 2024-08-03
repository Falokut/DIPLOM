import './app.css'
import App from './App.svelte'
import { initBackButton } from "@telegram-apps/sdk";

const backButtonRes = initBackButton();
let backButton = backButtonRes[0];
backButton.show();

const app = new App({
  target: document.getElementById('app'),
})

export default app
