<script>
  import { Route } from "svelte-router-spa";
  import TopHeader from "./top_header.svelte";
  import Dishes from "./dishes.svelte";

  export let currentRoute;
</script>

<div class="app">
  <TopHeader />
  <section class="section">
    <Route {currentRoute}/>
    <Dishes></Dishes>
  </section>
</div>
