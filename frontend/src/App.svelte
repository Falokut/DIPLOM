<script>
  import { Router } from "svelte-router-spa";
  import { routes } from "./routes";
  import { onMount } from "svelte";
  
  let tg = Telegram.WebApp;
  onMount(() => {
    tg.enableClosingConfirmation();
    tg.expand();
  });
  
</script>

<Router {routes} />
