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
  let baseUrl = import.meta.env.VITE_API_URL
  return baseUrl ? baseUrl : "/api/dish_as_a_service";
  //return "https://falokut.ru/api/dish_as_a_service";
}

export default app
