<script lang="ts">
  import { RenameRestaurant, DeleteRestaurant } from "../../client/restaurant";
  import TextInput from "../components/text_input.svelte";

  export let restaurant = {
    name: "",
    id: 0,
  };
  export let remove = function (id) {};

  async function updateRestaurant() {
    RenameRestaurant(restaurant.name, restaurant.id);
  }

  async function deleteRestaurant() {
    await DeleteRestaurant(restaurant.id);
    remove(restaurant.id);
  }
</script>

<section class="restaurant-container">
  <TextInput bind:value={restaurant.name} />
  <div class="buttons-container">
    <button class="apply-button" on:click={() => updateRestaurant()}>✓</button>
    <button class="remove-button" on:click={() => deleteRestaurant()}>✕</button>
  </div>
</section>

<style>
  .restaurant-container {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-around;
    padding: 10px;
  }

  .buttons-container {
    text-align: center;
    vertical-align: middle;
    display: flex;
    flex-direction: row;
    font-weight: 100;
    align-items: center;
    justify-content: space-around;
  }

  .remove-button {
    background-color: rgb(182, 19, 19);
  }

  .remove-button:hover {
    border: 1px solid var(--tg-theme-bg-color);
  }

  .apply-button {
    background-color: rgb(17, 151, 17);
  }

  .apply-button:hover {
    border: 1px solid var(--tg-theme-bg-color);
  }
</style>
