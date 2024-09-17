<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { navigate } from "svelte-routing";

  import Dish from "./dish.svelte";
  import { GetDishes } from "../../client/dish";
  import { LoadCart } from "../../client/cart";

  let dishes = [];
  let mainButton = LoadCart();
  let removeListFn;

  onMount(async () => {
    dishes = await GetDishes(null);

    removeListFn = mainButton.on("click", () => {
      navigate("/cart");
      mainButton.hide();
    });
  });

  onDestroy(() => {
    removeListFn();
  });
</script>

<div class="grid grid-cols-1 medium-grid-cols-3 grid-gap" id="dishes-grid">
  {#each dishes as dish}
    <Dish bind:dish />
  {/each}
</div>

<style>
  .grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1rem;
  }
</style>
