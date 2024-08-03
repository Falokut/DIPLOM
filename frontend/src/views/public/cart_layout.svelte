<script>
  import { Route } from "svelte-router-spa";
  import TopHeader from "./top_header.svelte";
  import Cart from "./cart.svelte";

  export let currentRoute;
</script>

<div class="app">
  <TopHeader />
  <section class="section">
    <Route {currentRoute}/>
    <Cart></Cart>
  </section>
</div>

