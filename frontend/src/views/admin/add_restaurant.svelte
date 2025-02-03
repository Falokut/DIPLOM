<script lang="ts">
  import {
    AddRestaurant,
    Restaurant,
  } from "../../client/restaurant";

  export let restaurantName = "";
  export let OnAdd = (Restaurant: Restaurant) => {};

  async function addRestaurant() {
    if (restaurantName == "") {
      return;
    }

    let restaurant = await AddRestaurant(restaurantName);
    if (restaurant.id != undefined) {
      restaurant.name = restaurantName;
      OnAdd(restaurant);
      restaurantName = "";
      return;
    }
  }
</script>

<section class="restaurant">
  <input
    class="input-area restaurant-input"
    type="text"
    bind:value={restaurantName}
  />
  <button class="add-button" on:click={addRestaurant}>Добавить</button>
</section>

<style>
  .restaurant {
    border-radius: 5px;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-around;
    background-color: var(--secondary-bg-color);
    align-items: center;
    margin: auto;
    padding: 5px;
    width: 100vw;
  }
  .restaurant-input {
    border-radius: 3px;
  }

  .add-button {
    display: block;
    border: 1px solid var(--primary-bg-color);
    margin-right: 5px;
  }
</style>
