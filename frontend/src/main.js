import './app.css'
import App from './App.svelte'
import { initBackButton } from "@telegram-apps/sdk";

const backButtonRes = initBackButton();
let backButton = backButtonRes[0];
backButton.show();

const app = new App({
  target: document.getElementById('app'),
})

export function GetBackendBasePath() {
  return import.meta.env.VITE_API_URL + "/api/dish_as_a_service";
}

export default app
