<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { GetAllRestaurants } from "../../client/restaurant";
  import { navigate } from "svelte-routing";

  import { initBackButton } from "@telegram-apps/sdk";
  import AddRestaurant from "./add_restaurant.svelte";
  import Restaurant from "./restaurant.svelte";

  const backButtonRes = initBackButton();
  var backButton = backButtonRes[0];
  let restaurants = [];
  var removeBackButtonListFn;
  onMount(async () => {
    removeBackButtonListFn = backButton.on("click", () => {
      navigate("/admin", { replace: true });
    });

    restaurants = await GetAllRestaurants();
  });
  onDestroy(() => {
    removeBackButtonListFn();
  });
  function remove(id: number) {
    restaurants = restaurants.filter((v) => v.id != id);
  }
</script>

<main>
  <h3>Рестораны</h3>
  {#key restaurants}
    <div class="restaurants">
      {#each restaurants as restaurant}
        <Restaurant restaurant={restaurant} {remove} />
      {/each}
    </div>
  {/key}

  <div class="add-restaurant">
    <AddRestaurant
      OnAdd={(restaurant) => {
        restaurants.push(restaurant);
        restaurants = restaurants;
      }}
    />
  </div>
</main>

<style>
  .restaurants {
    display: flex;
    flex-flow: column;
    width: 90vw;
    background-color: var(--secondary-bg-color);
    border-radius: 5px;
    padding: 10px;
  }

  .add-restaurant {
    display: flex;

    padding: 10px;
    width: 90vw;
    height: auto;
    min-height: 30vh;
  }
</style>
