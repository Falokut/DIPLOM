<script>
  import { routes } from "./routes";
  import { onMount } from "svelte";
  import { Router, Route } from "svelte-routing";

  let tg = Telegram.WebApp;
  onMount(() => {
    tg.expand();
  });
</script>

<div id="app">
  <Router>
    {#each routes as route}
      <Route path={route.name} component={route.component}></Route>
    {/each}
  </Router>
</div>