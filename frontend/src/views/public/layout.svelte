<script>
  import { Route } from "svelte-router-spa";

  export let currentRoute;
  const params = {};
</script>

<div>
  <Route {currentRoute} {params} />
</div>
