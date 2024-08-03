<script lang="ts">
  import { onMount } from "svelte";
  import { navigateTo } from "svelte-router-spa";

  import Dish from "./dish.svelte";
  import { GetDishes } from "../../client/dish";
  import { LoadCart } from "../../client/cart";

  let dishes = [];
  onMount(async () => {
    dishes = await GetDishes();
    let mainButton = LoadCart();
    mainButton.on("click", () => {
      navigateTo("/cart");
      mainButton.hide();
    });
  });
</script>

<div class="grid grid-cols-1 medium-grid-cols-3 grid-gap" id="dishes-grid">
  {#each dishes as dish}
    <Dish bind:dishInst={dish} />
  {/each}
</div>

<style>
  .grid {
    display: grid;
    grid-template-columns: minmax(var(--dish-size), 1fr);
    grid-auto-flow: column;
    gap: 0.5rem;
  }
</style>
